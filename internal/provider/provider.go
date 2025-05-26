package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure RailsCredentialProvider satisfies various provider interfaces.
var _ provider.Provider = &RailsCredentialProvider{}
var _ provider.ProviderWithFunctions = &RailsCredentialProvider{}
var _ provider.ProviderWithEphemeralResources = &RailsCredentialProvider{}

// RailsCredentialProvider defines the provider implementation.
type RailsCredentialProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// RailsCredentialProviderModel describes the provider data model.
type RailsCredentialProviderModel struct{}

func (p *RailsCredentialProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "railscred"
	resp.Version = p.version
}

func (p *RailsCredentialProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage Rails credentials files.",
	}
}

func (p *RailsCredentialProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data RailsCredentialProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (p *RailsCredentialProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewRailsMasterKeyResource,
	}
}

func (p *RailsCredentialProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (p *RailsCredentialProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewRailsCredentialsFileDataSource,
		NewRailsCredentialsInlineDataSource,
	}
}

func (p *RailsCredentialProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &RailsCredentialProvider{
			version: version,
		}
	}
}
