# Injection

Injection is a simple Python3 utility that consumes and processes [Ignition](https://coreos.github.io/ignition/) configs and writes out directories, files, and systemd units (and optionally enables/starts them!) on an **already-running Linux host**, differing from Ignition, which is only intended to configure immutable Fedora/RedHat CoreOS hosts on first boot.

In conjunction with [Butane](https://coreos.github.io/butane/), this allows you to build human-readable configuration files that can be used to easily (re)configure existing hosts.

## Prerequisites

- Python 3
- [Butane](https://coreos.github.io/butane/getting-started/#getting-butane), if you're planning on using Butane to generate Ignition configs (recommended)

## Installation

```
sudo wget https://github.com/andygeorge/injection/releases/download/v1.0.1/injection -O /usr/local/bin/injection
```

Or, manually:
- Download the `injection` script (either [a release](https://github.com/andygeorge/injection/releases) or by cloning this repo)
- Place `injection` somewhere in your `$PATH`

## Example

- Create a Butane config file, `example_files/example.bu`:

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

    - path: /tmp/example/hello_world_gzip.txt
      mode: 0755
      contents:
        inline: |
          This is a much larger Hello World file that will use gzip compression in the Ignition. Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.

systemd:
  units:
    - name: hello-world.service
      enabled: true
      contents: |
        [Unit]
        Description=Hello world service

        [Service]
        Type=oneshot
        ExecStart=/usr/bin/echo "hello world"
        StandardOutput=journal

        [Install]
        WantedBy=multi-user.target default.target
```

- Create an Ignition config using Butane:

```shell
$ butane --strict example_files/example.bu > example_files/example.ign
$ cat example_files/example.ign
{"ignition":{"version":"3.3.0"},"storage":{"directories":[{"path":"/tmp/example","mode":493}],"files":[{"path":"/tmp/example/hello_world.txt","contents":{"compression":"","source":"data:,Hello%2C%20world!%0A"},"mode":493},{"path":"/tmp/example/hello_world_gzip.txt","contents":{"compression":"gzip","source":"data:;base64,H4sIAAAAAAAC/zSQQc4UOwyE9+8UdYBRn+IhgcQSxNoknu6Skjjj2MOI06MwP7soicv1fd8uLnBB0LNcaOKnOj5ra4Yf5q3izqaISwK/2BpyKc7fnCjWp+tatAEOxKX4cg4GbRz4aq4dnCs7qjVzLAaka9xQbCwtoZEOqZxcLBwntDFuWFpRDcpc3SpC+zQHR2FlzRHIQJOf5gqNd7aiyzkE0vhIOfA9oIMdUtG5D08dlH7DI7kwbIVnhb7UC0N2YWRr0ou9k/cnLu5NfyM5oS+obORu1d4Ej5Q48P+OlAwFPV0/YDngOl0vHVWdsS+e1nKGhOK5SaFrKQpb+6dIoYl7npTA2IUwxSmRfuDTq+gMze1xBKwU0SKBkpNVYk/YwHRj1bEtblMcKNmmbG7Y/c5CQdWlvl+7tV1DtiBW6Prwmv34708AAAD//9sohz4WAgAA"},"mode":493}]},"systemd":{"units":[{"contents":"[Unit]\nDescription=Hello world service\n\n[Service]\nType=oneshot\nExecStart=/usr/bin/echo \"hello world\"\nStandardOutput=journal\n\n[Install]\nWantedBy=multi-user.target default.target\n","enabled":true,"name":"hello-world.service"}]}}
```

- Run Injection:

```shell
$ sudo injection example_files/example.ign
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
