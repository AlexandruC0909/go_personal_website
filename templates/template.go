package templates

import (
	"embed"
)

//go:embed ui/*
//go:embed auth/*
//go:embed user/*
//go:embed posts/*

var Templates embed.FS
