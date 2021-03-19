package structs

import (
	"context"
	"io"
)

type ShellResult struct {
	Success      bool
	Stdout       string
	Stderr       string
	ExitCode     int
	Signal       int
	ErrorMessage string
}

type ShellWriters struct {
	Input  io.Reader
	Output io.Writer
	Error  io.Writer
}

type ShellOptions struct {
	Context   context.Context
	Name      string
	Args      []string
	StdWriter *ShellWriters
	OnStart   func(io.Writer) error
}
