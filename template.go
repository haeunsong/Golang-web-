package main

import (
	"os"
	"text/template"
)

type Student struct {
	// exported field since it begins
	// with a capital letter
	Name string
}

func main() {
	s := Student{"Satish"}

	// create a new template with some name
	tmpl := template.New("test")

	// parse some content and generate a template
	tmpl, err := tmpl.Parse("Hello {{.Name}}!")

	if err != nil {
		panic(err)
	}
	// merge template 'tmpl' with content of 's'
	err1 := tmpl.Execute(os.Stdout, s)

	if err1 != nil {
		panic(err)
	}
}
