package main

import (
	"log"
	"os"

	"github.com/KyberNetwork/reserve-data/cmd/configuration"
)

func main() {
	var config *configuration.Config
	switch os.Getenv("KYBER_ENV") {
	case "mainnet", "production":
		log.Printf("Running in production mode")
		config = configuration.GetConfigForMainnet()
		break
	case "staging":
		log.Printf("Running in staging mode")
		config = configuration.GetConfigForStaging()
		break
	case "simulation":
		log.Printf("Running in simulation mode")
		config = configuration.GetConfigForSimulation()
		break
	case "kovan":
		log.Printf("Running in kovan mode")
		config = configuration.GetConfigForKovan()
		break
	case "ropsten":
		log.Printf("Running in ropsten mode")
		config = configuration.GetConfigForRopsten()
		break
	case "dev":
		log.Printf("Running in dev mode")
		config = configuration.GetConfigForDev()
	default:
		log.Printf("Running in dev mode")
		config = configuration.GetConfigForDev()
	}
	if config.AuthEngine == nil {
		Warning.Println("Current environment setting does not enable authentication. Please check again!!!")
	}
	verify := NewVerification(config.AuthEngine)
	verify.RunVerification()
}
