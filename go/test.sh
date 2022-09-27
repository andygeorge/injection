#!/usr/bin/env bash
/usr/bin/butane --strict ../example_files/example.bu > ../example_files/example.ign
rm -rf /tmp/example/
/usr/local/go/bin/go mod tidy && /usr/local/go/bin/go fmt *.go && /usr/local/go/bin/go build -o build/
./build/injection ../example_files/example.ign
