package structs

import (
	"context"
	"io"
)

// ShellResult  shell result for exec.Command
type ShellResult struct {
	Success      bool
	Stdout       string
	Stderr       string
	ExitCode     int
	Signal       int
	ErrorMessage string
}

// ShellWriters defind shell result writers
type ShellWriters struct {
	Input  io.Reader
	Output io.Writer
	Error  io.Writer
}

// ShellOptions defind shell options for exec.Command
type ShellOptions struct {
	Context   context.Context
	Name      string
	Args      []string
	StdWriter *ShellWriters
	OnStart   func(io.Writer) error
}
