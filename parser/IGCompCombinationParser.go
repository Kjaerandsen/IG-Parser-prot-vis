package parser

import (
	"IG-Parser/tree"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func main() {

	// Dummy input for convenient selection
	input := ""
	// Very simple input
	input = "(inspect and [OR] party)"
	// Simple input
	//input = "((inspect and [OR] party) [AND] sing)"
	// Proper complex input
	//input = "((inspect and [OR] party) [AND] ((review [XOR] muse) [AND] pray))"
	// Imbalanced parentheses will fail (whatever the direction)
	//input = "((inspect and [OR] party) [AND] ((review [XOR] muse) [AND] pray)"
	//input = "(inspect and [OR] party) [AND] ((review [XOR] muse) [AND] pray))"
	// Missing outer parentheses will lead to processing of right side only
	//input = "(inspect and [OR] party) [AND] ((review [XOR] muse) [AND] pray)"
	// Invalid operators lead to ignoring element in processing
	//input = "((inspect and [OR] party) [AND] ((review [XOR] muse) AND pray))"
	// Invalid operators even apply to parenthesized combinations --> FIX and make leaf node
	//input = "((inspect and OR party) [AND] ((review [XOR] muse) [AND] pray))"
	// Invalid operators without parentheses around operands will be treated as leaf node
	//input = "(inspect and OR party [AND] ((review [XOR] muse) [AND] pray))"
	// Excessive parentheses are acceptable (will be flattened in parsing process)
	input = "((((inspect and [OR] party) [AND] ((review [XOR] muse) [AND] pray))))"
	// Non-combinations are ignored (i.e., parentheses within components)
	input = "(French (and) [AND] (accredited certifying agents))"
	// Repeated AND combinations (e.g., "expr1 AND expr2 AND expr3") are collapsed into nested structures (e.g., "(expr1 AND expr2) AND expr3")
	input = "(French (and) [AND] (certified production and [XOR] handling operations) and [AND] (accredited certifying agents))"
	// Repeated logical operators will break (empty leaf value in the AND case)
	//input = "(French (and) [AND] [AND] (certified production and [XOR] handling operations) and [AND] (accredited certifying agents))"
	// Repeated logical operators will break (multiple non-AND operators (or mix thereof))
	//input = "(French (and) [AND] [OR] (certified production and [XOR] handling operations) and [AND] (accredited certifying agents))"

	// Create root node
	node := tree.Node{}
	// Parse provided expression
	_, modifiedInput, _ := ParseDepth(input, &node)
	// Print resulting tree
	fmt.Println("Final tree: \n" + node.String())
	fmt.Println("Corresponding (potentially modified) input string: " + modifiedInput)

	fmt.Println(node.Stringify())
}

/*
Parses combinations in string. The syntactic form of input is:
"( leftSide [OPERATOR] rightSide )", where [OPERATOR] is one
of the logical operators [AND], [OR], [XOR] (including brackets),
and left and right side are either text or combinations themselves.
For [AND] operators, an arbitrary number of expressions can be combined;
in this case the function will decompose those into nested structures
(e.g., expanding "( expr1 [AND] expr2 [AND] expr3 )" into
"(( expr1 [AND] expr2 ) [AND] expr3)"), with precedence for left combinations.
Note that expressions are trimmed prior storing in tree structure.
The parsing further supports shared values outside of the combination (e.g.,
'(shared left value (left element [AND] right element) shared right value)',
and returns those as part of the node that holds the logical operator.

Hint: Call Stringify() on the returned node to reconstruct string

The function returns
- a node tree of the structure, as well as
- the potentially modified input string corresponding to the node tree
  Note: Shared elements are stripped from the modified output string (but
  included in the node instance

Note:
- The entire expression must be surrounded with parentheses, else only
the right-most outer combination (and combinations nested therein) is parsed.
- Parsing checks for matching parentheses and stops otherwise
- Invalid combinations (e.g., missing logical operator) are discarded
in the processing.
 */
func ParseDepth(input string, node *tree.Node) (*tree.Node, string, tree.ParsingError) {
	// Return detected combinations, alongside potentially modified input string
	combinations, input, err := detectCombinations(input)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		return node, input, err
	}

	/*if len(combinations) == 0 {
		fmt.Println("No combinations detected.")
	} else {
		//TODO to be reviewed if issues arise with unwanted elements; else remove
		/*ct := 0
		toBeDeleted := []int{}
		for k, v := range combinations {
			if v.Complete {
				ct++
			} else {
				toBeDeleted = append(toBeDeleted, k)
			}
		}
		errorMsg := ""
		if len(toBeDeleted) > 0 {
			errorMsg = ", with one partial " + strconv.Itoa(len(toBeDeleted)) + " to be removed"
		}
		fmt.Println(strconv.Itoa(ct) + " valid combination detected" + errorMsg + "!")*/


		// Clean up invalid entries
		/*i := 0
		for i < len(toBeDeleted) {
			delete(combinations, toBeDeleted[i])
			i++
		}
	}*/

	// if no valid combinations are left, the parsing is finished
	if len(combinations) == 0 {
		node.Entry = input
		return node, input, tree.ParsingError{ErrorCode: tree.PARSING_NO_COMBINATIONS,
			ErrorMessage: "The input does not contain combinations"}
	}

	// Now the parsing of nested combinations starts
	fmt.Print("STARTING TREE CONSTRUCTION: Detected combinations: ")
	fmt.Println(combinations)

	// Depth first
	v := 0
	// Shared content framing combinations (e.g., "(shared info (left [AND] right))").
	sharedLeft := ""
	sharedRight := ""
	for { // infinite loop - breaks out eventually
		fmt.Println("Level " + strconv.Itoa(v))
		if _, ok := combinations[v]; ok {
			fmt.Print("Combinations on level " + strconv.Itoa(v) + ": ")
			fmt.Println(combinations[v])

			// Test for shared values
			if combinations[v].Operator == 0 &&
				combinations[v].OperatorVal == "" &&
				!combinations[v].Complete {
				fmt.Println("Detected shared string: ")
				// Must contain embedded combination (e.g., "(left shared info (left [AND] right) right shared info)").

				// Extract potential left shared element (without leading bracket)
				leftShared := input[combinations[v].Left:combinations[v+1].Left-1]
				leftShared = strings.Trim(leftShared, " ")
				// Extract potential right shared element (without trailing brackets)
				rightShared := input[combinations[v+1].Right+1:combinations[v].Right]
				rightShared = strings.Trim(rightShared, " ")
				// Store in shared variables
				if len(leftShared) > 0 {
					if sharedLeft != "" {
						// Append
						sharedLeft += INHERITANCE_DELIMITER + leftShared
					} else {
						// Overwrite
						sharedLeft = leftShared
					}
					fmt.Println("Left shared: " + sharedLeft)
				}
				if len(rightShared) > 0 {
					if sharedRight != "" {
						// Append
						sharedRight += INHERITANCE_DELIMITER + rightShared
					} else {
						// Overwrite
						sharedRight = rightShared
					}
					fmt.Println("Right shared: " + sharedRight)
				}
				// Move to next round by deleting entry
				delete(combinations, v)
				fmt.Println("Moving to next round after shared string processing ...")
				continue
			}
			left := input[combinations[v].Left:combinations[v].Operator]
			right := input[combinations[v].Operator+len(combinations[v].OperatorVal)+2:combinations[v].Right]
			fmt.Println("==Left value: " +  left)
			fmt.Println("==Operator: " + combinations[v].OperatorVal)
			fmt.Println("==Right value: " +  right)
			// Scan through top level only and break out afterwards ...

			// Assign left shared value if existing
			if len(sharedLeft) > 0 {
				fmt.Println("Assigning left shared element '" + sharedLeft + "'.")
				node.SharedLeft = sharedLeft
				// Reset shared
				sharedLeft = ""
			}
			// Assign shared value if existing
			if len(sharedRight) > 0 {
				fmt.Println("Assigning right shared element '" + sharedRight + "'.")
				node.SharedRight = sharedRight
				// Reset shared
				sharedRight = ""
			}

			// Assign logical operator
			node.LogicalOperator = combinations[v].OperatorVal

			fmt.Println("Tree after adding logical operator: " + node.String())


			// Left side (potentially modifying input string)
			leftCombos, left, err := detectCombinations(left)
			if err.ErrorCode != tree.PARSING_NO_ERROR && err.ErrorCode != tree.PARSING_ERROR_IGNORED_ELEMENTS {
				log.Println("Error when parsing left side: " + err.ErrorMessage)
				return node, input, err
			}
			if err.ErrorCode != tree.PARSING_ERROR_IGNORED_ELEMENTS {
				log.Print("Warning: Discarded elements during deep left parsing: " + tree.PrintArray(err.ErrorIgnoredElements))
			}
			// Assign nested nodes either way
			leftNode := tree.Node{}
			// Link both nodes
			node.Left = &leftNode
			leftNode.Parent = node

			if len(leftCombos) == 0 {
				// Trim content first
				left = strings.Trim(left, " ")
				if left != "" {
					// If no combinations exist, assign as left leaf
					fmt.Println("Found leaf on left side: " + left)
					node.Left.Entry = left
				} else {
					msg := "Empty leaf value on left side: " + left +
						" (Corresponding right value and operator: " + right + "; " + node.LogicalOperator +
						");\n processed expression: " + input
					log.Println(msg)
					return nil, input, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_EMPTY_LEAF,
						ErrorMessage: msg}
				}
			} else {
				// If combinations exist, delegate
				fmt.Println("Go deep on left side: " +  left)
				_, left, err := ParseDepth(left, &leftNode)
				if err.ErrorCode != tree.PARSING_NO_ERROR {
					return nil, left, err
				}

				// Check for inheriting shared elements on AND nodes
				inheritSharedElements(&leftNode)

				fmt.Println("Tree after processing left deep: " + node.String())
			}

			// Right side (potentially modifying input string)
			rightCombos, right, err := detectCombinations(right)
			if err.ErrorCode != tree.PARSING_NO_ERROR && err.ErrorCode != tree.PARSING_ERROR_IGNORED_ELEMENTS {
				log.Println("Error when parsing right side: " + err.ErrorMessage)
				return node, input, err
			}
			if err.ErrorCode != tree.PARSING_ERROR_IGNORED_ELEMENTS {
				log.Print("Warning: Discarded elements during deep right parsing: " + tree.PrintArray(err.ErrorIgnoredElements))
			}

			// Assign nested nodes either way
			rightNode := tree.Node{}
			// Link both nodes
			node.Right = &rightNode
			rightNode.Parent = node
			if len(rightCombos) == 0 {
				// Trim content first
				right = strings.Trim(right, " ")
				if right != "" {
					// If no combinations exist, assign as left leaf
					fmt.Println("Found leaf on right side: " + right)
					node.Right.Entry = right
				} else {
					msg := "Empty leaf value on right side: " + right + " (Corresponding left value and operator: " + left + "; " + node.LogicalOperator +
						");\n processed expression: " + input
					log.Println(msg)
					return nil, input, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_EMPTY_LEAF,
						ErrorMessage: msg}
				}
			} else {
				// If combinations exist, delegate
				fmt.Println("Go deep on right side: " +  right)
				_, right, err := ParseDepth(right, node.Right)
				if err.ErrorCode != tree.PARSING_NO_ERROR {
					return nil, right, err
				}

				// Check for inheriting shared elements on AND nodes
				inheritSharedElements(&rightNode)

				fmt.Println("Tree after processing right deep: " + node.String())
			}
			// break out if combination has been found and processed - must be top-level combination
			return node, "(" + left + " [" + node.LogicalOperator + "] " + right + ")", tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
		} else {
			fmt.Println("==No combination for key/level " + strconv.Itoa(v))
		}
		v++
	}

	//fmt.Println("Should not really reach here; probably empty node: " + node.String())
	return node, input, tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
}

