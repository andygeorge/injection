# Injection

Injection is a simple Python3 utility that consumes and processes [Ignition](https://coreos.github.io/ignition/) configs and writes out directories, files, and systemd units (and optionally enables/starts them!) on an **already-running host**, differing from Ignition, which is only intended to configure immutable hosts on first boot.

In conjunction with [Butane](https://coreos.github.io/butane/), this allows you to build human-readable configuration files that can be used to easily (re)configure existing hosts.

## Prerequisites

- Python 3
- [Butane](https://coreos.github.io/butane/getting-started/#getting-butane), if you're planning on using Butane to generate Ignition configs (recommended)

## Installation
- Download `injection` (either [a release](https://github.com/andygeorge/injection/releases) or by cloning this repo)
- Place `injection` somewhere in your `$PATH`

## Usage

```
$ injection
usage: injection [-h] path

$ injection path/to/your/ignition.ign
```

## Example

- Create a Butane config file, `hello_world.bu`:

```yaml
variant: fcos
version: 1.4.0

storage:
  directories:
    - path: /tmp/example
      mode: 0755

  files:
    - path: /tmp/example/hello_world.txt
      mode: 0755
      contents:
        inline: |
          Hello, world!
```

- Create an Ignition config using Butane:

```bash
$ butane --strict hello_world.bu > hello_world.ign
$ cat hello_world.ign
{"ignition":{"version":"3.3.0"},"storage":{"directories":[{"path":"/tmp/example","mode":493}],"files":[{"path":"/tmp/example/hello_world.txt","contents":{"compression":"","source":"data:,Hello%2C%20world!%0A"},"mode":493}]}}
```

- Run Injection:

```bash
$ sudo injection hello_world.ign
Creating directory: /tmp/example
Creating file: /tmp/example/hello_world.txt

$ cat /tmp/example/hello_world.txt
Hello, world!
```

## Ignition support

This only has basic Ignition support for a few specific fields:

- `storage:files`
- `storage:directories`
- `systemd:units`

Support for `passwd:users` is [planned](https://github.com/andygeorge/injection/issues/1).
