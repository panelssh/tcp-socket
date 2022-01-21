# Panel SSH - TCP Socket

Simple TCP connection for server command execution

## Installation

Download `.zip` or `.tar.gz` file, you can see all release on [this link](https://github.com/panelssh/tcp-socket/releases)

Example download with curl on Linux:

```bash
VERSION=v1.0.1
KERNEL=amd64
curl -LJ https://github.com/panelssh/tcp-socket/releases/download/${VERSION}/panelssh-tcp-socket-${VERSION}-linux-${KERNEL}.tar.gz | tar -xz
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
