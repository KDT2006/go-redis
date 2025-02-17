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

		var cmd Command
		if v.Type() == resp.Array {
			rawCMD := v.Array()[0]
			switch rawCMD.String() {
			case CommandClient:
				cmd = ClientCommand{
					value: v.Array()[1].String(),
				}
			case CommandSET:
				cmd = SetCommand{
					key: v.Array()[1].Bytes(),
					val: v.Array()[2].Bytes(),
				}

				fmt.Printf("got SET command %+v\n", cmd)

			case CommandGET:
				cmd = GetCommand{
					key: v.Array()[1].Bytes(),
				}

				fmt.Printf("got GET command %+v\n", cmd)

			case CommandHello:
				cmd = HelloCommand{
					value: v.Array()[1].String(),
				}

			default:
				fmt.Printf("got unknown command: %+v\n", v.Array())
			}

			p.msgCh <- Message{
				cmd:  cmd,
				peer: p,
			}
		}
	}

	return nil
}
