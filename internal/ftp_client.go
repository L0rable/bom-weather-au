package internal

import (
	"log"
	"time"

	"github.com/jlaffaye/ftp"
)

const FTP_URL = "ftp.bom.gov.au:21"

func OpenFtpServer() *ftp.ServerConn {
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

func CloseFtpServer(conn *ftp.ServerConn) {
	err := conn.Quit()
	if err != nil {
		log.Fatal(err)
	}
}
