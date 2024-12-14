package main

type ConnectionState int

const (
	Disconnected ConnectionState = iota
	Connecting
	Connected
	Failed
)

type UserRole int

const (
	Admin UserRole = iota + 1
	Editor
	Viewer
)
