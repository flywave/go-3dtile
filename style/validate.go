package style

import (
	"fmt"
	"strings"

	"github.com/flywave/go-3dtile/expr"
)

// ValidationError describes a single validation error in a style field.
type ValidationError struct {
	Field string
	Index int // -1 if not a condition index
	Err   error
}

func (ve ValidationError) Error() string {
	if ve.Index >= 0 {
		return fmt.Sprintf("%s[%d]: %v", ve.Field, ve.Index, ve.Err)
	}
	return fmt.Sprintf("%s: %v", ve.Field, ve.Err)
}

// ValidateExpression checks whether a 3D Tiles style expression string is syntactically valid.
func ValidateExpression(exprStr string) error {
	trimmed := strings.TrimSpace(exprStr)
	if trimmed == "" {
		return fmt.Errorf("empty expression")
	}
	_, err := expr.Parse(trimmed)
	if err != nil {
		return fmt.Errorf("syntax error: %w", err)
	}
	return nil
}

func validateRawExpr(raw BooleanExpression) error {
	if raw == nil {
		return nil
	}
	s := string(raw)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		inner := s[1 : len(s)-1]
		return ValidateExpression(inner)
	}
	return nil
}

func validateConditions(field string, conds *[][2]string) []ValidationError {
	if conds == nil {
		return nil
	}
	var errs []ValidationError
	for i, pair := range *conds {
		if len(pair) >= 1 {
			if err := ValidateExpression(pair[0]); err != nil {
				errs = append(errs, ValidationError{Field: field, Index: i, Err: fmt.Errorf("condition: %w", err)})
			}
		}
		if len(pair) >= 2 {
			if err := ValidateExpression(pair[1]); err != nil {
				errs = append(errs, ValidationError{Field: field, Index: i, Err: fmt.Errorf("value: %w", err)})
			}
		}
	}
	return errs
}

// ValidateStyle validates all expression fields in a Style, returning a list of
// validation errors. Each error identifies the specific field and condition index.
func ValidateStyle(s *Style) []ValidationError {
	if s == nil {
		return nil
	}
	var errs []ValidationError

	errs = append(errs, validateBooleanField("show", &s.Show)...)
	errs = append(errs, validateColorField("color", &s.Color)...)
	errs = append(errs, validateNumberField("pointSize", &s.PointSize)...)
	errs = append(errs, validateStringField("pointOutlineColor", s.PointOutlineColor)...)
	errs = append(errs, validateStringField("pointOutlineWidth", s.PointOutlineWidth)...)
	errs = append(errs, validateStringField("labelColor", s.LabelColor)...)
	errs = append(errs, validateStringField("labelOutlineColor", s.LabelOutlineColor)...)
	errs = append(errs, validateStringField("labelOutlineWidth", s.LabelOutlineWidth)...)
	errs = append(errs, validateStringField("font", s.Font)...)
	errs = append(errs, validateStringField("labelStyle", s.LabelStyle)...)
	errs = append(errs, validateStringField("backgroundColor", s.BackgroundColor)...)
	errs = append(errs, validateStringField("backgroundPadding", s.BackgroundPadding)...)
	errs = append(errs, validateStringField("backgroundEnabled", s.BackgroundEnabled)...)
	errs = append(errs, validateStringField("scaleByDistance", s.ScaleByDistance)...)
	errs = append(errs, validateStringField("translucencyByDistance", s.TranslucencyByDistance)...)
	errs = append(errs, validateStringField("distanceDisplayCondition", s.DistanceDisplayCondition)...)
	errs = append(errs, validateStringField("heightOffset", s.HeightOffset)...)
	errs = append(errs, validateStringField("anchorLineEnabled", s.AnchorLineEnabled)...)
	errs = append(errs, validateStringField("anchorLineColor", s.AnchorLineColor)...)
	errs = append(errs, validateStringField("image", s.Image)...)
	errs = append(errs, validateStringField("disableDepthTestDistance", s.DisableDepthTestDistance)...)
	errs = append(errs, validateStringField("horizontalOrigin", s.HorizontalOrigin)...)
	errs = append(errs, validateStringField("verticalOrigin", s.VerticalOrigin)...)
	errs = append(errs, validateStringField("labelHorizontalOrigin", s.LabelHorizontalOrigin)...)
	errs = append(errs, validateStringField("labelVerticalOrigin", s.LabelVerticalOrigin)...)

	if len(errs) == 0 {
		return nil
	}
	return errs
}

func validateBooleanField(field string, b *BooleanExpressionOrConditions) []ValidationError {
	if b == nil {
		return nil
	}
	if b.Conditions != nil && b.Conditions.Conditions != nil {
		return validateConditions(field, b.Conditions.Conditions)
	}
	if b.BooleanExpression != nil {
		if err := validateRawExpr(b.BooleanExpression); err != nil {
			return []ValidationError{{Field: field, Err: err}}
		}
	}
	return nil
}

func validateColorField(field string, c *ColorExpressionOrConditions) []ValidationError {
	if c == nil {
		return nil
	}
	if c.Conditions != nil && c.Conditions.Conditions != nil {
		return validateConditions(field, c.Conditions.Conditions)
	}
	if c.ColorExpression != "" {
		if err := ValidateExpression(c.ColorExpression); err != nil {
			return []ValidationError{{Field: field, Err: err}}
		}
	}
	return nil
}

func validateNumberField(field string, n *NumberExpressionOrConditions) []ValidationError {
	if n == nil {
		return nil
	}
	if n.Conditions != nil && n.Conditions.Conditions != nil {
		return validateConditions(field, n.Conditions.Conditions)
	}
	if n.NumberExpression != nil {
		if err := validateRawExpr(n.NumberExpression); err != nil {
			return []ValidationError{{Field: field, Err: err}}
		}
	}
	return nil
}

func validateStringField(field string, s *StringExpressionOrConditions) []ValidationError {
	if s == nil {
		return nil
	}
	if s.Conditions != nil && s.Conditions.Conditions != nil {
		return validateConditions(field, s.Conditions.Conditions)
	}
	if s.StringExpression != nil {
		if err := ValidateExpression(*s.StringExpression); err != nil {
			return []ValidationError{{Field: field, Err: err}}
		}
	}
	return nil
}
