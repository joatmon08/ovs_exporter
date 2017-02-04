package utils

import (
	"testing"
)

type test struct {
	Foo string
	Bar int
	Hi  map[string]string
}

type smallOvs struct {
	capabilities string
	Bar int
	Hi  map[string]string
}

func (t *test) Fill(m map[string]interface{}) error {
	for k, v := range m {
		err := SetField(t, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func TestGetField(t *testing.T) {
	t.Log("Testing getting fields from struct")
	testStructure := test{}
	values := GetFields(testStructure)
	if len(values) != 3 {
		t.Errorf("Expected %d, got %d for length", 3, len(values))
	}
	if _, exists := values["Bar"]; !exists {
		t.Errorf("Expected %s but doesn't exist", "Bar")
	}
	if _, exists := values["NOT"]; exists {
		t.Errorf("Expected %s not to exist but it is there", "NOT")
	}
}

func TestMapStringToInterface(t *testing.T) {
	t.Log("Testing mapping of a string to an interface")
	testString := "ignoreme: yes\n" +
		"Foo: hello\n" +
		"Bar: 0\n" +
		"alsoignoreme: yes\n" +
		"Hi: {foo:bar}"
	result := MapStringToInterface(test{}, testString)
	if len(result) != 3 {
		t.Errorf("Expected %d, got %d for length", 3, len(result))
	}
	if _, exists := result["Bar"]; !exists {
		t.Errorf("Expected %s but doesn't exist", "Bar")
	}
	if _, exists := result["Foo"]; !exists {
		t.Errorf("Expected %s but doesn't exist", "Foo")
	}
	if _, exists := result["alsoignoreme"]; exists {
		t.Errorf("Expected %s not to exist but it is there", "alsoignoreme")
	}
}

func TestSetField(t *testing.T) {
	t.Log("Testing setting of field")
	testStructure := &test{}
	if err := SetField(testStructure, "Bar", 20); err != nil {
		t.Error(err)
	}
	if testStructure.Bar != 20 {
		t.Errorf("Expected %d for %s, got %d", 20, "Bar", testStructure.Bar)
	}
}

func TestFillStructure(t *testing.T) {
	t.Log("Testing filling of structure")
	testStructure := &test{}
	testInterface := map[string]interface{}{
		"Foo": "bar",
		"Bar": 50,
		"Hi": map[string]string{
			"test":"this",
			"that":"works",
		},
	}
	testStructure.Fill(testInterface)
	if testStructure.Bar != 50 {
		t.Errorf("Expected %d for %s, got %d", 50, "Bar", testStructure.Bar)
	}
	if testStructure.Foo != "bar" {
		t.Errorf("Expected %s for %s, got %s", "bar", "Foo", testStructure.Foo)
	}
	if len(testStructure.Hi) != 2 {
		t.Errorf("Expected %d for %s, got %d", 2, "Hi", testStructure.Hi)
	}
}