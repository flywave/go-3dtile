// Package tile3d provides the legacy 3D Tiles 1.0 types.
// Styling types are re-exported from the unified style package.
package tile3d

import "github.com/flywave/go-3dtile/style"

// Styling defines how features are displayed using 3D Tiles styling expressions.
type Styling = style.Style

// Expression is a 3D Tiles style expression string.
type Expression = style.Expression

// ExpressionType defines the result type of a style expression.
type ExpressionType = style.ExpressionType

// ExpressionNodeType defines the AST node type for a parsed style expression.
type ExpressionNodeType = style.ExpressionNodeType

// Unary operator constants
const (
	OP_NEGATIVE = style.OP_NEGATIVE
	OP_POSITIVE = style.OP_POSITIVE
	OP_NOT      = style.OP_NOT
)

// Binary operator constants
const (
	OP_ADD        = style.OP_ADD
	OP_SUB        = style.OP_SUB
	OP_MUL        = style.OP_MUL
	OP_DIV        = style.OP_DIV
	OP_MOD        = style.OP_MOD
	OP_EQ         = style.OP_EQ
	OP_NEQ        = style.OP_NEQ
	OP_GREATER    = style.OP_GREATER
	OP_GREATER_EQ = style.OP_GREATER_EQ
	OP_LESS       = style.OP_LESS
	OP_LESS_EQ    = style.OP_LESS_EQ
	OP_AND        = style.OP_AND
	OP_OR         = style.OP_OR
	OP_REGEXP_NOT = style.OP_REGEXP_NOT
	OP_REGEXP     = style.OP_REGEXP
)

// Single-argument built-in function constants
const (
	FUNC_IS_NAN             = style.FUNC_IS_NAN
	FUNC_IS_FINITE          = style.FUNC_IS_FINITE
	FUNC_IS_EXACTCLASS      = style.FUNC_IS_EXACTCLASS
	FUNC_IS_CLASS           = style.FUNC_IS_CLASS
	FUNC_GET_EXACTCLASSNAME = style.FUNC_GET_EXACTCLASSNAME
	FUNC_BOOLEAN            = style.FUNC_BOOLEAN
	FUNC_NUMBER             = style.FUNC_NUMBER
	FUNC_STRING             = style.FUNC_STRING
	FUNC_ABS                = style.FUNC_ABS
	FUNC_SQRT               = style.FUNC_SQRT
	FUNC_COS                = style.FUNC_COS
	FUNC_SIN                = style.FUNC_SIN
	FUNC_TAN                = style.FUNC_TAN
	FUNC_ACOS               = style.FUNC_ACOS
	FUNC_ASIN               = style.FUNC_ASIN
	FUNC_ATAN               = style.FUNC_ATAN
	FUNC_RADIANS            = style.FUNC_RADIANS
	FUNC_DEGREES            = style.FUNC_DEGREES
	FUNC_SIGN               = style.FUNC_SIGN
	FUNC_FLOOR              = style.FUNC_FLOOR
	FUNC_CEIL               = style.FUNC_CEIL
	FUNC_ROUND              = style.FUNC_ROUND
	FUNC_EXP                = style.FUNC_EXP
	FUNC_EXP2               = style.FUNC_EXP2
	FUNC_LOG                = style.FUNC_LOG
	FUNC_LOG2               = style.FUNC_LOG2
	FUNC_FRACT              = style.FUNC_FRACT
	FUNC_LENGTH             = style.FUNC_LENGTH
	FUNC_NORMALIZE          = style.FUNC_NORMALIZE
)

// Two-argument built-in function constants
const (
	FUNC_ATAN2    = style.FUNC_ATAN2
	FUNC_POW      = style.FUNC_POW
	FUNC_MIN      = style.FUNC_MIN
	FUNC_MAX      = style.FUNC_MAX
	FUNC_DISTANCE = style.FUNC_DISTANCE
	FUNC_DOT      = style.FUNC_DOT
	FUNC_CROSS    = style.FUNC_CROSS
)

// Three-argument built-in function constants
const (
	FUNC_CLAMP = style.FUNC_CLAMP
	FUNC_MIX   = style.FUNC_MIX
)

// AST node type constants (re-exported)
const (
	VAR_BOOLEAN           = style.VAR_BOOLEAN
	VAR_NULL              = style.VAR_NULL
	VAR_UNDEFINED         = style.VAR_UNDEFINED
	VAR_NUMBER            = style.VAR_NUMBER
	VAR_STRING            = style.VAR_STRING
	VAR_ARRAY             = style.VAR_ARRAY
	VAR_VEC2              = style.VAR_VEC2
	VAR_VEC3              = style.VAR_VEC3
	VAR_VEC4              = style.VAR_VEC4
	VAR_REGEXP            = style.VAR_REGEXP
	EXP_VARIABLE          = style.EXP_VARIABLE
	EXP_UNARY             = style.EXP_UNARY
	EXP_BINARY            = style.EXP_BINARY
	EXP_TERNARY           = style.EXP_TERNARY
	EXP_CONDITIONAL       = style.EXP_CONDITIONAL
	EXP_MEMBER            = style.EXP_MEMBER
	EXP_FUNCTION_CALL     = style.EXP_FUNCTION_CALL
	EXP_ARRAY             = style.EXP_ARRAY
	EXP_REGEX             = style.EXP_REGEX
	EXP_VARIABLE_IN_STRING = style.EXP_VARIABLE_IN_STRING
	EXP_LITERAL_NULL      = style.EXP_LITERAL_NULL
	EXP_LITERAL_BOOLEAN   = style.EXP_LITERAL_BOOLEAN
	EXP_LITERAL_NUMBER    = style.EXP_LITERAL_NUMBER
	EXP_LITERAL_STRING    = style.EXP_LITERAL_STRING
	EXP_LITERAL_COLOR     = style.EXP_LITERAL_COLOR
	EXP_LITERAL_VECTOR    = style.EXP_LITERAL_VECTOR
	EXP_LITERAL_REGEX     = style.EXP_LITERAL_REGEX
	EXP_LITERAL_UNDEFINED = style.EXP_LITERAL_UNDEFINED
	EXP_BUILTIN_VARIABLE  = style.EXP_BUILTIN_VARIABLE
)
