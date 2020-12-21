package begundal

import (
	"flag"
	"log"
	"testing"
)

var testToken = flag.String("token", "", "Token for vault")
var testHost = flag.String("host", "http://127.0.0.1:8200", "Vault host")

func init() {
	var _ = func() bool {
		testing.Init()
		return true
	}()
	flag.Parse()
	err := Init(*testToken, WithHost(*testHost))
	if err != nil {
		log.Fatal(err)
	}
}
