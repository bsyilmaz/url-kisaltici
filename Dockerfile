# Derleme asamasi: kaynak kodu derleyip tek bir ikili dosya uretiyoruz.
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Once bagimliliklari indir, boylece kod degisince cache bozulmaz.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /url-kisaltici ./cmd/server

# Calistirma asamasi: sadece derlenmis dosyayi tasiyan kucuk bir imaj.
FROM alpine:3.19

WORKDIR /

COPY --from=builder /url-kisaltici /url-kisaltici

EXPOSE 8080

ENTRYPOINT ["/url-kisaltici"]
