package tree

import (
	"IG-Parser/core/shared"
	"log"
	"strconv"
	"strings"
)

// Entry delimiters
const TREE_PRINTER_OPEN_BRACE = "{"
const TREE_PRINTER_CLOSE_BRACE = "}"
const TREE_PRINTER_LINEBREAK = "\n"
const TREE_PRINTER_SEPARATOR = "," + TREE_PRINTER_LINEBREAK

// Key name
const TREE_PRINTER_KEY_NAME = "\"name\""
const TREE_PRINTER_KEY_COMPONENT = "\"comp\""
const TREE_PRINTER_KEY_NESTING_LEVEL = "\"level\""
const TREE_PRINTER_KEY_POSITION = "\"pos\""
const TREE_PRINTER_KEY_CHILDREN = "\"children\""
const TREE_PRINTER_KEY_PROPERTIES = "\"prop\""
const TREE_PRINTER_KEY_ANNOTATIONS = "\"anno\""
const TREE_PRINTER_KEY_COMPLEXITY = "\"dov\""

// Separator
const TREE_PRINTER_EQUALS = ": "

// Values
const TREE_PRINTER_VAL_POSITION_BELOW = "\"b\""

// Collection delimiters
const TREE_PRINTER_COLLECTION_OPEN = "["
const TREE_PRINTER_COLLECTION_CLOSE = "]"

