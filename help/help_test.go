package help

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

func TestGenerateHtml(t *testing.T) {
	dir, _ := os.Getwd()
	file := dir + "/test.md"
	bytes, err := GenerateHtml(file)
	assert.NoError(t, err)
	fmt.Println(string(bytes))
}

func TestHelp_HelpHandler(t *testing.T) {
	dir, _ := os.Getwd()
	file := dir + "/README.md"
	h := &Help{
		MarkdownFile: file,
	}
	mux := http.NewServeMux()
	mux.Handle("/getHelp", http.HandlerFunc(h.HelpHandler))
}
