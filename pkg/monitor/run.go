package monitor

//Monitor is
type Monitor struct{}

//Info is
func (e *Monitor) Info() (string, string) {
	return "monitor", nil
}

//Gen is
func (e *Monitor) Gen() error {
	return nil
}

//Run is
func (e *Monitor) Run() error {
	return nil
}

//Clean is
func (e *Monitor) Clean() error {
	return nil
}