/*
Prints JSON output format compatible with tree visualization in D3.
Allows indication for flat printing (nested property tree structure vs. flat listing of properties).
Allows indication for printing of binary trees, as opposed to tree aggregated by logical operators for given component.
Allows indication whether annotations should be included in output (as labels).
Allows indication whether degree of variability should be included in output (as labels).
Allows indication as to whether activation conditions should be moved to the beginning of the visual tree output
Requires specification of nesting level the nodes exists on (Default: 0).
This function is tested in TabularOutputGenerator_test.go, i.e., tests with focus on visual tree output.
*/
func (s Statement) PrintTree(parent *Node, printFlat bool, printBinary bool, includeAnnotations bool, includeDegreeOfVariability bool, moveActivationConditionsToFront bool, nestingLevel int) (strings.Builder, ParsingError) {

	// Default name if statement does not have root node
	rootName := ""

	if includeDegreeOfVariability {
		// Use root node label for complexity output
		rootName = "DoV: " + strconv.Itoa(s.CalculateComplexity().TotalStateComplexity)
	}

	if parent != nil {
		// if it is a nested statement, use component name it nests on as name
		rootName = parent.GetComponentName()
	}

	// Opening tree
	out := strings.Builder{}
	out.WriteString("{")
	out.WriteString(TREE_PRINTER_LINEBREAK)

	// Print statement-level node
	out.WriteString(TREE_PRINTER_KEY_NAME)
	out.WriteString(TREE_PRINTER_EQUALS)
	out.WriteString("\"")
	// Root node name
	out.WriteString(rootName)
	out.WriteString("\"")
	out.WriteString(TREE_PRINTER_SEPARATOR)

	// Append nesting level for every node (includes parent node of potential nested statement)
	out.WriteString(TREE_PRINTER_KEY_NESTING_LEVEL)
	out.WriteString(TREE_PRINTER_EQUALS)
	out.WriteString(strconv.Itoa(nestingLevel))
	out.WriteString(", ")

	// TODO CHECK ON ROOT NODE - Append annotations for root node if activated (and existing)
	if includeAnnotations {
		out.WriteString(parent.appendAnnotations("", false, true))
	}

	// TODO CHECK ON ROOT NODE - Append Degree of Variability
	if includeDegreeOfVariability {
		out.WriteString(parent.appendDegreeOfVariability("", false, true))
	}

	// Line break to separate children visually
	out.WriteString(TREE_PRINTER_LINEBREAK)

	// Indicates whether children have already been added below the top-level string
	childrenPresent := false

	components := []*Node{}

	// Move activation conditions to front if selected
	if moveActivationConditionsToFront {
		components = append(components, s.ActivationConditionSimple, s.ActivationConditionComplex)
	}

	// Add individual nodes (but suppress properties, since those are handled alongside corresponding components)
	components = append(components,
		// Regulative Side
		s.Attributes,
		//s.AttributesPropertySimple,
		//s.AttributesPropertyComplex,
		s.Deontic,
		s.Aim,
		s.DirectObject,
		s.DirectObjectComplex,
		//s.DirectObjectPropertySimple,
		//s.DirectObjectPropertyComplex,
		s.IndirectObject,
		s.IndirectObjectComplex,
		//s.IndirectObjectPropertySimple,
		//s.IndirectObjectPropertyComplex,

		// Constitutive Side
		s.ConstitutedEntity,
		//s.ConstitutedEntityPropertySimple,
		//s.ConstitutedEntityPropertyComplex,
		s.Modal,
		s.ConstitutiveFunction,
		s.ConstitutingProperties,
		s.ConstitutingPropertiesComplex,
		//s.ConstitutingPropertiesPropertySimple,
		//s.ConstitutingPropertiesPropertyComplex,
	)

	// Move activation conditions in the idiomatic ADIBCO position
	if !moveActivationConditionsToFront {
		components = append(components, s.ActivationConditionSimple, s.ActivationConditionComplex)
	}

	// Remainder of components
	components = append(components,
		// Shared elements
		s.ExecutionConstraintSimple,
		s.ExecutionConstraintComplex,

		s.OrElse)

	for _, v := range components {
		// only print components that have content, and for property components (the one whose name ends on ",p"),
		// only print if they have complex children - simple values are printed alongside first-order component

		if v != nil && (!strings.HasSuffix(v.GetComponentName(), PROPERTY_SYNTAX_SUFFIX) ||
			(strings.HasSuffix(v.GetComponentName(), PROPERTY_SYNTAX_SUFFIX) && !v.HasPrimitiveEntry())) {
			prepend := ""
			if childrenPresent {
				// If in next round (not first entry), prepend separator in case output is produced
				prepend = TREE_PRINTER_SEPARATOR
			}
			// Generate actual entry
			componentString, err := v.PrintNodeTree(s, printFlat, printBinary, includeAnnotations, includeDegreeOfVariability, moveActivationConditionsToFront, nestingLevel)
			if err.ErrorCode != TREE_NO_ERROR {
				// Special case of NodeError being embedded in ParsingError
				return strings.Builder{}, ParsingError{ErrorCode: PARSING_ERROR_EMBEDDED_NODE_ERROR, ErrorMessage: err.ErrorMessage}
			}

			Println("Output for " + v.GetComponentName() + ": " + componentString)
			if !childrenPresent && componentString != "" {
				// Print children prefix if components are present
				out.WriteString(TREE_PRINTER_KEY_CHILDREN)
				out.WriteString(TREE_PRINTER_EQUALS)
				out.WriteString(TREE_PRINTER_COLLECTION_OPEN)
				out.WriteString(TREE_PRINTER_LINEBREAK)
				childrenPresent = true
			}
			// Add the actual content
			out.WriteString(prepend)
			out.WriteString(componentString)
		}
	}
	// Close children
	if childrenPresent {
		out.WriteString(TREE_PRINTER_LINEBREAK)
		out.WriteString(TREE_PRINTER_COLLECTION_CLOSE)
	}
	// Close entire tree
	out.WriteString(TREE_PRINTER_LINEBREAK)
	out.WriteString(TREE_PRINTER_CLOSE_BRACE)

	return out, ParsingError{ErrorCode: PARSING_NO_ERROR}
}

