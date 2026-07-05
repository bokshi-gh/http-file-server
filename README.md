# HTTP File Server (httpfs)

A lightweight HTTP file server written in Go.

## Installation

### Install from source (build on your machine)

Clone the repository:

```bash
git clone https://github.com/bokshi-gh/http-file-server.git
cd http-file-server
```

Build the server:

```bash
go build -o httpfs ./cmd/httpfs
```

### Install using build scripts

For Unix-like systems (Linux, macOS):

```bash
curl -fsSL https://raw.githubusercontent.com/bokshi-gh/http-file-server/main/scripts/install.sh | bash
```

For Windows:

```powershell
irm https://raw.githubusercontent.com/bokshi-gh/http-file-server/main/scripts/install.ps1 | iex
```

## Usage

Run the server:

```bash
httpfs --root ./public --host 0.0.0.0 --port 8080 --v
```
