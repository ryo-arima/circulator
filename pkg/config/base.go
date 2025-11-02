package config

import (
	"os"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type BaseConfig struct {
	DBConnection *gorm.DB
	YamlConfig   YamlConfig
	Logger       LoggerInterface // Dependency injection for logger
}

type YamlConfig struct {
	Application Application  `yaml:"Application"`
	MySQL       MySQL        `yaml:"MySQL"`
	Pulsar      Pulsar       `yaml:"Pulsar"`
	Logger      LoggerConfig `yaml:"Logger"`
	// Redis は任意
	// Redis      Redis       `yaml:"Redis"`
}

type Application struct {
	Common Common `yaml:"Common"`
	Server Server `yaml:"Server"`
	Client Client `yaml:"Client"`
	Agent  Agent  `yaml:"Agent"`
}

type Common struct {
	Port string `yaml:"port"`
}

type Server struct {
	Base      Base   `yaml:"base"`
	JWTSecret string `yaml:"jwt_secret"`
}

type Base struct {
	Emails []string `yaml:"emails"`
}

type Client struct {
	ServerEndpoint string `yaml:"ServerEndpoint"`
	UserEmail      string `yaml:"UserEmail"`
}

type Agent struct {
	ServerEndpoint            string `yaml:"ServerEndpoint"`
	LoginEmail                string `yaml:"LoginEmail"`
	LoginPassword             string `yaml:"LoginPassword"`
	TokenCachePath            string `yaml:"TokenCachePath"`
	RefreshIntervalMinutes    int    `yaml:"RefreshIntervalMinutes"`
	RegistrationRetryInterval int    `yaml:"RegistrationRetryInterval"`
	HealthCheckInterval       int    `yaml:"HealthCheckInterval"`
}

type MySQL struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type Pulsar struct {
	URL               string         `yaml:"url"`
	ConnectionTimeout int            `yaml:"connection_timeout"` // seconds
	OperationTimeout  int            `yaml:"operation_timeout"`  // seconds
	Topics            PulsarTopics   `yaml:"topics"`
	Consumer          PulsarConsumer `yaml:"consumer"`
	Producer          PulsarProducer `yaml:"producer"`
}

type PulsarTopics struct {
	ExternalSensorData  string `yaml:"external_sensor_data"`
	ProcessedSensorData string `yaml:"processed_sensor_data"`
	SystemMetrics       string `yaml:"system_metrics"`
	AlertData           string `yaml:"alert_data"`
	ProcessingResults   string `yaml:"processing_results"`
}

type PulsarConsumer struct {
	SubscriptionName string `yaml:"subscription_name"`
	Type             string `yaml:"type"` // Shared, Exclusive, Failover, KeyShared
}

type PulsarProducer struct {
	SendTimeout int `yaml:"send_timeout"` // seconds
}

// type Redis struct {
//   Host string `yaml:"host"`
//   Port int    `yaml:"port"`
//   User string `yaml:"user"`
//   Pass string `yaml:"pass"`
//   DB   int    `yaml:"db"`
// }

func NewBaseConfig() BaseConfig {
	yamlConfig := loadYamlConfig()
	db := NewDBConnection(yamlConfig)

	baseConfig := BaseConfig{
		DBConnection: db,
		YamlConfig:   yamlConfig,
	}

	// Initialize logger with dependency injection
	logger := NewLogger(yamlConfig.Logger, &baseConfig)
	baseConfig.Logger = logger

	return baseConfig
}

func NewClientConfig() BaseConfig {
	yamlConfig := loadYamlConfig()

	baseConfig := BaseConfig{
		DBConnection: nil, // Client doesn't need DB connection
		YamlConfig:   yamlConfig,
	}

	// Initialize logger with dependency injection
	logger := NewLogger(yamlConfig.Logger, &baseConfig)
	baseConfig.Logger = logger

	return baseConfig
}

func loadYamlConfig() YamlConfig {
	var config YamlConfig

	data, err := os.ReadFile("etc/app.yaml")
	if err != nil {
		// Default configuration
		return YamlConfig{
			Application: Application{
				Common: Common{Port: "8080"},
				Server: Server{
					JWTSecret: "your-secret-key",
					Base:      Base{Emails: []string{"base@example.com"}},
				},
				Client: Client{
					ServerEndpoint: "http://localhost:8080",
				},
			},
			MySQL: MySQL{
				Host:     "localhost",
				Port:     "3306",
				User:     "root",
				Password: "password",
				Database: "circulator",
			},
			Pulsar: Pulsar{
				URL:               "pulsar://localhost:6650",
				ConnectionTimeout: 10,
				OperationTimeout:  5,
				Topics: PulsarTopics{
					ExternalSensorData:  "external-sensor-data",
					ProcessedSensorData: "processed-sensor-data",
					SystemMetrics:       "system-metrics",
					AlertData:           "alert-data",
					ProcessingResults:   "processing-results",
				},
				Consumer: PulsarConsumer{
					SubscriptionName: "agent-processor",
					Type:             "Shared",
				},
				Producer: PulsarProducer{
					SendTimeout: 30,
				},
			},
			Logger: LoggerConfig{
				Component:    "unknown",
				Service:      "unknown",
				Level:        "INFO",
				Structured:   true,
				EnableCaller: true,
				Output:       "stdout",
			},
		}
	}

	yaml.Unmarshal(data, &config)
	return config
}

func NewDBConnection(conf YamlConfig) *gorm.DB {
	dsn := conf.MySQL.User + ":" + conf.MySQL.Password + "@tcp(" +
		conf.MySQL.Host + ":" + conf.MySQL.Port + ")/" +
		conf.MySQL.Database + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		// For client-only usage, return nil instead of panic
		return nil
	}

	return db
}

