package domain

// SecretFileExt represents supported file extensions.
type SecretFileExt string

const (
	SecretFileExtEnv  SecretFileExt = "env"
	SecretFileExtYAML SecretFileExt = "yaml"
	SecretFileExtJSON SecretFileExt = "json"
	SecretFileExtTXT  SecretFileExt = "txt"
)

// Namespace represents a folder in the secret tree.
type Namespace struct {
	ID        string       `qwery:"id"`
	ParentID  *string      `qwery:"parent_id"`
	Name      string       `qwery:"name"`
	CreatedAt string       `qwery:"created_at"`
	UpdatedAt string       `qwery:"updated_at"`
	DeletedAt int          `qwery:"deleted_at"`
	Children  []*Namespace `qwery:"-"`
}

// NamespaceFilter holds filter criteria for namespace queries.
type NamespaceFilter struct {
	ParentID *string
}

// SecretFile represents a file inside a namespace.
type SecretFile struct {
	ID          string        `qwery:"id"`
	NamespaceID string        `qwery:"namespace_id"`
	Name        string        `qwery:"name"`
	Ext         SecretFileExt `qwery:"ext"`
	CreatedAt   string        `qwery:"created_at"`
	UpdatedAt   string        `qwery:"updated_at"`
	DeletedAt   int           `qwery:"deleted_at"`
}

// Secret represents a key-value pair inside an env file.
type Secret struct {
	ID        string `qwery:"id"`
	FileID    string `qwery:"file_id"`
	KeyName   string `qwery:"key_name"`
	ValueEnc  string `qwery:"value_enc"`
	CreatedAt string `qwery:"created_at"`
	UpdatedAt string `qwery:"updated_at"`
	DeletedAt int    `qwery:"deleted_at"`
}

// FileContent represents the raw encrypted content of a yaml/json/txt file.
type FileContent struct {
	ID        string `qwery:"id"`
	FileID    string `qwery:"file_id"`
	Content   string `qwery:"content"`
	UpdatedAt string `qwery:"updated_at"`
	DeletedAt int    `qwery:"deleted_at"`
}

// ExportFile is the decrypted export payload for a single file.
type ExportFile struct {
	Name    string            `json:"name"`
	Ext     SecretFileExt     `json:"ext"`
	Secrets map[string]string `json:"secrets,omitempty"`
	Content string            `json:"content,omitempty"`
}
