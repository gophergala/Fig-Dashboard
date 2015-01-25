// main.go for figdash web tool
// as no API in Fig, built up Fig commands with commend line calls to Fig
//

package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"

	"github.com/codegangsta/cli"
)

var addr = flag.String("addr", ":1984", "http service address")

var templ = template.Must(template.New("dash").Parse(templateStr))

func fixProjectName(c *cli.Context) string {
	if c.String("projectname") == "" {
		var cwd, err = os.Getwd()
		if err == nil {
			var wd = path.Base(cwd)
			fmt.Printf("cwd = %s\n", cwd)
			fmt.Printf("wd = %s\n", wd)
			return wd
		}
		return ""
	}
	return c.String("projectname")
}

func dash(w http.ResponseWriter, req *http.Request) {
	templ.Execute(w, req.FormValue("s"))
}

func main() {
	app := cli.NewApp()
	app.Name = "figdash"
	app.Usage = "fig dashboard"
	app.Version = "0.0.1"
	app.Email = "mkobar@rkosecurity.com"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "Show more output",
		},
		// version flag support is builtin
		//cli.BoolFlag{
		//	Name:  "version",
		//	Usage: "Print version and exit",
		//},
		cli.StringFlag{
			Name:   "file, f",
			Value:  "fig.yml",
			Usage:  "Specify an alternate fig file",
			EnvVar: "FIG_FILE",
		},
		cli.StringFlag{
			Name:   "projectname, p",
			Value:  "notset",
			Usage:  "Specify an alternate project name",
			EnvVar: "FIG_PROJECT_NAME",
		},
	}
	app.Commands = []cli.Command{
		// build - NOT supported yet
		// help - NOT supported yet
		{
			Name:  "kill",
			Usage: "Force stop service containers.",
			Action: func(c *cli.Context) {
				cmd := exec.Command("fig", "kill")
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Run()
			},
		},
		// logs - not supported
		// port - NOT supported
		{
			Name:  "ps",
			Usage: "List containers",
			Action: func(c *cli.Context) {
				var pn = fixProjectName(c)
				if pn != "" {
					fmt.Printf("ProjectName: %s\n", pn)
				} else {
					fmt.Printf("ProjectName: %s\n", "unknown")
				}
				cmd := exec.Command("fig", "ps")
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Run()
			},
		},
		// pull - NOT supported
		{
			Name:  "rm",
			Usage: "Remove stopped service containers.",
			Action: func(c *cli.Context) {
				cmd := exec.Command("fig", "rm")
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Run()
			},
		},
		// run - NOT supported
		// scale - NOT supported
		{
			Name:  "start",
			Usage: "Start existing containers for a service.",
			Action: func(c *cli.Context) {
				cmd := exec.Command("fig", "start")
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Run()
			},
		},
		{
			Name:  "stop",
			Usage: "Stop existing containers without removing them.",
			Action: func(c *cli.Context) {
				cmd := exec.Command("fig", "stop")
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Run()
			},
		},
		// up - NOT supported - due to logs
		{
			Name:  "web",
			Usage: "Enable web monitoring on http://localhost/1984/PROJECTNAME.",
			Action: func(c *cli.Context) {
				flag.Parse()
				var pn = fixProjectName(c)
				if pn == "" {
					pn = "unknown"
				}
				fmt.Printf("Starting web Fig Dashboard at http://localhost:1984/%s\n", pn)
				fmt.Printf("use Ctrl-C to exit\n")
				http.Handle("/"+pn, http.HandlerFunc(dash))
				err := http.ListenAndServe(*addr, nil)
				if err != nil {
					log.Fatal("ListenAndServe:", err)
				}


				//x.Run()
			},
		},
		// up - NOT supported - due to logs

	}
	app.Run(os.Args)
}

const templateStr = `
<html>
<head>
<title>Fig Dashboard</title>
</head>
<body>
<h2>Fig Dashboard v 0.0.1</h2>
<h3> for project </h3>
{{if .}}
<img src="http://chart.apis.google.com/chart?chs=300x300&cht=qr&choe=UTF-8&chl={{.}}" />
<br>
{{.}}
<br>
<br>
{{end}}
<form action="/" name=f method="GET"><input maxLength=1024 size=70
name=s value="" title="Text to QR Encode"><input type=submit
value="Refresh" name=dash>
</form>
</body>
</html>
`
