package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/av-belyakov/enricher_geoip/constants"
	"github.com/av-belyakov/enricher_sensor_information/internal/appname"
	"github.com/av-belyakov/enricher_sensor_information/internal/appversion"
	"github.com/av-belyakov/enricher_sensor_information/internal/confighandler"
)

func getInformationMessage(cfgNats *confighandler.CfgNats) string {
	version, err := appversion.GetVersion()
	if err != nil {
		log.Println(err)
	}

	appStatus := fmt.Sprintf("%vproduction%v", constants.Ansi_Bright_Blue, constants.Ansi_Reset)
	envValue, ok := os.LookupEnv("GO_ENRICHERSENSORINFO_MAIN")
	if ok && (envValue == "development" || envValue == "test") {
		appStatus = fmt.Sprintf("%v%s%v", constants.Ansi_Bright_Red, envValue, constants.Ansi_Reset)
	}

	msg := fmt.Sprintf("Application '%s' v%s was successfully launched", appname.GetName(), strings.Replace(version, "\n", "", -1))

	fmt.Printf("\n%v%v%s%v\n", constants.Bold_Font, constants.Ansi_Bright_Green, msg, constants.Ansi_Reset)
	fmt.Printf(
		"%v%vApplication status is '%s'%v\n",
		constants.Underlining,
		constants.Ansi_Bright_Green,
		appStatus,
		constants.Ansi_Reset,
	)
	fmt.Printf(
		"%vConnect to NATS with address %v%s:%d%v%v, subscription %v'%s'%v\n",
		constants.Ansi_Bright_Green,
		constants.Ansi_Dark_Gray,
		cfgNats.Host,
		cfgNats.Port,
		constants.Ansi_Reset,
		constants.Ansi_Bright_Green,
		constants.Ansi_Dark_Gray,
		cfgNats.Subscription,
		constants.Ansi_Reset,
	)

	return msg
}
