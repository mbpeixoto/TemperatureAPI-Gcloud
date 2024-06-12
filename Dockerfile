# Etapa de construção
FROM golang:1.22.2 as build
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cloudrun

# Etapa final para a aplicação
FROM scratch as runtime
WORKDIR /app
COPY --from=build /app/cloudrun .
ENTRYPOINT ["./cloudrun"]
