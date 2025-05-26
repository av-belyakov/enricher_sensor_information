package ncirccinteraction_test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/av-belyakov/enricher_sensor_information/internal/ncirccinteractions"
	"github.com/stretchr/testify/assert"
	"github.com/subosito/gotenv"
)

const (
	NICRCC_URL = "https://10.0.227.10/api/v2/companies"
	INN        = "6905011546"
)

var (
	client *ncirccinteractions.ClientNICRCC

	err error
)

func TestMain(t *testing.M) {
	if err = gotenv.Load("../../.env"); err != nil {
		log.Fatalln(err)
	}

	client, err = ncirccinteractions.NewClient(
		NICRCC_URL,
		os.Getenv("GO_ENRICHERSENSORINFO_SINCIRCCTOKEN"),
		15*time.Second)
	if err != nil {
		log.Fatalln(err)
	}

	os.Exit(t.Run())
}

func TestGetInformation(t *testing.T) {
	res, err := client.GetFullNameOrganizationByINN(t.Context(), INN)
	assert.NoError(t, err)

	fmt.Printf("Response:\n%+v\n", res)

	assert.NotEmpty(t, res, res.Data)

	t.Cleanup(func() {
		os.Unsetenv("GO_ENRICHERSENSORINFO_SIPASSWD")
		os.Unsetenv("GO_ENRICHERSENSORINFO_DBWLOGPASSWD")
		os.Unsetenv("GO_ENRICHERSENSORINFO_SINCIRCCTOKEN")
	})
}
