package core

// Settings ...
type Settings struct {
	Port    int32          `yaml:"port"`
	JWT     *JWT           `yaml:"jwt"`
	SQS     *SessionConfig `yaml:"sqs"`
	SNS     *SessionConfig `yaml:"sns"`
	MongoDB *MongoDBConfig `yaml:"mongodb"`
	Events  *Events        `yaml:"events"`
}

// JWT ...
type JWT struct {
	Secret string `yaml:"secret"`
}

// SessionConfig ...
type SessionConfig struct {
	Region   string `yaml:"region"`
	Endpoint string `yaml:"endpoint"`
}

// MongoDBConfig ...
type MongoDBConfig struct {
	Database         string `yaml:"database"`
	ConnectionString string `yaml:"connection_string"`
}

// Events ...
type Events struct {
	ItemsCreated       string `yaml:"items-created"`
	ItemsUpdated       string `yaml:"items-updated"`
	ItemsLockRequested string `yaml:"items-lock-requested"`
	ItemsLockCompleted string `yaml:"items-lock-completed"`
}
