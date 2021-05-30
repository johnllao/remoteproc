package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/johnllao/remoteproc/pkg/client"
)

const (
	token = "t6P8kt4QZpnaDoviGsOLKcfYitR1AjIoHjvF/LxPHC4n/4EDAQEFQ2xhaW0B/4IAAQIBBE5hbWUBDAABBkV4cGlyeQEEAAAAFf+CAQZzYW1wbGUB+C0Ffwv1TVwwAA=="
)

type RootHandler struct {
	csspath string
	jqpath  string

	cli *client.Client

	router http.Handler

	rootTempl *template.Template
}

func (h *RootHandler) Setup() error {
	var err error

	var mux = http.NewServeMux()
	mux.Handle("/css", NewCSSHandler(h.csspath))
	mux.Handle("/jquery", NewJQHandler(h.jqpath))
	mux.Handle("/company", NewCompanyHandler(h.cli))
	mux.Handle("/limits", NewLimitsAndUtilHandler(h.cli))
	h.router = mux

	h.rootTempl, err = template.New("root").Parse(roothtml)
	if err != nil {
		return err
	}
	return nil
}

func (h *RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.router.ServeHTTP(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	h.rootTempl.Execute(w, nil)
}

func main() {
	var err error
	var cli = &client.Client{
		Addr:  "localhost:6060",
		Token: token,
	}
	err = cli.Connect()
	if err != nil {
		log.Fatalf("ERR: main() %s", err.Error())
	}
	defer cli.Close()

	var rootHandler = &RootHandler{
		cli:     cli,
		csspath: "./cchk.css",
		jqpath:  "./jquery.js",
	}
	rootHandler.Setup()
	log.Printf("starting web server")
	log.Printf("css path: %s", rootHandler.csspath)
	log.Printf("jquery path: %s", rootHandler.jqpath)
	http.ListenAndServe("localhost:8080", rootHandler)
}

var roothtml = `<!DOCTYPE html>
<html>
<head>
	<title>Credit Check</title>
	<link rel="stylesheet" type="text/css" href="/css" />
	<script type="text/javascript" src="/jquery"></script>
</head>
<body>
	<div class="container">
		<div class="row">
			<div class="cell width1">Symbol</div>
			<div class="cell width3">
				<input id="findtext" name="findtext" type="text" />
			</div>
			<div class="cell width1">
				<button id="find">Find</button>
			</div>
		</div>
		<div id="codetails" class="row">
		</div>
		<div id="limits" class="row">
		</div>
	</div>
	<script type="text/javascript">
	(function(){
		function onfind() {
			var findtext = $('#findtext').val()
			if (findtext === '') {
				return
			}
			$('#codetails').load('/company', { n: findtext})
			$('#limits').load('/limits', { n: findtext})
		}

		$(document).ready(function(){
			$('#find').click(onfind)	
		})
	})()
	</script>
</body>
</html>`
