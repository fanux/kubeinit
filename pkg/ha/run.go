package ha

//Ha is
type Ha struct{}

//Info is
func (e *Ha) Info() (string, string) {
	return "ha", nil
}

//Gen is
func (e *Ha) Gen() error {
	return nil
}

//Run is
func (e *Ha) Run() error {
	return nil
}

//Clean is
func (e *Ha) Clean() error {
	return nil
}
