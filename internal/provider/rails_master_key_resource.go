package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jamesits/go-rails-credentials/pkg/credentials"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &RailsMasterKeyResource{}
var _ resource.ResourceWithImportState = &RailsMasterKeyResource{}

func NewRailsMasterKeyResource() resource.Resource {
	return &RailsMasterKeyResource{}
}

// RailsMasterKeyResource defines the resource implementation.
type RailsMasterKeyResource struct{}

// RailsMasterKeyResourceModel describes the resource data model.
type RailsMasterKeyResourceModel struct {
	MasterKey types.String `tfsdk:"master_key"`
}

func (r *RailsMasterKeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_master_key"
}

func (r *RailsMasterKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Random-generated Rails master key",

		Attributes: map[string]schema.Attribute{
			"master_key": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The master key",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *RailsMasterKeyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
}

func (r *RailsMasterKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RailsMasterKeyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	m, err := credentials.RandomMasterKey()
	if err != nil {
		resp.Diagnostics.AddError("unable to generate master key", err.Error())
		return
	}
	data.MasterKey = types.StringValue(m)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RailsMasterKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	//var data RailsMasterKeyResourceModel
	//
	//// Read Terraform prior state data into the model
	//resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	//
	//if resp.Diagnostics.HasError() {
	//	return
	//}

	// Save updated data into Terraform state
	//resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RailsMasterKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data RailsMasterKeyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RailsMasterKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	//var data RailsMasterKeyResourceModel

	// Read Terraform prior state data into the model
	//resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	//
	//if resp.Diagnostics.HasError() {
	//	return
	//}
}

func (r *RailsMasterKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("master_key"), req, resp)
}
