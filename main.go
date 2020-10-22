package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)


func readName(path string, name string) string{
	roamDir := os.Getenv("ROAMDIR")
	if roamDir == ""{
		roamDir = "/Users/hjertnes/txt/roam"
	}

	file, err := ioutil.ReadFile(fmt.Sprintf("%s%s/%s", roamDir, strings.Replace(path, ".", "/", 1), strings.ReplaceAll(name, ".html", ".org")))
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(file), "\n");
	for _, l := range lines{
		if strings.HasPrefix(strings.ToLower(l), "#+title: "){
			return strings.Split(l, ": ")[1]
		}
	}

	return ""
}



func crawl(path string, builder *strings.Builder){

	files, _ := ioutil.ReadDir(path)
	for _, i := range files {
		if i.Name() == "index.html"{
			continue
		}
		if i.IsDir(){
			builder.WriteString(fmt.Sprintf(`<li><a href="%s%s">%s</a></li>`,path, i.Name(), i.Name()))
			builder.WriteString("<ul>")
			crawl(fmt.Sprintf("%s/%s", path, i.Name()), builder)
			builder.WriteString("</ul>")
		} else {
			name := readName(path, i.Name())
			if name == ""{
				name = i.Name()
			}
			builder.WriteString(fmt.Sprintf(`<li><a href="%s/%s">%s</a></li>`, path, i.Name(), name))
		}

	}
}

func main(){
	path := "/"

	if len(os.Args) > 1{
		path = os.Args[1]
	}

	builder := strings.Builder{}
	builder.WriteString(`<ul>`)
	crawl(path, &builder)
	builder.WriteString(`</ul>`)

	f, err := ioutil.ReadFile("../index-template.html")
	if err != nil{
		panic(err)
	}
	template := string(f)

	fmt.Print(strings.Replace(template, "%BODY%", builder.String(), 1))
}
