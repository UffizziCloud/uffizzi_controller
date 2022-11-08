package types

type ComposeFileSourceKind string

const (
	ComposeFileSourceKindLocal  ComposeFileSourceKind = "local"
	ComposeFileSourceKindGithub ComposeFileSourceKind = "github"
)

type ComposeFile struct {
	Branch       string                `json:"branch"`
	RepoName     string                `json:"repo_name"`
	RepoUsername string                `json:"repo_username"`
	RepoPassword string                `json:"repo_password"`
	Path         string                `json:"path"`
	SourceKind   ComposeFileSourceKind `json:"source_kind"`
}

func (cf ComposeFile) IsGithubSourceKind() bool {
	return cf.SourceKind == ComposeFileSourceKindGithub
}

func (cf ComposeFile) IsLocalSourceKind() bool {
	return cf.SourceKind == ComposeFileSourceKindLocal
}
