package atomicfile_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/creachadair/atomicfile"
)

// withTempDir invokes f with the current working directory as a freshly
// created temp directory, whose path is given. When f returns, the temp
// directory is cleaned up.
func withTempDir(t *testing.T, f func(path string)) {
	t.Helper()

	dir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Creating temp directory: %v", err)
	}
	defer os.RemoveAll(dir)
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Chdir failed: %v", err)
	}
	f(dir)
}

func mustMkdirAll(t *testing.T, path string) string {
	t.Helper()

	if err := os.MkdirAll(path, 0700); err != nil {
		t.Fatalf("Mkdir: %v", err)
	}
	return path
}

func checkFile(t *testing.T, path string, perm os.FileMode, want string) {
	t.Helper()

	got, err := os.ReadFile(path)
	if err != nil {
		t.Errorf("Reading target: %v", err)
	}
	if s := string(got); s != want {
		t.Errorf("Target contents: got %q, want %q", s, want)
	}

	if fi, err := os.Stat(path); err != nil {
		t.Errorf("Stat failed: %v", err)
	} else if m := fi.Mode().Perm(); m != perm {
		t.Errorf("Target mode: got %v, want %v", m, perm)
	}
}

func TestFile(t *testing.T) {
	withTempDir(t, func(tmp string) {
		path := filepath.Join(mustMkdirAll(t, "a/b/c"), "target.txt")

		f, err := atomicfile.New(path, 0623)
		if err != nil {
			t.Fatalf("New %q failed: %v", path, err)
		}

		const message = "Hello, world\n"
		fmt.Fprint(f, message)
		if err := f.Close(); err != nil {
			t.Errorf("Close failed: %v", err)
		}
		checkFile(t, path, 0623, message)
	})
}

func TestCancel(t *testing.T) {
	withTempDir(t, func(tmp string) {
		path := filepath.Join(tmp, "target.txt")

		f, err := atomicfile.New(path, 0600)
		if err != nil {
			t.Fatalf("New %q failed: %v", path, err)
		}

		fmt.Fprintln(f, "Some of what a fool thinks often remains")
		f.Cancel()

		// After cancellation, a close should report an error.
		if err := f.Close(); err == nil {
			t.Error("Closing f should have reported an error")
		}

		// The target file should not exist, since it did not already.
		if fi, err := os.Stat(path); err == nil {
			t.Errorf("Stat %q should have failed, but found %d bytes", path, fi.Size())
		}

		// No temporary files should be left around.
		dc, err := ioutil.ReadDir(tmp)
		if err != nil {
			t.Errorf("ReadDir %q: unexpected error: %v", tmp, err)
		}
		for _, fi := range dc {
			t.Errorf("Unexpected file %q in output directory", fi.Name())
		}
	})
}

func TestNoClobber(t *testing.T) {
	withTempDir(t, func(tmp string) {
		path := filepath.Join(tmp, "target.txt")

		const oldMessage = "If I keep my eyes closed he looks just like you"
		if err := os.WriteFile(path, []byte(oldMessage), 0400); err != nil {
			t.Fatalf("Writing target file: %v", err)
		}

		f, err := atomicfile.New(path, 0644)
		if err != nil {
			t.Fatalf("New %q failed: %v", path, err)
		}

		fmt.Fprintln(f, "You should never see this")
		f.Cancel()

		if err := f.Close(); err == nil {
			t.Error("Closing f should have reported an error")
		}

		// After cancellation, the existing target should be unchanged.
		checkFile(t, path, 0400, oldMessage)
	})
}

func TestDeferredCancel(t *testing.T) {
	withTempDir(t, func(tmp string) {
		path := filepath.Join(tmp, "target.txt")

		f, err := atomicfile.New(path, 0640)
		if err != nil {
			t.Fatalf("New %q failed: %v", path, err)
		}

		// A cancel that happens after a successful close does not interfere with
		// the output.
		const message = "There's a place way down in Bed-Stuy"
		func() {
			defer f.Cancel()
			fmt.Fprint(f, message)

			if err := f.Close(); err != nil {
				t.Errorf("Close failed: %v", err)
			}
		}()

		checkFile(t, path, 0640, message)
	})
}
