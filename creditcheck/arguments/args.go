package arguments

import (
	"github.com/johnllao/remoteproc/creditcheck/models"
)

type NilArgs struct{}

type CompaniesArg struct {
	Companies []models.Company
}

type CompaniesReply struct {
	Companies []models.Company
	Status    int
}

type LoadFileArg struct {
	Path string
}
