package pkg

//Initer is a interface to define modules to install
type Initer interface {
	//generate and render config files
	Gen() error
	//install module
	Run() error
	//clean module
	Clean() error
	//query module name
	Name() string
}
