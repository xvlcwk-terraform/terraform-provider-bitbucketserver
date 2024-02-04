package bitbucket

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/xvlcwk-terraform/terraform-provider-bitbucketserver/bitbucket/util"
	"io"
	"net/http"
)

type repositoryAccessTokenModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Permissions types.Set    `tfsdk:"permissions"`
	ExpireIn    types.Int64  `tfsdk:"expire_in"`
	Project     types.String `tfsdk:"project"`
	Repository  types.String `tfsdk:"repository"`
	Token       types.String `tfsdk:"token"`
	CreatedDate types.Int64  `tfsdk:"created_date"`
}

type repositoryAccessTokenResource struct {
	resourceHelper *util.AccessTokenResourceHelper
}

func newRepositoryAccessTokenResource() resource.Resource {
	return &repositoryAccessTokenResource{
		resourceHelper: util.NewAccessTokenResourceHelper(),
	}
}

// Ensure the implementation satisfies the desired interfaces.
var _ resource.ResourceWithConfigure = &repositoryAccessTokenResource{}

// Metadata should return the full name of the resource.
func (r *repositoryAccessTokenResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository_access_token"
}

// Schema should return the schema for this resource.
func (r *repositoryAccessTokenResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "An HTTP-Access token limited to the given repository",
		Attributes: r.resourceHelper.Schema(map[string]schema.Attribute{
			"repository": schema.StringAttribute{
				Required:    true,
				Description: "The repository slug",
			},
		}),
	}
}

func (r *repositoryAccessTokenResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data repositoryAccessTokenModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload, diagnostics := r.createRequestData(ctx, data)

	if diagnostics != nil {
		resp.Diagnostics.Append(diagnostics...)
	}

	tokenResponse, tokenErrorResponse := r.resourceHelper.Client.Put(r.getUrlForProject(data), bytes.NewBuffer(payload))
	response, convertingResponseDiagnostics := r.readResponse(tokenErrorResponse, tokenResponse, &data)
	if convertingResponseDiagnostics != nil {
		resp.Diagnostics.Append(convertingResponseDiagnostics)
		return
	}

	data.Token = types.StringValue(response.Token)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
func (r *repositoryAccessTokenResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data repositoryAccessTokenModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tokenResponse, tokenErrorResponse := r.resourceHelper.Client.Get(r.getUrlForId(data))
	_, diagnostic := r.readResponse(tokenErrorResponse, tokenResponse, &data)
	if diagnostic != nil {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *repositoryAccessTokenResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data repositoryAccessTokenModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload, diagnostics := r.createRequestData(ctx, data)

	if diagnostics != nil {
		resp.Diagnostics.Append(diagnostics...)
	}
	tokenResponse, tokenErrorResponse := r.resourceHelper.Client.Post(r.getUrlForId(data), bytes.NewBuffer(payload))
	_, convertingResponseDiagnostics := r.readResponse(tokenErrorResponse, tokenResponse, &data)
	if convertingResponseDiagnostics != nil {
		resp.Diagnostics.Append(convertingResponseDiagnostics)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *repositoryAccessTokenResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data repositoryAccessTokenModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tokenResponse, tokenErrorResponse := r.resourceHelper.Client.Delete(r.getUrlForId(data))
	if tokenErrorResponse != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic(
			"Unable to Delete Resource",
			"An unexpected error occurred while deleting the resource. "+
				"Please report this issue to the provider developers.\n\n"+
				"Error: "+tokenErrorResponse.Error()))
	}

	if tokenResponse.StatusCode != 204 {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic(
			"Unable to Delete Resource",
			"An unexpected statusCode occurred while deleting the resource. "+
				"Please report this issue to the provider developers.\n\n"+
				"Status: "+tokenResponse.Status))
	}
}

func (r *repositoryAccessTokenResource) Configure(ctx context.Context, configureRequest resource.ConfigureRequest, configureResponse *resource.ConfigureResponse) {
	r.resourceHelper.Configure(ctx, configureRequest, configureResponse)
}

func (r *repositoryAccessTokenResource) createRequestData(ctx context.Context, data repositoryAccessTokenModel) ([]byte, diag.Diagnostics) {
	var permissions []string
	permissionConversionDiagnostics := data.Permissions.ElementsAs(ctx, &permissions, false)
	if permissionConversionDiagnostics != nil {
		return nil, permissionConversionDiagnostics
	}
	tokenRequest := &util.CreateAccessTokenRequest{
		ExpiryDays:  data.ExpireIn.ValueInt64(),
		Name:        data.Name.ValueString(),
		Permissions: permissions,
	}
	payload, jsonEncodingError := json.Marshal(tokenRequest)
	if jsonEncodingError != nil {
		return nil, diag.Diagnostics{diag.NewErrorDiagnostic("Marshalling error", fmt.Sprintf("Failed to encode %v. Is it valid Json?", tokenRequest))}
	}
	return payload, nil
}

func (r *repositoryAccessTokenResource) readResponse(tokenErrorResponse error, tokenResponse *http.Response, data *repositoryAccessTokenModel) (*util.AccessTokenResponse, *diag.ErrorDiagnostic) {
	if tokenErrorResponse != nil {
		diagnostic := diag.NewErrorDiagnostic("http error", tokenErrorResponse.Error())
		return nil, &diagnostic
	}
	if tokenResponse.StatusCode != 200 {
		diagnostic := diag.NewErrorDiagnostic("http error", fmt.Sprintf("Response Status: %d", tokenResponse.StatusCode))
		return nil, &diagnostic
	}

	body, readBodyError := io.ReadAll(tokenResponse.Body)
	if readBodyError != nil {
		diagnostic := diag.NewErrorDiagnostic("Error reading response", "Failed to read response. Is it valid Json?")
		return nil, &diagnostic
	}

	response := &util.AccessTokenResponse{}
	unmarshallError := json.Unmarshal(body, response)
	if unmarshallError != nil {
		diagnostic := diag.NewErrorDiagnostic("Error reading response", fmt.Sprintf("Failed to read response %v. Is it valid Json?, %v", string(body), unmarshallError))
		return nil, &diagnostic
	}

	data.Id = types.StringValue(response.Id)
	data.CreatedDate = types.Int64Value(response.CreatedDate)
	data.Name = types.StringValue(response.Name)
	return response, nil
}

func (r *repositoryAccessTokenResource) getUrlForId(data repositoryAccessTokenModel) string {
	return fmt.Sprintf("%v/%v", r.getUrlForProject(data), data.Id.ValueString())
}

func (r *repositoryAccessTokenResource) getUrlForProject(data repositoryAccessTokenModel) string {
	projectKey := data.Project.ValueString()
	repositorySlug := data.Repository.ValueString()
	repositoryUrl := fmt.Sprintf("/rest/access-tokens/latest/projects/%v/repos/%v", projectKey, repositorySlug)
	return repositoryUrl
}
