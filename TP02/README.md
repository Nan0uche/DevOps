# WIK-DPS-TP02

Dockerisation de l'API du TP01 (le serveur /ping). Il y a deux images : une en mono-stage et une en multi-stage.

## Build

Mono-stage :
```
docker build -t wik-dps-tp02:single .
```

Multi-stage :
```
docker build -f Dockerfile.multistage -t wik-dps-tp02:multi .
```

## Lancer

```
docker run -p 8080:8080 wik-dps-tp02:multi
```

Ensuite : http://localhost:8080/ping
Le port peut se changer avec la variable PING_LISTEN_PORT.

Avec docker compose :
```
docker compose up --build
```

## Optimisations

- Utilisateur non-root (appuser) dans les deux images.
- On copie go.mod et on fait go mod download avant de copier le code, comme ça modifier le code ne relance pas le téléchargement des dépendances.
- La version multi-stage build dans une image golang puis copie seulement le binaire et le dossier web dans une image alpine. Du coup elle fait environ 15 Mo au lieu de 305 Mo pour la mono-stage.

## Scan

Image scannée avec trivy :
```
docker save wik-dps-tp02:multi -o image.tar
docker run --rm -v "${PWD}:/work" aquasec/trivy image --input /work/image.tar
```

Il reste 15 vulnérabilités (14 HIGH, 1 CRITICAL) sur l'image multi, toutes dans la bibliothèque standard de Go 1.22. Elles seraient corrigées en passant à une version de Go plus récente (1.25+).
