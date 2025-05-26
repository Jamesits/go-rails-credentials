package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jamesits/go-rails-credentials/pkg/credentials"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &RailsCredentialsInlineDataSource{}

func NewRailsCredentialsInlineDataSource() datasource.DataSource {
	return &RailsCredentialsInlineDataSource{}
}

// RailsCredentialsInlineDataSource defines the data source implementation.
type RailsCredentialsInlineDataSource struct{}

// RailsCredentialsInlineDataSourceModel describes the data source data model.
type RailsCredentialsInlineDataSourceModel struct {
	MasterKey        types.String `tfsdk:"master_key"`
	EncryptedContent types.String `tfsdk:"encrypted_content"`
	DecryptedContent types.String `tfsdk:"content"`
}

func (d *RailsCredentialsInlineDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_inline"
}

func (d *RailsCredentialsInlineDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
				Computed:            true,
			},
			"content": schema.StringAttribute{
				MarkdownDescription: "Raw credentials in YAML format",
				Required:            true,
				Sensitive:           true,
			},
		},
	}
}

func (d *RailsCredentialsInlineDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
}

func (d *RailsCredentialsInlineDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RailsCredentialsInlineDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rawObject, err := credentials.MarshalSingleString(data.DecryptedContent.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Credentials marshal failed", err.Error())
		return
	}
	encryptedString, err := credentials.Encrypt(credentials.SanitizeMasterKey(data.MasterKey.ValueString()), rawObject)
	if err != nil {
		resp.Diagnostics.AddError("Credentials encryption failed", err.Error())
		return
	}
	data.EncryptedContent = types.StringValue(encryptedString)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
