package client

import (
    "github.com/urfave/cli/v2"
)

var RunFlags = []cli.Flag{
    &cli.BoolFlag{
        Name:  "no-clean",
        Value: false,
        Usage: "Don't delete session directory after judge",
    },
    &cli.StringFlag{
        Name:    "language",
        Aliases: []string{"l"},
        Value:   "auto",
        Usage:   "Code language name",
    },
    &cli.BoolFlag{
        Name:  "debug",
        Value: false,
        Usage: "Print debug log",
    },
    &cli.IntFlag{
        Name:  "benchmark",
        Value: 0,
        Usage: "Start benchmark",
    },
    &cli.StringFlag{
        Name:    "persistence",
        Aliases: []string{"p"},
        Value:   "",
        Usage:   "Persistent judge result to file (support: gzip, none)",
    },
    &cli.BoolFlag{
        Name:    "save-ac-data",
        Value:   false,
        Usage:   "Persistent an ACCEPTED test case's output data, will increase the file size",
    },
    &cli.StringFlag{
        Name:  "compress",
        Value: "gzip",
        Usage: "Persistent compressor type",
    },
    &cli.BoolFlag{
        Name:    "sign",
        Aliases: []string{"s"},
        Value:   false,
        Usage:   "Enable digital sign (GPG)",
    },
    &cli.BoolFlag{
        Name:  "detail",
        Value: false,
        Usage: "Show test-cases details",
    },
    &cli.StringFlag{
        Name:    "gpg-key",
        Aliases: []string{"key"},
        Value:   "",
        Usage:   "GPG private key file",
    },
    &cli.StringFlag{
        Name:    "passphrase",
        Aliases: []string{"password", "pwd"},
        Value:   "",
        Usage:   "GPG private key passphrase",
    },
    &cli.StringFlag{
        Name:    "work-dir",
        Aliases: []string{"w"},
        Value:   "",
        Usage:   "Working dir, using to unpack problem package",
    },
    &cli.StringFlag{
        Name:  "session-id",
        Value: "",
        Usage: "setup session id",
    },
    &cli.StringFlag{
        Name:  "session-root",
        Value: "",
        Usage: "setup session root dir",
    },
    &cli.StringFlag{
        Name:  "library",
        Value: "./lib",
        Usage: "library root for special judge, contains \"testlib.h\" and \"bits/stdc++.h\" etc.",
    },
    &cli.StringFlag{
        Name:  "log-level",
        Value: "",
        Usage: "set logs level (debug|info|warn|error)",
    },
    &cli.BoolFlag{
        Name:  "log",
        Value: false,
        Usage: "output logs to stdout",
    },
}