/*
Returns JSON output for visual tree rendering of individual nodes using D3.
Allows indication for flat printing (nested property tree structure vs. flat listing of properties).
Allows indication for printing of binary trees, as opposed to tree aggregated by logical operators for given component.
Allows indication for including annotations in output (as labels).
Allows indication of Degree of Variability in output (as labels).
Allows indication as to whether activation conditions should be moved to the beginning of the visual tree output
Requires specification of nesting level for node exists on (Default: 0).
*/
func (n *Node) PrintNodeTree(stmt Statement, printFlat bool, printBinary bool, includeAnnotations bool, includeDegreeOfVariability bool, moveActivationConditionsToFront bool, nestingLevel int) (string, NodeError) {
	out := strings.Builder{}

	if !n.IsNil() && !n.IsEmptyOrNilNode() {
		if n.HasPrimitiveEntry() || n.IsCombination() {

			// Indicates whether full entry should be printed
			printFullEntry := false

			// If the entry is not a statement but either leaf or combination
			if n.HasPrimitiveEntry() {

				// Produce output for simple entry
				out.WriteString(TREE_PRINTER_OPEN_BRACE)
				out.WriteString(TREE_PRINTER_KEY_NAME)
				out.WriteString(TREE_PRINTER_EQUALS)
				// Actual content (including escaping particular symbols)
				out.WriteString("\"")
				out.WriteString(shared.EscapeSymbolsForExport(n.Entry.(string)))
				out.WriteString("\"")

				// Ensure that entry is closed
				printFullEntry = true
			} else if n.IsCombination() {
				// Produce output for combination

				if printBinary {
					// Fall back to full entry parsing either way - resolving full binary tree structure
					printFullEntry = true
				} else {
					// If non-binary, print only leaf entries linked via same logical operator on same component
					// without considering logical operators in output
					if n.Parent != nil && n.LogicalOperator == n.Parent.LogicalOperator &&
						n.GetComponentName() == n.Parent.GetComponentName() {

						// Print left side
						outTmpL, err := n.Left.PrintNodeTree(stmt, printFlat, printBinary, includeAnnotations, includeDegreeOfVariability, moveActivationConditionsToFront, nestingLevel)
						if err.ErrorCode != TREE_NO_ERROR {
							return "", err
						}
						// Append if successful parsing
						out.WriteString(outTmpL)

						// Append separator to collapsed entries (i.e., on same level)
						out.WriteString(", ")

						// Print right side
						outTmpR, err := n.Right.PrintNodeTree(stmt, printFlat, printBinary, includeAnnotations, includeDegreeOfVariability, moveActivationConditionsToFront, nestingLevel)
						if err.ErrorCode != TREE_NO_ERROR {
							return "", err
						}
						// Append if successful parsing
						out.WriteString(outTmpR)

						// Suppress printing of closing parts of entry, since further nodes of same operator on same component may be appended
						printFullEntry = false
					} else {
						// Fall back to print full entry if logical operators or components don't match
						printFullEntry = true
					}
				}
				// Prints full entry as binary tree element (either applies if binary tree structure is activated,
				// or if no nested logical operators for given component were detected (e.g., multiple nested AND linkages)
				if printFullEntry {
					// Create logical node, two children, and delegate node entry generation to children
					out.WriteString(TREE_PRINTER_OPEN_BRACE)
					out.WriteString(TREE_PRINTER_KEY_NAME)
					out.WriteString(TREE_PRINTER_EQUALS)
					// Logical operator
					out.WriteString("\"")
					out.WriteString(n.LogicalOperator)
					out.WriteString("\"")
					out.WriteString(TREE_PRINTER_SEPARATOR)
					// Children
					out.WriteString(TREE_PRINTER_KEY_CHILDREN)
					out.WriteString(TREE_PRINTER_EQUALS)
					out.WriteString(TREE_PRINTER_COLLECTION_OPEN)

					// Left child
					outTmpL, err := n.Left.PrintNodeTree(stmt, printFlat, printBinary, includeAnnotations, includeDegreeOfVariability, moveActivationConditionsToFront, nestingLevel)
					if err.ErrorCode != TREE_NO_ERROR {
						return "", err
					}
					// Append if successful parsing
					out.WriteString(outTmpL)

					// Add separator
					out.WriteString(TREE_PRINTER_SEPARATOR)

					// Right child
					outTmpR, err := n.Right.PrintNodeTree(stmt, printFlat, printBinary, includeAnnotations, includeDegreeOfVariability, moveActivationConditionsToFront, nestingLevel)
					if err.ErrorCode != TREE_NO_ERROR {
						return "", err
					}
					// Append if successful parsing
					out.WriteString(outTmpR)

					// Closing collection
					out.WriteString(TREE_PRINTER_COLLECTION_CLOSE)
				}
			}

			// Continue and close full entry (with component, property and annotation information) only if entry is complete,
			// not if branches of logical operators are collapsed
			if printFullEntry {
				// Append component name as link label for any entry
				out.WriteString(", ")
				out.WriteString(TREE_PRINTER_KEY_COMPONENT)
				out.WriteString(TREE_PRINTER_EQUALS)
				out.WriteString("\"")
				out.WriteString(n.GetComponentName())
				out.WriteString("\"")

				// Append nesting level for every node
				out.WriteString(", ")
				out.WriteString(TREE_PRINTER_KEY_NESTING_LEVEL)
				out.WriteString(TREE_PRINTER_EQUALS)
				out.WriteString(strconv.Itoa(nestingLevel))

				// Print private properties
				outTmp, err := n.appendPropertyNodes("", stmt, printFlat, printBinary, includeAnnotations, includeDegreeOfVariability, moveActivationConditionsToFront, nestingLevel)
				if err.ErrorCode != TREE_NO_ERROR {
					return "", err
				}

				// Append annotations (if existing)
				if includeAnnotations {
					outTmp = n.appendAnnotations(outTmp, true, false)
				}

				// Append complexity
				if includeDegreeOfVariability {
					outTmp = n.appendDegreeOfVariability(outTmp, true, false)
				}

				out.WriteString(outTmp)

				// Close entry
				out.WriteString(TREE_PRINTER_CLOSE_BRACE)
			}
		} else {
			// Produce output for nested statement
			outTmp, err := n.Entry.(Statement).PrintTree(n, printFlat, printBinary, includeAnnotations, includeDegreeOfVariability, moveActivationConditionsToFront, nestingLevel+1)
			if err.ErrorCode != PARSING_NO_ERROR { // Important check on return value - different from all others in the function
				// Special case of NodeError embedding a ParsingError produced in nested invocation.
				return "", NodeError{ErrorCode: TREE_ERROR_EMBEDDED_PARSING_ERROR, ErrorMessage: err.ErrorMessage}
			}
			out.WriteString(outTmp.String())
		}
	}
	return out.String(), NodeError{ErrorCode: TREE_NO_ERROR}
}

