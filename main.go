package main

import (
	"context"
	"log"
	"runtime"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/cache"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
)

type exportReplaceLogger struct {
}

func (a *exportReplaceLogger) Replace(ctx context.Context, cache cache.Unmarshaler, cacheHints cache.ReplaceHints) error {
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		log.Printf("Replace(%s): called from %s\n", cacheHints.PartitionKey, details.Name())
	}

	return nil
}

func (a *exportReplaceLogger) Export(ctx context.Context, cache cache.Marshaler, cacheHints cache.ExportHints) error {
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		log.Printf("Export(%s): called from %s\n", cacheHints.PartitionKey, details.Name())
	}

	return nil
}

func main() {
	cache := &exportReplaceLogger{}

	options := []public.Option{
		public.WithCache(cache),
		public.WithAuthority("https://login.microsoftonline.com/organizations"),
	}

	publicClientApp, err := public.New("04b07795-8ddb-461a-bbee-02f9e1bf7b46", options...)
	if err != nil {
		panic(err)
	}

	scopes := []string{"https://management.azure.com//.default"}
	ctx := context.Background()
	_, err = publicClientApp.AcquireTokenInteractive(ctx, scopes)
	if err != nil {
		panic(err)
	}

	acc, err := publicClientApp.Accounts(ctx)
	if err != nil {
		panic(err)
	}

	if len(acc) != 1 {
		panic("expected 1 account")
	}

	_, err = publicClientApp.AcquireTokenSilent(ctx, scopes, public.WithSilentAccount(acc[0]))
	if err != nil {
		panic(err)
	}
}
