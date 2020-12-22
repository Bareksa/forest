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
	client := forest.NewClient("s.9bho6AeRSjyfObBNpgHUDH1Q") // This token is an example
	config, err := client.GetConfig(context.Background(), "/kv/foo")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(config))
}
```

Also you can use custom http clients

```go
httpClient := &http.Client{Timeout: 10 * time.Second}
client := forest.NewClient("s.9bho6AeRSjyfObBNpgHUDH1Q", forest.WithHttpClient(httpClient))
```

# Integration with Viper Example

```go
package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/Bareksa/forest"
	"github.com/spf13/viper"
)

func main() {
	token := flag.String("token", "", "-token [token]")
	host := flag.String("host", "http://127.0.0.1:8200", "-host [host]")
	flag.Parse()
	err := forest.Init(*token, forest.WithHost(*host))
	if err != nil {
		log.Fatal(err)
	}
	config, err := forest.GetKeyValue(context.Background(), "microservice-config")
	if err != nil {
		log.Fatal(err)
	}
	viper.SetConfigType("json")
	viper.ReadConfig(bytes.NewBuffer(config))
	ver := viper.GetString("app_version")
	fmt.Println(ver)
}
```

# Running Test

Please note integration test have to be modified for own use until vault dev is ready.

```bash
$ go test -token [token] -host [host] -v ./...
```

### [API Usage Documentation](./api.md)

### [Package Method Documentation](./package.md)

## Contributor

-   [Tigor Hutasuhut](https://gitlab.bareksa.com/tigor)
-   [Arif Rakhman](https://gitlab.bareksa.com/arif_rachman)
