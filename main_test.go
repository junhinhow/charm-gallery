package main

import (
	"bytes"
	"os"
	"os/exec"
	"sort"
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

func runGallery(t *testing.T, args ...string) (int, string, string) {
	t.Helper()
	var out, errOut bytes.Buffer
	cmd := exec.Command("go", append([]string{"run", "."}, args...)...)
	cmd.Dir = "."
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	cmd.Env = append(os.Environ(), "GOWORK=off", "GOPROXY=direct", "GONOSUMDB=github.com/junhinhow/*")
	err := cmd.Run()
	code := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			code = exitErr.ExitCode()
		} else {
			t.Fatalf("exec: %v", err)
		}
	}
	return code, out.String(), errOut.String()
}

func TestListContemTodosOsComponentes(t *testing.T) {
	code, stdout, stderr := runGallery(t, "-list")
	if code != 0 {
		t.Fatalf("code=%d stderr=%s", code, stderr)
	}
	names := strings.Fields(stdout)
	if len(names) != len(componentNames()) {
		t.Fatalf("listou %d, want %d", len(names), len(componentNames()))
	}
	seen := map[string]bool{}
	for _, n := range names {
		if seen[n] {
			t.Fatalf("duplicado: %s", n)
		}
		seen[n] = true
	}
	sorted := append([]string(nil), names...)
	sort.Strings(sorted)
	for i := range names {
		if names[i] != sorted[i] {
			t.Fatalf("nao ordenado: %v", names)
		}
	}
}

func TestRenderBanner(t *testing.T) {
	for _, args := range [][]string{
		{"-render", "banner", "-w", "40", "-h", "10"},
		{"-render", "banner", "-w", "80", "-h", "20"},
	} {
		code, stdout, stderr := runGallery(t, args...)
		if code != 0 {
			t.Fatalf("args=%v code=%d stderr=%s", args, code, stderr)
		}
		if strings.TrimSpace(stdout) == "" {
			t.Fatalf("saida vazia para %v", args)
		}
		for _, line := range strings.Split(strings.TrimRight(stdout, "\n"), "\n") {
			if lipgloss.Width(line) > 40 && args[3] == "40" {
				t.Fatalf("linha excede 40: %d %q", lipgloss.Width(line), line)
			}
		}
	}
}

func TestRenderComponenteInexistente(t *testing.T) {
	code, _, stderr := runGallery(t, "-render", "naoexiste", "-w", "40", "-h", "10")
	if code == 0 {
		t.Fatalf("deveria falhar")
	}
	if !strings.Contains(stderr, "nao encontrado") {
		t.Fatalf("stderr: %s", stderr)
	}
}

func TestRenderLargurasDegeneradas(t *testing.T) {
	for _, args := range [][]string{
		{"-render", "banner", "-w", "0", "-h", "0"},
		{"-render", "banner", "-w", "-5", "-h", "-3"},
		{"-render", "banner", "-w", "1", "-h", "1"},
	} {
		code, _, stderr := runGallery(t, args...)
		if code != 0 {
			t.Fatalf("args=%v code=%d stderr=%s", args, code, stderr)
		}
	}
}

func TestRenderTodosOsComponentes(t *testing.T) {
	for _, name := range componentNames() {
		code, stdout, stderr := runGallery(t, "-render", name, "-w", "60", "-h", "15")
		if code != 0 {
			t.Fatalf("%s code=%d stderr=%s", name, code, stderr)
		}
		if strings.TrimSpace(stdout) == "" {
			t.Fatalf("%s saida vazia", name)
		}
	}
}

func TestRenderDeterminismo(t *testing.T) {
	var first string
	for i := 0; i < 5; i++ {
		_, stdout, _ := runGallery(t, "-render", "banner", "-w", "40", "-h", "10")
		if first == "" {
			first = stdout
		} else if stdout != first {
			t.Fatalf("saida nao deterministica na iteracao %d", i)
		}
	}
}
