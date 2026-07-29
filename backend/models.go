package backend

// Note đại diện cho cấu trúc 1 file markdown
type Note struct {
	ID           string   `json:"id"`           // Đường dẫn file (tương đối hoặc tuyệt đối)
	Title        string   `json:"title"`        // Tiêu đề từ frontmatter hoặc tên file/heading
	Content      string   `json:"content"`      // Nội dung markdown
	Tags         []string `json:"tags"`         // Danh sách tag từ frontmatter
	Folder       string   `json:"folder"`       // ID thư mục cha
	LastModified string   `json:"lastModified"` // Thời gian cập nhật gần nhất (chuẩn ISO)
	Excerpt      string   `json:"excerpt"`      // Tóm tắt nội dung (đoạn đầu)
}

// Folder đại diện cho thư mục chứa ghi chú hoặc tệp tin
type Folder struct {
	ID       string    `json:"id"`                 // Đường dẫn chuẩn hoá
	Name     string    `json:"name"`               // Tên thư mục
	Icon     string    `json:"icon"`               // Icon biểu trưng (vd: 📁, 📚)
	Children []*Folder `json:"children,omitempty"` // Các thư mục con
	Level    int       `json:"level"`              // Độ sâu cây thư mục
	Type     string    `json:"type,omitempty"`     // Loại ('markdown' | 'file')
}

// WorkspaceData lưu trữ toàn bộ dữ liệu quét được từ workspace
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
