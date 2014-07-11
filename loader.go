package settings

import (
	"encoding/json"
	//"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

/*
func main() {
	jsonConfig := GetConfig()
	redis := jsonConfig.Redis
	fmt.Println(redis)
	fmt.Println(redis[0].Flight.Host)
	fmt.Println(redis[1].Hotel.Host)
	fmt.Println(redis[2].Bus.Host)
}
*/

type redisDialParams struct {
	Host string
	Port int
}

type mysqlDialParams struct {
	Host     string
	Port     int
	Username string
	Password string
}

type Config struct {
	Redis []struct {
		Flight redisDialParams
		Hotel  redisDialParams
		Bus    redisDialParams
	}
	Mysql []struct {
		Flight mysqlDialParams
		Hotel  mysqlDialParams
	}
}

var (
	configFilePath = "config.json"
	config         *Config
	configLock     = new(sync.RWMutex)
)

func GetConfig() *Config {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

func init() {
	loadConfig(true)
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGUSR2)
	go func() {
		for {
			<-s
			loadConfig(false)
			log.Println("Reloaded config file")
		}
	}()
}

func loadConfig(fail bool) {
	file, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Println("Error in opening config: ", err)
		if fail {
			os.Exit(1)
		}
	}

	temp := new(Config)
	if err = json.Unmarshal(file, &temp); err != nil {
		log.Println("Error in parsing config: ", err)
		if fail {
			os.Exit(1)
		}
	}
	configLock.Lock()
	config = temp
	configLock.Unlock()
}
