package im

import (
	"bufio"
	"bytes"

	"io"

	"github.com/pkg/errors"
)

type Parser struct {
	r            *bufio.Reader
	line         int //not useful
	bufferedItem *item
}

func (p *Parser) readUnfoldedLine() (string, error) {
	buf, e := p.r.ReadBytes('\r')
	if e != nil {
		return "", e
	}

	b, err := p.r.ReadByte()
	if err != nil {
		return "", err
	}
	if b != '\n' {
		return "", errors.New("Expected CRLF:" + string(buf))
	}
	b1, err2 := p.r.Peek(1)

	if err2 != nil {
		return string(buf[:len(buf)-1]), err2
	}
	if bytes.Equal(b1, []byte(" ")) {
		p.r.ReadByte()
		s, e := p.readUnfoldedLine()
		if s == "" {
			return "", e
		}
		return string(buf[:len(buf)-1]) + s, e
	}
	return string(buf[:len(buf)-1]), nil
}

func InitParser(reader io.Reader) Parser {
	return Parser{bufio.NewReader(reader), 0}
}

func (p *Parser) ParseComponent() (component *Component, err error) {

}
