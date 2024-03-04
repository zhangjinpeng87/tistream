package utils

// GlobalConfig is the global configuration of the system.
// Please refer to etc/config.toml for the details of each field.
type GlobalConfig struct {
	LogLevel   string
	Dispatcher DispatcherConfig
	Sorter     SorterConfig
	MetaServer MetaServerConfig
	ApiServer  ApiServerConfig
}

type StorageConfig struct {
	Region      string
	Bucket      string
	Prefix      string
	SecurityKey string
	AccessKey   string
}

type DispatcherConfig struct {
	Addr string
	Port int

	MetaServerEndpoints []string

	// If the dispatcher is isolated from the meta server for a long time,
	// it will be paused.
	SuicideDur int

	// storage configuration
	Storage StorageConfig

	// How often to check the store changes, in seconds.
	// Store changes caused by upstream scale out or scale in.
	// Default to 10 seconds.
	CheckStoreInterval int
	// How often to check the new file changes, in seconds.
	// Default to 5 seconds.
	CheckFileInterval int
}

type SorterConfig struct {
	Addr    string
	Port    int
	Storage StorageConfig
}

type MetaServerConfig struct {
	Addr          string
	Port          int
	MysqlHost     string
	MysqlPort     int
	MysqlUser     string
	MysqlPassword string

	// The interval to update the lease of the master role, in seconds.
	// Default to 5 seconds.
	UpdateLeaseInterval int
	// Lease duration, in seconds, default to 10 seconds.
	LeaseDuration int
}

type ApiServerConfig struct {
	Addr string
	Port int
}

type SchemaRegistryConfig struct {
	Addr string
	Port int
}
