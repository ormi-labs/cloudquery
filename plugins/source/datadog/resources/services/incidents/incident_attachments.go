// Code generated by codegen; DO NOT EDIT.

package incidents

import (
	"context"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
	"github.com/cloudquery/cloudquery/plugins/source/datadog/client"
	"github.com/cloudquery/plugin-sdk/v2/schema"
)

func IncidentAttachments() *schema.Table {
	return &schema.Table{
		Name:     "datadog_incident_attachments",
		Resolver: fetchIncidentAttachments,
		Columns: []schema.Column{
			{
				Name:     "account_name",
				Type:     schema.TypeString,
				Resolver: client.ResolveAccountName,
			},
			{
				Name:     "attributes",
				Type:     schema.TypeJSON,
				Resolver: schema.PathResolver("Attributes"),
			},
			{
				Name:     "id",
				Type:     schema.TypeString,
				Resolver: schema.PathResolver("Id"),
			},
			{
				Name:     "relationships",
				Type:     schema.TypeJSON,
				Resolver: schema.PathResolver("Relationships"),
			},
			{
				Name:     "type",
				Type:     schema.TypeString,
				Resolver: schema.PathResolver("Type"),
			},
		},
	}
}


func fetchIncidentAttachments(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
	p := parent.Item.(datadogV2.IncidentResponseData)
	c := meta.(*client.Client)
	ctx = c.BuildContextV2(ctx)
	resp, _, err := c.DDServices.IncidentsAPI.ListIncidentAttachments(ctx, p.Id)
	if err != nil {
		return err
	}
	res <- resp.GetData()
	return nil
}

