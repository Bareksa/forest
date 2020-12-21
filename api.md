# forest API Documentation

There are two ways to initiate forest. Returning a client instance or start a global
initialization.

Use global initialization method for simple use everywhere in the code, but less flexible. If using
multiple instances (like connecting to multiple Vault hosts) please use the Client instance method,
and pass the instance(s) around.

### [Method Documentations Here](./package.md)

## Global Initialization

```go
// main.go

package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"gitlab.bareksa.com/backend/forest"
)

func main() {
	token := "abc"
	host := "http://127.0.0.1"
	// This initializes forest and allow usage for
	// other methods
	forest.Init(token, forest.WithHost(host))
}
```

## Client Instance Initialization With Transit Engine different than default

```go
// main.go

package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"gitlab.bareksa.com/backend/forest"
)

func main() {
	token := "abc"
	host := "http://127.0.0.1"
	transitEngine := "sayap-krispy"
	forest.NewClient(
		token,
		forest.WithHost(host),
		forest.WithTransitEngine(transitEngine),
	)
}
```

## Global Initialization Example With Different KV Engine and Transit Engine different than default

```go
// main.go

package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"gitlab.bareksa.com/backend/forest"
)

func main() {
	var token string
	var host string
	flag.StringVar(&token, "token", "", "vaultex -token [token]")
	flag.StringVar(&host, "host", "", "vaultex -host [vault host]")
	flag.Parse()
	if host == "" {
		log.Fatal("Host is empty")
	}
	err := forest.Init(
		token,
		forest.WithHost(host),
		forest.WithTransitEngine("ayam-kuning"),
		forest.WithKeyValueEngine("bukan-consul"),
	)
	if err != nil {
		log.Fatal(err)
	}
	cipher, err := forest.TransitEncrypt(context.Background(), "awawa", []byte(`1234`))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cipher)
}
```

## Client Instance Initialization Example

```go
// main.go

package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"gitlab.bareksa.com/backend/forest"
)

func main() {
	var token string
	var host string
	flag.StringVar(&token, "token", "", "vaultex -token [token]")
	flag.StringVar(&host, "host", "", "vaultex -host [vault host]")
	flag.Parse()
	if host == "" {
		log.Fatal("Host is empty")
	}
	vault, err := forest.NewClient(
		token,
		forest.WithHost(host),
		forest.WithTransitEngine("ayam-kuning"),
		forest.WithKeyValueEngine("bukan-consul"),
	)
	if err != nil {
		log.Fatal(err)
	}
	config, err := vault.GetKeyValue(context.Background(), "ms-order-conf")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(config))
}
```
