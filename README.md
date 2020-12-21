# forest

![forthebadge made-with-Go](http://ForTheBadge.com/images/badges/made-with-go.svg)

## What?

A simplified vault request. Taylored to Bareksa's paradigm.

# What Does it Do?

1. Gets key value store from vault.

# Why?

1. You (the user) don't need to research how to use Vault much, and example is given directly in code description.

# How to Install

```bash
$ GOPRIVATE=gitlab.bareksa.com go get -u gitlab.bareksa.com/backend/forest
```

Since this is a private repo, we may need to override the _https_ protocol for `go get` to _ssh_.
Obviously this need SSH private key for your pc registered in your gitlab account.

Run this in your machine.

```bash
$ git config --global --add url."git@gitlab.bareksa.com:".insteadOf "https://gitlab.bareksa.com/"
```

Do note for every `go get` or `git clone` request to `gitlab.bareksa.com` will now enforced to use SSH. Delete
the corresponding line in `~/.gitconfig` to disable this.

```go
package main

import (
	"context"
	"fmt"
	"gitlab.bareksa.com/backend/forest"
)

func main() {
	client := forest.NewClient("s.9bho6AeRSjyfObBNpgHUDH1Q")
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
