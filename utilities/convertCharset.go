package utilities

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

func ConvertEncoding(inputfile *os.File) *transform.Reader {
	b, err := ioutil.ReadAll(inputfile)
	if err != nil {
		log.Panicln(err)
	}

	e, _, _, err := DetermineEncodingFromReader(bytes.NewReader(b))
	if err != nil {
		log.Panicln(err)
	}

	r := transform.NewReader(bytes.NewReader(b), e.NewDecoder())

	return r
}

func DetermineEncodingFromReader(r io.Reader) (e encoding.Encoding, name string, certain bool, err error) {
	b, err := bufio.NewReader(r).Peek(1024)
	if err != nil {
		return
	}

	e, name, certain = charset.DetermineEncoding(b, "")
	return
}
