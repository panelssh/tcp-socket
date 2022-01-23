# Panel SSH - TCP Socket

Simple TCP connection for server command execution

## Installation

Download `.zip` or `.tar.gz` file, you can see all release on [this link](https://github.com/panelssh/tcp-socket/releases)

Example download with curl on Linux:

```bash
VERSION=v1.0.1
KERNEL=amd64
curl -LJ https://github.com/panelssh/tcp-socket/releases/download/${VERSION}/panelssh-tcp-socket-${VERSION}-linux-${KERNEL}.tar.gz | tar -xz -C /usr/bin
```

## Usage

Environment Variables:

| Key               | Default   | Description                              |
|-------------------|-----------|------------------------------------------|
| `HOST`            | `0.0.0.0` | Bind Host/IP Address                     |
| `PORT`            | `3000`    | Listening Port                           |
| `SECRET_KEY`      | `test`    | Secret key validation                    |
| `ALLOWED_ADDRESS` | `%` (any) | List IP Address (v4), separate by comma. |

Example Command:

```bash
PORT=5555 SECRET_KEY=1234 panelssh-tcp-socket
```

Flag / Argument Options:

> Flags / arguments will be replace environment variable

| Flags/Args         | Default               | Description                              |
|--------------------|-----------------------|------------------------------------------|
| `-host`            | env `HOST`            | Bind Host/IP Address                     |
| `-port`            | env `PORT`            | Listening Port                           |
| `-secret-key`      | env `SECRET_KEY`      | Secret key validation                    |
| `-allowed-address` | env `ALLOWED_ADDRESS` | List IP Address (v4), separate by comma. |

Example Command:

```bash
panelssh-tcp-socket -allowed-address="192.168.0.2,192.168.0.3,192.168.0.4"
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.
