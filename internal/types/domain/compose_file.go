package types

type ComposeFile struct {
	Branch       string `json:"branch"`
	RepoName     string `json:"repo_name"`
	RepoUsername string `json:"repo_username"`
	RepoPassword string `json:"repo_password"`
}
