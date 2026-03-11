package config

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type AppConfig struct {
	Server     Server     `yaml:"server"`
	API        API        `yaml:"api"`
	Kafka      Kafka      `yaml:"kafka"`
	Prometheus Prometheus `yaml:"prometheus"`
	Repository Repository `yaml:"repository"`
	Tracing    Tracing    `yaml:"tracing"`
}

type Tracing struct {
	Endpoint    string `yaml:"endpoint" env:"TRACING_EXPORTER_OTLP_ENDPOINT" env_default:"http://localhost:4318"`
	ServiceName string `yaml:"service_name" env:"TRACING_SERVICE_NAME" env_default:"hotels"`
	Enabled     bool   `yaml:"enabled" env:"TRACING_ENABLED" env_default:"true"`
}

type Server struct {
	Address string `yaml:"address" env:"SERVER_ADDRESS" env_default:":8080"`
}

type API struct {
	Timeout           time.Duration `yaml:"timeout" env:"API_TIMEOUT" env_default:"30s"`
	ReadHeaderTimeout time.Duration `yaml:"read_header_timeout" env:"API_READ_HEADER_TIMEOUT" env_default:"5s"`
}

type Kafka struct {
	URL      string   `yaml:"url" env:"KAFKA_URL" env_default:"localhost:9092"`
	Brokers  []string `yaml:"brokers" env:"KAFKA_BROKERS" env_default:"localhost:9092"`
	Version  string   `yaml:"version" env:"KAFKA_VERSION" env_default:""`
	Topics   []string `yaml:"topics" env:"KAFKA_TOPICS"`
	GroupID  string   `yaml:"group_id" env:"KAFKA_GROUP_ID" env_default:"default-consumer-group"`
	Producer Producer `yaml:"producer"`
	Consumer Consumer `yaml:"consumer"`
}

type Producer struct {
	AckPolicy int16 `yaml:"ack_policy" env:"KAFKA_ACK_POLICY" env_default:"1"`
	RetryMax  int   `yaml:"retry_max" env:"KAFKA_RETRY_MAX" env_default:"5"`
}

type Consumer struct {
	AutoCommitEnable bool          `yaml:"auto_commit_enable" env:"KAFKA_AUTO_COMMIT" env_default:"false"`
	MaxWaitTime      time.Duration `yaml:"max_wait_time" env:"KAFKA_MAX_WAIT_TIME" env_default:"500ms"`
}

type Prometheus struct {
	Host string `yaml:"host" env:"PROM_HOST" env_default:"localhost:2112"`
}

type Repository struct {
	Hosts    string `yaml:"hosts" env:"DB_HOSTS" env_default:"localhost:3306"`
	Dbname   string `yaml:"dbname" env:"DB_NAME" env_default:"orders_db"`
	User     string `yaml:"user" env:"DB_USER" env_default:"root"`
	Password string `yaml:"password" env:"DB_PASSWORD" env_default:""`
	Schema   string `yaml:"schema" env:"DB_SCHEMA"`

	MaxIdleConn int `yaml:"max_idle_conn" env:"DB_MAX_IDLE_CONN" env_default:"5"`
	MaxOpenConn int `yaml:"max_open_conn" env:"DB_MAX_OPEN_CONN" env_default:"10"`

	MaxIdleLifetime time.Duration `yaml:"max_idle_lifetime" env:"DB_MAX_IDLE_LIFETIME" env_default:"5m"`
	MaxOpenLifetime time.Duration `yaml:"max_open_lifetime" env:"DB_MAX_OPEN_LIFETIME" env_default:"30m"`
}

func Load(path string) (*AppConfig, error) {
	var cfg AppConfig
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, fmt.Errorf("reading config %s: %w", path, err)
	}
	return &cfg, nil
}

func MustLoad(path string) *AppConfig {
	cfg, err := Load(path)
	if err != nil {
		log.Fatalf("loading config: %v", err)
	}
	return cfg
}

func (r *Repository) DSN() string {
	return r.DSNMySQL()
}

// DSNMySQL возвращает DSN для подключения к MySQL.
// Hosts: "host:port" или "host:3306", один или несколько через запятую (берётся первый).
func (r *Repository) DSNMySQL() string {
	host := "localhost:3306"
	if r.Hosts != "" {
		parts := strings.Split(r.Hosts, ",")
		host = strings.TrimSpace(parts[0])
	}
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&charset=utf8mb4&multiStatements=true",
		r.User, r.Password, host, r.Dbname)
}
