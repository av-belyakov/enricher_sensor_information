package information

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/av-belyakov/enricher_geoip/constants"
	"github.com/av-belyakov/enricher_sensor_information/internal/appname"
	"github.com/av-belyakov/enricher_sensor_information/internal/appversion"
	"github.com/av-belyakov/enricher_sensor_information/internal/dicontainer"
)

func GetInformationMessage(cfg dicontainer.Configer) string {
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
	fmt.Printf("%vSensor information database settings:%v\n", constants.Ansi_Bright_Green, constants.Ansi_Reset)
	fmt.Printf("  %vZabbix host:%v%s%v\n", constants.Ansi_Bright_Green, constants.Ansi_Bright_Blue, cfg.GetSensorInformationDB().ZabbixHost, constants.Ansi_Reset)
	fmt.Printf("  %vZabbix user:%v%s%v\n", constants.Ansi_Bright_Green, constants.Ansi_Bright_Blue, cfg.GetSensorInformationDB().ZabbixUser, constants.Ansi_Reset)
	fmt.Printf("  %vNetbox host:%v%s%v\n", constants.Ansi_Bright_Green, constants.Ansi_Bright_Blue, cfg.GetSensorInformationDB().NetboxHost, constants.Ansi_Reset)
	fmt.Printf("  %vNetbox port:%v%d%v\n", constants.Ansi_Bright_Green, constants.Ansi_Bright_Magenta, cfg.GetSensorInformationDB().NetboxPort, constants.Ansi_Reset)
	fmt.Printf("  %vNCIRCC URL:%v%s%v\n", constants.Ansi_Bright_Green, constants.Ansi_Bright_Blue, cfg.GetSensorInformationDB().NCIRCCURL, constants.Ansi_Reset)
	fmt.Printf(
		"%vConnect to NATS with address %v%s:%d%v%v, subscription %v'%s'%v\n",
		constants.Ansi_Bright_Green,
		constants.Ansi_Dark_Gray,
		cfg.GetNATS().Host,
		cfg.GetNATS().Port,
		constants.Ansi_Reset,
		constants.Ansi_Bright_Green,
		constants.Ansi_Dark_Gray,
		cfg.GetNATS().Subscription,
		constants.Ansi_Reset,
	)

	return msg
}
