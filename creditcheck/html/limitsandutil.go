package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/johnllao/remoteproc/creditcheck/arguments"
	"github.com/johnllao/remoteproc/pkg/client"
)

type CoLimitsAndUtil struct {
	Name   string
	Limits string
	Util   string
}

type LimitsAndUtilHandler struct {
	cli          *client.Client
	limsTempl    *template.Template
	limsErrTempl *template.Template
}

func NewLimitsAndUtilHandler(cli *client.Client) *LimitsAndUtilHandler {
	var t, _ = template.New("limits").Parse(limsTempl)
	var errt, _ = template.New("limits_error").Parse(limsErrTempl)
	return &LimitsAndUtilHandler{
		cli:          cli,
		limsTempl:    t,
		limsErrTempl: errt,
	}
}

func (h *LimitsAndUtilHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var symbol = r.PostForm.Get("n")
	if symbol == "" {
		h.limsErrTempl.Execute(w, "symbol is missing from the parameter or arument")
		return
	}
	var args arguments.LimitsAndUtilizationArg
	args.Name = symbol

	var reply arguments.LimitsAndUtilizationReply
	var err = h.cli.Call("CustomerOp.CompanyLimitAndUtilization", &args, &reply)
	if err != nil {
		h.limsErrTempl.Execute(w, reply.ErrorMessage)
		return
	}

	var co CoLimitsAndUtil
	co.Name = symbol
	co.Limits = fmt.Sprintf("%.2f", reply.Limit)
	co.Util = fmt.Sprintf("%.2f", reply.Utilization)
	h.limsTempl.Execute(w, co)
}

var limsTempl = `
<div class="row">
	<h1>Limits and Utilization</h1>
</div>
<div class="row">
	<div class="cell width2">Limits</div>
	<div class="cell width5">{{ .Limits }}</div>
</div>
<div class="row">
	<div class="cell width2">Utilization</div>
	<div class="cell width5">{{ .Util }}</div>
</div>
`

var limsErrTempl = `
<div class="row">
	Error retrieving the company limits. {{ . }}
</div>`
