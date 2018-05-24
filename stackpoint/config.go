package stackpoint

import (
	"github.com/StackPointCloud/stackpoint-sdk-go/stackpointio"
)

type Config struct {
	Token     string
	EndPoint  string
	Client    *stackpointio.APIClient
	OrgID     int
	SSHKeyset int
}
