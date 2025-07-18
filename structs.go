package main

type conf struct {
	environmentFile   string
	templateFilePath  string
	stdinTemplateFile string
	valuesFile        string
	strictMode        bool
	customDelimiters  string
	setValues         []string
	setStringValues   []string
}
