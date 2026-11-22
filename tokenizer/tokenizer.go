package tokenizer

import (
	"bytes"
)

type Tokenizer struct {
	input  *bytes.Reader
	output []string
}

func not(c byte) func(x byte) bool {
	return func(x byte) bool { return x != c }
}

func (t *Tokenizer) skipWhile(cond func(c byte) bool) error {
	for {
		c, err := t.input.ReadByte()
		if err != nil {
			return err
		}

		if cond(c) {
			break
		}
	}

	return nil
}

func (t *Tokenizer) readWhile(cond func(c byte) bool, hasCloser bool) error {
	t.input.UnreadByte()

	c, err := t.input.ReadByte()
	content := bytes.NewBuffer([]byte{c})
	for {
		c, err = t.input.ReadByte()
		if err != nil {
			return err
		}
		if !cond(c) {
			if !hasCloser {
				t.input.UnreadByte()
			} else {
				content.WriteByte(c)
			}
			break
		}
		content.WriteByte(c)
	}
	t.output = append(t.output, content.String())
	return nil
}

func (t *Tokenizer) process() error {
	c, err := t.input.ReadByte()
	if err != nil {
		return err
	}

	switch {
	case isOpener(c), isCloser(c):
		t.output = append(t.output, string(c))
	case isSpace(c):
	case isVarChar(c):
		return t.readWhile(isVarChar, false)
	case c == '"':
		return t.readWhile(not('"'), true)
	case c == ';':
		return t.skipWhile(not('\n'))
	}

	return nil
}

func Tokenize(code []byte) ([]string, error) {
	t := Tokenizer{}
	t.input = bytes.NewReader(code)

	for {
		err := t.process()
		if err != nil {
			return t.output, err
		}
	}
}

func isWordChar(c byte) bool {
	return c >= 'a' && c <= 'z' ||
		c >= 'A' && c <= 'Z' ||
		c >= '0' && c <= '9' ||
		c == '.' || c == ',' || c == '-' ||
		c == '>' || c == '<' || c == '=' ||
		c == '+' || c == '*' || c == '/'
}

func isVarChar(c byte) bool {
	return isWordChar(c) || c == '-'
}

func isOpener(c byte) bool {
	return c == '(' || c == '[' || c == '{'
}

func isCloser(c byte) bool {
	return c == ')' || c == ']' || c == '}'
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\n' || c == '\r' || c == '\t'
}
