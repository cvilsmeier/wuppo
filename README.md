
# wuppo

A web framework for Go, ridiculously simple.

## Usage

A very basic usage is this:


```go
func main() {
    handler := wuppo.DefaultHandler(func (req *wuppo.Req) {
        req.Html = fmt.Sprintf("<html>Hello %s</html>", req.Path)
    })
    http.Handle("/", handler)
    log.Panic(http.ListenAndServe(":8080", nil))
}
```

See the examples folder for more samples

## Licence

Hell, it's free! Do whatevery you like.

