package main

import ()

type BuildStatus struct {
	ClientID    string
	JobID       string
	BuildNumber int
	BuildStatus bool
}

type LightConfig struct {
	ClientID    string
	LightID     string
	Description string
	JobIDRegEx  string
}
