package main

import (
	"bytes"
	"flag"
	"io"
	slog "log"
	"os"
	"time"
)

var log = NewLogger(Info)

func main() {
	flag.Parse()
	err := Exec(bytes.NewReader([]byte(initData)))
	if err != nil {
		slog.Fatalln("Failed to load 'core.cj'", err)
	}
	if len(flag.Args()) < 1 {
		StartRepl()
		return
	}

	fname := flag.Args()[0]
	err = LoadFile(fname)
	if err != nil {
		slog.Fatalln("Failed to load '%s'", fname, err)
	}
}

func LoadFile(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	return Exec(f)
}

func Exec(r io.Reader) error {
	lexer := NewLexer(r)
	parser := NewParser(lexer)
	start := time.Now()
	sexpr, err := parser.Parse()
	if err != nil {
		return err
	}
	log.Debug("Parsed in ", time.Since(start))
	// Evaluate
	for _, s := range sexpr {
		s.Eval(coreContext)
	}
	return nil
}
