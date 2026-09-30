# Archive

Earlier experiments, kept because they are part of how the project got here.
Nothing in this folder is built, tested, or maintained.

The leading underscore is deliberate: the Go tool skips directories whose names
begin with `_`, so `go build ./...` and `go test ./...` ignore everything in
here.

## backend/

A Go HTTP server from the project's first pass at the idea, when the plan was
to run the language on a server and talk to it from the browser. It contains a
`gorilla/mux` hello-world endpoint and a copy of the `prose` natural-language
library's example, from when parsing English sentences directly looked like a
promising route.

It does not compile — `app.go` and `api.go` each declare `func main()` in the
same package — and its dependencies were never recorded in a module. The
language ended up as a plain interpreter you run locally instead, which needs
no server at all.
