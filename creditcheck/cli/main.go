package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/johnllao/remoteproc/creditcheck/arguments"
	"github.com/johnllao/remoteproc/pkg/client"
	"github.com/johnllao/remoteproc/pkg/hmac"
)

func main() {
	var cmd = os.Args[1]

	if cmd == "gentoken" {
		gentoken(os.Args[2:])
	} else if cmd == "list_companies" {
		listCompanies(os.Args[2:])
	} else if cmd == "find_company" {
		findCompany(os.Args[2:])
	} else if cmd == "update_limits" {
		updateLimits(os.Args[2:])
	} else if cmd == "load_file" {
		loadFromFile(os.Args[2:])
	} else {
		fmt.Printf("ERR: main() invalid command argument")
	}
}

func gentoken(args []string) {
	var err error

	var key, name string
	var expiry time.Duration

	var flagset = flag.NewFlagSet("gentoken", flag.ContinueOnError)
	flagset.StringVar(&key, "key", "", "secret key")
	flagset.StringVar(&name, "name", "", "name of the token")
	flagset.DurationVar(&expiry, "expiry", 5, "no of days to expire")
	err = flagset.Parse(args)
	if err != nil {
		fmt.Printf("ERR: gentoken() %s \n", err.Error())
		return
	}

	var claim = &hmac.Claim{
		Name:   name,
		Expiry: time.Now().Add(expiry * 24 * time.Hour).UnixNano(),
	}

	var token string
	token, err = hmac.Token(key, claim)
	if err != nil {
		fmt.Printf("ERR: gentoken() %s \n", err.Error())
		return
	}

	var valid bool
	valid, err = hmac.VerifyToken(key, token)
	if err != nil {
		fmt.Printf("ERR: gentoken() %s \n", err.Error())
		return
	}

	if valid {
		fmt.Println("TOKEN:")
		fmt.Println(token)
	}
}

func listCompanies(args []string) {
	var token = args[0]

	var err error

	var cli = &client.Client{
		Addr:  "localhost:6060",
		Token: token,
	}

	err = cli.Connect()
	if err != nil {
		fmt.Printf("ERR: listCompanies() %s \n", err.Error())
		return
	}
	defer cli.Close()

	var a arguments.NilArgs
	var r arguments.CompaniesReply
	err = cli.Call("CustomerOp.Companies", &a, &r)
	if err != nil {
		fmt.Printf("ERR: listCompanies() %s \n", err.Error())
		return
	}
	for _, co := range r.Companies {
		fmt.Printf("%7s %s \n", co.Symbol, co.Name)
	}
	fmt.Println("bye!")
}

func findCompany(args []string) {
	var token = args[0]
	var symbol = args[1]

	var err error

	var cli = &client.Client{
		Addr:  "localhost:6060",
		Token: token,
	}

	err = cli.Connect()
	if err != nil {
		fmt.Printf("ERR: findCompany() %s \n", err.Error())
		return
	}
	defer cli.Close()

	var a arguments.FindCompanyArg
	var r arguments.FindCompanyReply
	a.Name = symbol
	err = cli.Call("CustomerOp.FindCompany", &a, &r)
	if err != nil {
		fmt.Printf("ERR: findCompany() %s \n", err.Error())
		return
	}
	if r.Status > 0 {
		fmt.Printf("Symbol:   %s \n", r.Co.Symbol)
		fmt.Printf("Name:     %s \n", r.Co.Name)
		fmt.Printf("Industry: %s \n", r.Co.Industry)
		fmt.Printf("Sector:   %s \n", r.Co.Sector)
	}

	var aa arguments.LimitsAndUtilizationArg
	aa.Name = symbol
	var limReply arguments.LimitsAndUtilizationReply
	err = cli.Call("CustomerOp.CompanyLimitAndUtilization", &aa, &limReply)
	if err != nil {
		fmt.Printf("ERR: findCompany() %s \n", err.Error())
		return
	}
	if limReply.Status > 0 {
		fmt.Printf("Limit:       %f \n", limReply.Limit)
		fmt.Printf("Utilization: %f \n", limReply.Utilization)
	}
	fmt.Println("bye!")
}

func updateLimits(args []string) {
	const defaultLimit = 2000000
	const defaultUtilization = 20000

	var token = args[0]
	var symbol = args[1]

	var err error

	var cli = &client.Client{
		Addr:  "localhost:6060",
		Token: token,
	}

	err = cli.Connect()
	if err != nil {
		fmt.Printf("ERR: updateLimits() %s \n", err.Error())
		return
	}
	defer cli.Close()

	var limArg arguments.UpdateLimitArg
	var utilArg arguments.UpdateUtilizationArg
	var r int

	limArg.Symbol = symbol
	limArg.Limit = defaultLimit
	err = cli.Call("CustomerOp.UpdateLimit", &limArg, &r)
	if err != nil {
		fmt.Printf("ERR: updateLimits() %s \n", err.Error())
		return
	}

	utilArg.Symbol = symbol
	utilArg.Utilization = defaultUtilization
	err = cli.Call("CustomerOp.UpdateUtilization", &utilArg, &r)
	if err != nil {
		fmt.Printf("ERR: updateLimits() %s \n", err.Error())
		return
	}

}

func loadFromFile(args []string) {
	var err error

	var token = args[0]

	var cli = &client.Client{
		Addr:  "localhost:6060",
		Token: token,
	}

	err = cli.Connect()
	if err != nil {
		fmt.Printf("ERR: testmsg() %s \n", err.Error())
		return
	}
	defer cli.Close()

	var a arguments.LoadFileArg
	a.Path = "./nasdaq_comapnies.csv"

	var r int
	err = cli.Call("CustomerOp.LoadFromFile", &a, &r)
	if err != nil {
		fmt.Printf("ERR: loadFromFile() %s \n", err.Error())
		return
	}
}
