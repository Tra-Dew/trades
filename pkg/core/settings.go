package core

// Settings ...
type Settings struct {
	Port             int32          `yaml:"port"`
	JWT              *JWT           `yaml:"jwt"`
	MongoDB          *MongoDBConfig `yaml:"mongodb"`
	InventoryService *GRPCService   `yaml:"inventory_service"`
}

// JWT ...
type JWT struct {
	Secret string `yaml:"secret"`
}

// MongoDBConfig ...
type MongoDBConfig struct {
	Database         string `yaml:"database"`
	ConnectionString string `yaml:"connection_string"`
}

// GRPCService ...
type GRPCService struct {
	URL string `yaml:"url"`
}
