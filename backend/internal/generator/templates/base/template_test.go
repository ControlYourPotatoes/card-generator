package base

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestNewBaseTemplate(t *testing.T) {
	template := NewBaseTemplate()

	if template == nil {
		t.Fatal("NewBaseTemplate returned nil")
	}

	if template.framesPath == "" {
		t.Error("template frames path is empty")
	}

	if template.artBounds.Empty() {
		t.Error("template art bounds are empty")
	}

	expectedArtBounds := GetDefaultArtBounds()
	if template.artBounds != expectedArtBounds {
		t.Errorf("art bounds = %v, want %v", template.artBounds, expectedArtBounds)
	}
}

func TestLoadFrame(t *testing.T) {
	template := NewBaseTemplate()

	// Create a temporary test image
	tmpDir := t.TempDir()
	testImagePath := filepath.Join(tmpDir, "test.png")

	// Create a 1x1 black PNG for testing
	f, err := os.Create(testImagePath)
	if err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}
	defer f.Close()

	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	if err := png.Encode(f, img); err != nil {
		t.Fatalf("Failed to encode test image: %v", err)
	}

	// Test cases
	tests := []struct {
		name      string
		imageName string
		wantErr   bool
	}{
		{
			name:      "Valid image",
			imageName: "test.png",
			wantErr:   false,
		},
		{
			name:      "Non-existent image",
			imageName: "nonexistent.png",
			wantErr:   true,
		},
	}

	// Override template's frames path for testing
	template.framesPath = tmpDir

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frame, err := template.LoadFrame(tt.imageName)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if frame == nil {
				t.Error("frame is nil")
			}
		})
	}
}

func TestGetTemplateDir(t *testing.T) {
	dir := getTemplateDir()

	if dir == "" {
		t.Error("getTemplateDir returned empty string")
	}

	// Verify the directory structure
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Logf("Template directory not found at: %s", dir)
		t.Log("Note: This is not necessarily an error if running tests before directory setup")
	}
}

func TestGetDefaultArtBounds(t *testing.T) {
	bounds := GetDefaultArtBounds()

	if bounds.Empty() {
		t.Error("default art bounds are empty")
	}

	// Test specific dimensions
	expectedBounds := image.Rect(170, 240, 1330, 1000)
	if bounds != expectedBounds {
		t.Errorf("art bounds = %v, want %v", bounds, expectedBounds)
	}

	// Verify bounds are reasonable
	width := bounds.Dx()
	height := bounds.Dy()

	if width <= 0 {
		t.Error("art bounds width should be positive")
	}
	if height <= 0 {
		t.Error("art bounds height should be positive")
	}
}

func TestGetArtBounds(t *testing.T) {
	template := NewBaseTemplate()
	bounds := template.GetArtBounds()

	if bounds.Empty() {
		t.Error("GetArtBounds returned empty bounds")
	}

	expectedBounds := GetDefaultArtBounds()
	if bounds != expectedBounds {
		t.Errorf("art bounds = %v, want %v", bounds, expectedBounds)
	}
}
