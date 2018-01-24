package kubecore

//Kubecore is
type Kubecore struct{}

//Info is
func (e *Kubecore) Info() (string, string) {
	return "kubecore", nil
}

//Gen is
func (e *Kubecore) Gen() error {
	return nil
}

//Run is
func (e *Kubecore) Run() error {
	return nil
}

//Clean is
func (e *Kubecore) Clean() error {
	return nil
}
