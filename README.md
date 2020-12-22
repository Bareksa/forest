# forest

![forthebadge made-with-Go](http://ForTheBadge.com/images/badges/made-with-go.svg)

## What?

A simplified vault request. Taylored to Bareksa's paradigm.

# What Does it Do?

1. Gets key value store from vault.
2. Creates token from vault
3. Handles policy creation

# Why?

1. You (the user) don't need to research how to use Vault much, and example is given directly in code description.

```go
package main

import (
	"context"
	"fmt"
	"gitlab.bareksa.com/backend/forest"
)

func main() {
	client := forest.NewClient("s.token") // This token is an example
	conf, err := forest.GetKeyValue(context.Background(), "some-conf")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(conf))
}
```

Also you can use custom http clients

```go
httpClient := &http.Client{Timeout: 10 * time.Second}
client := forest.NewClient("s.token", forest.WithHttpClient(httpClient))
```

# Integration with Viper Example

```go
package main

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/Bareksa/forest"
	"github.com/spf13/viper"
)

func main() {
	vaultHost := os.Getenv("VAULT_HOST")
	vaultToken := os.Getenv("VAULT_TOKEN")

	forest.Init(vaultToken, forest.WithHost(vaultHost))

	conf, err := forest.GetKeyValue(context.Background(), "some-conf")
	if err != nil {
		fmt.Printf("Failed to read configuration from vault : %v", err)
		os.Exit(1)
	}

	viper.SetConfigType("json") // Need to explicitly set this to json
	viper.ReadConfig(bytes.NewBuffer(conf))

	fmt.Printf("Using configuration file from : %s \n", vaultHost)

}
```

# Running Test for this Library Package

Please note integration test have to be modified for own use until vault dev is ready.

```bash
$ go test -token [token] -host [host] -v ./...
```

### [API Usage Documentation](./api.md)

### [Package Method Documentation](./package.md)

## Contributor

-   [Tigor Hutasuhut](https://gitlab.bareksa.com/tigor)
-   [Arif Rakhman](https://gitlab.bareksa.com/arif_rachman)
