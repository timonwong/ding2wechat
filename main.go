// Copyright (c) Timon Wong
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package main

import (
	"html/template"
	"net/http"

	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"github.com/timonwong/ding2wechat/config"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	configFile    = kingpin.Flag("config.file", "Path to configuration file.").Default("ding2wechat.yml").String()
	listenAddress = kingpin.Flag("web.listen-address", "Address to listen on for web interface.").Default(":8080").String()
	dryRun        = kingpin.Flag("dry-run", "Only verify configuration is valid and exit.").Default("false").Bool()
)

func main() {
	log.AddFlags(kingpin.CommandLine)
	kingpin.Version(version.Print("ding2wechat"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	log.Infoln("Starting ding2wechat", version.Info())
	log.Infoln("Build context", version.BuildContext())

	cfg, err := config.LoadFile(*configFile)
	if err != nil {
		log.Fatalf("Error parsing config file: %s", err)
	}

	// Exit if in dry-run mode.
	if *dryRun {
		log.Infoln("Configuration parsed successfully.")
		return
	}

	const indexTemplateString = `<html>
		<head><title>DingTalk To WeChat</title></head>
		<body>
			<h1>DingTalk To WeChat</h1>
			<h2>Receivers</h2>
			<ul>
			{{range .Receivers}}<li><code>{{ "{{ url_base }}" }}/receiver?name={{ .Name }}</code></li>{{end}}
			</ul>
		</body>
	</html>`
	indexTpl := template.Must(template.New("index").Parse(indexTemplateString))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := indexTpl.Execute(w, cfg)
		if err != nil {
			log.Errorf("unable to execute template: %s", err)
		}
	})
	http.HandleFunc("/receiver", ReceiverHandler(cfg))

	log.Infof("Listening on %s", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
