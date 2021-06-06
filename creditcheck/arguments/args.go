package arguments

import (
	"github.com/johnllao/remoteproc/creditcheck/models"
)

type NilArgs struct{}

type Reply struct {
	Status       int
	ErrorMessage string
}

type UpsertCompaniesArg struct {
	Companies []models.Company
}

type CompaniesReply struct {
	Companies    []models.Company
	Status       int
	ErrorMessage string
}

type LoadFileArg struct {
	Path string
}

type FindCompanyArg struct {
	Name string
}

type FindCompanyReply struct {
	Co           models.Company
	Status       int
	ErrorMessage string
}

type UpdateLimitArg struct {
	Symbol string
	Limit  float64
}

type UpdateUtilizationArg struct {
	Symbol      string
	Utilization float64
}
type LimitsAndUtilizationArg struct {
	Name string
}

type LimitsAndUtilizationReply struct {
	Limit        float64
	Utilization  float64
	Status       int
	ErrorMessage string
}

type BookDealArg struct {
	Deal models.Deal
}