/*
Appends shared and private nodes to D3-consumable JSON output string based on related properties, as well as own private nodes.
The shared and private property nodes are combined in the order "shared, private".
Note: In flat output mode only primitive private properties are included in the rendered output.
Flat output implies the printing of private properties as labels for component nodes, rather than an own node hierarchy.
Allows indication for printing of binary trees, as opposed to tree aggregated by logical operators for given component.
Includes indication whether annotations are to be included in output.
Includes indication whether Degree of Variability is to be included in output.
Allows indication as to whether activation conditions should be moved to the beginning of the visual tree output
Requires specification of nesting level the property node lies on (Lowest level: 0)
*/
func (n *Node) appendPropertyNodes(stringToPrepend string, stmt Statement, printFlat bool, printBinary bool, includeAnnotations bool, includeDegreeOfVariability bool, moveActivationConditionsToFront bool, nestingLevel int) (string, NodeError) {

	stringToAppendTo := strings.Builder{}
	stringToAppendTo.WriteString(stringToPrepend)

	// Append potential private and shared property nodes under the condition that those nodes are leaf nodes, or if flat printing is activated
	if n != nil && (n.IsLeafNode() || printFlat) &&
		// Check for shared and private properties
		len(stmt.GetPropertyComponent(n, true)) > 0 || (len(n.PrivateNodeLinks) > 0 && n.PrivateNodeLinks[0] != nil) {

		// Retrieve relevant component property and combine with private nodes before printing
		allNodes := stmt.GetPropertyComponent(n, true)
		Println("Properties associated with component node "+n.GetComponentName()+":", allNodes)
		includeAllNodes := true
		// Test whether shared property nodes exist
		if len(allNodes) == 0 || (len(allNodes) > 0 && allNodes[0] == nil) {
			includeAllNodes = false
		}
		Println("Append nodes are property nodes:", includeAllNodes)

		// Check whether private nodes are populated
		if len(n.PrivateNodeLinks) > 0 && n.PrivateNodeLinks[0] != nil {

			mergedNodes := n.PrivateNodeLinks[0]
			// Infer AND-linkage of private nodes, de facto forming tree structure - this may be decomposed later depending on flat printing setting
			if len(n.PrivateNodeLinks) > 1 {
				// Start with second node if there are actually multiple, and merge using implicit between-component AND linkage
				for _, v := range n.PrivateNodeLinks[1:] {
					err := NodeError{}
					mergedNodes, err = Combine(mergedNodes, v, SAND_BETWEEN_COMPONENTS)
					if err.ErrorCode != TREE_NO_ERROR {
						log.Println("Invalid merge of private nodes when combining private nodes (this should never happen). Nodes to be merged:", mergedNodes, "and", v)
						return "", err
					}
				}
			}

			if includeAllNodes {
				// Append private nodes to shared nodes
				allNodes = append(allNodes, mergedNodes)
			} else {
				// Override potential shared nodes
				allNodes = []*Node{mergedNodes}
			}
		}

		Println("Property nodes to process:", allNodes)

		// Only add output if properties exist
		if len(allNodes) > 0 && allNodes[0] != nil {

			// keeps track whether any element has been extracted
			elementPrinted := false

			// Add individual items
			for _, privateNode := range allNodes {

				// Initiate entry structure
				if !elementPrinted {
					// Start general output for property only if nothing is printed yet
					if printFlat {
						// Initiate flat output for properties
						stringToAppendTo.WriteString(", ")
						stringToAppendTo.WriteString(TREE_PRINTER_KEY_PROPERTIES)
						stringToAppendTo.WriteString(TREE_PRINTER_EQUALS)
					} else {
						// Add position information to ensure offset printing of leading node content (since it is followed by nested property structure)
						stringToAppendTo.WriteString(", ")
						stringToAppendTo.WriteString(TREE_PRINTER_KEY_POSITION)
						stringToAppendTo.WriteString(TREE_PRINTER_EQUALS)
						stringToAppendTo.WriteString(TREE_PRINTER_VAL_POSITION_BELOW)
						// Initiate tree structure for tree output of properties
						stringToAppendTo.WriteString(", ")
						stringToAppendTo.WriteString(TREE_PRINTER_KEY_CHILDREN)
						stringToAppendTo.WriteString(TREE_PRINTER_EQUALS)
						stringToAppendTo.WriteString(TREE_PRINTER_COLLECTION_OPEN)
					}
				}

				// Add separators, or open new entry if needed
				if elementPrinted {
					// Add separator if element has been added
					stringToAppendTo.WriteString(", ")
				} else if printFlat {
					// Prepend quotation
					stringToAppendTo.WriteString("\"")
				}

				// Print per-property entry
				if printFlat {
					// Consider distinct treatment for complex or primitive properties in the case of flat printing
					if privateNode.IsCombination() {
						// Decompose and print combinations
						nodes := privateNode.GetLeafNodes(false)
						entryAdded := false
						// Check outer element
						for _, v1 := range nodes {
							// Check inner element
							for _, v := range v1 {
								// Add separator if previous entry exists
								if entryAdded {
									stringToAppendTo.WriteString(", ")
								}
								// Append each entry individually
								stringToAppendTo.WriteString(shared.EscapeSymbolsForExport(v.Entry.(string)))
								entryAdded = true
							}
						}
					} else if !privateNode.HasPrimitiveEntry() {
						// Embedded statement (is printed as flat string, e.g., A: actor I: action, Cac: context)
						stringToAppendTo.WriteString(privateNode.Entry.(Statement).StringFlatStatement(true))
					} else {
						// Primitive properties
						stringToAppendTo.WriteString(shared.EscapeSymbolsForExport(privateNode.Entry.(string)))
					}
				} else {
					// If no flat printing, append complete nested tree structure (property tree)
					stringTmp, err := privateNode.PrintNodeTree(stmt, printFlat, printBinary, includeAnnotations, includeDegreeOfVariability, moveActivationConditionsToFront, nestingLevel)
					if err.ErrorCode != TREE_NO_ERROR {
						return "", err
					}
					stringToAppendTo.WriteString(stringTmp)
				}
				// Mark if initial item has been performed
				elementPrinted = true
			}

			// Handle termination of entries
			if elementPrinted {
				if printFlat {
					// Close flat entry
					stringToAppendTo.WriteString("\"")
				} else {
					// Close collection
					stringToAppendTo.WriteString(TREE_PRINTER_COLLECTION_CLOSE)
				}
			}
		}
	}

	// Return extended string
	return stringToAppendTo.String(), NodeError{ErrorCode: TREE_NO_ERROR}
}

