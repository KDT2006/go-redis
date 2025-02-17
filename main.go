package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"

	"github.com/tidwall/resp"
)

const (
	defaultListenAddress = ":5000"
)

type Config struct {
	ListenAddress string
}

type Message struct {
	cmd  Command
	peer *Peer
}

type Server struct {
	Config
	peers     map[*Peer]bool
	addPeerCh chan *Peer
	quitch    chan struct{}
	msgCh     chan Message
	delPeerCh chan *Peer
	ln        net.Listener

	kv *KV
}

func NewServer(cfg Config) *Server {
	return &Server{
		Config:    cfg,
		peers:     make(map[*Peer]bool),
		addPeerCh: make(chan *Peer),
		quitch:    make(chan struct{}),
		msgCh:     make(chan Message),
		delPeerCh: make(chan *Peer),
		kv:        NewKV(),
	}
}

func (s *Server) Start() error {
	if len(s.ListenAddress) == 0 {
		s.ListenAddress = defaultListenAddress
	}
	ln, err := net.Listen("tcp", s.ListenAddress)
	if err != nil {
		return err
	}
	s.ln = ln

	go s.loop()

	slog.Info("server runnning", "listenaddr", s.ListenAddress)

	return s.acceptLoop()
}

func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			slog.Error("accept error", "err", err)
			continue
		}

		go s.handleConn(conn)
	}
}

func (s *Server) handleMessage(msg Message) error {
	switch cmd := msg.cmd.(type) {
	case ClientCommand:
		err := resp.NewWriter(msg.peer.conn).WriteString("OK")
		if err != nil {
			return err
		}
	case SetCommand:
		if err := s.kv.Set(cmd.key, cmd.val); err != nil {
			return err
		}

		err := resp.NewWriter(msg.peer.conn).WriteString("OK")
		if err != nil {
			return err
		}
	case GetCommand:
		val, ok := s.kv.Get(cmd.key)
		if !ok {
			return fmt.Errorf("key %s not found", cmd.key)
		}

		err := resp.NewWriter(msg.peer.conn).WriteString(string(val))
		if err != nil {
			return err
		}
	case HelloCommand:
		spec := map[string]string{
			"server": "redis",
		}

		_, err := msg.peer.Write(writeRespMap(spec))
		if err != nil {
			return fmt.Errorf("error writing to peer: %s", err)
		}
	}

	return nil
}

// loop for accepting peers and messages
func (s *Server) loop() {
	for {
		select {
		case msg := <-s.msgCh:
			if err := s.handleMessage(msg); err != nil {
				slog.Error("handleRawMessage() error", "err", err)
			}

			// log.Println(rawMsg)
		case <-s.quitch:
			return
		case peer := <-s.addPeerCh:
			slog.Info("new peer connected", "remoteAddr", peer.conn.RemoteAddr())
			s.peers[peer] = true
		case peer := <-s.delPeerCh:
			slog.Info("peer disconnected", "remoteAddr", peer.conn.RemoteAddr())
			delete(s.peers, peer)
		}
	}
}

func (s *Server) handleConn(conn net.Conn) {
	peer := NewPeer(conn, s.msgCh, s.delPeerCh)
	s.addPeerCh <- peer
	// slog.Info("new peer connected", "remoteAddr", conn.RemoteAddr())
	if err := peer.readLoop(); err != nil {
		slog.Error("peer read error", "err", err, "remoteAddr", conn.RemoteAddr())
	}
}

func main() {
	listenAddr := flag.String("listenAddr", defaultListenAddress, "listen address for the go-redis server")
	flag.Parse()
	if *listenAddr == "" {
		*listenAddr = defaultListenAddress
	}
	cfg := Config{
		ListenAddress: *listenAddr,
	}

	server := NewServer(cfg)
	log.Fatal(server.Start())
}
