# 🎯 Démonstration Dual API Crypto

## Vue d'ensemble

Ce projet démontre une architecture "Dual API" où une seule logique métier sert deux types de clients :
- **API REST** : Pour les interfaces web et intégrations simples
- **API gRPC** : Pour la communication machine-à-machine haute performance

## 🚀 Démarrage Rapide

```bash
# Cloner et démarrer
git clone https://github.com/xAPT42/dual-api-crypto.git
cd dual-api-crypto
./start.sh up

# Tester les APIs
./start.sh test
```

## 📡 Tests des APIs

### API REST (Port 8080)

#### Health Check
```bash
curl http://localhost:8080/health
# Réponse: OK
```

#### Valeur du Portefeuille
```bash
curl http://localhost:8080/portfolio | jq
```

**Réponse :**
```json
{
  "value": 72619.897,
  "currency": "EUR",
  "timestamp": "2025-10-18T20:01:21.142082557Z"
}
```

### API gRPC (Port 9090)

#### Lister les Services
```bash
grpcurl -plaintext localhost:9090 list
```

**Réponse :**
```
grpc.reflection.v1.ServerReflection
grpc.reflection.v1alpha.ServerReflection
portfolio.PortfolioService
```

#### Appeler GetPortfolioValue
```bash
grpcurl -plaintext -d '{}' localhost:9090 portfolio.PortfolioService/GetPortfolioValue
```

**Réponse :**
```json
{
  "value": 72611.124,
  "currency": "EUR",
  "timestamp": "1760817624"
}
```

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
│        CoinGecko API                     │
│     (Prix crypto temps réel)            │
└─────────────────────────────────────────┘
```

## 💼 Portefeuille par Défaut

Le portefeuille contient :
- **Bitcoin** : 0.5 BTC
- **Ethereum** : 5.0 ETH  
- **Cardano** : 1000 ADA
- **Solana** : 50 SOL
- **Chainlink** : 100 LINK

## 🔧 Commandes Utiles

```bash
# Démarrer le service
./start.sh up

# Voir les logs
./start.sh logs

# Tester les APIs
./start.sh test

# Arrêter le service
./start.sh down

# Nettoyer
./start.sh clean
```

## 📊 Comparaison REST vs gRPC

| Aspect | REST API | gRPC API |
|--------|----------|----------|
| **Format** | JSON | Protocol Buffers |
| **Performance** | Bonne | Excellente |
| **Taille** | Plus volumineux | Compact |
| **Type Safety** | Faible | Fort |
| **Outils** | curl, Postman | grpcurl, clients générés |
| **Cas d'usage** | Web, intégrations | Microservices, IoT |

## 🎓 Points d'Apprentissage

1. **Architecture Dual API** : Une logique métier, deux interfaces
2. **gRPC vs REST** : Quand utiliser chaque technologie
3. **Microservices** : Design et implémentation
4. **Docker** : Conteneurisation d'applications Go
5. **API externes** : Intégration avec CoinGecko

## 🔍 Monitoring

Le service expose un endpoint de santé sur `/health` pour le monitoring.

## 🛠️ Développement

```bash
# Développement local
go mod tidy
go run cmd/server/main.go

# Générer le code gRPC
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       api/portfolio.proto
```

---

**Développé avec ❤️ pour démontrer les architectures microservices modernes**
