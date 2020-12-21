# Begundal API Documentation

There are two ways to initiate Begundal. Returning a client instance or start a global
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

	"gitlab.bareksa.com/backend/begundal"
)

func main() {
	token := "abc"
	host := "http://127.0.0.1"
	// This initializes begundal and allow usage for
	// other methods
	begundal.Init(token, begundal.WithHost(host))
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

	"gitlab.bareksa.com/backend/begundal"
)

func main() {
	token := "abc"
	host := "http://127.0.0.1"
	transitEngine := "sayap-krispy"
	begundal.NewClient(
		token,
		begundal.WithHost(host),
		begundal.WithTransitEngine(transitEngine),
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

	"gitlab.bareksa.com/backend/begundal"
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
	err := begundal.Init(
		token,
		begundal.WithHost(host),
		begundal.WithTransitEngine("ayam-kuning"),
		begundal.WithKeyValueEngine("bukan-consul"),
	)
	if err != nil {
		log.Fatal(err)
	}
	cipher, err := begundal.TransitEncrypt(context.Background(), "awawa", []byte(`1234`))
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

	"gitlab.bareksa.com/backend/begundal"
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
	vault, err := begundal.NewClient(
		token,
		begundal.WithHost(host),
		begundal.WithTransitEngine("ayam-kuning"),
		begundal.WithKeyValueEngine("bukan-consul"),
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
