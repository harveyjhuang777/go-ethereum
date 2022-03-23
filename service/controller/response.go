package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"

	"github.com/harveyjhuang777/go-ethereum/service/util/codebook"
)

func respondError(ctx *gin.Context, err error) {
	if multipleErrorsIs(err, codebook.ErrDatabase, codebook.ErrServer) {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusBadRequest, err.Error())
}

func multipleErrorsIs(target error, expects ...error) bool {
	for _, err := range expects {
		if xerrors.Is(err, target) {
			return true
		}
	}

	return false
}
