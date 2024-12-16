package config

var G_Config = &Config{}

type Config struct {
	Db struct {
		Name string
	}
}

func SetDbName(name string) {
	G_Config.Db.Name = name
}
