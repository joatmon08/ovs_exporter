package ovsofctl

type OvsPorts struct {
	id        string
	name      string
	addr      string
	config    int
	state     int
	current   string
	speed     int
	max_speed int
}

type OpenvSwitch struct {
	capabilities string
	actions      string
	ports        map[string]OvsPorts
}
