package main

import (
	"io"
	"log"
	"time"

	"github.com/jlaffaye/ftp"
)

const FTP_URL = "ftp.bom.gov.au:21"

func openFtpServer() *ftp.ServerConn {
	conn, err := ftp.Dial(FTP_URL, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		log.Fatal(err)
	}

	err = conn.Login("anonymous", "")
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

func closeFtpServer(conn *ftp.ServerConn) {
	err := conn.Quit()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	conn := openFtpServer()
	obsURL := "/anon/gen/fwo/IDV60920.xml"
	resp, err := conn.Retr(obsURL)
	if err != nil {
		log.Fatal(err)
	}

	data, err := io.ReadAll(resp)
	if err != nil {
		log.Println(data)
	}
	log.Println(string(data))

	closeFtpServer(conn)
}
