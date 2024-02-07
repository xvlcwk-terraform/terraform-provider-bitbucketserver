package bitbucket

import (
	"context"
	"github.com/hashicorp/go-cty/cty"
	"github.com/xvlcwk-terraform/terraform-provider-bitbucketserver/bitbucket/util/client"
	"net/http"
	"strings"

	"github.com/xvlcwk-terraform/terraform-provider-bitbucketserver/bitbucket/marketplace"
	bitbucketTypes "github.com/xvlcwk-terraform/terraform-provider-bitbucketserver/bitbucket/util/types"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"server": {
				Optional:    true,
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("BITBUCKET_SERVER", nil),
				Description: "The url of your bitbucket instance. For the docker compose instance this is http://localhost:7990",
			},
			"username": {
				Optional:    true,
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("BITBUCKET_USERNAME", nil),
				Description: "The username for authentication. If you're using a personal access token use your normal username.",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("BITBUCKET_PASSWORD", nil),
				Description: "the password for authentication. Personal access tokens are allowed, but http access token aren't yet",
			},
			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("BITBUCKET_TOKEN", nil),
				Description: "Token as alternative to the Password. Only use for repository access tokens. Personal access tokens can use the normal basic authentication",
			},
		},
		ConfigureContextFunc: providerConfigure,
		DataSourcesMap: map[string]*schema.Resource{
			"bitbucketserver_application_properties":        dataSourceApplicationProperties(),
			"bitbucketserver_cluster":                       dataSourceCluster(),
			"bitbucketserver_global_permissions_groups":     dataSourceGlobalPermissionsGroups(),
			"bitbucketserver_global_permissions_users":      dataSourceGlobalPermissionsUsers(),
			"bitbucketserver_groups":                        dataSourceGroups(),
			"bitbucketserver_group_users":                   dataSourceGroupUsers(),
			"bitbucketserver_plugin":                        dataSourcePlugin(),
			"bitbucketserver_project_hooks":                 dataSourceProjectHooks(),
			"bitbucketserver_project_permissions_groups":    dataSourceProjectPermissionsGroups(),
			"bitbucketserver_project_permissions_users":     dataSourceProjectPermissionsUsers(),
			"bitbucketserver_repository_hooks":              dataSourceRepositoryHooks(),
			"bitbucketserver_repository_permissions_groups": dataSourceRepositoryPermissionsGroups(),
			"bitbucketserver_repository_permissions_users":  dataSourceRepositoryPermissionsUsers(),
			"bitbucketserver_user":                          dataSourceUser(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"bitbucketserver_banner":                       resourceBanner(),
			"bitbucketserver_default_reviewers_condition":  resourceDefaultReviewersCondition(),
			"bitbucketserver_global_permissions_group":     resourceGlobalPermissionsGroup(),
			"bitbucketserver_global_permissions_user":      resourceGlobalPermissionsUser(),
			"bitbucketserver_group":                        resourceGroup(),
			"bitbucketserver_license":                      resourceLicense(),
			"bitbucketserver_mail_server":                  resourceMailServer(),
			"bitbucketserver_plugin":                       resourcePlugin(),
			"bitbucketserver_plugin_config":                resourcePluginConfig(),
			"bitbucketserver_project":                      resourceProject(),
			"bitbucketserver_project_branch_permissions":   resourceBranchPermissions(),
			"bitbucketserver_project_hook":                 resourceProjectHook(),
			"bitbucketserver_project_permissions_group":    resourceProjectPermissionsGroup(),
			"bitbucketserver_project_permissions_user":     resourceProjectPermissionsUser(),
			"bitbucketserver_repository":                   resourceRepository(),
			"bitbucketserver_repository_deploy_key":        resourceRepositoryDeployKey(),
			"bitbucketserver_repository_hook":              resourceRepositoryHook(),
			"bitbucketserver_repository_permissions_group": resourceRepositoryPermissionsGroup(),
			"bitbucketserver_repository_permissions_user":  resourceRepositoryPermissionsUser(),
			"bitbucketserver_repository_webhook":           resourceRepositoryWebhook(),
			"bitbucketserver_user":                         resourceUser(),
			"bitbucketserver_user_access_token":            resourceUserAccessToken(),
			"bitbucketserver_user_group":                   resourceUserGroup(),
		},
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	serverSanitized := d.Get("server").(string)
	if strings.HasSuffix(serverSanitized, "/") {
		serverSanitized = serverSanitized[0 : len(serverSanitized)-1]
	}

	username := d.Get("username").(string)
	password := d.Get("password").(string)
	token := d.Get("token").(string)

	configErrors := diag.Diagnostics{}

	if serverSanitized == "" {
		configErrors = append(configErrors,
			diag.Diagnostic{
				Severity:      diag.Error,
				AttributePath: cty.Path{}.GetAttr("server"),
				Detail:        "server is required and must be provided in the provider config or the BITBUCKET_SERVER environment variable",
			})
	}
	if username == "" && token == "" {
		configErrors = append(configErrors,
			diag.Diagnostic{

				Severity:      diag.Error,
				AttributePath: cty.Path{}.GetAttr("username"),
				Detail:        "username is required and must be provided in the provider config or the BITBUCKET_USERNAME environment variable",
			})
	}
	if password == "" && token == "" {
		configErrors = append(configErrors,
			diag.Diagnostic{
				Severity:      diag.Error,
				AttributePath: cty.Path{}.GetAttr("password"),
				Detail:        "password is required and must be provided in the provider config or the BITBUCKET_PASSWORD environment variable",
			})
	}

	if configErrors.HasError() {
		return nil, configErrors
	}

	b := &client.BitbucketClient{
		Server:     serverSanitized,
		Username:   username,
		Password:   password,
		Token:      token,
		HTTPClient: &http.Client{},
	}

	m := &marketplace.Client{
		HTTPClient: &http.Client{},
	}

	return &bitbucketTypes.BitbucketServerProvider{
		BitbucketClient:   b,
		MarketplaceClient: m,
	}, nil
}
