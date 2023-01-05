// Code generated by codegen; DO NOT EDIT.

package plugin

import (
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugins/source/gcp/resources/services/aiplatform"
	"github.com/cloudquery/plugins/source/gcp/resources/services/apigateway"
	"github.com/cloudquery/plugins/source/gcp/resources/services/apikeys"
	"github.com/cloudquery/plugins/source/gcp/resources/services/appengine"
	"github.com/cloudquery/plugins/source/gcp/resources/services/artifactregistry"
	"github.com/cloudquery/plugins/source/gcp/resources/services/baremetalsolution"
	"github.com/cloudquery/plugins/source/gcp/resources/services/batch"
	"github.com/cloudquery/plugins/source/gcp/resources/services/beyondcorp"
	"github.com/cloudquery/plugins/source/gcp/resources/services/bigquery"
	"github.com/cloudquery/plugins/source/gcp/resources/services/bigtableadmin"
	"github.com/cloudquery/plugins/source/gcp/resources/services/billing"
	"github.com/cloudquery/plugins/source/gcp/resources/services/binaryauthorization"
	"github.com/cloudquery/plugins/source/gcp/resources/services/certificatemanager"
	"github.com/cloudquery/plugins/source/gcp/resources/services/clouddeploy"
	"github.com/cloudquery/plugins/source/gcp/resources/services/clouderrorreporting"
	"github.com/cloudquery/plugins/source/gcp/resources/services/cloudiot"
	"github.com/cloudquery/plugins/source/gcp/resources/services/compute"
	"github.com/cloudquery/plugins/source/gcp/resources/services/container"
	"github.com/cloudquery/plugins/source/gcp/resources/services/containeranalysis"
	"github.com/cloudquery/plugins/source/gcp/resources/services/dns"
	"github.com/cloudquery/plugins/source/gcp/resources/services/domains"
	"github.com/cloudquery/plugins/source/gcp/resources/services/functions"
	"github.com/cloudquery/plugins/source/gcp/resources/services/iam"
	"github.com/cloudquery/plugins/source/gcp/resources/services/kms"
	"github.com/cloudquery/plugins/source/gcp/resources/services/logging"
	"github.com/cloudquery/plugins/source/gcp/resources/services/monitoring"
	"github.com/cloudquery/plugins/source/gcp/resources/services/redis"
	"github.com/cloudquery/plugins/source/gcp/resources/services/resourcemanager"
	"github.com/cloudquery/plugins/source/gcp/resources/services/run"
	"github.com/cloudquery/plugins/source/gcp/resources/services/secretmanager"
	"github.com/cloudquery/plugins/source/gcp/resources/services/serviceusage"
	"github.com/cloudquery/plugins/source/gcp/resources/services/sql"
	"github.com/cloudquery/plugins/source/gcp/resources/services/storage"
)

func PluginAutoGeneratedTables() []*schema.Table {
	return []*schema.Table{
		aiplatform.JobLocations(),
		aiplatform.DatasetLocations(),
		aiplatform.EndpointLocations(),
		aiplatform.FeaturestoreLocations(),
		aiplatform.IndexendpointLocations(),
		aiplatform.IndexLocations(),
		aiplatform.MetadataLocations(),
		aiplatform.ModelLocations(),
		aiplatform.PipelineLocations(),
		aiplatform.SpecialistpoolLocations(),
		aiplatform.VizierLocations(),
		aiplatform.TensorboardLocations(),
		aiplatform.Operations(),
		apigateway.Apis(),
		apigateway.Gateways(),
		apikeys.Keys(),
		appengine.Apps(),
		appengine.Services(),
		appengine.AuthorizedCertificates(),
		appengine.AuthorizedDomains(),
		appengine.DomainMappings(),
		appengine.FirewallIngressRules(),
		artifactregistry.Locations(),
		baremetalsolution.Instances(),
		baremetalsolution.Networks(),
		baremetalsolution.NfsShares(),
		baremetalsolution.Volumes(),
		batch.Jobs(),
		batch.TaskGroups(),
		beyondcorp.AppConnections(),
		beyondcorp.AppConnectors(),
		beyondcorp.AppGateways(),
		beyondcorp.ClientConnectorServices(),
		beyondcorp.ClientGateways(),
		bigquery.Datasets(),
		bigtableadmin.Instances(),
		billing.BillingAccounts(),
		billing.Services(),
		binaryauthorization.Assertors(),
		certificatemanager.CertificateIssuanceConfigs(),
		certificatemanager.CertificateMaps(),
		certificatemanager.Certificates(),
		certificatemanager.DnsAuthorizations(),
		clouddeploy.Targets(),
		clouddeploy.DeliveryPipelines(),
		clouderrorreporting.ErrorGroupStats(),
		cloudiot.DeviceRegistries(),
		compute.Addresses(),
		compute.Autoscalers(),
		compute.BackendServices(),
		compute.DiskTypes(),
		compute.Disks(),
		compute.ForwardingRules(),
		compute.Instances(),
		compute.SslCertificates(),
		compute.Subnetworks(),
		compute.TargetHttpProxies(),
		compute.UrlMaps(),
		compute.VpnGateways(),
		compute.InstanceGroups(),
		compute.Images(),
		compute.Firewalls(),
		compute.Networks(),
		compute.SslPolicies(),
		compute.Interconnects(),
		compute.TargetSslProxies(),
		compute.Projects(),
		container.Clusters(),
		containeranalysis.Occurrences(),
		dns.Policies(),
		dns.ManagedZones(),
		domains.Registrations(),
		functions.Functions(),
		iam.Roles(),
		iam.ServiceAccounts(),
		iam.DenyPolicies(),
		kms.Keyrings(),
		logging.Metrics(),
		logging.Sinks(),
		monitoring.AlertPolicies(),
		redis.Instances(),
		resourcemanager.Folders(),
		resourcemanager.Projects(),
		resourcemanager.ProjectPolicies(),
		run.Locations(),
		secretmanager.Secrets(),
		serviceusage.Services(),
		sql.Instances(),
		storage.Buckets(),
	}
}
