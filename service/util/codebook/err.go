package codebook

import "golang.org/x/xerrors"

var (
	//500
	ErrDatabase = xerrors.New("db error")
	ErrServer   = xerrors.New("server error")

	//400
	ErrNoData         = xerrors.New("no data")
	ErrInvalidRequest = xerrors.New("invalid request")
	ErrWrongSort      = xerrors.New("sort should be field.orderDirection, orderDirection should be asc or desc")
)
