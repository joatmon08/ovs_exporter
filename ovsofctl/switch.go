package ovsofctl

type OpenvSwitch struct {
	capabilities string
	actions      string
	ports        map[string]OvsPort
}