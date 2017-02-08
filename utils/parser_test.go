package utils

import (
	"testing"
	"reflect"
)

type test struct {
	Foo string
	Bar int
	Hi  string
}

type smallOvs struct {
	capabilities string
	Bar int
	Hi  string
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
		"foo: hello\n" +
		"bar: 20\n" +
		"alsoignoreme: yes\n" +
		"hi: world"
	expectedTestStruct := map[string]interface{}{
		"Foo": "hello",
		"Bar": 20,
		"Hi": "world",
	}
	result := MapStringToInterface(test{}, testString)
	if reflect.DeepEqual(result, expectedTestStruct) {
		t.Errorf("Expected %v, got %v", expectedTestStruct, result)
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
		"Hi": "test",
	}
	if err := testStructure.Fill(testInterface); err != nil {
		t.Error(err)
	}
	if testStructure.Bar != 50 {
		t.Errorf("Expected %d for %s, got %d", 50, "Bar", testStructure.Bar)
	}
	if testStructure.Foo != "bar" {
		t.Errorf("Expected %s for %s, got %s", "bar", "Foo", testStructure.Foo)
	}
	if testStructure.Hi != "test" {
		t.Errorf("Expected %s for %s, got %s", "test", "Hi", testStructure.Hi)
	}

}