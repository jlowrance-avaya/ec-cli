package main

type Job struct {
	Kind        string      `json:"kind"`
	Type        string      `json:"type"`
	Subtype     string      `json:"subtype"`
	Version     string      `json:"version"`
	Metadata    JobMetadata `json:"metadata"`
	Payload     JobPayload  `json:"payload"`
	ID          string      `json:"id"`
	RID         string      `json:"_rid"`
	Self        string      `json:"_self"`
	ETag        string      `json:"_etag"`
	Attachments string      `json:"_attachments"`
	Timestamp   int64       `json:"_ts"`
}

type JobMetadata struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Created     string `json:"created"`
}

type JobPayload struct {
	Workspace  string     `json:"workspace"`
	CodeRepo   Repo       `json:"code_repo"`
	ConfigRepo ConfigRepo `json:"config_repo"`
}

type Repo struct {
	URL     string `json:"url"`
	Branch  string `json:"branch"`
	Version string `json:"version"`
}

type ConfigRepo struct {
	URL            string `json:"url"`
	Branch         string `json:"branch"`
	Version        string `json:"version"`
	ConfigFileName string `json:"configfile_name"`
	ConfigFilePath string `json:"configfile_path"`
}
