package utils

// GlobalConfig is the global configuration of the system.
// Please refer to etc/config.toml for the details of each field.
type GlobalConfig struct {
	LogLevel        string
	ExternalStorage StorageConfig
	Dispatcher      DispatcherConfig
	Sorter          SorterConfig
	MetaServer      MetaServerConfig
	ApiServer       ApiServerConfig
}

type StorageConfig struct {
	Url         string
	SecurityKey string
	AccessKey   string
}

type DispatcherConfig struct {
	Addr string
	Port int

	// Prefix of the data change buffer.
	Prefix string

	// How often to check the store changes, in seconds.
	// Store changes caused by upstream scale out or scale in.
	// Default to 10 seconds.
	CheckStoreInterval int
	// How often to check the new file changes, in seconds.
	// Default to 5 seconds.
	CheckFileInterval int
}

type SorterConfig struct {
	Addr string
	Port int
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
