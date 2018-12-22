package nks

import (
	"github.com/StackPointCloud/nks-sdk-go/nks"
)

type Config struct {
	Token    string
	EndPoint string
	Client   *nks.APIClient
}
