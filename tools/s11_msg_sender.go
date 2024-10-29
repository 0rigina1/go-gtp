package main

import (
	"flag"
	"log"
	"net"
	"os"

	"github.com/wmnsk/go-gtp/gtpv2/ie"
	"github.com/wmnsk/go-gtp/gtpv2/message"
)

func createConnection(mmeAddr string) (*net.UDPConn, error) {
	addr, err := net.ResolveUDPAddr("udp", mmeAddr)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func generate_pgwRestartNotification_message(pgw, sgw string) *message.PGWRestartNotification {
	return message.NewPGWRestartNotification(
		0x00000000,
		0x0000cb,
		ie.NewIPAddress(pgw).WithInstance(0),
		ie.NewIPAddress(sgw).WithInstance(1),
		// ie.NewCause(gtpv2.CausePGWNotResponding, 0, 0, 0, nil),
	)
}

func generate_echoRequest_message() *message.EchoRequest {
	return message.NewEchoRequest(
		0x000002,
	)
}

func send_echoRequest_message(conn *net.UDPConn, echo *message.EchoRequest) error {
	serialized, err := echo.Marshal()
	if err != nil {
		return err
	}

	_, err = conn.Write(serialized)
	return err
}

func send_pgwRestartNotification_message(conn *net.UDPConn, prn *message.PGWRestartNotification) error {
	serialized, err := prn.Marshal()
	if err != nil {
		return err
	}

	_, err = conn.Write(serialized)
	return err
}

func main() {
	mme := flag.String("mme", "", "mme ip:host")
	pgw := flag.String("pgw", "", "pgw ip")
	sgw := flag.String("sgw", "", "sgw ip")
	msg_type := flag.String("type", "echo", "message type. (prn,echo)")

	flag.Parse()

	switch *msg_type {
	case "prn":
		if *mme == "" || *pgw == "" || *sgw == "" {
			log.Fatalf("For 'prn' message type, MME, PGW, and SGW params are required.")
			flag.Usage()
			os.Exit(1)
		}
	case "echo":
		if *mme == "" {
			log.Println("For 'echo' message type, MME params is required.")
			flag.Usage()
			os.Exit(1)
		}
	default:
		log.Println("Invalid message type. Choose either 'prn' or 'echo'.")
		flag.Usage()
		os.Exit(1)
	}

	conn, err := createConnection(*mme)
	if err != nil {
		log.Fatalf("failed to connect to mme: %s", err)
	}
	defer conn.Close()

	switch *msg_type {
	case "prn":
		prn := generate_pgwRestartNotification_message(*pgw, *sgw)
		if err := send_pgwRestartNotification_message(conn, prn); err != nil {
			log.Fatalf("failed to send PGW Restart Notification message: %s", err)
		}
		log.Println("PGW Restart Notification message sent successfully!")

	case "echo":
		echo := generate_echoRequest_message()

		if err := send_echoRequest_message(conn, echo); err != nil {
			log.Fatalf("failed to send Echo Request message: %s", err)
		}

		log.Println("Echo Request message sent successfully!")
	}
}
