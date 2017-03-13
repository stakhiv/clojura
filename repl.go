package main

import (
	"fmt"
	"github.com/peterh/liner"
	"io"
	"strings"
)

func StartRepl() {
	line := liner.NewLiner()
	defer line.Close()

	r := strings.NewReplacer("(", "", ")", "")
	line.SetCompleter(func(line string) (c []string) {
		word := line
		words := strings.Split(line, " ")
		if len(words) > 0 {
			word = words[len(words)-1]
		}
		word = r.Replace(word)
		for name := range coreContext.vars {
			if strings.HasPrefix(string(name), strings.ToLower(word)) {
				c = append(c, line[:len(line)-len(word)]+string(name))
			}
		}
		return
	})

	for {
		if text, err := line.Prompt(">"); err == nil {
			line.AppendHistory(text)
			lex := NewLexer(strings.NewReader(text))
			sexpr, err := NewParser(lex).Parse()
			if err != nil {
				fmt.Println("Failed parsing line: ", err)
				continue
			}
			for _, s := range sexpr {
				fmt.Println(">>", s.Eval(coreContext))
			}
		} else if err == io.EOF {
			fmt.Println("\nExiting...")
			return
		} else if err == liner.ErrPromptAborted {
			fmt.Println("Aborted")
			break
		} else {
			fmt.Println("Error reading line: ", err)
		}
	}
}
