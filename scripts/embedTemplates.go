/*
 * Copyright © 2018 Rasmus Hansen
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */

package main

import (
	"fmt"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
	"github.com/tdewolff/minify/json"
	"github.com/tdewolff/minify/svg"
	"github.com/tdewolff/minify/xml"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

const (
	templateLocation = "web/template"
	embedFileName    = "internal/pkg/templates/templates.go"

	outputTemplate = `// Code Generated by go run scripts/embedTemplates DO NOT EDIT.

/*
 * Copyright © 2018 Rasmus Hansen
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */

package templates

import "html/template"

type templateContent struct{
	Filename string
	Content string
}

var contents = []templateContent{
	{{range .TemplateFiles}}
	templateContent{
		Filename:"{{.Filename}}", 
		Content:{{.Content}},
	},
	{{end}}
}


func getTemplates() (*template.Template, error) {
	var t *template.Template
	
	for _, file := range contents {

		filename := file.Filename
		var tmpl *template.Template
		if t == nil {
			t = template.New(filename)
		}
		
		if filename == t.Name() {
			tmpl = t
		} else {
			tmpl = t.New(filename)
		}
		_, err := tmpl.Parse(file.Content)
		if err != nil {
			return nil, err
		}
	}

	return t, nil
}`
)

func main() {

	err := createEmbedTemplates()
	if err != nil {
		log.Fatal(err)
	}
}

type templateFile struct {
	Filename string
	Content  string
}

type templateArgs struct {
	TemplateFiles []templateFile
}

func getTemplateFiles() ([]string, error) {
	matches := make([]string, 0)
	err := filepath.Walk(templateLocation, func(path string, info os.FileInfo, err error) error {
		isMatch, err := filepath.Match("*.go.html", info.Name())
		if err != nil {
			return err
		}

		if isMatch {
			matches = append(matches, path)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return matches, err
}

func createEmbedTemplates() error {

	minifer := getMinifier()

	matches, err := getTemplateFiles()
	if err != nil {
		log.Fatal(err)
	}

	tmpl := template.Must(template.New("templates").Parse(outputTemplate))

	ta := templateArgs{
		TemplateFiles: []templateFile{},
	}

	fmt.Println(matches)

	for _, match := range matches {
		b, err := ioutil.ReadFile(match)
		if err != nil {
			return err
		}

		minified, err := minifer.Bytes("text/html", b)
		if err != nil {
			return err
		}

		content := strconv.Quote(string(minified))

		filename := strings.Replace(filepath.Base(match), ".go.html", "", 1)

		ta.TemplateFiles = append(ta.TemplateFiles, templateFile{
			Filename: filename,
			Content:  content,
		})
	}

	log.Println(len(ta.TemplateFiles))

	file, err := os.Create(embedFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, ta)
}

func getMinifier() *minify.M {
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("text/javascript", js.Minify)
	m.AddFunc("image/svg+xml", svg.Minify)
	m.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	m.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)

	return m
}
