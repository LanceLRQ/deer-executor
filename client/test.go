package client

import (
	"github.com/urfave/cli/v2"
	"time"
)

// Test for cli command 'test'
func Test(c *cli.Context) error {
	time.Sleep(time.Second * 10)
	return nil
}
