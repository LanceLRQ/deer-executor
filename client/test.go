package client

import (
    "github.com/pkg/errors"
    "github.com/urfave/cli/v2"
)

func Test(c *cli.Context) error {
    return errors.Errorf("AAAA")
}