/*
This function detects levels of combinations present in the expression
and returns the boundary indices as well logical operator (where present).
It further returns the input string in potentially modified form to reflect
changes performed during processing.
To signal incomplete combinations, it contains a Complete flag that signals
completeness for further postprocessing.
Note: This function does not extract all combinations present in the expression,
since combinations on the same level will not be detected, but overwritten.
In essence, the function provides the depth of the nesting in the expression.

Default syntactic form of input: "( leftSide [OPERATOR] rightSide )", where
[OPERATOR] is one of the logical operators (including brackets), and left
and right side are either text or combinations themselves.
 */
func detectCombinations(expression string) (map[int]tree.Boundaries, string, tree.ParsingError) {

	// Tracks current parsing level
	level := 0

	// Parentheses count to check for balance
	parCount := 0

	// Initial test run for parentheses
	for i, letter := range expression {

		switch string(letter) {
		case "(":
			parCount++
		case ")":
			parCount--
		}
		i++
	}
	if parCount != 0 {
		msg := "Uneven number of parentheses (positive --> too many left; negative --> too many right): " + strconv.Itoa(parCount)
		log.Println(msg)
		return nil, expression, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_IMBALANCED_PARENTHESES, ErrorMessage: msg}
	}
	// Passed parentheses count

	// Map of mode states across levels (to recover during parsing)
	modeMap := make(map[int]string)

	// Maintain map of boundaries across different levels // note: only records one entry per level; good for top-level identification
	levelMap := make(map[int]tree.Boundaries)

	// Collection of found operators (with operator as key, followed by level, and count as value)
	foundOperators := make(map[string]map[int]int)

	fmt.Println("Testing expression " + expression)
	for i, letter := range expression {

		switch string(letter) {
		case "(":
			// Increase level
			level++
			fmt.Println("Expression start detected (Level " + strconv.Itoa(level) + ")")
			// Configure mode
			modeMap[level] = tree.PARSING_MODE_LEFT
			fmt.Println("Mode: " + modeMap[level])
			//Test of existing entries
			if _, ok := levelMap[level]; ok {
				fmt.Println("Key already defined - not added")
			} else {
				// Store index reference (incremented with 1 to avoid left parenthesis - balanced)
				levelMap[level] = tree.Boundaries{Left: i + 1, Complete: false}
			}
			// Count parentheses to detect uneven matching
			parCount++
		case ")":
			fmt.Println("Expression end detected (Level " + strconv.Itoa(level) + ")")
			// Check if there are repetitions
			for k, _ := range foundOperators { // key: operator, value:
				for k2, v2 := range foundOperators[k] { // key: level, value: number of occurrences
					if v2 > 1 {
						log.Println("Found " + strconv.Itoa(v2) + " occurrences of operator " + k + " on level " + strconv.Itoa(k2) + " in map.")
					}
				}

			}
			// Store index reference
			fmt.Println("Level before saving: " + strconv.Itoa(level))
			if b, ok := levelMap[level]; ok {
				b.Right = i
				levelMap[level] = b
				fmt.Println("Level map for level " + strconv.Itoa(level) + " (after adding right value): " + b.String())
				// Test whether indices are identical or immediately following - suggesting gaps in values
				if (b.Left == b.Operator) || ((b.Operator + len(b.OperatorVal) + 2) == b.Right) {
					msg := "Input contains invalid combination expression in the range '" + expression[b.Left:b.Right] + "'."
					return levelMap, expression, tree.ParsingError{ErrorCode: tree.PARSING_INVALID_COMBINATION,
						ErrorMessage: msg}
				}
				if b.Left != 0 && b.OperatorVal != "" {
					// Update complete marker
					b := levelMap[level]
					b.Complete = true
					levelMap[level] = b
				} else {
					fmt.Println("Detected end, but combination incomplete (Missing operator or left parenthesis). " +
						"Discarding combination. (Input: '" + expression + "')")
					fmt.Println(levelMap[level])
					// Retrieve higher level
					_, ok := levelMap[level+1]
					if !ok {
						fmt.Print("No nested complete combination, so deleted this incomplete one: ")
						fmt.Println(levelMap[level])
						// If no higher-level exists (within the lower level), delete this incomplete entry
						delete(levelMap, level)
					} else {
						// else retain
						fmt.Println("Nested complete combination, so incomplete higher-level combination is retained.")
					}
				}
			}

			// Configure mode
			modeMap[level] = tree.PARSING_MODE_OUTSIDE_EXPRESSION
			fmt.Println("Mode: " + modeMap[level])

			// Reset operator count for given level
			for op := range foundOperators {
				fmt.Println("Deleting operator " + op + " for level " + strconv.Itoa(level))
				delete(foundOperators[op], level)
			}

			// Reduce level
			level--
			fmt.Println("Moving back to level " + strconv.Itoa(level) + ", Mode: " + modeMap[level])
			// Count parentheses to detect uneven matching
			parCount--
		case "[":
			//fmt.Println("Checking for logical operator ... " + expression[i:i+5])
			foundOperator := ""
			switch expression[i:i+5] {
			case tree.AND_BRACKETS:
				fmt.Println("Detected " + tree.AND_BRACKETS)
				foundOperator = tree.AND
			case tree.XOR_BRACKETS:
				fmt.Println("Detected " + tree.XOR_BRACKETS)
				foundOperator = tree.XOR
			}
			// Separately test for OR due to differing length
			if foundOperator == "" && expression[i:i+4] == tree.OR_BRACKETS {
				fmt.Println("Detected " + tree.OR_BRACKETS)
				foundOperator = tree.OR
			}
			if foundOperator != "" {

				fmt.Println("Found logical operator " + foundOperator + " on level " + strconv.Itoa(level))

				// Check whether the logical operator is immediately adjacent to left parenthesis (e.g., ... ([AND] ... - invalid combination
				if levelMap[level].Left == i {
					msg := "Input contains invalid combination expression in the range '" + expression[levelMap[level].Left:] + "'."
					log.Println(msg)
					return levelMap, expression, tree.ParsingError{ErrorCode: tree.PARSING_INVALID_COMBINATION,
						ErrorMessage: msg}
				}

				// Store operators
				if _, ok := foundOperators[foundOperator]; ok {
					if _, ok2 := foundOperators[foundOperator][level]; ok2 {
						// if entry exists, increment
						foundOperators[foundOperator][level] = foundOperators[foundOperator][level] + 1
						fmt.Println(" -> Added. Count: " + strconv.Itoa(foundOperators[foundOperator][level]))
					} else {
						// else create new level entry with default value of 1
						foundOperators[foundOperator][level] = 1
						fmt.Println(" -> Created. Count: " + strconv.Itoa(foundOperators[foundOperator][level]))
					}
				} else {
					// if no operator entry exists, else create new operator entry with default value of 1
					foundOperators[foundOperator] = make(map[int]int)
					foundOperators[foundOperator][level] = 1
					fmt.Println(" -> Created level and value. Count: " + strconv.Itoa(foundOperators[foundOperator][level]))
				}

				// If already in right parsing mode, there should be no operator
				if modeMap[level] == tree.PARSING_MODE_RIGHT {
					log.Println("Found additional operator [" + foundOperator + "] (now " + strconv.Itoa(foundOperators[foundOperator][level]) +
						" times on level " + strconv.Itoa(level) + "), even though looking for terminating parenthesis.")
					if foundOperator == tree.AND && foundOperators[foundOperator][level] > 1 { // if AND operator and multiple on the same level
						// Consider injecting a left parenthesis before the expression and add mixfix ") " before logical operator, e.g., "( left ... [AND] right ... ) [AND] ..."
						expression = expression[:levelMap[level].Left] + "(" + expression[levelMap[level].Left:i-1] + ")" + expression[i:]
						log.Println("Multiple [AND] operators found. Reconstructed nested structure by introducing parentheses, now: " + expression)
						log.Println("Rerunning all parsing on combination to capture nested AND combinations")
						return detectCombinations(expression)
					} else {
						return levelMap, expression, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_INVALID_OPERATOR_COMBINATIONS,
							ErrorMessage: "Error: Duplicate non-[AND] operators (or mix of [AND] and non-[AND] operators) on level " + strconv.Itoa(level) +
								" in single expression (Expression: " + expression + ")"}
					}
				}

				// Configure mode
				modeMap[level] = tree.PARSING_MODE_RIGHT
				fmt.Println("Mode: " + modeMap[level])

				// Store index reference and value
				if o, ok := levelMap[level]; ok {
					o.Operator = i
					o.OperatorVal = foundOperator
					levelMap[level] = o
				}
			}
		}
	}

	if parCount != 0 {
		msg := "Uneven number of parentheses (positive --> too many left; negative --> too many right): " + strconv.Itoa(parCount)
		log.Println(msg)
		return nil, expression, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_IMBALANCED_PARENTHESES, ErrorMessage: msg}
	}

	// Check for non-parsed prefix or suffix of input string
	/*i := 0
	firstIdx := -1
	lastIdx := -1
	for i < len(levelMap) {
		if _, ok := levelMap[i]; ok {
			if firstIdx == -1 {
				// Assign first value
				firstIdx = levelMap[i].Left
				fmt.Println("Prefix pos: " + strconv.Itoa(firstIdx))
			}
			if levelMap[i].Right > lastIdx {
				// Assign highest last index
				lastIdx = levelMap[i].Right
				fmt.Println("Suffix pos: " + strconv.Itoa(lastIdx))
			}
		}
		i++
	}
	prefix := ""
	suffix := ""
	if firstIdx > 0 {
		prefix = strings.Trim(expression[:firstIdx], " (")
	}
	if lastIdx != -1 {
		suffix = strings.Trim(expression[lastIdx+1:], ") ")
	}
	if prefix != "" || suffix != "" {
		fmt.Println("Prefix: " + prefix)
		fmt.Println("Suffix: " + suffix)
		ignoredElements := []string{}
		errorString := ""
		if prefix != "" {
			ignoredElements = append(ignoredElements, prefix)
			errorString += prefix
		}
		if suffix != "" {
			ignoredElements = append(ignoredElements, suffix)
			if errorString != "" {
				errorString += ", "
			}
			errorString += suffix
		}
		fmt.Println("Returning expression (ignored elements: " + errorString + "): " + expression)
		return levelMap, expression, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_IGNORED_ELEMENTS,
			ErrorMessage: "Parsing was successful, but expression parts were ignored during coding (" + errorString + "). " +
			"This commonly occurs when logical operators between simple strings and combinations are omitted " +
			"(e.g., ... some string (left [AND] right) ...) and not wrapped by parentheses to signal shared elements. " +
			"In this case, simple strings are ignored in the parsing process.",
			ErrorIgnoredElements: ignoredElements}
	}*/

	fmt.Println("Returning expression (complete parsing): " + expression)
	// if no omitted elements during parsing, regular return without error
	return levelMap, expression, tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
}

