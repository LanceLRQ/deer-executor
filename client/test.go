package client

import (
    deer_common "github.com/LanceLRQ/deer-common"
    "github.com/urfave/cli/v2"
)

func Test(c *cli.Context) error {
    deer_common.Test()
    return nil
}
