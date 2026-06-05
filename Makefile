.PHONY: build run test test-coverage coverage-html vet fmt up down clean

# Uygulamayi derler ve bin/ altina koyar.
build:
	go build -o bin/url-kisaltici ./cmd/server

# Uygulamayi dogrudan calistirir (lokalde Redis gerekir).
run:
	go run ./cmd/server

# Tum testleri ayrintili calistirir.
test:
	go test ./... -v

# Testleri calistirip kapsam ozetini ekrana basar.
test-coverage:
	go test ./... -coverprofile=coverage.out -covermode=atomic
	go tool cover -func=coverage.out

# Kapsam raporunu html olarak uretir.
coverage-html: test-coverage
	go tool cover -html=coverage.out -o coverage.html

# Kodda supheli durumlari arar.
vet:
	go vet ./...

# Kodu bicimlendirir.
fmt:
	go fmt ./...

# Redis ve uygulamayi docker ile ayaga kaldirir.
up:
	docker compose up --build

# Docker servislerini kapatir.
down:
	docker compose down

# Uretilen dosyalari temizler.
clean:
	rm -rf bin coverage.out coverage.html
