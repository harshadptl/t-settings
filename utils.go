package settings

func ConstructRedisPath(params map[string]string) string {
	return params["host"] + ":" + params["port"]
}

func ConstructMysqlPath(params map[string]string) string {
	return params["username"] + ":" + params["password"] +
		"@tcp(" + params["host"] + ":" + params["port"] +
		")/" + params["database"]
}
