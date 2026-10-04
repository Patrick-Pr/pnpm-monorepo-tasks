package delegate

import (
	"errors"
	"os"
	"os/exec"
	"runtime"
	"testing"
)

func TestScriptCommandWaitsForDescendantOutput(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("test fixture requires a POSIX shell")
	}

	for _, tc := range []struct {
		name     string
		script   string
		exitCode int
		stdout   string
	}{
		{"success", "(sleep 0.1; printf child; printf child-error >&2) & printf parent; exit 0", 0, "parentchild"},
		{"failure", "(sleep 0.1; printf child; printf child-error >&2) & printf parent; exit 7", 7, "parentchild"},
		{"stderr only", "(sleep 0.1; printf child-error >&2) >/dev/null & exit 0", 0, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			stdout, err := os.CreateTemp(t.TempDir(), "stdout")
			if err != nil {
				t.Fatal(err)
			}
			defer stdout.Close()
			stderr, err := os.CreateTemp(t.TempDir(), "stderr")
			if err != nil {
				t.Fatal(err)
			}
			defer stderr.Close()

			cmd := &scriptCommand{Cmd: exec.Command("sh", "-c", tc.script)}
			cmd.SetStdout(stdout)
			cmd.SetStderr(stderr)
			err = cmd.Run()
			if tc.exitCode == 0 {
				if err != nil {
					t.Fatal(err)
				}
			} else {
				var exitErr *exec.ExitError
				if !errors.As(err, &exitErr) || exitErr.ExitCode() != tc.exitCode {
					t.Fatalf("expected exit code %d, got %v", tc.exitCode, err)
				}
			}

			for _, output := range []struct {
				file *os.File
				want string
			}{{stdout, tc.stdout}, {stderr, "child-error"}} {
				got, err := os.ReadFile(output.file.Name())
				if err != nil {
					t.Fatal(err)
				}
				if string(got) != output.want {
					t.Errorf("output at command completion = %q, want %q", got, output.want)
				}
			}
		})
	}
}
