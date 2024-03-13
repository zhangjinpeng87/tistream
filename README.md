# TiStream

A cloud-native data change stream service for TiDB Cloud. See [design](./docs/design/20240220-tistream.md) for more details.

![image](docs/resources/tistream-arch.png)

## Functionality

- [x] Dispatchers read tenants' raft log/watermarks changes files(S3) and dispatcher to corresponding remote sorter.
- [x] Remote Sorters sort committed transaction changes and flush into commited data pool(S3).
- [x] Remote Sorters support persistent unflushed data to S3.
- [x] Meta Server manage and store tenants' data changes stream tasks into backend RDS.
- [x] Meta Server campaign leader.
- [x] Schema Registry support MVCC schema storage.
- [x] API sever supports APIs to pull committed data in order.
- [ ] Upstream TiDB cluster push full snapshot data.
- [ ] Balance tasks dynamically across dispatchers and sorters according to workload.
- [ ] Dispatchers and sorters failover.
- [ ] Tenant level data througput metering.
- [ ] Security: tenant level CMEK.
- [ ] Integrate with CDC.

## Code Structure

- docs/design: feature or architecture design docs.
- pkg/codec: encoding and decoding of different type of data files
- pkg/dispatcher: remote data dispatcher workers, dispatch upstream changes to corresponding sorters.
- pkg/remotesorter: remote sorter workers, organize upstream data changes in commit-ts order and flush them as incremental data change files.
- pkg/schemaregistry: upstream clusters/tenants' table schema registry, contains all upstream table schema history for last n days where n is configurable.
- pkg/metaserver: meta server is responsible for cluster health monitoring, tasks persistent, workload balancing, etc.
- pkg/apiserver: providing api for downstream system to query and subscribe database or table's data changes stream.
- pkg: and other reusable modules.
- cmd: target binary files.
- proto: protobuf definition files
