package main

import "time"

type ManifestTemplate struct {
	CreatedAt   *time.Time                `json:"createdAt"`
	UpdatedAt   *time.Time                `json:"updatedAt"`
	ID          string                    `json:"id"`
	Version     string                    `json:"version"`
	Kind        string                    `json:"kind"`
	Metadata    ManifestTemplateMetadata  `json:"metadata"`
	Variables   ManifestTemplateVariables `json:"variables"`
	Blocks      []ManifestTemplateBlock   `json:"blocks"`
	RID         string                    `json:"_rid"`
	Self        string                    `json:"_self"`
	ETag        string                    `json:"_etag"`
	Attachments string                    `json:"_attachments"`
	Timestamp   int64                     `json:"_ts"`
}

type ManifestTemplateMetadata struct {
	UUID        string                `json:"uuid"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Labels      interface{}           `json:"labels"`
	Tags        []ManifestTemplateTag `json:"tags"`
	Selector    interface{}           `json:"selector"`
	Created     string                `json:"created"`
	Creator     string                `json:"creator"`
	Modified    string                `json:"modified"`
	Modifier    string                `json:"modifier"`
}

type ManifestTemplateTag struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ManifestTemplateVariables map[string]string

type ManifestTemplateBlock struct {
	Kind     string                        `json:"kind"`
	Metadata ManifestTemplateBlockMetadata `json:"metadata"`
	Jobs     []ManifestTemplateBlockJob    `json:"jobs"`
}

type ManifestTemplateBlockMetadata struct {
	UUID        string      `json:"uuid"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Labels      interface{} `json:"labels"`
	Tags        interface{} `json:"tags"`
	Selector    interface{} `json:"selector"`
	Created     interface{} `json:"created"`
	Creator     interface{} `json:"creator"`
	Modified    interface{} `json:"modified"`
	Modifier    interface{} `json:"modifier"`
}

type ManifestTemplateBlockJob struct {
	Version  string      `json:"version"`
	Type     string      `json:"type"`
	UUID     string      `json:"uuid"`
	Name     string      `json:"name"`
	Sequence int         `json:"sequence"`
	Target   string      `json:"target"`
	Labels   []string    `json:"labels"`
	Tags     interface{} `json:"tags"`
}
