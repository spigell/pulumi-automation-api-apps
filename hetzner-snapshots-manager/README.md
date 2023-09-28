# hetzner Snapshots Manager
The tool for creating snapshots of Hetzner Cloud servers while deleting server via Pulumi program.

It is very useful if you need some *cloud* server for a couple of hours or days for very cheap price and you would like to continue your work with same data after a week. It creates a snapshot for your server before deleting it and you can restore it later. And it would be cheaper than taken from big cloud providers.

# Motivation
## Why Hetzner?
Hetzner Cloud is a great service for hosting your projects. It has a pretty decent set of features (firewall, loadbalancer) and it's a quite simple. And it's very cheap.

## limitation
The only limitation is number of allowed snapshots per account. Only 30 snapshots per account are allowed by default. Setting `max-keep` to 1 (default) will allow you to keep only last shapshot.

# Installation
Grab the binary from releases page or build it from source.

```
go install github.com/spigell/pulumi-automation-api-apps/hetzner-snapshots-manager
```

```
Usage:
  hetzner-snapshot-manager [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  destroy     The 'destroy' command destroy pulumi stack
  help        Help about any command
  pre         The 'pre' only make preview
  up          The 'up' command runs pulumi

Flags:
      --api-server-port int    default is random
      --cleaner-max-keep int   default is keepling only the last snaphot (default 1)
  -c, --config string          /path/to/file for config. Required
  -d, --diff                   Enable the diff option for pulumi command
      --hcloud-token string    Hetzner Cloud token
  -h, --help                   help for hetzner-snapshot-manager
      --only-api-server        Run only api server and do not stop it. For testing purposes.
  -v, --verbose                verbose

```

# Usage
## Configuration
The tool uses a configuration file. For example:
```
---
token: 
max-keep: 2 # keep only 2 last snapshot (sorted by creation date)
stack:
  name: test # name of the pulumi stack
  path: examples/go # path to the pulumi program
```

Also you can use environment variables for configuration:
`HCLOUD_TOKEN` - use it instead of hardcoding token in config file or passing via flags


## Run
```
hetzner-snapshots-manager up -c config.example.yaml -d
```

## API
For golang Pulumi program the one can use `github.com/spigell/pulumi-automation-api-apps/hetzner-snaphots-manager/sdk` package. For example please see [example](examples/go.main.go) file.
