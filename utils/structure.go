package utils

type Structure interface {
	Fill(m map[string]interface{}) error
}