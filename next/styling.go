// Package next provides the 3D Tiles 1.1 types.
// Styling types are re-exported from the unified style package.
package next

import "github.com/flywave/go-3dtile/style"

// Style represents a 3D Tiles style.
type Style = style.Style

// DefaultStyle returns a new Style with default values.
func DefaultStyle() style.Style {
	return style.DefaultStyle()
}

// Expression represents a valid 3D Tiles style expression.
type Expression = style.Expression

// BooleanExpression represents a boolean or string expression.
type BooleanExpression = style.BooleanExpression

// ColorExpression represents a color expression string.
type ColorExpression = style.ColorExpression

// NumberExpression represents a number expression.
type NumberExpression = style.NumberExpression

// BooleanExpressionOrConditions represents either a boolean expression or conditions.
type BooleanExpressionOrConditions = style.BooleanExpressionOrConditions

// NumberExpressionOrConditions represents either a number expression or conditions.
type NumberExpressionOrConditions = style.NumberExpressionOrConditions

// ColorExpressionOrConditions represents either a color expression or conditions.
type ColorExpressionOrConditions = style.ColorExpressionOrConditions

// StringExpressionOrConditions represents either a string expression or conditions.
type StringExpressionOrConditions = style.StringExpressionOrConditions

// Conditions represents a series of conditional expressions.
type Conditions = style.Conditions

// ValidationError describes a single validation error in a style field.
type ValidationError = style.ValidationError

// ValidateExpression checks whether a 3D Tiles style expression string is syntactically valid.
func ValidateExpression(expr string) error {
	return style.ValidateExpression(expr)
}

// ValidateStyle validates all expression fields in a Style.
func ValidateStyle(s *style.Style) []style.ValidationError {
	return style.ValidateStyle(s)
}
