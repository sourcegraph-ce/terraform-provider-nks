package nks

import (
	"github.com/NetApp/nks-sdk-go/nks"
)

type Config struct {
	Token    string
	EndPoint string
	Client   *nks.APIClient
}
