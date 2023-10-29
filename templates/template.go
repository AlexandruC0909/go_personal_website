package templates

import (
	"embed"
)

//go:embed ui/*
//go:embed auth/*
//go:embed user/*
//go:embed posts/*
//go:embed chat/*
//go:embed workspace/*

var Templates embed.FS
