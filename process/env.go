package process

import (
	"path/filepath"

	"github.com/Maycon-Santos/test-brand-monitor-backend/internal/projectpath"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Env struct {
	ServerPort                int    `mapstructure:"server_port"`
	MongoDbConnUri            string `mapstructure:"mongo_db_conn_uri"`
	AccessControlAllowOrigin  string `mapstructure:"access_control_allow_origin"`
	AccessControlAllowHeaders string `mapstructure:"access_control_allow_headers"`
}

func NewEnv() (*Env, error) {
	env := &Env{}

	defaultEnvs, err := godotenv.Read(filepath.Join(projectpath.Root, ".env.sample"))
	if err != nil {
		return nil, err
	}

	for key, value := range defaultEnvs {
		viper.SetDefault(key, value)
	}

	viper.SetConfigType("env")
	viper.SetConfigFile(filepath.Join(projectpath.Root, ".env"))
	viper.AutomaticEnv()

	_ = viper.ReadInConfig()

	if err := viper.Unmarshal(env); err != nil {
		return nil, err
	}

	return env, nil
}
