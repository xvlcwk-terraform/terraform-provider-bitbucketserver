package bitbucket

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/xvlcwk-terraform/terraform-provider-bitbucketserver/bitbucket/util/client"
	types2 "github.com/xvlcwk-terraform/terraform-provider-bitbucketserver/bitbucket/util/types"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/xvlcwk-terraform/terraform-provider-bitbucketserver/bitbucket/marketplace"
)

func New() provider.Provider {
	return &BitbucketServerProviderFramework{}
}

var _ provider.Provider = (*BitbucketServerProviderFramework)(nil)

type BitbucketServerProviderFramework struct {
	BitbucketClient   *client.BitbucketClient
	MarketplaceClient *marketplace.Client
}

type BitbucketServerProviderModel struct {
	Server   types.String `tfsdk:"server"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	Token    types.String `tfsdk:"token"`
}

func (p *BitbucketServerProviderFramework) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		newRepositoryAccessTokenResource,
	}
}

func (p *BitbucketServerProviderFramework) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *BitbucketServerProviderFramework) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"server": schema.StringAttribute{
				Optional:    true,
				Description: "The url of your bitbucket instance. For the docker compose instance this is http://localhost:7990",
				Validators: []validator.String{
					// Validate string value satisfies the regular expression for alphanumeric characters
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^https?://.*$`),
						"must contain the http schema (http:// or https://)",
					),
				},
				Sensitive: false,
			},
			"username": schema.StringAttribute{
				Optional:    true,
				Description: "The username for authentication. If you're using a personal access token use your normal username.",
				Sensitive:   false,
			},
			"password": schema.StringAttribute{
				Optional:    true,
				Description: "the password for authentication. Personal access tokens are allowed, but http access token aren't yet",
				Sensitive:   true,
			},
			"token": schema.StringAttribute{
				Optional:    true,
				Description: "Token as alternative to the Password. Only use for repository access tokens. Personal access tokens can use the normal basic authentication",
				Sensitive:   true,
			},
		},
	}
}

func (p *BitbucketServerProviderFramework) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "bitbucketserver"
}

func (p *BitbucketServerProviderFramework) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data BitbucketServerProviderModel

	// Read configuration data into model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	server := os.Getenv("BITBUCKET_SERVER")
	if data.Server.ValueString() != "" {
		server = data.Server.ValueString()
	}

	if strings.HasSuffix(server, "/") {
		server = server[0 : len(server)-1]
	}

	if server == "" {
		resp.Diagnostics.AddError(
			"server is required",
			"server is required and must be provided in the provider config or the BITBUCKET_SERVER environment variable",
		)
	}

	token := os.Getenv("BITBUCKET_TOKEN")
	if data.Token.ValueString() != "" {
		token = data.Token.ValueString()
	}

	username := os.Getenv("BITBUCKET_USERNAME")
	if data.Username.ValueString() != "" {
		username = data.Username.ValueString()
	}
	if username == "" && token == "" {
		resp.Diagnostics.AddError(
			"username is required",
			"username is required and must be provided in the provider config or the BITBUCKET_USERNAME environment variable",
		)
	}

	password := os.Getenv("BITBUCKET_PASSWORD")
	if data.Password.ValueString() != "" {
		password = data.Password.ValueString()
	}
	if password == "" && token == "" {
		resp.Diagnostics.AddError(
			"password is required",
			"password is required and must be provided in the provider config or the BITBUCKET_PASSWORD environment variable",
		)
	}

	b := &client.BitbucketClient{
		Server:     server,
		Username:   username,
		Password:   password,
		Token:      token,
		HTTPClient: &http.Client{},
	}

	m := &marketplace.Client{
		HTTPClient: &http.Client{},
	}

	resp.ResourceData = &types2.BitbucketServerProvider{
		BitbucketClient:   b,
		MarketplaceClient: m,
	}

	resp.DataSourceData = &types2.BitbucketServerProvider{
		BitbucketClient:   b,
		MarketplaceClient: m,
	}
}
