package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/xAPT42/dual-api-crypto/api"
	"github.com/xAPT42/dual-api-crypto/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// PortfolioServer implémente l'interface gRPC
type PortfolioServer struct {
	api.UnimplementedPortfolioServiceServer
}

// GetPortfolioValue implémente la méthode gRPC
func (s *PortfolioServer) GetPortfolioValue(ctx context.Context, req *api.GetPortfolioValueRequest) (*api.GetPortfolioValueResponse, error) {
	value, currency, timestamp, err := service.CalculatePortfolioValue()
	if err != nil {
		return nil, fmt.Errorf("erreur lors du calcul du portefeuille: %v", err)
	}

	return &api.GetPortfolioValueResponse{
		Value:     value,
		Currency:  currency,
		Timestamp: timestamp.Unix(),
	}, nil
}

// handlePortfolioValue gère les requêtes REST
func handlePortfolioValue(w http.ResponseWriter, r *http.Request) {
	// Vérifier la méthode HTTP
	if r.Method != http.MethodGet {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Calculer la valeur du portefeuille
	value, currency, timestamp, err := service.CalculatePortfolioValue()
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors du calcul: %v", err), http.StatusInternalServerError)
		return
	}

	// Créer la réponse JSON
	response := service.PortfolioValue{
		Value:     value,
		Currency:  currency,
		Timestamp: timestamp,
	}

	// Définir les headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Encoder et envoyer la réponse
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Erreur lors de l'encodage JSON", http.StatusInternalServerError)
		return
	}
}

// startGRPCServer démarre le serveur gRPC
func startGRPCServer() {
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("Erreur lors de l'écoute gRPC: %v", err)
	}

	grpcServer := grpc.NewServer()
	api.RegisterPortfolioServiceServer(grpcServer, &PortfolioServer{})

	// Activer la réflexion gRPC pour grpcurl
	reflection.Register(grpcServer)

	log.Println("Serveur gRPC démarré sur le port 9090")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Erreur serveur gRPC: %v", err)
	}
}

// startRESTServer démarre le serveur REST
func startRESTServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/portfolio", handlePortfolioValue)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Println("Serveur REST démarré sur le port 8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Erreur serveur REST: %v", err)
	}
}

func main() {
	log.Println("Démarrage du microservice Dual API Crypto...")

	// Démarrer les serveurs en parallèle
	go startGRPCServer()
	go startRESTServer()

	// Attendre indéfiniment
	select {}
}
