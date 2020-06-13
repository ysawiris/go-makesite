package main

import (
	"flag"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

//
type post struct {
	Description string
}

func readFile(filename string) string {
	// Reading a File
	fileContents, err := ioutil.ReadFile(filename)
	if err != nil {
		// A common use of `panic` is to abort if a function returns an error
		// value that we don’t know how to (or want to) handle. This example
		// panics if we get an unexpected error when creating a new file.
		panic(err)
	}
	return string(fileContents)
}

func writeFile(fileContents string) {
	//Writing a File
	bytesToWrite := []byte(fileContents)
	err := ioutil.WriteFile("new-file1.txt", bytesToWrite, 0644)
	if err != nil {
		panic(err)
	}
}
func renderTemplate(filename string, data string) {
	// var that holds content
	c := post{Description: data}
	// render given content in the template
	t := template.Must(template.New("template.tmpl").ParseFiles(filename))

	var err error
	// Print out using Stdout
	err = t.Execute(os.Stdout, c)
	if err != nil {
		panic(err)
	}
}

func changeExtHTML(filename string) string {
	// var to hold ext ending
	ext := ".html"
	// Deletes extension ending and adds new one
	withExt := strings.Split(filename, ".")[0] + ext
	return withExt
}

func writeTemplateToFile(tmplName string, data string) {
	// var to hold content
	c := post{Description: readFile(data)}
	// render given content in the template
	t := template.Must(template.New("template.tmpl").ParseFiles(tmplName))
	// change file extension
	file := changeExtHTML(data)
	// Create new file
	f, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	// render template
	err = t.Execute(f, c)
	if err != nil {
		panic(err)
	}
}

func main() {
	fileParse := flag.String("file", "", "txt file will be converted to html file")
	flag.Parse()
	if *fileParse != "" {
		renderTemplate("template.tmpl", readFile(*fileParse))
		writeTemplateToFile("template.tmpl", *fileParse)
	} else {
		renderTemplate("template.tmpl", readFile("first-post.txt"))
		writeTemplateToFile("template.tmpl", "test.txt")
	}
}
