// Package help provides a Go tool of markdown file covert to html
// By：dyy
package help

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gomarkdown/markdown"

	"github.com/gin-gonic/gin"
)

// GetHandlerHelp implements http.ServeHTTP.Interface
func GetHandlerHelp(markdownFile string, mode ...string) http.Handler {
	return &Help{
		MarkdownFile: markdownFile,
		Mode:         mode,
	}
}

// MarkdownFile is the absolute path to markdown
// Response （[]byte） is no parameters
// Debug mode return data or return nil
func GenerateHtml(markdownFile string, mode ...string) ([]byte, error) {

	if len(mode) == 0 {
		if modeStr, ok := os.LookupEnv("MICRO_LOG_LEVEL"); ok {
			mode = append(mode, modeStr)
		}
		if modeStr, ok := os.LookupEnv("MODE"); ok {
			mode = append(mode, modeStr)
		}
		if modeStr, ok := os.LookupEnv("DEBUG"); ok {
			mode = append(mode, modeStr)
		}
		if modeStr, ok := os.LookupEnv("debug"); ok {
			mode = append(mode, modeStr)
		}
		if modeStr, ok := os.LookupEnv("Debug"); ok {
			mode = append(mode, modeStr)
		}
	}
	if len(mode) == 0 || (len(mode) != 0 && strings.ToLower(mode[0]) != "debug") {
		return nil, nil
	}
	f, err := os.Open(markdownFile)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	content, _ := ioutil.ReadAll(f)
	output := markdown.ToHTML(content, nil, nil)
	tmpl, err := template.New("help").Parse(templates(style(), string(output)))
	if err != nil {
		return nil, err
	}
	var doc bytes.Buffer
	tmpl.Execute(&doc, nil)
	return doc.Bytes(), nil
}

// GetHelp provide gin Handler
func GetHelp(markdownFile string, mode ...string) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		bytes, err := GenerateHtml(markdownFile, mode...)
		if err != nil {
			c.JSON(500, nil)
			return
		}
		c.Writer.Write(bytes)
	}
	return gin.HandlerFunc(fn)
}

type Help struct {
	MarkdownFile string
	Mode         []string
}

func (h *Help) ServeHTTP(rep http.ResponseWriter, req *http.Request) {
	h.helpHandler(rep, req)
}
func (h *Help) helpHandler(w http.ResponseWriter, r *http.Request) {
	bytes, err := GenerateHtml(h.MarkdownFile, h.Mode...)
	if err != nil {
		return
	}
	w.Write(bytes)
}
