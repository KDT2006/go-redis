package main

import (
	"log"
	"log/slog"
	"net"
)

const (
	defaultListenAddress = ":5000"
)

type Config struct {
	ListenAddress string
}

type Server struct {
	Config
	peers     map[*Peer]bool
	addPeerCh chan *Peer
	quitch    chan struct{}
	msgCh     chan []byte
	ln        net.Listener
}

func NewServer(cfg Config) *Server {
	return &Server{
		Config:    cfg,
		peers:     make(map[*Peer]bool),
		addPeerCh: make(chan *Peer),
		quitch:    make(chan struct{}),
		msgCh:     make(chan []byte),
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

func (s *Server) handleRawMessage(rawMsg []byte) error {
	cmd, err := parseCommand(string(rawMsg))
	if err != nil {
		return err
	}

	switch cmd := cmd.(type) {
	case SetCommand:
		slog.Info("SET command", "key", cmd.key, "value", cmd.val)
	}

	return nil
}

// loop for accepting peers
func (s *Server) loop() {
	for {
		select {
		case rawMsg := <-s.msgCh:
			if err := s.handleRawMessage(rawMsg); err != nil {
				slog.Error("handleRawMessage() error", "err", err)
			}

			log.Println(rawMsg)
		case <-s.quitch:
			return
		case peer := <-s.addPeerCh:
			s.peers[peer] = true
		}
	}
}

func (s *Server) handleConn(conn net.Conn) {
	peer := NewPeer(conn, s.msgCh)
	s.addPeerCh <- peer
	slog.Info("new peer connected", "remoteAddr", conn.RemoteAddr())
	if err := peer.readLoop(); err != nil {
		slog.Error("peer read error", "err", err, "remoteAddr", conn.RemoteAddr())
	}
}

func main() {
	cfg := Config{
		ListenAddress: defaultListenAddress,
	}
	server := NewServer(cfg)
	log.Fatal(server.Start())
}
