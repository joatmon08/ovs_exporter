package utils

import "io/ioutil"

func ReadTestData(filename string) (string, error) {
	file, err := ioutil.ReadFile("testdata/" + filename)
	if err != nil {
	return "", err
	}
	return string(file), nil
}