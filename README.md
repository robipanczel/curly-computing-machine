# Project curly-computing-machine

Curly project

## Getting Started

Simple RestAPI. The API requires a mongodb connection, please create a .env file based on the .env.example file. A makefile is available to get started easily.
Or you can start the API and the database by running the docker compose file with `docker compose up`, it will build and start both services.

## MakeFile

Run build make command with tests

```bash
make all
```

Build the application

```bash
make build
```

Run the application

```bash
make run
```

Create DB container

```bash
make docker-run
```

Shutdown DB Container

```bash
make docker-down
```

DB Integrations Test:

```bash
make itest
```

Live reload the application:

```bash
make watch
```

Run the test suite:

```bash
make test
```

Clean up binary from the last build:

```bash
make clean
```
