package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/tidwall/resp"
)

type Peer struct {
	conn      net.Conn
	msgCh     chan<- Message
	delPeerCh chan<- *Peer
}

func NewPeer(conn net.Conn, msgCh chan<- Message, delPeerCh chan<- *Peer) *Peer {
	return &Peer{
		conn:      conn,
		msgCh:     msgCh,
		delPeerCh: delPeerCh,
	}
}

func (p *Peer) Write(msg []byte) (int, error) {
	return p.conn.Write(msg)
}

func (p *Peer) readLoop() error {
	rd := resp.NewReader(p.conn)
	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			// the conn is closed by the peer, so remove the peer from the server
			p.delPeerCh <- p
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Printf("Read %s\n", v.Type())
		if v.Type() == resp.Array {
			for _, value := range v.Array() {
				switch value.String() {
				case CommandSET:
					// fmt.Println(len(v.Array()))
					if len(v.Array()) != 3 {
						return fmt.Errorf("invalid number of arguments for SET command")
					}

					cmd := SetCommand{
						key: v.Array()[1].Bytes(),
						val: v.Array()[2].Bytes(),
					}

					p.msgCh <- Message{
						cmd:  cmd,
						peer: p,
					}

					fmt.Printf("got SET command %+v\n", cmd)

				case CommandGET:
					if len(v.Array()) != 2 {
						return fmt.Errorf("invalid number of arguments for GET command")
					}

					cmd := GetCommand{
						key: v.Array()[1].Bytes(),
					}

					p.msgCh <- Message{
						cmd:  cmd,
						peer: p,
					}

					fmt.Printf("got GET command %+v\n", cmd)
				}
			}
		}
	}

	return nil
}
