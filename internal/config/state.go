package config

type State string

const (
	StateLocal    State = "local"
	StateLoadTest State = "loadtest"
	StateDEV      State = "dev"
	StateSIT      State = "sit"
	StateUAT      State = "uat"
	StateProd     State = "prod"
)
