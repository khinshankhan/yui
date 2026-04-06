package clipboard

import (
	"os"
	"testing"
)

func TestDetectBackendPrefersWaylandWhenAvailable(t *testing.T) {
	t.Setenv("WAYLAND_DISPLAY", "wayland-1")
	t.Setenv("DISPLAY", ":0")

	origGOOS := goos
	origLookPath := lookPath
	t.Cleanup(func() {
		goos = origGOOS
		lookPath = origLookPath
	})

	goos = "linux"
	lookPath = func(file string) (string, error) {
		switch file {
		case "wl-copy", "wl-paste", "xclip", "xsel":
			return "/usr/bin/" + file, nil
		default:
			return "", os.ErrNotExist
		}
	}

	backend, err := detectBackend()
	if err != nil {
		t.Fatalf("detectBackend() error = %v", err)
	}

	if backend.Name != "wl-clipboard" {
		t.Fatalf("detectBackend() backend = %q, want wl-clipboard", backend.Name)
	}
}

func TestDetectBackendFallsBackToXclip(t *testing.T) {
	t.Setenv("WAYLAND_DISPLAY", "")
	t.Setenv("DISPLAY", ":0")

	origGOOS := goos
	origLookPath := lookPath
	t.Cleanup(func() {
		goos = origGOOS
		lookPath = origLookPath
	})

	goos = "linux"
	lookPath = func(file string) (string, error) {
		switch file {
		case "xclip":
			return "/usr/bin/" + file, nil
		default:
			return "", os.ErrNotExist
		}
	}

	backend, err := detectBackend()
	if err != nil {
		t.Fatalf("detectBackend() error = %v", err)
	}

	if backend.Name != "xclip" {
		t.Fatalf("detectBackend() backend = %q, want xclip", backend.Name)
	}
}

func TestDetectBackendOnWindowsUsesPowerShell(t *testing.T) {
	origGOOS := goos
	origLookPath := lookPath
	t.Cleanup(func() {
		goos = origGOOS
		lookPath = origLookPath
	})

	goos = "windows"
	lookPath = func(file string) (string, error) {
		if file == "powershell.exe" {
			return "C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe", nil
		}
		return "", os.ErrNotExist
	}

	backend, err := detectBackend()
	if err != nil {
		t.Fatalf("detectBackend() error = %v", err)
	}

	if backend.Name != "powershell" {
		t.Fatalf("detectBackend() backend = %q, want powershell", backend.Name)
	}
}
