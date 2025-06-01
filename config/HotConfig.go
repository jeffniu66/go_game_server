package config

var ConfigList []interface{}

type HotConfig interface {
	HotInit()
}
