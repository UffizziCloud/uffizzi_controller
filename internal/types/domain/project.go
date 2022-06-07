package types

type Project struct {
	IsPreviewsProtected bool   `json:"is_previews_protected"`
	PreviewsPassword    string `json:"previews_password"`
	PreviewsUserName    string `json:"previews_username"`
}
