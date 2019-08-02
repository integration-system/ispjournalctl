package config

type Configuration struct {
	GateHost string `valid:"required~Required"`
}
