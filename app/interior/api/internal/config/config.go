package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	Sqlite
}

type Sqlite struct {
	SqliteDsn string `json:"SqliteDsn"`
}
