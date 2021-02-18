package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/smartystreets/redact"
)

func main() {
	redactor := redact.New()

	buffer := bufio.NewReaderSize(os.Stdin, 1024*1024*16)

	for {
		line, _, err := buffer.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		items := bytes.Split(line, []byte("|"))
		address := fmt.Sprintf("%s %s %s %s %s %s", items[15], items[16], items[17], items[18], items[19], items[20])
		redacted := redactor.All(address)
		if address != redacted && !strings.Contains(address, "@"){
			fmt.Printf("[%s] [%s] \n", address, redacted)
		}
	}
}