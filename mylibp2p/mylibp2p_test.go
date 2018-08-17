package mylibp2p

import (
	"fmt"
	"testing"
)

func TestInitHost(t *testing.T) {
	host, err := InitHost()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(host.ID().Pretty())
}
