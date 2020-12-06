## About The Project

This repository contains a very basic example of a configurable reverse proxy made with go.

## Getting Started

### Prerequisites

You will need [Docker](https://www.docker.com/) and [Docker compose](https://docs.docker.com/compose/install/).

### Installation

The configuration file it is self-explanatory and includes all possible options.

```yaml
config:
  proxy address: "0.0.0.0:9090"

routes:
  /redirect:
    strategy: roundrobin
    root: true
    servers:
      - http://backend01:80
      - http://backend02:80
  /nginx:
    strategy: random
    root: false
    servers:
      - http://backend01:80
      - http://backend02:80
      - http://backend03:80
```

## Usage

The repository includes a working example which includes a basic `config.yml`.

```sh
$ git clone https://github.com/lopz82/reverse-proxy.git
$ cd reverse-proxy
$ docker-compose up 
```

## License

Distributed under the MIT License. See `LICENSE` for more information.
