package templates

import (
	"embed"
)

//go:embed ui/*
//go:embed auth/*
//go:embed user/*
//go:embed posts/*
//go:embed chat/*

var Templates embed.FS
