package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	if len(os.Args) == 1 {
		runPrompt()
	} else if len(os.Args) == 2 {
		err := runFile(os.Args[1])
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(65)
		}
	} else {
		fmt.Println("Usage: glox [script]")
		os.Exit(64)
	}
}

func runFile(filepath string) error {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	return run(string(bytes))
}

func run(source string) error {
	reader := strings.NewReader(source)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanRunes)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	return nil
}

func runPrompt() {
	for {
		var line string
		fmt.Print("> ")
		fmt.Scanln(&line)
		if len(line) == 0 {
			break
		}
		err := runFile(line)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
	}
}
