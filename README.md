# ReqViz

Visual OpenAPI explorer. Upload an API spec, view it as an interactive graph, and test endpoints with a built-in client.

## Features

- **Import** an OpenAPI spec — upload a file, paste JSON/YAML, or load a bundled demo
- **Visualize** endpoints and data models as an interactive graph
- **Try it out** — build and send requests to any endpoint, with response viewer and request history

## Tech Stack

- **Backend** — Go 1.26 + Chi
- **Frontend** — Angular 21 (standalone components, signals)
- **Contract-first** — `api/openapi.yaml` drives generated types on both sides

## Getting Started

Requires [Go](https://go.dev) 1.26+, [Bun](https://bun.sh), and [Task](https://taskfile.dev).

```sh
task dev      # frontend (:4200) + backend (:3000)
task build    # build both
task test     # run all tests
task lint     # lint both
```

## Contract-First Workflow

All backend–frontend communication goes through `api/openapi.yaml`. Edit the spec, then regenerate types for both sides:

```sh
task generate:api
```

See [CLAUDE.md](CLAUDE.md) for conventions and project structure.
