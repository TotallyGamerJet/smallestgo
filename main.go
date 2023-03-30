package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Running this command generates a json overlay file necessary to remove the Go runtime.
// It creates the following files: overlay.json, empty.go, notempty.go, empty.s
// Just simply `go run main.go` in whatever directory you want to generate the files.
func main() {
	cmd := exec.Command("go", "list", "-f", "{{ .Dir }}", "runtime")
	b, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	getwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	var (
		emptyGoPath    = filepath.Join(getwd, "empty.go")
		notEmptyGoPath = filepath.Join(getwd, "notempty.go")
		emptySPath     = filepath.Join(getwd, "empty.s")
	)
	{
		empty_go, err := os.Create(emptyGoPath)
		if err != nil {
			return
		}
		empty_go.WriteString("package runtime\n//this file is empty on purpose")
		empty_go.Close()
		notempty_go, err := os.Create(notEmptyGoPath)
		if err != nil {
			return
		}
		notempty_go.WriteString(
			`package runtime

type moduledata struct {
	_ [505]byte
}`)
		notempty_go.Close()
		empty_s, err := os.Create(emptySPath)
		if err != nil {
			return
		}
		empty_s.WriteString("//this file is empty on purpose")
		empty_s.Close()
	}
	path := strings.TrimSpace(string(b))
	dir, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}
	names := make(map[string]string, len(dir))
	for _, entry := range dir {
		if entry.IsDir() {
			continue
		}
		n := entry.Name()
		full := filepath.Join(path, n)
		if strings.HasSuffix(n, ".go") {
			names[full] = emptyGoPath
		} else if strings.HasSuffix(n, ".s") {
			names[full] = emptySPath
		}
	}
	// set the first go source file to the notempty.go file
	for k := range names {
		if !strings.HasSuffix(k, "runtime.go") {
			continue
		}
		names[k] = notEmptyGoPath
		break
	}
	type Overlay struct {
		Replace map[string]string
	}
	marshal, err := json.Marshal(Overlay{Replace: names})
	if err != nil {
		panic(err)
	}
	file, err := os.Create("overlay.json")
	if err != nil {
		panic(err)
	}
	if _, err = file.Write(marshal); err != nil {
		panic(err)
	}
	if err := file.Close(); err != nil {
		panic(err)
	}
}
