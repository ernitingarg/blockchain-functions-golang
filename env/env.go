package env

import (
	"os"
)

// Constants for project ids
const (
	DEVELOP    string = "black-stream-292507"
	PRODUCTION string = "soteria-production"
)

type globalEnv struct {
	ProjectID string
	Keypath   string
	BtcChain  string
}

// EnvVars container for global variables
var EnvVars *globalEnv

// InitEnvVars initialize env variables
func InitEnvVars() {
	gcpProject := os.Getenv("GCP_PROJECT")
	keyPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	btcChain := "btc_test3"
	var projectID string
	switch gcpProject {
	case DEVELOP:
		projectID = DEVELOP
	case PRODUCTION:
		projectID = PRODUCTION
		btcChain = "btc_main"
	default:
		panic("project id is invalid")
	}

	EnvVars = &globalEnv{
		ProjectID: projectID,
		Keypath:   keyPath,
		BtcChain:  btcChain,
	}
}
