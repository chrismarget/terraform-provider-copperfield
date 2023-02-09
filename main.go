package main

import (
	"context"
	"flag"
	"github.com/chrismarget/terraform-provider-copperfield/copperfield"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"log"
)

func NewProvider() provider.Provider {
	return &copperfield.Provider{}
}

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	err := providerserver.Serve(context.Background(), NewProvider, providerserver.ServeOpts{
		Address: "example.com/chrismarget/copperfield",
		Debug:   debug,
	})
	if err != nil {
		log.Fatal(err)
	}
}
