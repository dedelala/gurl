package main

import (
	"bufio"
	"bytes"
	"flag"
	"io"
	"log"
	"net/http"
	"net/textproto"
	"os"

	"github.com/dedelala/xflag"
)

func main() {
	log.SetFlags(0)

	var (
		method = flag.String("X", "GET", "the request `method`")
		header = xflag.Buffer("H", "add `header`s")
		body   = xflag.InFile("d", "file for `body`, \"-\" for stdin")
		url    = xflag.Pos().String("url", "", "the request `url`")
	)
	xflag.Parse()

	bb := bytes.NewBuffer([]byte{})
	io.Copy(bb, body)

	req, err := http.NewRequest(*method, *url, bb)
	if err != nil {
		log.Fatal(err)
	}

	header.WriteString("\n")
	r := textproto.NewReader(bufio.NewReader(header))
	h, err := r.ReadMIMEHeader()
	if err != nil {
		log.Fatal(err)
	}
	req.Header = http.Header(h)

	cli := &http.Client{}
	rsp, err := cli.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(rsp.Status)

	io.Copy(os.Stdout, rsp.Body)
	rsp.Body.Close()
}
