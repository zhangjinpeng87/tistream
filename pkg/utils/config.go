package utils

// GlobalConfig is the global configuration of the system.
// Please refer to etc/config.toml for the details of each field.
type GlobalConfig struct {
	LogLevel string
	ExternalStorage StorageConfig
	Dispatcher DispatcherConfig
	Sorter SorterConfig
	MetaServer MetaServerConfig
	ApiServer ApiServerConfig
}

type StorageConfig struct {
	Url string
	SecurityKey string
	AccessKey string
}

type DispatcherConfig struct {
	Addr string
	Port int
}

type SorterConfig struct {
	Addr string
	Port int
}

type MetaServerConfig struct {
	Addr string
	Port int
	MysqlHost string
	MysqlPort int
	MysqlUser string
	MysqlPassword string
}

type ApiServerConfig struct {
	Addr string
	Port int
}

type SchemaRegistryConfig struct {
	Addr string
	Port int
}
