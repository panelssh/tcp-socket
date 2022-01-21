# Panel SSH - TCP Socket

Simple TCP connection for server command execution

## Installation

Download Binary File

```bash
curl -LJ https://github.com/panelssh/tcp-socket/releases/download/v1.0.0/panelssh-tcp-socket-v1.0.0-linux-amd64.tar.gz | tar -xz
```

## Usage

Environment Variables:

| Key               | Default   | Description                              |
|-------------------|-----------|------------------------------------------|
| `HOST`            | `0.0.0.0` | Bind Hostname                            |
| `PORT`            | `3000`    | Listening Port                           |
| `SECRET_KEY`      | `test`    | Secret key validation                    |
| `ALLOWED_ADDRESS` | `%` (any) | List IP Address (v4), separate by comma. |

Command with default environment variable:

```bash
./panelssh-tcp-socket
```

Command override environment variable:

```bash
PORT=5555 SECRET_KEY=1234 ./panelssh-tcp-socket
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.
