<!-- markdownlint-disable MD033 -->
<h1 align="center"><img alt="cloudquery logo" width=75% src="https://github.com/cloudquery/cloudquery/raw/main/cli/docs/images/logo.png"/></h1>
<!-- markdownlint-enable MD033 -->

CloudQuery is an [open-source](https://github.com/cloudquery/cloudquery),
high-performance data integration framework built for developers.

CloudQuery extracts, transforms, and loads configuration from cloud APIs to
variety of supported destinations such as databases, data lakes, or streaming platforms
for further analysis.

### Ormi Labs Changes

To run it with changes made to `source/postgresql` and `destication/bigquery` build selected plugin with

```shell
./scripts/build.sh
```

Create yaml config files for source and destination respectively pointing them to the local bin directory.
Examples:

```yaml
kind: source
spec:
  name: "postgresql"
  # registry: "github"
  # path: "cloudquery/postgresql"
  registry: "local"
  path: "../bin/postgresql"
  version: "v3.0.0"
  tables: ["block","transaction","transaction_log","trace"]
  destinations: ["bigquery"]
  spec:
    # Provide connection string. For ex.: postgres://user:pass@host:port/db?sslmode=verify-ca&pool_max_conns=10
    # Available options:
    # - pool_max_conns: integer greater than 0
    # - pool_min_conns: integer 0 or greater
    # - pool_max_conn_lifetime: duration string
    # - pool_max_conn_idle_time: duration string
    # - pool_health_check_period: duration string
    # - pool_max_conn_lifetime_jitter: duration string
    connection_string: "postgresql://user:pass@host:port/db?sslmode=disable"
    # Optional parameters:
    # cdc_id: "postgresql" # Set to a unique string per source to enable Change Data Capture mode (logical replication, or CDC)
    pgx_log_level: error
    rows_per_record: 1
    # "block" or "table", default is table.
    # If "block" is set, 'block' section is read in.
    sync_mode: "block"
    report_dir: "."
    report_fmt: "csv" # "csv" or "json", default is csv
    block:
      start: 0 # block number from start, 0 points to the very first block
      limit: 0 # limit number of blocks, default 0 which is no limitation at all
      # Index of the table in 'tables' list starting from 0.
      # This is necessary because actual tables/views may be named differently.
      table_idx: 0
      # Other entities such as: transaction,log,token_transfer,trace,contract,token
      # see https://ethereum-etl.readthedocs.io/en/latest/commands/
      entities:
        - name: "transaction"
          table_idx: 1
          enable: true
        - name: "log"
          table_idx: 2
          enable: true
        - name: "trace"
          table_idx: 3
          enable: true
        - name: "token_transfer"
          # table_idx:
          enable: false
```

```yaml
kind: destination
spec:
  name: bigquery
  registry: "local"
  path: "../bin/bigquery"
  version: "v3.3.0"
  write_mode: "append"
  spec:
    project_id: temporal-wares-368022
    dataset_id: ethblocks_svr6_dev
```

Run from **_cloudquery/cli_** directory the following command:

```shell
go run main.go sync /path/to/dir/where/source/and/dest/yaml/conf/files/stored
```

### Installation

See the **[Quickstart guide](https://www.cloudquery.io/docs/quickstart)** for instructions how to start syncing data with CloudQuery.

## Why CloudQuery?

- **Open source**: Extensible plugin architecture. Contribute to our official plugins or develop your own with CloudQuery SDK.
- **Blazing fast**: CloudQuery is optimized for performance, utilizing the excellent Go concurrency model with light-weight goroutines.
- **Deploy anywhere**: CloudQuery plugins are single-binary executables and can be deployed and run anywhere.
- **Pre-built queries**: CloudQuery maintains a number of out-of-the-box security and compliance policies for cloud infrastructure.
- **Eliminate data silos**: Eliminate data silos across your organization, unifying data between security, infrastructure, marketing and finance teams.
- **Unlimited scale**: CloudQuery plugins are stateless and can be scaled horizontally on any platform, such as EC2, Kubernetes, batch jobs or any other compute.

## Use Cases

- **Cloud Security Posture Management**: Use as an open source CSPM solution to monitor and enforce security policies across your cloud infrastructure for AWS, GCP, Azure and many more.
- **Cloud Asset Inventory**: First-class support for major cloud infrastructure providers such as AWS, GCP and Azure allow you to collect and unify configuration data.
- **Cloud FinOps**: Collect and unify billing data from cloud providers to drive financial accountability.
- **ELT Platform**: With hundreds of plugin combinations and extensible architecture, CloudQuery can be used for reliable, efficient export from any API to any database, or from one database to another.
- **Attack Surface Management**: [Open source solution](https://www.cloudquery.io/how-to-guides/attack-surface-management-with-graph) for continuous discovery, analysis and monitoring of potential attack vectors that make up your organization's attack surface.

### Links

- Homepage: https://www.cloudquery.io
- Documentation: https://www.cloudquery.io/docs
- Plugins: https://www.cloudquery.io/docs/plugins/overview
- How-to Guides: https://www.cloudquery.io/how-to-guides
- Integrations: https://www.cloudquery.io/integrations
- Releases: https://github.com/cloudquery/cloudquery/releases
- Plugin SDK: https://github.com/cloudquery/plugin-sdk

## License

By contributing to CloudQuery you agree that your contributions will be licensed as defined on the LICENSE file.

## Hiring

If you are into Go, Backend, Cloud, GCP, AWS - ping us at jobs [at] our domain

## Contribution

Feel free to open a pull request for small fixes and changes. For bigger changes and new plugins, please open an issue first to prevent duplicated work and to have the relevant discussions first.
