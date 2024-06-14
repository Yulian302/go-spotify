package env

import "github.com/gofor-little/env"

const ENV_FILE_PATH = ".env"

var (
	ClusterName  string = ""
	UserPassword        = ""
	DbName       string = ""
)

func LoadEnv() {
	// environment variables
	if err := env.Load(ENV_FILE_PATH); err != nil {
		panic(err)
	}
	ClusterName = env.Get("CLUSTER_NAME", "")
	UserPassword = env.Get("CLUSTER_USER_PASSWD", "")
	DbName = env.Get("DB_NAME", "")
}
