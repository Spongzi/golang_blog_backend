package settings

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func Init() (err error) {
	viper.SetConfigFile("config.y")
	//viper.SetConfigName("config")
	//viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("viper.ReadInConfig() failed, err: %s \n", err)
		return
	}
	viper.WatchConfig()
	if err != nil {
		fmt.Printf("viper.WatchRemoteConfig() failed, err: %s \n", err)
		return
	}
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了...")
	})
	return
}