// NewServerConfig creates a BaseConfig instance for server component
func NewServerConfig() BaseConfig {
	yamlConfig := loadYamlConfig()
	// Override logger component for server
	yamlConfig.Logger.Component = "server"
	if yamlConfig.Logger.Service == "unknown" {
		yamlConfig.Logger.Service = "circulator-server"
	}

	db := NewDBConnection(yamlConfig)

	baseConfig := BaseConfig{
		DBConnection: db,
		YamlConfig:   yamlConfig,
	}

	// Initialize logger with dependency injection
	logger := NewLogger(yamlConfig.Logger, &baseConfig)
	baseConfig.Logger = logger

	return baseConfig
}

// NewAgentConfig creates a BaseConfig instance for agent component
func NewAgentConfig() BaseConfig {
	yamlConfig := loadYamlConfig()
	// Override logger component for agent
	yamlConfig.Logger.Component = "agent"
	if yamlConfig.Logger.Service == "unknown" {
		yamlConfig.Logger.Service = "circulator-agent"
	}

	baseConfig := BaseConfig{
		DBConnection: nil, // Agent may not need direct DB connection
		YamlConfig:   yamlConfig,
	}

	// Initialize logger with dependency injection
	logger := NewLogger(yamlConfig.Logger, &baseConfig)
	baseConfig.Logger = logger

	return baseConfig
}

// NewClientConfigWithComponent creates a BaseConfig instance for client component
func NewClientConfigWithComponent(service string) BaseConfig {
	yamlConfig := loadYamlConfig()
	// Override logger component for client
	yamlConfig.Logger.Component = "client"
	yamlConfig.Logger.Service = service

	baseConfig := BaseConfig{
		DBConnection: nil, // Client doesn't need DB connection
		YamlConfig:   yamlConfig,
	}

	// Initialize logger with dependency injection
	logger := NewLogger(yamlConfig.Logger, &baseConfig)
	baseConfig.Logger = logger

	return baseConfig
}

// IncomingAgentData represents incoming data for agent processing
type IncomingAgentData struct {
	UUID       string  `json:"uuid"`
	Source     string  `json:"source"`
	SensorType string  `json:"sensor_type"`
	Value      float64 `json:"value"`
	RawPayload []byte  `json:"raw_payload"`
}

// ProcessedAgentData represents processed agent data
type ProcessedAgentData struct {
	AgentUUID      string  `json:"agent_uuid"`
	OriginalValue  float64 `json:"original_value"`
	ProcessedValue float64 `json:"processed_value"`
	Anomaly        bool    `json:"anomaly"`
	Confidence     float64 `json:"confidence"`
	ProcessingTime int64   `json:"processing_time"` // microseconds
}

// ProcessingRule represents a processing rule configuration
type ProcessingRule struct {
	Name    string                 `json:"name"`
	Enabled bool                   `json:"enabled"`
	Params  map[string]interface{} `json:"params"`
}

// NewPulsarClient creates a new Pulsar client with configuration
func NewPulsarClient(conf YamlConfig) (pulsar.Client, error) {
	return pulsar.NewClient(pulsar.ClientOptions{
		URL:               conf.Pulsar.URL,
		ConnectionTimeout: time.Duration(conf.Pulsar.ConnectionTimeout) * time.Second,
		OperationTimeout:  time.Duration(conf.Pulsar.OperationTimeout) * time.Second,
	})
}

// GetPulsarConsumerOptions returns consumer options based on configuration
func (conf YamlConfig) GetPulsarConsumerOptions(topic string) pulsar.ConsumerOptions {
	var consumerType pulsar.SubscriptionType
	switch conf.Pulsar.Consumer.Type {
	case "Exclusive":
		consumerType = pulsar.Exclusive
	case "Failover":
		consumerType = pulsar.Failover
	case "KeyShared":
		consumerType = pulsar.KeyShared
	default:
		consumerType = pulsar.Shared
	}

	return pulsar.ConsumerOptions{
		Topic:            topic,
		SubscriptionName: conf.Pulsar.Consumer.SubscriptionName,
		Type:             consumerType,
	}
}

// GetPulsarProducerOptions returns producer options based on configuration
func (conf YamlConfig) GetPulsarProducerOptions(topic string) pulsar.ProducerOptions {
	return pulsar.ProducerOptions{
		Topic:       topic,
		SendTimeout: time.Duration(conf.Pulsar.Producer.SendTimeout) * time.Second,
	}
}

// GetTopicName returns the topic name for a given topic type
func (p *Pulsar) GetTopicName(topicType string) string {
	switch topicType {
	case "external_sensor_data":
		return p.Topics.ExternalSensorData
	case "processed_sensor_data":
		return p.Topics.ProcessedSensorData
	case "system_metrics":
		return p.Topics.SystemMetrics
	case "alert_data":
		return p.Topics.AlertData
	case "processing_results":
		return p.Topics.ProcessingResults
	default:
		return ""
	}
}
