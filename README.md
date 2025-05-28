# Konfiguracja Docker Build & Security Workflow

## Wyzwalacze
Workflow uruchamia się przy:
- Push na branch `main`
- Pull request do branch `main`
- Ręcznym uruchomieniu (`workflow_dispatch`)

## Zmienne środowiskowe
- **Registry**: GitHub Container Registry (`ghcr.io`)
- **Nazwa obrazu**: Automatycznie z nazwy repozytorium
- **Cache**: DockerHub cache image (`dalvy07/weather-app-cache`)

## Główne etapy workflow

### 1. Przygotowanie środowiska
- Checkout repozytorium
- Konfiguracja QEMU dla buildów multi-platform
- Setup Docker Buildx
- Logowanie do DockerHub (cache) i GitHub Container Registry

### 2. Skanowanie bezpieczeństwa
- **Build testowy**: Szybki build tylko dla `linux/amd64`
- **Trivy scan**: Skanowanie pod kątem krytycznych i wysokich luk w zabezpieczeniach
- **Analiza wyników**: Sprawdzenie czy znaleziono vulnerabilities
- **Upload do GitHub Security**: Automatyczne przesłanie wyników do Security tab

### 3. Warunki deploymentu
- **Blokada bezpieczeństwa**: Jeśli znajdą się krytyczne/wysokie luki, deployment zostaje zablokowany
- **Deployment**: Tylko jeśli skan bezpieczeństwa przejdzie pomyślnie

### 4. Finalna publikacja (jeśli bezpieczne)
- Build multi-platform (`linux/amd64`, `linux/arm64`)
- Wykorzystanie cache z DockerHub
- Push do GitHub Container Registry
- Tagowanie: `latest` + SHA commit