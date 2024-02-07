package types

import (
	"github.com/xvlcwk-terraform/terraform-provider-bitbucketserver/bitbucket/marketplace"
	"github.com/xvlcwk-terraform/terraform-provider-bitbucketserver/bitbucket/util/client"
)

type BitbucketServerProvider struct {
	BitbucketClient   *client.BitbucketClient
	MarketplaceClient *marketplace.Client
}
