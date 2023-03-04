package api

import (
	"net/http"

	"github.com/MinimalistMarkets/mrkts.fyi/storage"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
)

type Server struct {
	listenAddr   string
	alpacaClient *marketdata.Client
	storage      *storage.CloudStorage
}

func NewServer(listenAddr string, alpacaClient *marketdata.Client, storage *storage.CloudStorage) *Server {
	return &Server{
		listenAddr:   listenAddr,
		alpacaClient: alpacaClient,
		storage:      storage,
	}
}

func (s *Server) Start() error {
	http.HandleFunc("/", s.handleIndex)
	http.HandleFunc("/sync", s.handleSync)
	return http.ListenAndServe(s.listenAddr, nil)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (s *Server) handleSync(w http.ResponseWriter, r *http.Request) {
	if !verifyToken(w, r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sync := NewSync(s.alpacaClient, s.storage)
	if ok := sync.Start(); !ok {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("OK"))
}
