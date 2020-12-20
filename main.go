package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type (
	// Lox language instance
	Lox struct {
		hasError bool
	}
)

// ErrScanning returned when an error has been encountered while scanning the source code
var ErrScanning = errors.New("Error encountered while scanning")

func main() {
	lox := Lox{}
	if len(os.Args) == 1 {
		lox.runPrompt()
	} else if len(os.Args) == 2 {
		err := lox.runFile(os.Args[1])
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(65)
		}
	} else {
		fmt.Println("Usage: glox [script]")
		os.Exit(64)
	}
}

func (lox *Lox) runFile(filepath string) error {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	err = lox.run(string(bytes))
	if err != nil {
		return err
	}
	if lox.hasError {
		return ErrScanning
	}
	return nil
}

func (lox *Lox) run(source string) error {
	scanner := NewScanner(source, lox)
	tokens := scanner.ScanTokens()

	parser := NewParser(tokens, lox)
	expr := parser.Parse()
	if expr != nil {
		fmt.Println(NewAstPrinter().Print(expr))
	}

	return nil
}

func (lox *Lox) runPrompt() {
	for {
		var line string
		fmt.Print("> ")
		fmt.Scanln(&line)
		if len(line) == 0 {
			break
		}
		err := lox.runFile(line)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			lox.hasError = false
		}
	}
}

// Exception reports an error
func (lox *Lox) Exception(line int, message string) {
	lox.report(line, "", message)
}

func (lox *Lox) report(line int, where, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error %s: %s\n", line, where, message)
	lox.hasError = true
}
