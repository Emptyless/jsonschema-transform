package domain

// Source returns information about the path where the data was parsed from
type Source interface {
	Path() string
}

// FileSource is a simple Source representation where it's a single file source
type FileSource struct {
	FilePath string
}

// Path of the Source
func (f FileSource) Path() string {
	return f.FilePath
}
