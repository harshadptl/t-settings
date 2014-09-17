package settings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

type Config map[string]interface{}

func getConfigsFor(db string) (dbSettings Config) {
	jsonConfig := *config
	redisSpecificConfigs := jsonConfig[db]
	dbSettings = redisSpecificConfigs.(map[string]interface{})
	return
}

func getVerticalSpecificSettings(db string, vertical string, settings Config) Config {
	verticalSpecificConfigs := settings[vertical]
	if verticalSpecificConfigs == nil {
		panic(vertical + " doesn't exist. Use one of these dbs " + strings.Join(GetAllDbsFor(db), ", "))
		verticalSpecificConfigs = make(map[string]interface{})
	}
	config := verticalSpecificConfigs.(map[string]interface{})
	return config
}

func GetConfigsFor(db string, vertical string) (params map[string]string) {
	configLock.RLock()
	defer configLock.RUnlock()
	params = make(map[string]string)
	verticalSettings := getConfigsFor(db)
	configs := getVerticalSpecificSettings(db, vertical, verticalSettings)
	for config := range configs {
		params[config] = configs[config].(string)
	}
	return
}

func GetAllDbsFor(db string) (allDbs []string) {
	verticalSettings := getConfigsFor(db)
	for key := range verticalSettings {
		allDbs = append(allDbs, key)
	}
	return
}

var (
	configFilePath = "" //file path to a config file
	config         *Config
	configLock     = new(sync.RWMutex)
)

func Configure(filePath ...string) {
	num_config_files := len(filePath)
	if num_config_files == 0 {
		configFilePath = "src/github.com/vireshas/t-settings/config.json"
	} else if num_config_files == 1 {
		fmt.Println("we dont support than one config file at the moment")
		os.Exit(1)
	} else {
		configFilePath = filePath[0]
	}
	loadConfig(true)
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGUSR2)
	go func() {
		for {
			<-s
			loadConfig(false)
			fmt.Println("Reloaded config file")
		}
	}()
}

func loadConfig(fail bool) {
	file, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		fmt.Println("Error in opening config: ", err)
		if fail {
			os.Exit(1)
		}
	}

	temp := new(Config)
	if err = json.Unmarshal(file, &temp); err != nil {
		fmt.Println("Error in parsing config: ", err)
		if fail {
			os.Exit(1)
		}
	}
	configLock.Lock()
	config = temp
	configLock.Unlock()
}
