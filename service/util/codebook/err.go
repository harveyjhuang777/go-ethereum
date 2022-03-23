package codebook

import "golang.org/x/xerrors"

var (
	//500
	ErrDatabase = xerrors.New("db error")
	ErrServer   = xerrors.New("server error")

	//400
	ErrDataNotExist   = xerrors.New("data not exist")
	ErrInvalidRequest = xerrors.New("invalid request")
)
