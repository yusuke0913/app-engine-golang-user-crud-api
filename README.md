# app-engine-golang-user-crud-api
This package implements a CRUD operation APIs to manage User data on Google App Engine and Datastore.

# Example


## Directory layout
    .
    ├── app.yaml
    ├── main.go

- main.go
```go
package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yusuke0913/app-engine-golang-user-crud-api"
	"google.golang.org/appengine"
)

func main() {
	r := mux.NewRouter()
	users.Register(r)
	http.Handle("/", r)
	appengine.Main()
}
```

- app.yaml
```yaml
runtime: go
api_version: go1

handlers:
  # All URLs are handled by the Go application script
  - url: /.*
    script: _go_app

skip_files:
  - .*node_modules
  - .*vendor
````