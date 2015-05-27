package main

import (
	"github.com/gkiki90/envman/pathutil"
)

var DefaultPath string = pathutil.UserHomeDir() + ".envman/envlist.json"