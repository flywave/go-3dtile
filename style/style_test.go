package style

import (
	"encoding/json"
	"testing"
)

func TestDefaultStyle(t *testing.T) {
	s := DefaultStyle()
	if s.Show.BooleanExpression == nil {
		t.Error("Show should have BooleanExpression")
	}
	if s.Color.ColorExpression != "color('#FFFFFF')" {
		t.Errorf("Color = %q, want color('#FFFFFF')", s.Color.ColorExpression)
	}
}

func TestStyleJSONRoundTrip(t *testing.T) {
	s := DefaultStyle()
	data, err := json.Marshal(s)
	if err != nil {
		t.Fatal("Marshal failed:", err)
	}
	var decoded Style
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal("Unmarshal failed:", err)
	}
	if string(decoded.Show.BooleanExpression) != "true" {
		t.Errorf("Show = %s, want true", decoded.Show.BooleanExpression)
	}
}

func TestStyleWithConditions(t *testing.T) {
	s := Style{
		Show: BooleanExpressionOrConditions{
			Conditions: &Conditions{
				Conditions: &[][2]string{
					{"${Height} > 10", "true"},
					{"true", "false"},
				},
			},
		},
		Color: ColorExpressionOrConditions{
			ColorExpression: "color('red')",
		},
		PointSize: NumberExpressionOrConditions{
			NumberExpression: json.RawMessage(`2.0`),
		},
	}
	data, err := json.Marshal(s)
	if err != nil {
		t.Fatal("Marshal failed:", err)
	}
	var decoded Style
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal("Unmarshal failed:", err)
	}
	if decoded.Show.Conditions == nil {
		t.Fatal("Show.Conditions should not be nil")
	}
}

func TestStyleWithStringField(t *testing.T) {
	expr := "1.0"
	s := Style{
		Show: BooleanExpressionOrConditions{
			BooleanExpression: json.RawMessage(`true`),
		},
		Color: ColorExpressionOrConditions{
			ColorExpression: "color('white')",
		},
		PointSize: NumberExpressionOrConditions{
			NumberExpression: json.RawMessage(`1.0`),
		},
		Font: &StringExpressionOrConditions{
			StringExpression: &expr,
		},
	}
	data, err := json.Marshal(s)
	if err != nil {
		t.Fatal("Marshal failed:", err)
	}
	var decoded Style
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal("Unmarshal failed:", err)
	}
	if decoded.Font == nil || decoded.Font.StringExpression == nil {
		t.Fatal("Font should have StringExpression")
	}
	if *decoded.Font.StringExpression != "1.0" {
		t.Errorf("Font = %q, want 1.0", *decoded.Font.StringExpression)
	}
}

func TestValidateExpression(t *testing.T) {
	if err := ValidateExpression("${Height} > 50"); err != nil {
		t.Errorf("valid expression should pass: %v", err)
	}
	if err := ValidateExpression(""); err == nil {
		t.Error("empty expression should fail")
	}
}

func TestValidateStyle(t *testing.T) {
	s := DefaultStyle()
	errs := ValidateStyle(&s)
	if errs != nil {
		t.Errorf("default style should be valid: %v", errs)
	}
	errs = ValidateStyle(nil)
	if errs != nil {
		t.Errorf("nil style should return nil: %v", errs)
	}
}

func TestValidateStyleInvalid(t *testing.T) {
	s := DefaultStyle()
	s.Color.ColorExpression = "color(" // incomplete
	errs := ValidateStyle(&s)
	if len(errs) == 0 {
		t.Error("expected error for invalid expression")
	}
}

func TestConstants(t *testing.T) {
	if OP_ADD != "+" {
		t.Error("OP_ADD should be +")
	}
	if FUNC_SIN != "sin" {
		t.Error("FUNC_SIN should be sin")
	}
	if VAR_BOOLEAN != 0 {
		t.Error("VAR_BOOLEAN should be 0")
	}
	if EXP_VARIABLE != 0 {
		t.Error("EXP_VARIABLE should be 0")
	}
}

func TestConditionsMarshal(t *testing.T) {
	conds := Conditions{
		Conditions: &[][2]string{
			{"${H} > 0", "true"},
		},
	}
	data, err := json.Marshal(conds)
	if err != nil {
		t.Fatal("Marshal failed:", err)
	}
	var decoded Conditions
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal("Unmarshal failed:", err)
	}
	if decoded.Conditions == nil {
		t.Fatal("Conditions should not be nil")
	}
	if (*decoded.Conditions)[0][0] != "${H} > 0" {
		t.Errorf("condition = %q", (*decoded.Conditions)[0][0])
	}
}

func TestOrConditionsRoundTrip(t *testing.T) {
	// Simple expression
	b := BooleanExpressionOrConditions{
		BooleanExpression: json.RawMessage(`true`),
	}
	data, _ := json.Marshal(b)
	if string(data) != "true" {
		t.Errorf("marshaled = %s, want true", string(data))
	}

	// Conditions
	b2 := BooleanExpressionOrConditions{
		Conditions: &Conditions{
			Conditions: &[][2]string{
				{"${H} > 0", "false"},
			},
		},
	}
	data2, _ := json.Marshal(b2)
	var b3 BooleanExpressionOrConditions
	json.Unmarshal(data2, &b3)
	if b3.Conditions == nil {
		t.Error("Conditions should not be nil")
	}
}
