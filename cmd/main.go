package main

import (
	"context"
	"log"
	"os"

	functions "github.com/SoteriaTech/blockchain-functions"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
)

func main() {
	ctx := context.Background()
	funcframework.RegisterHTTPFunctionContext(ctx, "/SyncBtcBalance", functions.SyncBtcBalance)
	funcframework.RegisterHTTPFunctionContext(ctx, "/ScanBtcBlock", functions.ScanBtcBlock)
	funcframework.RegisterHTTPFunctionContext(ctx, "/test", functions.ScanBtcHead)

	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}
