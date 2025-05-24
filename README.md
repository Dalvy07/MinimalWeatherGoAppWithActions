# MinimalWeatherGoApp

Repository for PAwCHO Zad1

[weather-app dockerhub](https://hub.docker.com/repository/docker/dalvy07/weather-app/general)

[weather-app-cache dockerhub](https://hub.docker.com/repository/docker/dalvy07/weather-app-cache/general)

## Setup Instructions

### Getting an API Key

1. Sign up on [Weather API](https://www.weatherapi.com/)
2. After logging in, change API Response Fields in Dashboard:
   - Unmark all responses with Imperial units 🦅🦅🦅 (Fahrenheit, miles per hour, inches)
3. Create an `api_key.txt` file in the project root folder and paste your API key

## Docker Commands

### Base Building of Container

- **Build using source code from GitHub repo**
  ```bash
  docker build --no-cache --ssh github=~/.ssh/gh_lab6 --secret id=api_key,src=api_key.txt -t weather-app .
  ```
  ```bash
  docker build --build-arg CACHE_BUST=$(date +%s) --ssh github=~/.ssh/gh_lab6 --secret id=api_key,src=api_key.txt -t weather-app .
  ```

### Building Container with SBOM and Provenance for Checking Vulnerabilities

- **Build with SBOM and provenance**
  ```bash
  docker buildx build \
  --platform linux/amd64 \
  --build-arg CACHE_BUST=$(date +%s) \
  --ssh github=~/.ssh/gh_lab6 \
  --secret id=api_key,src=api_key.txt \
  --tag docker.io/dalvy07/weather-app:latest \
  --sbom=true \
  --provenance=mode=max \
  --push \
  .
  ```

- **To see SBOM:**
  ```bash
  docker buildx imagetools inspect docker.io/dalvy07/weather-app:latest --format '{{json .SBOM}}'
  ```

- **To see Provenance:**
  ```bash
  docker buildx imagetools inspect docker.io/dalvy07/weather-app:latest --format '{{json .Provenance}}'
  ```

- **Checking image for vulnerabilities**
  ```bash
  docker scout cves docker.io/dalvy07/weather-app:latest
  ```
![image-alt](https://github.com/Dalvy07/MinimalWeatherGoApp/blob/main/screenshots/vulnerabilities_check.png?raw=true)
### Building Multiplatform OCI Image Using Cache Type Registry

- **Building OCI image for 2 platforms using cache type registry**
  ```bash
  docker buildx build \
  --platform linux/amd64,linux/arm64 \
  --build-arg CACHE_BUST=$(date +%s) \
  --ssh github=~/.ssh/gh_lab6 \
  --secret id=api_key,src=api_key.txt \
  --tag docker.io/dalvy07/weather-app:latest \
  --sbom=true \
  --provenance=mode=max \
  --cache-to type=registry,ref=docker.io/dalvy07/weather-app-cache:latest,mode=max \
  --cache-from type=registry,ref=docker.io/dalvy07/weather-app-cache:latest \
  --output type=image,name=docker.io/dalvy07/weather-app:latest,push=true,oci-mediatypes=true \
  .
  ```
*First iteration of building using registry cache*
![image-alt](https://github.com/Dalvy07/MinimalWeatherGoApp/blob/main/screenshots/first_iteration_using_registry_cache.png?raw=true)

*Second iteration of building using registry cache*
![image-alt](https://github.com/Dalvy07/MinimalWeatherGoApp/blob/main/screenshots/second_iteration_using_registry_cache.png?raw=true)
- **To check manifest**
  ```bash
  docker manifest inspect docker.io/dalvy07/weather-app:latest
  ```

- **To check usage of OCI media-typesand multiplatforming **
  ```bash
  docker buildx imagetools inspect docker.io/dalvy07/weather-app:latest
  ```
*Check usage of OCI media-types and multiplatforming*
![image-alt](https://github.com/Dalvy07/MinimalWeatherGoApp/blob/main/screenshots/check_manifest_for_OCI_and_multiplatform.png?raw=true)
### Sending Image to DockerHub

- **You can send your image to DockerHub using:**
  ```bash
  docker tag weather-app docker.io/dalvy07/weather-app:latest
  ```
  ```bash
  docker push docker.io/dalvy07/weather-app:latest
  ```

### Running the Container

- **Simple run with:**
  ```bash
  docker run -p 3000:3000 docker.io/dalvy07/weather-app
  ```

### Analyzing Image

- **Basic information about image (can check layers, env, etc.):**
  ```bash
  docker inspect docker.io/dalvy07/weather-app
  ```

- **Build history:**
  ```bash
  docker history docker.io/dalvy07/weather-app
  ```

- **To check manifest**
  ```bash
  docker manifest inspect docker.io/dalvy07/weather-app:latest
  ```

- **To check usage of OCI media-types**
  ```bash
  docker buildx imagetools inspect docker.io/dalvy07/weather-app:latest
  ```

- **To see SBOM:**
  ```bash
  docker buildx imagetools inspect docker.io/dalvy07/weather-app:latest --format '{{json .SBOM}}'
  ```

- **To see Provenance:**
  ```bash
  docker buildx imagetools inspect docker.io/dalvy07/weather-app:latest --format '{{json .Provenance}}'
  ```

- **Checking image for vulnerabilities**
  ```bash
  docker scout cves docker.io/dalvy07/weather-app:latest
  ```

> **Note:** Unfortunately, you can only view the logs using Docker Desktop or the terminal in which the container was started. Since I tried to make the container as small as possible, the base image has nothing but scratch, which makes it impossible to connect to the container terminal to view logs, because there is no terminal there.
