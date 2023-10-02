# 0xdb source Plugin

The CloudQuery 0xdb source plugin replicates blockchain data to any supported CloudQuery destination.

## Links

- [User Guide](https://docs.cloudquery.io/docs/plugins/sources/0xdb/overview)

Use `bq` to extract schema for `syncs` table and load reporting files to it. For ex.: 

```shell
bq show --schema --format=prettyjson temporal-wares-368022:ethblocks_svr6_dev.syncs > sync_schema.json

bq load --skip_leading_rows=1 --source_format=CSV ethblocks_svr6_dev.syncs \
  sync_Oct-2-17:16:35.csv \
  sync_schema.json
```
