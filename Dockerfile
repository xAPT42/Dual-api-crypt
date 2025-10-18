# Multi-stage build pour optimiser la taille de l'image finale
FROM golang:1.24-alpine AS builder

# Installer les dépendances système nécessaires
RUN apk add --no-cache git ca-certificates tzdata

# Définir le répertoire de travail
WORKDIR /app

# Copier les fichiers de dépendances
COPY go.mod go.sum ./

# Télécharger les dépendances
RUN go mod download

# Copier le code source
COPY . .

# Compiler l'application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Image finale minimale
FROM alpine:latest

# Installer ca-certificates pour les requêtes HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Créer un utilisateur non-root
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Définir le répertoire de travail
WORKDIR /app

# Copier le binaire depuis l'étape builder
COPY --from=builder /app/main .

# Changer les permissions
RUN chown -R appuser:appgroup /app

# Passer à l'utilisateur non-root
USER appuser

# Exposer les ports
EXPOSE 8080 9090

# Variables d'environnement
ENV GIN_MODE=release

# Commande de démarrage
CMD ["./main"]
