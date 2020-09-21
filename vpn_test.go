package cuckoo

import (
	"context"
	"fmt"
	"testing"
)

func TestVPNStatus(t *testing.T) {
	c := getTestingClient()

	vpnStatus, err := c.VPNStatus(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(vpnStatus)
}
