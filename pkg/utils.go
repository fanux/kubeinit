package pkg

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
)

//RenderToStr is
func RenderToStr(t *template.Template, tp string, args interface{}) string {
	template.Must(t.Parse(tp))

	var buf []byte
	buffer := bytes.NewBuffer(buf)

	err = t.Execute(buffer, args)
	if err != nil {
		fmt.Println("exec template file error: %s", err)
	}
	str := buffer.String()
	return str
}

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

//WriteFile is
func WriteFile(fileName string, content string) {
	b := []byte(content)
	err := ioutil.WriteFile(fileName, b, 0644)
	if err != nil {
		fmt.Println("write file error", err)
	}
}

//ApplyShell is
func ApplyShell(sh string) {
	fmt.Println("+ ", sh)
	cmd := exec.Command("bash", "-c", sh)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

//ApplyShellOutput is
func ApplyShellOutput(sh string) string {
	fmt.Println("+ ", sh)
	s, err := exec.Command("bash", "-c", sh).Output()
	if err != nil {
		fmt.Println("exec shell failed: ", sh)
		return ""
	}
	return string(s)
}
