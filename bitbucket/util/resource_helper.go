package util

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	client2 "github.com/xvlcwk-terraform/terraform-provider-bitbucketserver/bitbucket/util/client"
	bitbucketTypes "github.com/xvlcwk-terraform/terraform-provider-bitbucketserver/bitbucket/util/types"
)

type (
	// ResourceHelper provides assistive snippets of logic to help reduce duplication in
	// each resource definition.
	ResourceHelper struct {
		client *client2.BitbucketClient
	}
)

func NewResourceHelper() *ResourceHelper {
	return &ResourceHelper{}
}

// Configure should register the client for the resource.
func (r *ResourceHelper) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	provider, ok := req.ProviderData.(*bitbucketTypes.BitbucketServerProvider)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *bitbucket.BitbucketServerProvider, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = provider.BitbucketClient
}

func (r *ResourceHelper) Schema(s map[string]schema.Attribute) map[string]schema.Attribute {
	return s
}
