package config

import (
	"fmt"
	"time"

	"github.com/sirjager/goreact/pkg/db"
	"github.com/spf13/viper"
)

type Config struct {
	LogErrors bool `mapstructure:"LOG_ERRORS"` //? for logging errors in console

	//? for internal
	StartTime   time.Time // StartTime is the timestamp when the application started.
	ServiceName string    // ServiceName is the name of the service.

	GrpcPort             string        `mapstructure:"GRPC_PORT"`              // GrpcPort is the port number for the gRPC server.
	GatewayPort          string        `mapstructure:"GATEWAY_PORT"`           // RestPort is the port number for the REST server.
	//? for pkg
	Database db.Config   // DBConfig holds the configuration for the database.

}

func LoadConfigs(path, name string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}
	if err = viper.Unmarshal(&config); err != nil {
		return
	}

	if err = viper.Unmarshal(&config.Database); err != nil {
		return
	}

	// Construct the DBUrl using the DBConfig values.
	config.Database.Url = fmt.Sprintf("%s://%s:%s@%s:%s/%s%s", config.Database.Driver, config.Database.User, config.Database.Pass, config.Database.Host, config.Database.Port, config.Database.Name, config.Database.Args)
	config.Database.Migrate = "file://" + config.Database.Migrate
	return
}
