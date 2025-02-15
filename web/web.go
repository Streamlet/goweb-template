package web

import "embed"

//go:embed js/* *.html
var RootFs embed.FS

//go:embed include/*
var IncludeFs embed.FS
