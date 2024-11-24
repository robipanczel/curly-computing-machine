# Project curly-computing-machine

Curly project

## Getting Started

1. Run docker compose up to start everything automatically, or
2. Create a .env file from .env.example and use the Makefile for manual setup (MongoDB required)

Documentation is available in openapi.yml or through our live OpenAPI interface.

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
