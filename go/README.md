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
