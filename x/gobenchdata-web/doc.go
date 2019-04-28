/*

Gobenchdata-web is a utility for generating a website template to visualize
your continuous benchmark. To install it:

	go get -u github.com/bobheadxi/gobenchdata/x/gobenchdata-web
	gobenchdata-web help

An example usage might be:

	git checkout gh-pages
	gobenchdata-web
	git commit -a -m "init gobenchdata website"

You can test the generated website using a static file server like `serve`
(https://www.npmjs.com/package/serve):

	serve .

*/
package main
