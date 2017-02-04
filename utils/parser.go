package utils

import (
	"reflect"
	"strings"
	"regexp"
	"github.com/Sirupsen/logrus"
	"errors"
)

func GetFields(structure interface{}) (map[string]int) {
	fields := map[string]int{}
	t := reflect.TypeOf(structure)
	for i := 0; i < t.NumField(); i++ {
		fields[t.Field(i).Name] = 1
	}
	return fields
}

func MapStringToInterface(expectedStruct interface{}, input string) (map[string]interface{}) {
	structure := map[string]interface{}{}
	fields := GetFields(expectedStruct)
	lines := strings.Split(input, "\n")
	for field := range fields {
		re := regexp.MustCompile(field + ":")
		for _, line := range lines {
			matchedIndices := re.FindStringIndex(line)
			if len(matchedIndices) < 2 {
				continue
			} else {
				logrus.Debugf("Adding %s:%s", field, line[matchedIndices[1]:])
				structure[field] = strings.TrimSpace(line[matchedIndices[1]:])
			}
		}
	}
	return structure
}

func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return errors.New("No such field: " + name + " in obj")
	}

	if !structFieldValue.CanSet() {
		return errors.New("Cannot set " + name + " field value")
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}

