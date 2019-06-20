package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

var i = flag.Int("i", 2, "group that represents variable name")
var g = flag.String("d", "*.conf", "gob pattern for selecting files to modify")
var p = flag.String("p", "", "prefix for env variables")
var e = flag.String("e", `^(\s*)set\s*\$(\w*)\s*("(.*?)")?\s*;\s*$`, "regular expression to search for")
var r = flag.String("r", `${1}set $$$2 "%s";`, "replacement template (use %s where you want value)")

func main() {
	flag.Parse()

	template.New("name")
	re := regexp.MustCompile(*e)

	// ioutil.ReadFile
	// for some directory
	files, err := filepath.Glob(*g)
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		fmt.Print(f)
		file, err := os.OpenFile(f, os.O_RDONLY, os.ModePerm)
		if err != nil {
			panic(err)
		}

		tmp, err := ioutil.TempFile("", "*.tmp")
		if err != nil {
			panic(err)
		}

		scanner := bufio.NewScanner(file)
		writter := bufio.NewWriter(tmp)

		var variables []string
		for scanner.Scan() {
			txt := scanner.Text()
			matches := re.FindAllStringSubmatchIndex(txt, -1)
			if len(matches) == 0 {
				writter.WriteString(txt)
				writter.WriteString("\n")
				continue
			}
			for _, match := range matches {
				varName := txt[match[*i*2]:match[*i*2+1]]

				e := envName(*p, varName)
				if v, ok := os.LookupEnv(e); ok {
					variables = append(variables, varName+"*")
					t := re.ReplaceAllString(txt, *r)
					writter.WriteString(fmt.Sprintf(t, v))
					writter.WriteString("\n")
				} else {
					variables = append(variables, varName)
					writter.WriteString(txt)
					writter.WriteString("\n")
				}
			}
		}

		fmt.Println("[" + strings.Join(variables, ", ") + "]")

		if err := scanner.Err(); err != nil {
			fmt.Println("error:", err)
		}

		if err := writter.Flush(); err != nil {
			panic(err)
		}

		if err := os.Rename(tmp.Name(), file.Name()); err != nil {
			panic(err)
		}

		file.Close()
		tmp.Close()

	}
}

var re = regexp.MustCompile(`\pP`)

func envName(prefix, s string) string {
	in := prefix + s
	// fmt.Println("envName:", s)
	upper := strings.ToUpper(in)
	return re.ReplaceAllString(upper, "_")
}
