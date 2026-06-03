# WIK-DPS-TP01 — Bonus

Reprend le serveur du sujet 1 et ajoute une route `/stats` avec un compteur de requêtes.

## Prérequis

- Go installé (https://go.dev/dl/)

## Lancer

Depuis ce dossier (`bonus/`) :

```
go run .
```

Le port est configurable avec `PING_LISTEN_PORT` (8080 par défaut), via l'environnement ou le `.env`.
L'identifiant d'instance vient de la variable `INSTANCE_ID`, et si elle n'est pas définie on prend le hostname de la machine.

## Routes

- `GET /` : page d'accueil
- `GET /ping` : renvoie les headers de la requête au format JSON
- `GET /stats` : renvoie en JSON le nombre total de requêtes, l'uptime (en secondes) et l'identifiant d'instance
- `GET /health` : renvoie "ok"

Le compteur est incrémenté à chaque requête (toutes routes confondues, `/ping` comprise), grâce à un middleware qui englobe le routeur. Le stockage du compteur est isolé derrière l'interface `CounterStore`, avec ici une implémentation en mémoire (`MemoryStore`). Comme ça on pourrait brancher une autre implémentation plus tard sans toucher au reste.

## Test rapide

```
curl http://localhost:8080/stats
curl http://localhost:8080/ping
curl http://localhost:8080/stats
```

Le `requests` du `/stats` augmente entre les deux appels.
