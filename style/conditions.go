package style

import "encoding/json"

// StringExpressionOrConditions represents either a string expression or conditions.
type StringExpressionOrConditions struct {
	StringExpression *string    `json:"-"`
	Conditions       *Conditions `json:"-"`
}

func (s StringExpressionOrConditions) MarshalJSON() ([]byte, error) {
	if s.Conditions != nil {
		return json.Marshal(s.Conditions)
	}
	if s.StringExpression != nil {
		return json.Marshal(*s.StringExpression)
	}
	return json.Marshal(nil)
}

func (s *StringExpressionOrConditions) UnmarshalJSON(data []byte) error {
	var conditions Conditions
	if err := json.Unmarshal(data, &conditions); err == nil && conditions.Conditions != nil {
		s.Conditions = &conditions
		return nil
	}
	var expr string
	if err := json.Unmarshal(data, &expr); err != nil {
		return err
	}
	s.StringExpression = &expr
	return nil
}

// BooleanExpressionOrConditions represents either a boolean expression or conditions.
type BooleanExpressionOrConditions struct {
	BooleanExpression BooleanExpression `json:"-"`
	Conditions        *Conditions       `json:"-"`
}

func (b BooleanExpressionOrConditions) MarshalJSON() ([]byte, error) {
	if b.Conditions != nil {
		return json.Marshal(b.Conditions)
	}
	return json.Marshal(b.BooleanExpression)
}

func (b *BooleanExpressionOrConditions) UnmarshalJSON(data []byte) error {
	var conditions Conditions
	if err := json.Unmarshal(data, &conditions); err == nil && conditions.Conditions != nil {
		b.Conditions = &conditions
		return nil
	}
	var expr BooleanExpression
	if err := json.Unmarshal(data, &expr); err != nil {
		return err
	}
	b.BooleanExpression = expr
	return nil
}

// NumberExpressionOrConditions represents either a number expression or conditions.
type NumberExpressionOrConditions struct {
	NumberExpression NumberExpression `json:"-"`
	Conditions       *Conditions      `json:"-"`
}

func (n NumberExpressionOrConditions) MarshalJSON() ([]byte, error) {
	if n.Conditions != nil {
		return json.Marshal(n.Conditions)
	}
	return json.Marshal(n.NumberExpression)
}

func (n *NumberExpressionOrConditions) UnmarshalJSON(data []byte) error {
	var conditions Conditions
	if err := json.Unmarshal(data, &conditions); err == nil && conditions.Conditions != nil {
		n.Conditions = &conditions
		return nil
	}
	var expr NumberExpression
	if err := json.Unmarshal(data, &expr); err != nil {
		return err
	}
	n.NumberExpression = expr
	return nil
}

// ColorExpressionOrConditions represents either a color expression or conditions.
type ColorExpressionOrConditions struct {
	ColorExpression ColorExpression `json:"-"`
	Conditions      *Conditions     `json:"-"`
}

func (c ColorExpressionOrConditions) MarshalJSON() ([]byte, error) {
	if c.Conditions != nil {
		return json.Marshal(c.Conditions)
	}
	return json.Marshal(c.ColorExpression)
}

func (c *ColorExpressionOrConditions) UnmarshalJSON(data []byte) error {
	var conditions Conditions
	if err := json.Unmarshal(data, &conditions); err == nil && conditions.Conditions != nil {
		c.Conditions = &conditions
		return nil
	}
	var expr ColorExpression
	if err := json.Unmarshal(data, &expr); err != nil {
		return err
	}
	c.ColorExpression = expr
	return nil
}

// Conditions represents a series of conditional expressions [[cond, value], ...].
type Conditions struct {
	Conditions *[][2]string `json:"conditions,omitempty"`
}
