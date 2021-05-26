package ops

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/boltdb/bolt"

	"github.com/johnllao/remoteproc/creditcheck/arguments"
	"github.com/johnllao/remoteproc/creditcheck/models"
	"github.com/johnllao/remoteproc/creditcheck/repositories"
)

type CustomerOp struct {
	repo *repositories.Repository
}

func NewCustomerOps(db *bolt.DB) *CustomerOp {
	return &CustomerOp{
		repo: repositories.NewRepository(db),
	}
}

func (o *CustomerOp) Ping(args *arguments.NilArgs, reply *int) error {
	*reply = 1
	return nil
}

func (o CustomerOp) UpsertCompanies(a *arguments.UpsertCompaniesArg, r *int) error {
	var err = o.repo.SaveCompanies(a.Companies)
	if err != nil {
		*r = -1
		log.Printf("WARN: UpsertCompanies() %s", err.Error())
		return err
	}
	*r = 1
	return nil
}

func (o CustomerOp) Companies(a *arguments.NilArgs, r *arguments.CompaniesReply) error {
	var err error
	var companies []models.Company
	companies, err = o.repo.Companies()
	if err != nil {
		r.Status = -1
		log.Printf("WARN: Companies() %s", err.Error())
		return err
	}
	r.Status = 1
	r.Companies = companies
	return nil
}

func (o CustomerOp) FindCompany(a *arguments.FindCompanyArg, r *arguments.FindCompanyReply) error {
	var err error
	var co *models.Company
	co, err = o.repo.FindCompany(a.Name)
	if co == nil && err == nil {
		r.Status = 0
		return nil
	}
	r.Status = 1
	r.Co = *co
	return nil
}

func (o CustomerOp) LoadFromFile(a *arguments.LoadFileArg, r *int) error {
	var err error

	*r = 1

	var filer *os.File
	filer, err = os.Open(a.Path)
	if err != nil {
		*r = -1
		return err
	}
	var companies = make([]models.Company, 0)
	var csvr = csv.NewReader(filer)
	for {
		var record []string
		record, err = csvr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			*r = -1
			return err
		}
		var co models.Company
		co.Symbol = record[0]
		co.Name = record[1]
		co.LastSale = record[2]
		co.NetChange = 0
		if netchg, err := strconv.ParseFloat(record[3], 64); err == nil {
			co.NetChange = netchg
		}
		co.PercentChange = record[4]
		co.MarketCap = 0
		if mktCap, err := strconv.ParseFloat(record[5], 64); err == nil {
			co.MarketCap = mktCap
		}
		co.Country = record[6]
		co.IPOYear = 0
		if ipoyr, err := strconv.ParseInt(record[7], 10, 32); err == nil {
			co.IPOYear = int(ipoyr)
		}
		co.Volume = 0
		if vol, err := strconv.ParseInt(record[8], 10, 32); err == nil {
			co.Volume = int(vol)
		}
		co.Sector = record[9]
		co.Industry = record[10]
		companies = append(companies, co)
	}
	err = o.repo.SaveCompanies(companies)
	if err != nil {
		*r = -1
		return err
	}
	return nil
}

func (o CustomerOp) CompanyLimitAndUtilization(a *arguments.FindCompanyArg, r *arguments.LimitsAndUtilizationReply) error {
	var err error
	var lim, util float64
	lim, util, err = o.repo.CompanyLimitsAndUtilization(a.Name)
	if err != nil {
		r.Limit = 0
		r.Utilization = 0
		r.Status = -1
		return err
	}
	r.Limit = lim
	r.Utilization = util
	r.Status = 1
	return nil
}

func (o CustomerOp) UpdateLimit(a *arguments.UpdateLimitArg, r *int) error {
	var err = o.repo.UpdateCompanyLimit(a.Symbol, a.Limit)
	if err != nil {
		*r = -1
		return err
	}
	*r = 1
	return nil
}

func (o CustomerOp) UpdateUtilization(a *arguments.UpdateUtilizationArg, r *int) error {
	var err = o.repo.UpdateCompanyUtilization(a.Symbol, a.Utilization)
	if err != nil {
		*r = -1
		return err
	}
	*r = 1
	return nil
}

func (o CustomerOp) BookDeal(a *arguments.BookDealArg, r *int) error {
	var err = o.repo.BookDeal(a.Deal)
	if err != nil {
		*r = -1
		return err
	}
	*r = 1
	return nil
}
