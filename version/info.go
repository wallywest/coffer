package version

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"
	"text/template"

	"github.com/codegangsta/cli"
	"github.com/gin-gonic/gin"
)

var (
	Version   string
	Revision  string
	Branch    string
	BuildUser string
	BuildDate string
)

var GoVersion = fmt.Sprintf("%s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)

var Map = map[string]string{
	"version":   Version,
	"revision":  Revision,
	"branch":    Branch,
	"buildDate": BuildDate,
	"goVersion": GoVersion,
}

var versionInfoTmpl = `
	bifrost, version {{.version}} (branch: {{.branch}}, revision: {{.revision}})
	build date:       {{.buildDate}}
	go version:       {{.goVersion}}
`

func printVersion() {
	t := template.Must(template.New("version").Parse(versionInfoTmpl))
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "version", Map); err != nil {
		panic(err)
	}
	fmt.Fprintln(os.Stdout, strings.TrimSpace(buf.String()))
}

func BuildInfo() func(c *cli.Context) {
	return func(c *cli.Context) {
		printVersion()
	}
}

func Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, Map)
	}
}
