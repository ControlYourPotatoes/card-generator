package manager

import (
	"path/filepath"
	"sync"
)

type FontManager struct {
	fontPaths   map[string]string
	fontCache   sync.Map
	defaultFont string
}

func NewFontManager() (*FontManager, error) {
	return &FontManager{
		fontPaths: map[string]string{
			"regular": filepath.Join("assets", "fonts", "regular.ttf"),
			"bold":    filepath.Join("assets", "fonts", "bold.ttf"),
			"italic":  filepath.Join("assets", "fonts", "italic.ttf"),
		},
		defaultFont: "regular",
	}, nil
}

func (fm *FontManager) GetFontPath(name string) string {
	if path, exists := fm.fontPaths[name]; exists {
		return path
	}
	return fm.fontPaths[fm.defaultFont]
}
