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

You can create a sort of database of benchmarks by appending new benchmarks to
an existing file:

```
go test -benchtime 10000x -bench . -benchmem ./... | gobenchdata --json benchmarks.json --append
```

You can also merge results:

```
gobenchdata merge file1.json file2.json file3.json
```
