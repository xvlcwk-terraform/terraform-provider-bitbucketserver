package bitbucket

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"bitbucketserver": func() (tfprotov6.ProviderServer, error) {
		ctx := context.Background()

		upgradedSdkServer, err := tf5to6server.UpgradeServer(
			ctx,
			//  terraform-plugin-sdk provider
			Provider().GRPCProvider, //nolint:staticcheck
		)

		if err != nil {
			return nil, err
		}

		providers := []func() tfprotov6.ProviderServer{
			providerserver.NewProtocol6(New()), // Example terraform-plugin-framework provider
			func() tfprotov6.ProviderServer {
				return upgradedSdkServer
			},
		}

		muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)

		if err != nil {
			return nil, err
		}

		return muxServer.ProviderServer(), nil
	},
}

func TestMuxServer(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					provider "bitbucketserver"{
						server = "http://localhost:7990"
						username = "admin"
						password = "admin"
					}
					resource "bitbucketserver_banner" "testbanner" {
						message = "testing muxing"
					}
				`,
			},
		},
	})
}