/*
Appends potentially existing annotations to node-specific output.
Input is the string to be appended to (stringToAppendTo), as well as a parameter indicating whether
termination separator (", ") should be added (either prepended, appended, or both) if annotations are added.
*/
func (n *Node) appendAnnotations(stringToPrepend string, prependSeparator bool, appendSeparator bool) string {

	stringToAppendTo := strings.Builder{}
	stringToAppendTo.WriteString(stringToPrepend)

	// Append potential annotations (while replacing specific conflicting symbols)
	if n != nil && n.GetAnnotations() != nil {
		if prependSeparator {
			stringToAppendTo.WriteString(", ")
		}
		stringToAppendTo.WriteString(TREE_PRINTER_KEY_ANNOTATIONS)
		stringToAppendTo.WriteString(TREE_PRINTER_EQUALS)
		stringToAppendTo.WriteString("\"")
		stringToAppendTo.WriteString(shared.EscapeSymbolsForExport(n.GetAnnotations().(string)))
		stringToAppendTo.WriteString("\"")
		if appendSeparator {
			stringToAppendTo.WriteString(", ")
		}
	}
	// Return potentially extended string
	return stringToAppendTo.String()
}

/*
Appends Degree of Variability metric to node-specific output.
Input is the string to be appended to (stringToAppendTo), as well as a parameter indicating whether
termination separator (", ") should be added (either prepended, appended, or both) if annotations are added.
*/
func (n *Node) appendDegreeOfVariability(stringToPrepend string, prependSeparator bool, appendSeparator bool) string {

	stringToAppendTo := strings.Builder{}
	stringToAppendTo.WriteString(stringToPrepend)

	// Append potential complexity output (while replacing specific conflicting symbols)
	retVal, err := n.CalculateStateComplexity()
	if n != nil && err.ErrorCode == TREE_NO_ERROR {
		if prependSeparator {
			stringToAppendTo.WriteString(", ")
		}
		stringToAppendTo.WriteString(TREE_PRINTER_KEY_COMPLEXITY)
		stringToAppendTo.WriteString(TREE_PRINTER_EQUALS)
		stringToAppendTo.WriteString("\"")
		stringToAppendTo.WriteString(shared.EscapeSymbolsForExport(strconv.Itoa(retVal)))
		stringToAppendTo.WriteString("\"")
		if appendSeparator {
			stringToAppendTo.WriteString(", ")
		}
	}
	// Return potentially extended string
	return stringToAppendTo.String()
}