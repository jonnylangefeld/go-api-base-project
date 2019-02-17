# GO API Base Project

This repository is my favorite way of kicking of a new API project in go.
Read [my blog post](https://jonnylangefeld.com/blog/how-to-write-a-go-api-part-1-webserver-with-iris) for all the details.

### Run it

* [Install go](https://golang.org/doc/install) (v1.11 or later)
* `go build`
* `./my-go-api`

### Try it

* `curl localhost:8080`
* `curl localhost:8080/Jonny`
* `curl localhost:8080/orders`

### Test it

* `go test -v`