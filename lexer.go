package main

import (
	"bufio"
	"errors"
	"io"
)

const (
	LBrace       = '('
	RBrace       = ')'
	RSquareBrace = ']'
	LSquareBrace = '['
	RCurlyBrace  = '}'
	LCurlyBrace  = '{'
	CommentStart = ';'
	Space        = ' '
	Tick         = '\''
	Newline      = '\n'
	Tab          = '\t'
)

var separators = map[rune]bool{
	LBrace:       true,
	RBrace:       true,
	RSquareBrace: true,
	LSquareBrace: true,
	RCurlyBrace:  true,
	LCurlyBrace:  true,
	Space:        true,
	Tab:          true,
	Newline:      true,
}

var tokens = map[rune]bool{
	LBrace:       true,
	RBrace:       true,
	RSquareBrace: true,
	LSquareBrace: true,
	RCurlyBrace:  true,
	LCurlyBrace:  true,
}

type Lexer struct {
	reader *bufio.Reader
}

func NewLexer(r io.Reader) *Lexer {
	return &Lexer{
		reader: bufio.NewReader(r),
	}
}

func (l *Lexer) ReadToken() (string, error) {
	r, _, err := l.reader.ReadRune()
	if err != nil {
		return "", err
	}
	for {
		if isToken(r) {
			return string(r), nil
		} else if isWhitespace(r) {
			r, err = l.drainWhitespace()
			if err != nil {
				return "", err
			}
			continue
		} else if r == '"' {
			s, err := l.readWhile('"')
			if err != nil {
				return "", err
			}
			return "\"" + s + "\"", err
		} else if r == ';' {
			r, err = l.drainWhile('\n')
			if err != nil {
				return "", err
			}
			continue
		} else {
			token, err := l.readToken()
			if err != nil {
				return "", err
			}
			return string(r) + token, nil
		}
	}

	return "", errors.New("how did we get here?")
}

func (l *Lexer) readToken() (string, error) {
	defer l.reader.UnreadRune()
	var res string
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			return "", err
		}
		if isSeparator(r) {
			break
		}
		res += string(r)
	}
	return res, nil
}

func (l *Lexer) readWhile(cr rune) (string, error) {
	var res string
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			return "", err
		}
		if r == cr {
			break
		}
		res += string(r)
	}
	return res, nil
}

func (l *Lexer) drainWhile(t rune) (r rune, err error) {
	for {
		r, _, err = l.reader.ReadRune()
		if err != nil {
			return 0, err
		}
		if r == t {
			break
		}
	}
	return
}

func (l *Lexer) drainWhitespace() (r rune, err error) {
	for {
		r, _, err = l.reader.ReadRune()
		if err != nil {
			return 0, err
		}
		if !isWhitespace(r) {
			break
		}
	}
	return
}

func isToken(r rune) bool {
	return tokens[r]
}

func isSeparator(r rune) bool {
	return separators[r]
}

func isWhitespace(r rune) bool {
	if r == ' ' || r == '\n' || r == '\t' {
		return true
	}
	return false
}
