#!/bin/bash

# Script de démarrage pour Dual API Crypto
# Usage: ./start.sh [build|up|down|logs|test]

set -e

PROJECT_NAME="dual-api-crypto"
COMPOSE_FILE="docker-compose.yml"

# Couleurs pour les messages
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Fonction d'affichage des messages
print_message() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Fonction d'aide
show_help() {
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  build    - Construire l'image Docker"
    echo "  up       - Démarrer les services"
    echo "  down     - Arrêter les services"
    echo "  logs     - Afficher les logs"
    echo "  test     - Tester les APIs"
    echo "  clean    - Nettoyer les images et volumes"
    echo "  help     - Afficher cette aide"
    echo ""
    echo "Exemples:"
    echo "  $0 up      # Démarrer le service"
    echo "  $0 test    # Tester les APIs"
    echo "  $0 logs    # Voir les logs"
}

# Fonction de test des APIs
test_apis() {
    print_message "Test des APIs..."
    
    # Attendre que le service soit prêt
    print_message "Attente du démarrage du service..."
    sleep 10
    
    # Test de l'API REST
    print_message "Test de l'API REST..."
    if curl -s http://localhost:8080/health > /dev/null; then
        print_success "API REST Health Check: OK"
    else
        print_error "API REST Health Check: FAILED"
        return 1
    fi
    
    # Test de l'endpoint portfolio
    print_message "Test de l'endpoint /portfolio..."
    PORTFOLIO_RESPONSE=$(curl -s http://localhost:8080/portfolio)
    if echo "$PORTFOLIO_RESPONSE" | jq . > /dev/null 2>&1; then
        print_success "API REST Portfolio: OK"
        print_message "Réponse:"
        echo "$PORTFOLIO_RESPONSE" | jq .
    else
        print_error "API REST Portfolio: FAILED"
        return 1
    fi
    
    # Test de l'API gRPC
    print_message "Test de l'API gRPC..."
    if command -v grpcurl > /dev/null; then
        if grpcurl -plaintext localhost:9090 list > /dev/null 2>&1; then
            print_success "API gRPC: OK"
            print_message "Services disponibles:"
            grpcurl -plaintext localhost:9090 list
        else
            print_error "API gRPC: FAILED"
            return 1
        fi
    else
        print_warning "grpcurl non installé, test gRPC ignoré"
        print_message "Pour installer grpcurl: go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest"
    fi
    
    print_success "Tous les tests sont passés !"
}

# Fonction de nettoyage
clean_up() {
    print_message "Nettoyage des ressources Docker..."
    docker-compose down --volumes --remove-orphans
    docker system prune -f
    print_success "Nettoyage terminé"
}

# Traitement des arguments
case "${1:-help}" in
    "build")
        print_message "Construction de l'image Docker..."
        docker-compose build
        print_success "Construction terminée"
        ;;
    "up")
        print_message "Démarrage des services..."
        docker-compose up -d
        print_success "Services démarrés"
        print_message "APIs disponibles:"
        print_message "  REST: http://localhost:8080/portfolio"
        print_message "  gRPC: localhost:9090"
        ;;
    "down")
        print_message "Arrêt des services..."
        docker-compose down
        print_success "Services arrêtés"
        ;;
    "logs")
        print_message "Affichage des logs..."
        docker-compose logs -f
        ;;
    "test")
        test_apis
        ;;
    "clean")
        clean_up
        ;;
    "help"|*)
        show_help
        ;;
esac
