package settings

func ConstructRedisPath(params Config) string {
	return params["host"] + ":" + params["port"]
}
