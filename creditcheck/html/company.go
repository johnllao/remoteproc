package main

import (
	"html/template"
	"net/http"

	"github.com/johnllao/remoteproc/creditcheck/arguments"
)

type CompanyHandler struct {
	cli      *Client
	coTempl  *template.Template
	errTempl *template.Template
}

func NewCompanyHandler(cli *Client) *CompanyHandler {
	var t, _ = template.New("company").Parse(coTempl)
	var errt, _ = template.New("company_error").Parse(coErrTempl)
	return &CompanyHandler{
		cli:      cli,
		coTempl:  t,
		errTempl: errt,
	}
}

func (h *CompanyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var symbol = r.PostForm.Get("n")
	if symbol == "" {
		h.errTempl.Execute(w, nil)
		return
	}
	var args arguments.FindCompanyArg
	args.Name = symbol

	var reply arguments.FindCompanyReply
	var err = h.cli.Call("CustomerOp.FindCompany", &args, &reply)
	if err != nil {
		h.errTempl.Execute(w, nil)
		return
	}
	var co = reply.Co
	h.coTempl.Execute(w, co)
}

var coTempl = `
<div class="row">
	<h1>{{ .Symbol }}</h1>
</div>
<div class="row">
	<p class="title"><strong>{{ .Name }}</strong></p>
</div>
<div class="row">
	<div class="cell width1">IPO Year</div>
	<div class="cell width5">{{ .IPOYear }}</div>
</div>
<div class="row">
	<div class="cell width1">Country</div>
	<div class="cell width5">{{ .Country }}</div>
</div>
<div class="row">
	<div class="cell width1">Sector</div>
	<div class="cell width5">{{ .Sector }}</div>
</div>
<div class="row">
	<div class="cell width1">Industry</div>
	<div class="cell width5">{{ .Industry }}</div>
</div>	
`

var coErrTempl = `
<div class="row">
	Error retrieving the company symbol
</div>`
