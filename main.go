package main

import (
	"flag"
	slog "log"
	"os"
	"time"
)

var log = NewLogger(Info)

func main() {
	flag.Parse()
	if len(flag.Args()) < 1 {
		slog.Fatalln("Missing file arg")
	}
	err := LoadFile("core.clj")
	if err != nil {
		slog.Fatalln("Failed to load 'core.cj'", err)
	}

	fname := flag.Args()[0]
	err = LoadFile(fname)
	if err != nil {
		slog.Fatalln("Failed to load '%s'", fname, err)
	}
}

func LoadFile(name string) error {
	start := time.Now()
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	sexpr := make([]Sexpr, 0)

	lexer := NewLexer(f)
	parser := NewParser(lexer)
	sexpr, err = parser.Parse()
	if err != nil {
		return err
	}
	log.Debugf("Parsed '%s' in ", name, time.Since(start))
	// Evaluate
	for _, s := range sexpr {
		s.Eval(coreContext)
	}
	return nil
}
