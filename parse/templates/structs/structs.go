package structs

type TemplateList struct {
	AllTemplates []TemplateInfo
}

type PlatformTemplateYAML struct {
	Version int `yaml:"version"`
	Info    struct {
		ID          string   `yaml:"id"`
		Name        string   `yaml:"name"`
		Description string   `yaml:"description"`
		Tags        []string `yaml:"tags"`
		Image       string   `yaml:"image"`
		Notes       []struct {
			Heading string `yaml:"heading"`
			Content string `yaml:"content"`
		} `yaml:"notes"`
	} `yaml:"info"`
	Initialize struct {
		Repository string        `yaml:"repository"`
		Config     interface{}   `yaml:"config"`
		Files      []interface{} `yaml:"files"`
		Profile    string        `yaml:"profile"`
	} `yaml:"initialize"`
}

type PlatformApplicationsYAML struct {
	Apps []PlatformAppYAML
}

type PlatformAppYAML struct {
	Type string `json:"type"`
}

type TemplateInfo struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Sha         string `json:"sha"`
	Size        int    `json:"size"`
	URL         string `json:"url"`
	HTMLURL     string `json:"html_url"`
	GitURL      string `json:"git_url"`
	DownloadURL string `json:"download_url"`
	Type        string `json:"type"`
	Content     string `json:"content"`
	Encoding    string `json:"encoding"`
	Links       struct {
		Self string `json:"self"`
		Git  string `json:"git"`
		HTML string `json:"html"`
	} `json:"_links"`
}

// https://api.github.com/search/code?q=filename:.platform.app.yaml+repo:platformsh-templates/gatsby-strapi
type PlatformAppYAMLSearchResults struct {
	Items []struct {
		Name string `json:"name"`
		Path string `json:"path"`
	} `json:"items"`
}
