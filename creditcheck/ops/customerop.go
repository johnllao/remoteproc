package ops

import (
	"log"

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

func (o CustomerOp) UpsertCompanies(a *arguments.CompaniesArg, r *int) error {
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
	var companies []*models.Company
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
