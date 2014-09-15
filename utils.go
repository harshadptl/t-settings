package settings

func ConstructRedisPath(params map[string]string) string {
	return params["host"] + ":" + params["port"]
}
