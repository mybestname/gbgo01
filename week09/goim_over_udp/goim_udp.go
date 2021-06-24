package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
	"goim_udp/message"
	"github.com/Terry-Mao/goim/api/protocol"
	proto_old_api "github.com/golang/protobuf/proto"
	proto_new_api "google.golang.org/protobuf/proto"
)

func usage() {
	fmt.Println("usage:[client|server]")
}

func main() {
	if len(os.Args[1:]) < 1 {
		usage();
		return
	}
	t := os.Args[1]
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	if t == "server" {
		startServer(ctx);
	} else if t == "client" {
		msg := "hello"
		if len(os.Args) == 4 && os.Args[2] == "-msg" {
		 	msg = fmt.Sprintf("%s", os.Args[3])
		}
		startClient(ctx, msg);
	} else {
		fmt.Printf("error arg %s\n", t)
		usage();
		return
	}
}

func check(err error, errChan chan error){
	if err != nil {
		errChan<-err
	}
}

const maxBufferSize = 1024
const timeout = 15 * time.Second

// AuthToken auth token.
type AuthToken struct {
	Mid      int64   `json:"mid"`
	Key      string  `json:"key"`
	RoomID   string  `json:"room_id"`
	Platform string  `json:"platform"`
	Accepts  []int32 `json:"accepts"`
}

const (
	opHeartbeat      = int32(2)
	opHeartbeatReply = int32(3)
	opAuth           = int32(7)
	opAuthReply      = int32(8)
)

func startServer(ctx context.Context) {
	errChan := make(chan error, 1)
	log.Println("start server")
	pc, err := net.ListenPacket("udp", ":2000")
	check(err,errChan)
	defer pc.Close()
	buffer := make([]byte, maxBufferSize)
	go func() {
		for {
			n, addr, err := pc.ReadFrom(buffer)
			check(err,errChan)
			if err == nil {
				log.Printf("packet received byte=%d from %s\n", n, addr)
			}
			log.Printf("%s\n",buffer)
			deadline := time.Now().Add(timeout)
			err = pc.SetWriteDeadline(deadline)
			check(err,errChan)
			cmd := fmt.Sprintf("%s",buffer[:n])
			var data []byte
			switch cmd {
			case "hello":
				messageProto := message.Message{Text: "Hello World", Timestamp: time.Now().Unix()}
				data, err = proto_new_api.Marshal(&messageProto)
				check(err,errChan)
				n,err = pc.WriteTo(data,addr)
			case "auth":
				seq := int32(0)
				authToken := &AuthToken{
					time.Now().Unix(),
					"",
					"test://1",
					"ios",
					[]int32{1000, 1001, 1002},
				}
				body, err := json.Marshal(authToken)
				check(err, errChan)
				msg := protocol.Proto{Ver: int32(1), Op: opAuth, Seq: seq, Body: body}
				data, err = proto_old_api.Marshal(&msg)
				check(err,errChan)
			default:
				data = buffer[:n]
			}
			n, err = pc.WriteTo(data, addr)
			check(err,errChan)
			if err == nil {
				log.Printf("packet writen byte=%d to %s\n", n, addr)
			}
		}
	}()
	select {
	case <-ctx.Done():
		log.Println("server canceled")
	case err = <-errChan:
		log.Printf("server stopped since error: %v\n",err)
	}
	log.Println("server done")
}

func startClient(ctx context.Context, msg string) {
	log.Println("start client")
	errChan := make(chan error, 1)
	conn, err := net.ListenPacket("udp", ":0")
	check(err,errChan)
	defer conn.Close()
	dst, err := net.ResolveUDPAddr("udp", "127.0.0.1:2000")
	check(err,errChan)

	go func() {
		n, err := conn.WriteTo([]byte(msg), dst)
		check(err,errChan)
		fmt.Printf("packet-written: bytes=%d\n", n)
		buffer := make([]byte, maxBufferSize)
		deadline := time.Now().Add(timeout)
		err = conn.SetReadDeadline(deadline)
		check(err,errChan)
		nRead, addr, err := conn.ReadFrom(buffer)
		check(err,errChan)
		fmt.Printf("packet-received: bytes=%d from=%s\n",
			nRead, addr.String())
		switch(msg) {
		case "hello":
			messagePb := message.Message{}
			err = proto_new_api.Unmarshal(buffer[:nRead], &messagePb)
			check(err,errChan)
			log.Printf("received message: %s, timestamp: %v", messagePb.Text, messagePb.Timestamp)
		case "auth":
			protoPb := protocol.Proto{}
			err = proto_old_api.Unmarshal(buffer[:nRead], &protoPb)
			check(err,errChan)
			var token AuthToken;
			json.Unmarshal(protoPb.Body,&token)
			log.Printf("received goim proto message: ver %d, op %d, seq %d, body %v\n", protoPb.Ver, protoPb.Op, protoPb.Seq, token)
		default:
			log.Printf("received raw message %s",buffer[:nRead])
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("client canceled")
	case err = <-errChan:
		log.Printf("client stopped since error: %s\n",err)
	}
	log.Println("client done")
}






