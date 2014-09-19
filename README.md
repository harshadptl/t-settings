t-settings
==========

####A go module to deal with config files; One can update a config file without restarting the service.

        package main

        import (
                "fmt"
                "github.com/vireshas/t-settings"
        )

        func main() {
                //pass filepath to config file
                //default : "src/github.com/vireshas/t-settings/config.json"
                settings.Configure()
                //or settings.Configure("src/github.com/vireshas/t-settings/config.json")
                configs := settings.GetConfigsFor("mysql", "m1")
                fmt.Println(configs)
                //construct mysql url from config map
                fmt.Println(settings.ConstructMysqlPath(configs))
                //all the dbs for redis
                fmt.Println(settings.GetAllDbsFor("redis"))
                //all the dbs for mysql
                fmt.Println(settings.GetAllDbsFor("mysql"))
        }
