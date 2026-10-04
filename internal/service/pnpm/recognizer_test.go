package pnpm

import (
	"path/filepath"
	"testing"
)

func TestParseWorkspace(t *testing.T) {
	testbed := filepath.Join("..", "..", "..", "/testbed")
	workspace, err := ParseWorkspace(testbed)
	if err != nil {
		t.Fatalf("Failed to parse workspace: %v", err)
	}

	if len(workspace.Packages) == 0 {
		t.Fatalf("Workspace has no packages")
	}

	for _, pkg := range workspace.Packages {
		if pkg.Name == "" {
			t.Errorf("Package has no name")
		}
		if pkg.Version == "" {
			t.Errorf("Package has no version")
		}
		if len(pkg.Scripts) == 0 {
			t.Errorf("Package has no scripts")
		}

		for name, script := range pkg.Scripts {
			switch pkg.Name {
			case "domain":
				if name != "hello" {
					t.Errorf("Script name should be 'hello'")
				}

				if script != "echo \"Hello\"" {
					t.Errorf("Script should be empty: Command: echo \"Hello\"")
				}
			case "server":
				switch name {
				case "dev":
					if script != "tsx src/main.ts" {
						t.Errorf("Script for server dev should be: tsx src/main.ts")
					}
				case "dev:watch":
					if script != "tsx --watch src/main.ts" {
						t.Errorf("Script for server dev:watch should be: tsx --watch src/main.ts")
					}
				default:
					t.Fatalf("Script name should be dev or dev:watch")
				}
			default:
				t.Fatalf("Package name should be domain or server")
			}
		}
	}
}
