# Speed Date

A simple HTTP API for a dating app

## Requirements 

* Go =< 1.23
* Docker

To run on bare metal, openssl and make are useful too

## Starting

Use docker compose to start everything 

```sh
docker compose up -d
```

Or you can start postgres and run the app without docker

```sh
make gen-rsa
docker compose up -d postgres
go run cmd/main.go
```

## Assumtions & Descisions

* There will only be two environments: local and docker
* There are no refresh tokens, one has to re-login. Done for brevity
* No SSL handler. More things to generate and more potential for a 'works on my machine'
* HTTP is a suitable. The logic is written in a way the implementing another protocol should be simple
* Locations are stored as fixed Coordinates. I assume location data is a lot less specific in prod, but again keeping it simple
* No logging collector. One could assume that a pod/container collector is enough. But if this were going to production hooking up OpenTelemetry for tracing and logging would be ideal
* The PII in the database is not encrypted. This is not suitable for production apps, but again keep it simple
* The vendor folder is not checked in. The need for reproducibility is not as important for a dummy service
