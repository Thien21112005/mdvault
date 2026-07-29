package model

// Note đại diện cho cấu trúc 1 file markdown trong Domain
type Note struct {
	ID           string   `json:"id"`           // Đường dẫn file chuẩn hóa
	Title        string   `json:"title"`        // Tiêu đề note
	Content      string   `json:"content"`      // Nội dung markdown
	Tags         []string `json:"tags"`         // Danh sách tag từ frontmatter
	Folder       string   `json:"folder"`       // ID thư mục cha
	LastModified string   `json:"lastModified"` // Chuẩn RFC3339 ISO String
	Excerpt      string   `json:"excerpt"`      // Tóm tắt nội dung
}

// Folder đại diện cho thư mục lưu trữ trong cây Workspace
type Folder struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Icon     string    `json:"icon"`
	Children []*Folder `json:"children,omitempty"`
	Level    int       `json:"level"`
	Type     string    `json:"type,omitempty"`
}

// WorkspaceData lưu trữ toàn bộ thông tin hệ thống tệp ghi chú
type WorkspaceData struct {
	RootPath string    `json:"rootPath"`
	Folders  []*Folder `json:"folders"`
	Notes    []Note    `json:"notes"`
	AllTags  []string  `json:"allTags"`
}

// Frontmatter mô tả cấu trúc YAML đầu file markdown
type Frontmatter struct {
	Title string   `yaml:"title"`
	Tags  []string `yaml:"tags"`
	Icon  string   `yaml:"icon,omitempty"`
}
