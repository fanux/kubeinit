package net

//Net is
type Net struct {
}

//Info is
func (e *Net) Info() (string, string) {
	return "net", nil
}

//Gen is
func (e *Net) Gen() error {
	return nil
}

//Run is
func (e *Net) Run() error {
	return nil
}

//Clean is
func (e *Net) Clean() error {
	return nil
}
