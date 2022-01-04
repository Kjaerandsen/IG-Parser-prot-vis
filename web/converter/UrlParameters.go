package converter

/*
Constants for URL parameter keys to pass UI input (e.g., "?rawStmt=Certifiers ...")
*/

// GENERIC CONSTANTS FOR PARSER

// Raw statement
const PARAM_RAW_STATEMENT = "rawStmt"

// Encoded statement
const PARAM_CODED_STATEMENT = "codedStmt"

// Indicates whether parser is directly invoked with passed information
const PARAM_EXECUTE_PARSER = "execute"

// SHEETS ONLY

// Statement ID
const PARAM_STATEMENT_ID = "stmtId"

// Dynamic schema output
const PARAM_DYNAMIC_SCHEMA = "dynamicSchema"

// Component-level nesting
const PARAM_EXTENDED_OUTPUT = "igExtended"

// SHARED AMONGST SHEETS AND VISUAL OUTPUT

// Annotations
const PARAM_LOGICO_OUTPUT = "annotations"

// VISUAL ONLY

// Properties as tree structure
const PARAM_PROPERTY_TREE = "propertyTree"

// Strictly binary tree structure (for debugging)
const PARAM_BINARY_TREE = "binaryTree"

// Canvas width for visual output
const PARAM_WIDTH = "canvasWidth"

// Canvas height for visual output
const PARAM_HEIGHT = "canvasHeight"

// CHECKBOX CONSTANTS

// Checkbox constant as read from form input
const CHECKBOX_ON = "on"

// Checked checkbox as represented in HTML
const CHECKBOX_CHECKED = "checked"

// Unchecked checkbox as represented in HTML
const CHECKBOX_UNCHECKED = "unchecked"