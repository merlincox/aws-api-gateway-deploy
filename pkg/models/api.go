// Code generated by schema-generator. DO NOT EDIT.

package models

// Empty: (No description)
type Empty struct {
}

// Status: API status information
type Status struct {
	Branch    string `json:"branch"`
	Platform  string `json:"platform"`
	Commit    string `json:"commit"`
	Release   string `json:"release"`
	Timestamp string `json:"timestamp"`
}
