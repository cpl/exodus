package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/mr-tron/base58"
	"golang.org/x/net/dns/dnsmessage"
)

var (
	port     = flag.Int("port", 53, "Port for listening")
	verbose  = flag.Bool("v", false, "Enable verbose logging")
	dataDirF = flag.String("data", "", "Data directory, default is temp dir")
)

type entry struct {
	token string
	count int
	data  []byte
}

func (e entry) save(dir string) error {
	dirPath := path.Join(dir, e.token)
	if err := os.MkdirAll(dirPath, 0744); err != nil {
		return fmt.Errorf("failed creating dir %s, %w", dirPath, err)
	}

	filename := path.Join(dir, e.token, fmt.Sprintf("%08d.out", e.count))
	if err := ioutil.WriteFile(filename, e.data, 0644); err != nil {
		return fmt.Errorf("failed writing file %s, %w", filename, err)
	}

	return nil
}

func main() {
	flag.Parse()

	dataDir := path.Join(os.TempDir(), "exodus")
	if *dataDirF != "" {
		dataDir = *dataDirF
	}

	conn, err := net.ListenUDP("udp4", &net.UDPAddr{Port: *port})
	if err != nil {
		log.Fatalf("failed network listen, %s", err.Error())
	}
	defer conn.Close()

	if *verbose {
		log.Println("started")
		log.Println("data dir:", dataDir)
	}

	start(conn, dataDir)
}

func extractData(msg dnsmessage.Message) (e entry, err error) {
	if len(msg.Questions) != 1 {
		return e, fmt.Errorf("dns query has no questions")
	}

	domains := strings.Split(msg.Questions[0].Name.String(), ".")
	if len(domains) < 5 {
		return e, fmt.Errorf("dns question is missing requiered information")
	}

	e.data, err = base58.Decode(domains[0])
	if err != nil {
		return e, fmt.Errorf("failed decoding data, %w", err)
	}

	e.count, _ = strconv.Atoi(domains[1])
	e.token = domains[2]

	return
}

func start(conn *net.UDPConn, dataDir string) {
	for {
		var (
			buffer [512]byte
			msg    dnsmessage.Message
		)

		read, addr, err := conn.ReadFromUDP(buffer[:])
		if err != nil {
			log.Printf("failed network read, %s\n", err.Error())
			continue
		}

		if err := msg.Unpack(buffer[:read]); err != nil {
			log.Printf("failed message unpacking, %s\n", err.Error())
			continue
		}

		e, err := extractData(msg)
		if err != nil {
			log.Println(err)
			continue
		}

		if *verbose {
			log.Printf("received [token: %16s] [count: %8d]\n", e.token, e.count)
		}

		if err := e.save(dataDir); err != nil {
			log.Println(err)
			continue
		}

		msg.Header.Response = true

		packed, err := msg.Pack()
		if err != nil {
			log.Printf("failed packing dns msg, %s\n", err.Error())
		}

		if _, err := conn.WriteToUDP(packed, addr); err != nil {
			log.Printf("failed network write, %s\n", err.Error())
		}
	}
}
