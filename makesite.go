package main

import (
	"bytes"
	"flag"
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/aofei/air"
	"github.com/aofei/cameron"
)

// global variable
type post struct {
	Content string
	Avatar  string
}

var a = air.Default

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

func renderTemplate(filename string, data string, url string) {
	// var that holds content
	c := post{Content: data, Avatar: url}
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
	addExt := strings.Split(filename, ".")[0] + ext
	return addExt
}

func writeTemplateToFile(tmplName string, data string, url string) {
	// var to hold content
	c := post{Content: readFile(data), Avatar: url}
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

func checkTextFile(filename string) bool {
	if strings.Contains(filename, ".") {
		return strings.Split(filename, ".")[1] == "txt"
	}
	return false
}

func identicon(req *air.Request, res *air.Response) error {
	buf := bytes.Buffer{}
	jpeg.Encode(
		&buf,
		cameron.Identicon(
			[]byte(req.Param("Name").Value().String()),
			540,
			60,
		),
		&jpeg.Options{
			Quality: 100,
		},
	)

	res.Header.Set("Content-Type", "image/jpeg")

	return res.Write(bytes.NewReader(buf.Bytes()))
}

func main() {
	fileParse := flag.String("file", "", "txt file will be converted to html file")
	directory := flag.String("dir", "", "search files in this directory")
	avatarFlag := flag.String("avatar", "Default Avatar", "add avatar to your profile")
	flag.Parse()
	if *directory != "" {
		textFiles, err := ioutil.ReadDir(*directory)
		if err != nil {
			panic(err)
		}
		var numFiles int
		for _, file := range textFiles {
			filename := file.Name()
			if checkTextFile(filename) == true {
				renderTemplate("template.tmpl", readFile(filename), *avatarFlag)
				writeTemplateToFile("template.tmpl", filename, *avatarFlag)
				numFiles++
			}
		}
	}

	if *fileParse != "" {
		renderTemplate("template.tmpl", readFile(*fileParse), *avatarFlag)
		writeTemplateToFile("template.tmpl", *fileParse, *avatarFlag)
	} else {
		renderTemplate("template.tmpl", readFile("first-post.txt"), *avatarFlag)
		writeTemplateToFile("template.tmpl", "test.txt", *avatarFlag)
	}

	if *avatarFlag != "Default Avatar" {
		fmt.Printf(*avatarFlag)
	}

	a.DebugMode = true
	a.GET("/identicons/:Name", identicon)
	a.Serve()
}
