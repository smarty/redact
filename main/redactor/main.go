package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/smartystreets/redact"
)

func main() {
	redactor := redact.New()

	//buffer := bufio.NewReaderSize(os.Stdin, 1024*1024*16)
	file, err := os.Open("/Users/savannah/SmartyStreets/production-data-redact/test-data_1.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	buffer := bufio.NewReader(file)
	for {
		line, _, err := buffer.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		items := bytes.Split(line, []byte("|"))
		address := fmt.Sprintf("%s %s %s %s %s %s", items[15], items[16], items[17], items[18], items[19], items[20])
		redacted := redactor.RedactAll([]byte(address))
		if address != string(redacted) && !strings.Contains(address, "@") {
			fmt.Printf("[%s] [%s] \n", address, redacted)
		}
	}
}
