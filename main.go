package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

func main() {
	var (
		in  = flag.String("in", "", "file to parse instead of stdin")
		out = flag.String("out", "", "file to save output to instead of stdout")
	)
	flag.Parse()

	outWriter := os.Stdout
	if len(*out) > 0 {
		err := os.MkdirAll(path.Dir(*out), 0755)
		if err != nil {
			panic(err)
		}

		outFile, err := os.Create(*out)
		if err != nil {
			panic(err)
		}
		defer outFile.Close()
		outWriter = outFile
	}

	inReader := os.Stdin
	if len(*in) > 0 {
		inFile, err := os.Open(*in)
		if err != nil {
			panic(err)
		}
		defer inFile.Close()
		inReader = inFile
	}

	// build replacement mapping
	repl := make(map[string][]string)
	var lb []string
	s := bufio.NewScanner(inReader)
	for s.Scan() {
		sl := s.Text()
		if strings.HasPrefix(sl, "//genify:") {
			for _, f := range strings.Fields(strings.TrimPrefix(sl, "//genify:")) {
				p := strings.SplitN(f, "=", 2)
				if len(p) < 2 {
					panic(fmt.Sprintf("Invalid comment: %v", sl))
				}
				repl[p[0]] = strings.Split(p[1], ",")
			}
			continue
		}
		//fmt.Println("repl: %v", repl)

		lb = append(lb, sl)
		if len(sl) == 0 && len(lb) > 1 && len(lb[len(lb)-2]) > 0 && !strings.HasPrefix(lb[len(lb)-2], "\t") { // this is the end of a block
			processLines(lb, repl, outWriter)
			lb = []string{}
		}
	}
	if len(lb) > 0 {
		processLines(lb, repl, outWriter)
	}
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

func uncapitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToLower(string(s[0])) + s[1:]
}

func processLines(lb []string, repl map[string][]string, outWriter io.Writer) {
	ls := strings.Join(lb, "\n") + "\n"
	var m bool
	for sw, rws := range repl {
		if !strings.Contains(ls, capitalize(sw)) && !strings.Contains(ls, uncapitalize(sw)) {
			continue
		}
		m = true
		for _, rw := range rws {
			ns := strings.Replace(ls, capitalize(sw), capitalize(rw), -1)
			ns = strings.Replace(ns, uncapitalize(sw), uncapitalize(rw), -1)
			if _, err := io.WriteString(outWriter, ns); err != nil {
				panic(err)
			}
		}
	}
	if !m {
		if _, err := io.WriteString(outWriter, ls); err != nil {
			panic(err)
		}
	}
}
