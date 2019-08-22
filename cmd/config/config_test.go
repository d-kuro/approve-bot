package config_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/d-kuro/approve-bot/cmd/config"
)

const yaml = `
repo: github.com/d-kuro/approve-bot
owners:
  - name: d-kuro
    patterns:
      - path/to/file
      - regexp
  - name: d-kuro-kuro
    patterns:
      - path/to/file
      - regexp
`

func TestLoadConfigFromFile(t *testing.T) {
	want := &config.ApproveConfig{
		Repo: "github.com/d-kuro/approve-bot",
		Owners: []config.OwnerConfig{
			{
				Name:     "d-kuro",
				Patterns: []string{"path/to/file", "regexp"},
			},
			{
				Name:     "d-kuro-kuro",
				Patterns: []string{"path/to/file", "regexp"},
			},
		},
	}

	f := createTempFile(t)
	defer removeTempFile(t, f)
	got, err := config.LoadConfigFromFile(f)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got: %#v, want: %#v", got, want)
	}
}

func createTempFile(t *testing.T) string {
	t.Helper()
	f, err := ioutil.TempFile("", "config")
	if err != nil {
		t.Fatalf("create temp file error: %v", err)
	}
	defer f.Close()
	if _, err := fmt.Fprint(f, yaml); err != nil {
		t.Fatalf("write to temp file error: %v", err)
	}
	return f.Name()
}

func removeTempFile(t *testing.T, file string) {
	t.Helper()
	if err := os.RemoveAll(file); err != nil {
		t.Fatalf("remove temp file error: %v", err)
	}
}