/*
Processes potential inheritance of shared element values from parent to child nodes,
where both parent and child nodes have AND operators
 */
func inheritSharedElements(node *tree.Node) {
	if node.LogicalOperator == tree.AND &&
		node.Parent.LogicalOperator == tree.AND &&
		SHARED_ELEMENT_INHERITANCE_MODE != SHARED_ELEMENT_INHERIT_NOTHING {

		switch SHARED_ELEMENT_INHERITANCE_MODE {
		case SHARED_ELEMENT_INHERIT_OVERRIDE:
			// Overwrite child with parent shared element values
			if node.Parent.SharedLeft != "" {
				node.SharedLeft = node.Parent.SharedLeft
			}
			if node.Parent.SharedRight != "" {
				node.SharedRight = node.Parent.SharedRight
			}
		case SHARED_ELEMENT_INHERIT_APPEND:
			if node.Parent.SharedLeft != "" && node.SharedLeft != "" {
				// Append child to parent values and assign to child
				node.SharedLeft = node.Parent.SharedLeft + INHERITANCE_DELIMITER + node.SharedLeft
			} else if node.Parent.SharedLeft != "" {
				//if child is empty, just overwrite
				node.SharedLeft = node.Parent.SharedLeft
			}
			if node.Parent.SharedRight != "" && node.SharedRight != "" {
				// Append child to parent values and assign to child
				node.SharedRight = node.Parent.SharedRight + INHERITANCE_DELIMITER + node.SharedRight
			} else if node.Parent.SharedRight != "" {
				//if child is empty, just overwrite
				node.SharedRight = node.Parent.SharedRight
			}
		}
		fmt.Println("Inherited shared component from parent component in mode " + SHARED_ELEMENT_INHERITANCE_MODE + ": " +
			"Left: " + node.SharedLeft + ", Right: " + node.SharedRight)
	}
}