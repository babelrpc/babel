package main

import "embed"

//go:embed media/*.png swagger/* swagger/*/*
var webContent embed.FS
