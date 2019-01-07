# app-engine-golang-user-crud-api
This is simple CRUD operation APIs to manage User data on Google App Engine and Datastore.

# Example
```
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