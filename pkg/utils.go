package pkg

import (
	"fmt"
	"html/template"
	"os"
)

//Render is
func Render(t *template.Template, tp string, args interface{}, outFile string) {
	template.Must(t.Parse(tp))

	file, err := os.Create(outFile)
	defer file.Close()
	if err != nil {
		fmt.Println("create out file error: %s", err)
		return
	}

	err = t.Execute(file, args)
	if err != nil {
		fmt.Println("exec template file error: %s", err)
	}
}
