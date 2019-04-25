# gobenchdata

a tool for manipulating `go test -bench` data.

## Usage

```
go get -u github.com/bobheadxi/gobenchdata
```

Then pipe your benchmark into the tool:

```
go test -bench . -benchmem ./... | gobenchdata --json bench.json
```
