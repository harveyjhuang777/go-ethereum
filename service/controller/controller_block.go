package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/harveyjhuang777/go-ethereum/service/util/codebook"
)

type IBlockController interface {
	GetBlocks(ctx *gin.Context)
	GetBlockDetail(ctx *gin.Context)
}

func newBlockController(in digIn) IBlockController {
	return &blockCtrl{
		in: in,
	}
}

type blockCtrl struct {
	in digIn
}

func (ctl *blockCtrl) GetBlocks(ctx *gin.Context) {
	var (
		limit int
		err   error
	)
	limitStr := ctx.Query("limit")
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			ctl.in.Logger.Error(ctx, err)
			respondError(ctx, codebook.ErrInvalidRequest)
			return
		}
	}

	resp, err := ctl.in.BlockListUseCase.Handle(ctx, limit)
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (ctl *blockCtrl) GetBlockDetail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	if idStr == "" {
		respondError(ctx, codebook.ErrInvalidRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctl.in.Logger.Error(ctx, err)
		respondError(ctx, codebook.ErrInvalidRequest)
		return
	}

	resp, err := ctl.in.BlockDetailUseCase.Handle(ctx, id)
	if err != nil {
		respondError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
