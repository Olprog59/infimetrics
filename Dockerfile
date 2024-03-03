# Étape de build
FROM golang:1.22-alpine3.19 AS builder

LABEL maintainer="Samuel MICHAUX <samuel.michaux@olprog.fr>"

ARG GOOS=linux
ARG GOARCH=amd64

ENV CGO_ENABLED 0

RUN apk update --no-cache && apk add --no-cache git tzdata ca-certificates

# Définis le répertoire de travail
WORKDIR /app

# Copie les fichiers de dépendances et télécharge les dépendances
COPY go.mod go.sum ./
RUN go mod download && go mod tidy && go mod verify

# Copie le code source dans l'image
COPY . .

# Compile l'application
#RUN GOOS=${GOOS} GOARCH=${GOARCH} go build -a -installsuffix cgo -o main .
RUN GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags="-s -w" -o main .

# Étape de création de l'image finale
FROM scratch

# Copie l'exécutable depuis l'étape de build
COPY --from=builder /app/main /app/main
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY web /web

# Copie les certificats CA (si nécessaire)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080/tcp

CMD ["/app/main"]
