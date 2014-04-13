package main

import (
	"bytes"
	"fmt"
	"github.com/laher/uggo"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)
const(
	VERSION = "0.0.1"
)

func main() {
	err := platego(os.Args)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}

func platego(call []string) error {
	flagSet := uggo.NewFlagSetDefault("platego", "[options] [key=val...]", VERSION)
	var tpl string
	var outputTempl string
	flagSet.AliasedStringVar(&tpl, []string{"t", "template"}, "template.tpl", "Template")
	flagSet.AliasedStringVar(&outputTempl, []string{"o", "output"}, "", "Output file (also a template)")
	err := flagSet.Parse(call[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Flag error:  %v\n\n", err.Error())
		flagSet.Usage()
		return err
	}
	if flagSet.ProcessHelpOrVersion() {
		return nil
	}

	args := flagSet.Args()
	t, err := template.ParseFiles(tpl)
	if err != nil {
		return err
	}

	data := map[string]string{}
	for _, a := range args {
		parts := strings.Split(a, "=")
		data[parts[0]] = parts[1]
		data[parts[0]+"UCF"] = strings.ToUpper(parts[1][0:1]) + parts[1][1:]
	}

	var out io.Writer
	if outputTempl != "" {
		ot, err := template.New("out").Parse(outputTempl)
		if err != nil {
			return err
		}
		var bout bytes.Buffer
		err = ot.Execute(&bout, data)
		if err != nil {
			return err
		}
		output := bout.String()
		//make dir(s) if necessary
		dirname := filepath.Dir(output)
		if _, err = os.Stat(dirname); os.IsNotExist(err) {
			err = os.MkdirAll(dirname, 0777)
			if err != nil {
				return err
			}
		}
		outputFile, err := os.Create(output)
		defer outputFile.Close()
		if err != nil {
			return err
		}
		out = outputFile
	} else {
		out = os.Stdout
	}
	return t.Execute(out, data)
}
