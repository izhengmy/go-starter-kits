package viperx

import (
	"fmt"

	"github.com/spf13/viper"
)

func Scan(path string, config any) error {
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %v \n", err)
	}

	if err := viper.Unmarshal(config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %v \n", err)
	}

	return nil
}
