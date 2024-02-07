package util

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/xvlcwk-terraform/terraform-provider-bitbucketserver/bitbucket/util/client"
)

type (
	// AccessTokenResourceHelper provides assistive snippets of logic to help reduce duplication in
	// each resource definition.
	AccessTokenResourceHelper struct {
		Client *client.BitbucketClient
		helper *ResourceHelper
	}
)

func NewAccessTokenResourceHelper() *AccessTokenResourceHelper {
	return &AccessTokenResourceHelper{
		helper: NewResourceHelper(),
	}
}

func (r *AccessTokenResourceHelper) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.helper.Configure(ctx, req, resp)
	r.Client = r.helper.client
}

func (r *AccessTokenResourceHelper) Schema(s map[string]schema.Attribute) map[string]schema.Attribute {
	s = r.helper.Schema(s)
	if _, ok := s["id"]; !ok {
		s["id"] = schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The id of the access token",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		}
	}
	if _, ok := s["name"]; !ok {
		s["name"] = schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The label for the access token",
		}
	}
	if _, ok := s["project"]; !ok {
		s["project"] = schema.StringAttribute{
			Required:    true,
			Description: "The project slug",
		}
	}
	if _, ok := s["token"]; !ok {
		s["token"] = schema.StringAttribute{
			Computed:    true,
			Description: "The token. Only set if created by Terraform",
		}
	}
	if _, ok := s["permissions"]; !ok {
		s["permissions"] = schema.SetAttribute{
			Required:    true,
			Description: fmt.Sprintf("The permissions this access token has for repositories."),
			ElementType: types.StringType,
		}
	}
	if _, ok := s["expire_in"]; !ok {
		s["expire_in"] = schema.Int64Attribute{
			Computed:    true,
			Description: "Expire in X Days. If not set it does not expire",
			Default:     int64default.StaticInt64(1095),
		}
	}
	if _, ok := s["created_date"]; !ok {
		s["created_date"] = schema.Int64Attribute{
			Computed:    true,
			Description: "Created Date",
		}
	}
	return s
}

type CreateAccessTokenRequest struct {
	ExpiryDays  int64    `json:"expiryDays" binding:"required"`
	Name        string   `json:"name" binding:"required"`
	Permissions []string `json:"permissions" binding:"required"`
}

type AccessTokenResponse struct {
	Token       string `json:"token,omitempty"` // Only available on creation
	Name        string `json:"name" binding:"required"`
	Id          string `json:"id" binding:"required"`
	CreatedDate int64  `json:"createdDate" binding:"required"`
}
