# WAX Blockchain Go SDK

This SDK is used to get data from the WAX Blockchain

### Basic usage

```go
package main

import (
	"context"
	"log"
	"time"

	wax "github.com/luktech-labs/wax-go-sdk"
)

func main() {
	waxSdk := wax.NewSdk("https://wax.greymass.com")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	info, err := waxSdk.GetInfoContext(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Chain info: %+v\n", info)

	// Let's suppose that we want to get the account info of the farmersworld's account `asdasdasd225`
	type AccountInfo struct {
		CurrentEnergy int      `json:"energy"`
		MaxEnergy     int      `json:"max_energy"`
		Balances      []string `json:"balances"`
	}

	fwAccountsPayload := wax.GetTableRowsPayload{
		Json:       true,
		Code:       "farmersworld",
		Scope:      "farmersworld",
		Table:      "accounts",
		LowerBound: "asdasdasd225",
		UpperBound: "asdasdasd225",
		Limit:      "1",
	}

	var accountInfo []AccountInfo

	err = waxSdk.GetTableRowsContext(ctx, fwAccountsPayload, &accountInfo)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%+v\n", accountInfo)
}
```

### Running with a pool of proxies
```go
package main

import (
	"context"
	"log"
	"time"

	wax "github.com/luktech-labs/wax-go-sdk"
)

func main() {
	proxies := []string{"https://proxy1", "https://proxy2"}
	
	// this will randomly select one of the proxies from the pool at each request.
	proxyOpt := wax.WithProxies(proxies)
	waxSdk := wax.NewSdk("https://wax.greymass.com", proxyOpt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	info, err := waxSdk.GetInfoContext(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Chain info: %+v\n", info)
}
```

### Running tests

```shell
[luukkk@wax]$ go test -v ./...
=== RUN   TestCreateRequestErrorMessage
=== RUN   TestCreateRequestErrorMessage/response_nil
    http_test.go:58: ✓   Messages should be identical. 
=== RUN   TestCreateRequestErrorMessage/response_empty_body
    http_test.go:58: ✓   Messages should be identical. 
=== RUN   TestCreateRequestErrorMessage/response_nil_body
    http_test.go:58: ✓   Messages should be identical. 
=== RUN   TestCreateRequestErrorMessage/response_with_body
    http_test.go:58: ✓   Messages should be identical. 
--- PASS: TestCreateRequestErrorMessage (0.00s)
    --- PASS: TestCreateRequestErrorMessage/response_nil (0.00s)
    --- PASS: TestCreateRequestErrorMessage/response_empty_body (0.00s)
    --- PASS: TestCreateRequestErrorMessage/response_nil_body (0.00s)
    --- PASS: TestCreateRequestErrorMessage/response_with_body (0.00s)
PASS
ok      github.com/luktech-labs/wax-go-sdk 0.127s
```

