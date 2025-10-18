# 🚀 Dual API Crypto

[![Docker Build](https://github.com/xAPT42/Dual-api-crypt/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/xAPT42/Dual-api-crypt/actions/workflows/docker-publish.yml)
[![Docker Hub](https://img.shields.io/docker/pulls/xapt42/dual-api-crypt)](https://hub.docker.com/r/xapt42/dual-api-crypt)

Un microservice Go démontrant une architecture "Dual API" avec une API REST et une API gRPC pour le suivi de portefeuille crypto en temps réel.

## 🎯 Objectif

Ce projet illustre comment une seule logique métier peut servir deux types de clients :
- **API REST** : Pour les interfaces web et les intégrations simples
- **API gRPC** : Pour la communication machine-à-machine haute performance

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐
│   Client REST   │    │   Client gRPC   │
└─────────┬───────┘    └─────────┬───────┘
          │                      │
          ▼                      ▼
┌─────────────────────────────────────────┐
│           Dual API Server                │
│  ┌─────────────┐  ┌─────────────────┐   │
│  │ REST API    │  │   gRPC API      │   │
│  │ Port 8080   │  │   Port 9090     │   │
│  └──────┬──────┘  └─────────┬───────┘   │
│         │                   │           │
│         └─────────┬─────────┘           │
│                   ▼                     │
│         ┌─────────────────┐             │
│         │  Business Logic  │             │
│         │  (Portfolio)    │             │
│         └─────────────────┘             │
└─────────────────────────────────────────┘
          │
          ▼
┌─────────────────────────────────────────┐
│        CoinGecko API                    │
│     (Prix crypto temps réel)           │
└─────────────────────────────────────────┘
```

## 🚀 Démarrage Rapide

### Prérequis

- Docker et Docker Compose
- Go 1.24+ (pour le développement local)

### Lancement avec Docker Hub

```bash
# Utiliser l'image pré-construite depuis Docker Hub
docker pull xapt42/dual-api-crypt:latest

# Lancer le conteneur
docker run -d -p 8080:8080 -p 9090:9090 --name dual-api-crypt xapt42/dual-api-crypt:latest
```

### Lancement avec Docker Compose

```bash
# Cloner le projet
git clone https://github.com/xAPT42/Dual-api-crypt.git
cd Dual-api-crypt

# Construire et lancer
docker-compose up --build
```

### Développement Local

```bash
# Installer les dépendances
go mod tidy

# Générer le code gRPC
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       api/portfolio.proto

# Lancer l'application
go run cmd/server/main.go
```

## 📡 API Endpoints

### REST API (Port 8080)

#### Obtenir la valeur du portefeuille
```bash
curl http://localhost:8080/portfolio
```

**Réponse :**
```json
{
  "value": 12543.67,
  "currency": "EUR",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

#### Health Check
```bash
curl http://localhost:8080/health
```

### gRPC API (Port 9090)

#### Lister les services
```bash
grpcurl -plaintext localhost:9090 list
```

#### Appeler GetPortfolioValue
```bash
grpcurl -plaintext localhost:9090 portfolio.PortfolioService/GetPortfolioValue
```

## 💼 Portefeuille par Défaut

Le portefeuille est défini dans le code et contient :
- **Bitcoin** : 0.5 BTC
- **Ethereum** : 5.0 ETH  
- **Cardano** : 1000 ADA
- **Solana** : 50 SOL
- **Chainlink** : 100 LINK

## 🛠️ Structure du Projet

```
dual-api-crypto/
├── api/                    # Contrats gRPC
│   ├── portfolio.proto     # Définition du service
│   ├── portfolio.pb.go     # Code généré
│   └── portfolio_grpc.pb.go
├── cmd/server/             # Point d'entrée
│   └── main.go            # Serveurs REST et gRPC
├── internal/service/      # Logique métier
│   └── portfolio.go       # Calcul du portefeuille
├── Dockerfile             # Image Docker
├── docker-compose.yml     # Orchestration
├── go.mod                 # Dépendances Go
└── README.md              # Documentation
```

## 🔧 Technologies Utilisées

- **Go 1.24** : Langage principal
- **gRPC** : Communication haute performance
- **Protocol Buffers** : Sérialisation efficace
- **Docker** : Conteneurisation
- **CoinGecko API** : Données crypto temps réel

## 🧪 Tests

### Test REST API
```bash
# Test de base
curl -X GET http://localhost:8080/portfolio

# Test avec jq pour un affichage formaté
curl -s http://localhost:8080/portfolio | jq
```

### Test gRPC API
```bash
# Installer grpcurl
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Lister les services
grpcurl -plaintext localhost:9090 list

# Appeler la méthode
grpcurl -plaintext localhost:9090 portfolio.PortfolioService/GetPortfolioValue
```

## 📊 Monitoring

Le service expose un endpoint de santé sur `/health` pour le monitoring.

## 🔄 CI/CD et Docker Hub

### Workflow GitHub Actions

Le projet utilise GitHub Actions pour automatiser le build et le push des images Docker vers Docker Hub :

- **Déclenchement** :
  - Push sur la branche `main` → tag `latest`
  - Tags Git (ex: `v1.0.0`) → tags versionnés
  - Pull Requests → build de test uniquement

- **Multi-architecture** : Les images sont buildées pour `linux/amd64` et `linux/arm64`

- **Cache** : Utilise GitHub Actions cache pour accélérer les builds

### Configuration des Secrets GitHub

Pour activer le workflow, configurez les secrets suivants dans votre dépôt GitHub :

1. Allez dans `Settings` → `Secrets and variables` → `Actions`
2. Ajoutez les secrets suivants :
   - `DOCKERHUB_USERNAME` : Votre nom d'utilisateur Docker Hub
   - `DOCKERHUB_TOKEN` : Votre token d'accès Docker Hub (créer un token sur hub.docker.com)

### Images Docker Hub

Les images sont disponibles sur : **[xapt42/dual-api-crypt](https://hub.docker.com/r/xapt42/dual-api-crypt)**

```bash
# Dernière version
docker pull xapt42/dual-api-crypt:latest

# Version spécifique
docker pull xapt42/dual-api-crypt:v1.0.0
```

## 🔒 Sécurité

- Utilisateur non-root dans le conteneur
- Certificats SSL pour les requêtes HTTPS
- Headers CORS configurés
- Timeouts configurés

## 🤝 Contribution

1. Fork le projet
2. Créer une branche feature (`git checkout -b feature/AmazingFeature`)
3. Commit les changements (`git commit -m 'Add some AmazingFeature'`)
4. Push vers la branche (`git push origin feature/AmazingFeature`)
5. Ouvrir une Pull Request

## 📝 Licence

Ce projet est sous licence MIT. Voir le fichier `LICENSE` pour plus de détails.

## 🎓 Apprentissage

Ce projet démontre :
- **Architecture Dual API** : Une logique métier, deux interfaces
- **gRPC vs REST** : Quand utiliser chaque technologie
- **Microservices** : Design et implémentation
- **Docker** : Conteneurisation d'applications Go
- **API externes** : Intégration avec CoinGecko

---

**Développé avec ❤️ pour démontrer les architectures microservices modernes**
