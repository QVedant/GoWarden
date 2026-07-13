# GoWarden

A Containerised Go sandbox to run untrusted code and return output. Sandbox using nsjail. All in a singular docker image.

## Features
- Added basic yaml structure
- Setup a temp http server

## Architecture
Work In progress

## Getting Started

### Prerequisites
- Go 1.22+
- Docker
- nsjail (or: built automatically via Dockerfile)

### Installation
```bash
git clone https://github.com/QVedant/GoWarden.git
cd GoWarden
go mod tidy
```

### Running
```bash
make run
```

## API

### `GET /healthz`
Returns `200 OK` if the server is running.

### `GET /languages`
Returns supported languages as JSON.

## Supported Languages
Table or list of what's in `config/languages.yaml`.

## Project Structure
