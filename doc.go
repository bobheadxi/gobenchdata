/*

Gobenchdata is a tool for inspecting golang benchmark outputs. To install it,
you must have Go installed:

	go get -u go.bobheadxi.dev/gobenchdata
	gobenchdata help

Then pipe your benchmark into the tool:

	go test -bench . -benchmem ./... | gobenchdata --json bench.json

You can create a sort of database of benchmarks by appending new benchmarks to
an existing file:

	go test -benchtime 10000x -bench . -benchmem ./... | gobenchdata --json benchmarks.json --append

You can also merge results:

	gobenchdata merge file1.json file2.json file3.json

*/
package main
