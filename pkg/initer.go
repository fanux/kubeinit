package pkg

//Initer is a interface to define modules to install
type Initer interface {
	//query module name
	Info() (string, string)
	//generate and render config files
	Gen() error
	//install module
	Run() error
	//clean module
	Clean() error
	//install offline
	InstallOffline() error
	//Install online
	InstallOnline() error
	//Save image, config file or bin
	Save() error
}
