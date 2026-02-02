# Publishing ListenBucket Container Images

This document describes how to publish ListenBucket container images to Docker Hub and GitHub Container Registry (GHCR).

## Prerequisites

- Docker installed and running
- Account on Docker Hub and/or GitHub
- Logged in to the respective registries

## Docker Hub

### Initial Setup

1. Create a Docker Hub account at https://hub.docker.com
2. Create a repository named `listenbucket`
3. Login to Docker Hub:
   ```bash
   docker login
   ```

### Manual Publishing

```bash
# Build the image
docker build -t listenbucket/listenbucket:latest .

# Tag with version
docker tag listenbucket/listenbucket:latest listenbucket/listenbucket:v1.0.0

# Push to Docker Hub
docker push listenbucket/listenbucket:latest
docker push listenbucket/listenbucket:v1.0.0
```

### Multi-Architecture Build (Recommended)

For broader compatibility (amd64, arm64), use Docker Buildx:

```bash
# Create and use a new builder
docker buildx create --name multiarch --use

# Build and push multi-arch image
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t listenbucket/listenbucket:latest \
  -t listenbucket/listenbucket:v1.0.0 \
  --push \
  .
```

## GitHub Container Registry (GHCR)

### Initial Setup

1. Create a Personal Access Token (PAT) with `write:packages` scope:
   - Go to GitHub Settings → Developer settings → Personal access tokens
   - Generate new token with `write:packages` permission

2. Login to GHCR:
   ```bash
   echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin
   ```

### Manual Publishing

```bash
# Build the image
docker build -t ghcr.io/listenbucket/listenbucket:latest .

# Tag with version
docker tag ghcr.io/listenbucket/listenbucket:latest ghcr.io/listenbucket/listenbucket:v1.0.0

# Push to GHCR
docker push ghcr.io/listenbucket/listenbucket:latest
docker push ghcr.io/listenbucket/listenbucket:v1.0.0
```

### Multi-Architecture Build

```bash
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t ghcr.io/listenbucket/listenbucket:latest \
  -t ghcr.io/listenbucket/listenbucket:v1.0.0 \
  --push \
  .
```

## Automated Publishing with GitHub Actions

The repository includes GitHub Actions workflows for automated publishing:

- `.github/workflows/release.yml` - Publishes on version tags (v*)
- `.github/workflows/ci.yml` - Builds and tests on every push

### Creating a Release

1. Tag the release:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. The release workflow will automatically:
   - Build multi-architecture images
   - Push to both Docker Hub and GHCR
   - Create a GitHub Release

### Required Secrets

Configure these secrets in your GitHub repository settings:

| Secret | Description |
|--------|-------------|
| `DOCKERHUB_USERNAME` | Docker Hub username |
| `DOCKERHUB_TOKEN` | Docker Hub access token |

GHCR uses the built-in `GITHUB_TOKEN` automatically.

## Version Tagging Strategy

- `latest` - Most recent build from main branch
- `v1.0.0` - Specific version release
- `v1.0` - Latest patch for v1.0.x
- `v1` - Latest minor/patch for v1.x.x

## Verifying Published Images

```bash
# Docker Hub
docker pull listenbucket/listenbucket:latest
docker run --rm listenbucket/listenbucket:latest --version

# GHCR
docker pull ghcr.io/listenbucket/listenbucket:latest
docker run --rm ghcr.io/listenbucket/listenbucket:latest --version
```

## Troubleshooting

### Permission Denied on GHCR

Ensure your PAT has the `write:packages` scope and you're logged in:
```bash
echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin
```

### Multi-arch Build Fails

Ensure buildx is properly configured:
```bash
docker buildx ls
docker buildx create --name multiarch --use
docker buildx inspect --bootstrap
```

### Rate Limiting

Docker Hub has pull rate limits for anonymous users. Consider using authenticated pulls or GHCR for CI/CD pipelines.
