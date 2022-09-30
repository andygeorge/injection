# Injection

Injection is a simple utility that consumes and processes [Ignition](https://coreos.github.io/ignition/) configs and writes out directories, files, and systemd units (and enables/disables them!) on an **already-running Linux host**, differing from Ignition, which is only intended to configure immutable Fedora/RedHat CoreOS hosts on first boot.

In conjunction with [Butane](https://coreos.github.io/butane/), this allows you to build human-readable configuration files that can be used to easily (re)configure existing hosts.

## Prerequisites

- Linux/macOS
- [Butane](https://coreos.github.io/butane/getting-started/#getting-butane), if you're planning on using Butane to generate Ignition configs (recommended)

## Installation

```shell
sudo wget https://github.com/andygeorge/injection/releases/download/v0.1.1-beta/injection-amd64-linux -O /usr/local/bin/injection && sudo chmod a+x /usr/local/bin/injection
```

Or, manually:
- Download an [`injection` binary release](https://github.com/andygeorge/injection/releases)
- Place `injection` somewhere in your `$PATH`

## Example

- Clone and enter this repo
```shell
gh repo clone andygeorge/injection
cd injection
```

- Use the example Butane config file:
```shell
cp example_butane.bu example_files/example.bu
```

- The file `example_files/example.bu` should be populated

- Create an Ignition config using Butane:
```shell
butane --strict example_files/example.bu > example_files/example.ign
```

- The file `example_files/example.ign` should be populated and look like this:
```json
{"ignition":{"version":"3.3.0"},"storage":{"directories":[{"path":"/tmp/example","mode":493}],"files":[{"path":"/tmp/example/hello_world.txt","contents":{"compression":"","source":"data:,Hello%2C%20world!%0A"},"mode":420},{"path":"/tmp/example/hello_world_gzip.txt","contents":{"compression":"gzip","source":"data:;base64,H4sIAAAAAAAC/zSQQc4UOwyE9+8UdYBRn+IhgcQSxNoknu6Skjjj2MOI06MwP7soicv1fd8uLnBB0LNcaOKnOj5ra4Yf5q3izqaISwK/2BpyKc7fnCjWp+tatAEOxKX4cg4GbRz4aq4dnCs7qjVzLAaka9xQbCwtoZEOqZxcLBwntDFuWFpRDcpc3SpC+zQHR2FlzRHIQJOf5gqNd7aiyzkE0vhIOfA9oIMdUtG5D08dlH7DI7kwbIVnhb7UC0N2YWRr0ou9k/cnLu5NfyM5oS+obORu1d4Ej5Q48P+OlAwFPV0/YDngOl0vHVWdsS+e1nKGhOK5SaFrKQpb+6dIoYl7npTA2IUwxSmRfuDTq+gMze1xBKwU0SKBkpNVYk/YwHRj1bEtblMcKNmmbG7Y/c5CQdWlvl+7tV1DtiBW6Prwmv34708AAAD//9sohz4WAgAA"},"mode":420}]},"systemd":{"units":[{"contents":"[Unit]\nDescription=Hello world service\n\n[Service]\nType=oneshot\nExecStart=/usr/bin/echo \"hello world\"\nStandardOutput=journal\n\n[Install]\nWantedBy=multi-user.target default.target\n","enabled":true,"name":"hello-world.service"}]}}
```

- Run Injection:
```shell
sudo injection example_files/example.ign
```

This will write out the directories, files, and systemd units defined in the example:
```shell
/tmp/example
/tmp/example2
/tmp/example3/example4
/tmp/example/hello_world.txt
/tmp/example/hello_world_gzip.txt
/etc/systemd/system/hello-world.service
```

## Ignition support

This only has basic Ignition support for a few specific fields:

- `storage:files`
- `storage:directories`
- `systemd:units`

Support for `passwd:users` is [planned](https://github.com/andygeorge/injection/issues/1).


----
# injection-go
`injection` built in Go.

```shell
# generate ignition
butane --strict ../example_files/example.bu > ../example_files/example.ign

# build injection
go mod tidy && go fmt *.go && go build -o build/

# run injection
./build/injection ../example_files/example.ign
```
