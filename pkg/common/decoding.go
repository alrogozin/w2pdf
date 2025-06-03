package common

import (
	"bytes"
	"io"
	"log"

	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/transform"
)

func DoProcess(p_inp string, p_charset string) string {
	// text := "\x8f\xe0\xae⮪\xae\xab ࠧ\xad\xae\xa3\xab\xa0ᨩ \xaa \x84\xae\xa3\xae\xa2\xae\xe0\xe3 \xfc19609 \x84\x8f-\x82 \xe1 \xee\xe0\xa8\xe1⠬\xa8.docx"
	// charset := "866"
	// us-ascii
	// latin1 == iso-8859-1
	// windows-1252.
	e, err := ianaindex.MIME.Encoding(p_charset)
	if err != nil {
		log.Fatal(err)
	}
	r := transform.NewReader(bytes.NewBufferString(p_inp), e.NewDecoder())
	result, err := io.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("%s\n", result)
	return string(result)
}
