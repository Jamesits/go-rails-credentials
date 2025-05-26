package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jamesits/go-rails-credentials/pkg/credentials"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &RailsCredentialsFileDataSource{}

func NewRailsCredentialsFileDataSource() datasource.DataSource {
	return &RailsCredentialsFileDataSource{}
}

// RailsCredentialsFileDataSource defines the data source implementation.
type RailsCredentialsFileDataSource struct{}

// RailsCredentialsFileDataSourceModel describes the data source data model.
type RailsCredentialsFileDataSourceModel struct {
	MasterKey        types.String `tfsdk:"master_key"`
	EncryptedContent types.String `tfsdk:"encrypted_content"`
	DecryptedContent types.String `tfsdk:"content"`
}

func (d *RailsCredentialsFileDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file"
}

func (d *RailsCredentialsFileDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Reads and decrypts a Rails credentials file.",

		Attributes: map[string]schema.Attribute{
			"master_key": schema.StringAttribute{
				MarkdownDescription: "The master key",
				Required:            true,
				Sensitive:           true,
			},
			"encrypted_content": schema.StringAttribute{
				MarkdownDescription: "The credentials file content",
				Required:            true,
			},
			"content": schema.StringAttribute{
				MarkdownDescription: "Decrypted credentials in YAML format",
				Computed:            true,
				Sensitive:           true,
			},
		},
	}
}

func (d *RailsCredentialsFileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
}

func (d *RailsCredentialsFileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RailsCredentialsFileDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rawObject, err := credentials.Decrypt(data.MasterKey.ValueString(), data.EncryptedContent.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Credentials decryption failed", err.Error())
		return
	}
	rawString, err := credentials.UnmarshalSingleString(rawObject)
	if err != nil {
		resp.Diagnostics.AddError("Credentials unmarshal failed", err.Error())
		return
	}
	data.DecryptedContent = types.StringValue(rawString)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
