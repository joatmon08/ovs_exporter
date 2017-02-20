package utils

import "io/ioutil"

func ReadTestDataToString(filename string) (string, error) {
	file, err := ioutil.ReadFile("testdata/" + filename)
	if err != nil {
	return "", err
	}
	return string(file), nil
}


func ReadTestDataToBytes(filename string) ([]byte, error) {
	file, err := ioutil.ReadFile("testdata/" + filename)
	if err != nil {
		return nil, err
	}
	return []byte(file), nil
}