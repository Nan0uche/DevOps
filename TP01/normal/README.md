# WIK-DPS-TP01

Petit serveur HTTP en Go. Il expose une route `/ping` qui renvoie les headers de la requête au format JSON.

## Prérequis

- Go installé (https://go.dev/dl/)

## Lancer

Depuis la racine du projet :

```
go run .
```

Par défaut le serveur écoute sur le port 8080. Le port est configurable avec la variable d'environnement `PING_LISTEN_PORT`, qu'on peut mettre directement ou dans le fichier `.env` :

```
PING_LISTEN_PORT=8080
```

## Routes

- `GET /` : page d'accueil
- `GET /ping` : renvoie les headers de la requête au format JSON
- `GET /health` : renvoie "ok"

Toute autre route renvoie un 404. Sur `/ping`, si la méthode n'est pas GET, le serveur répond un 404 vide.

## Test rapide

```
curl http://localhost:8080/ping
```
