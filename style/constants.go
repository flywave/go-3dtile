package style

// Unary operators
const (
	OP_NEGATIVE = "-"
	OP_POSITIVE = "+"
	OP_NOT      = "!"
)

// Binary operators
const (
	OP_ADD        = "+"
	OP_SUB        = "-"
	OP_MUL        = "*"
	OP_DIV        = "/"
	OP_MOD        = "%"
	OP_EQ         = "==="
	OP_NEQ        = "!=="
	OP_GREATER    = ">"
	OP_GREATER_EQ = ">="
	OP_LESS       = "<"
	OP_LESS_EQ    = "<="
	OP_AND        = "&&"
	OP_OR         = "||"
	OP_REGEXP_NOT = "!~"
	OP_REGEXP     = "=~"
)

// Single-argument built-in functions
const (
	FUNC_IS_NAN             = "isNaN"
	FUNC_IS_FINITE          = "isFinite"
	FUNC_IS_EXACTCLASS      = "isExactClass"
	FUNC_IS_CLASS           = "isClass"
	FUNC_GET_EXACTCLASSNAME = "getExactClassName"
	FUNC_BOOLEAN            = "Boolean"
	FUNC_NUMBER             = "Number"
	FUNC_STRING             = "String"
	FUNC_ABS                = "abs"
	FUNC_SQRT               = "sqrt"
	FUNC_COS                = "cos"
	FUNC_SIN                = "sin"
	FUNC_TAN                = "tan"
	FUNC_ACOS               = "acos"
	FUNC_ASIN               = "asin"
	FUNC_ATAN               = "atan"
	FUNC_RADIANS            = "radians"
	FUNC_DEGREES            = "degrees"
	FUNC_SIGN               = "sign"
	FUNC_FLOOR              = "floor"
	FUNC_CEIL               = "ceil"
	FUNC_ROUND              = "round"
	FUNC_EXP                = "exp"
	FUNC_EXP2               = "exp2"
	FUNC_LOG                = "log"
	FUNC_LOG2               = "log2"
	FUNC_FRACT              = "fract"
	FUNC_LENGTH             = "length"
	FUNC_NORMALIZE          = "normalize"
)

// Two-argument built-in functions
const (
	FUNC_ATAN2    = "atan2"
	FUNC_POW      = "pow"
	FUNC_MIN      = "min"
	FUNC_MAX      = "max"
	FUNC_DISTANCE = "distance"
	FUNC_DOT      = "dot"
	FUNC_CROSS    = "cross"
)

// Three-argument built-in functions
const (
	FUNC_CLAMP = "clamp"
	FUNC_MIX   = "mix"
)

// ExpressionType defines the result type of a style expression.
type ExpressionType int32

const (
	VAR_BOOLEAN   ExpressionType = 0
	VAR_NULL      ExpressionType = 1
	VAR_UNDEFINED ExpressionType = 2
	VAR_NUMBER    ExpressionType = 3
	VAR_STRING    ExpressionType = 4
	VAR_ARRAY     ExpressionType = 5
	VAR_VEC2      ExpressionType = 6
	VAR_VEC3      ExpressionType = 7
	VAR_VEC4      ExpressionType = 8
	VAR_REGEXP    ExpressionType = 9
)

// ExpressionNodeType defines the AST node type for a parsed style expression.
type ExpressionNodeType int32

const (
	EXP_VARIABLE           ExpressionNodeType = 0  // Variable reference: ${Height}
	EXP_UNARY              ExpressionNodeType = 1  // Unary operation: !x, -x
	EXP_BINARY             ExpressionNodeType = 2  // Binary operation: a + b
	EXP_TERNARY            ExpressionNodeType = 3  // Ternary: a ? b : c
	EXP_CONDITIONAL        ExpressionNodeType = 4  // Conditional: [[...], [...]]
	EXP_MEMBER             ExpressionNodeType = 5  // Member access: a.b
	EXP_FUNCTION_CALL      ExpressionNodeType = 6  // Function call: sin(x)
	EXP_ARRAY              ExpressionNodeType = 7  // Array literal: [1, 2, 3]
	EXP_REGEX              ExpressionNodeType = 8  // Regex literal: /pattern/
	EXP_VARIABLE_IN_STRING ExpressionNodeType = 9  // Variable in string: "${x}"
	EXP_LITERAL_NULL       ExpressionNodeType = 10
	EXP_LITERAL_BOOLEAN    ExpressionNodeType = 11
	EXP_LITERAL_NUMBER     ExpressionNodeType = 12
	EXP_LITERAL_STRING     ExpressionNodeType = 13
	EXP_LITERAL_COLOR      ExpressionNodeType = 14 // color('#FF0000')
	EXP_LITERAL_VECTOR     ExpressionNodeType = 15 // vec3(1,2,3)
	EXP_LITERAL_REGEX      ExpressionNodeType = 16
	EXP_LITERAL_UNDEFINED  ExpressionNodeType = 17
	EXP_BUILTIN_VARIABLE   ExpressionNodeType = 18 // tiles3d_* variables
)
