package block

import (
	"context"
	"encoding/hex"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"gorm.io/gorm"

	"github.com/harveyjhuang777/go-ethereum/service/model"
)

type IBlockInsert interface {
	Handle(ctx context.Context, sub ethereum.Subscription, headers chan *types.Header)
}

func newBlockInsert(in digIn) IBlockInsert {
	return &blockInsert{
		in: in,
	}
}

type blockInsert struct {
	in digIn
}

func (uc *blockInsert) Handle(ctx context.Context, sub ethereum.Subscription, headers chan *types.Header) {
	db := uc.in.DB.Session()
	for {
		select {
		case err := <-sub.Err():
			uc.in.Logger.Error(ctx, err)
			continue
		case header := <-headers:
			block, err := uc.in.EthApiCli.GetBlockByNumber(context.Background(), header.Number)
			if err != nil {
				uc.in.Logger.Panic(ctx, err)
			}

			tx := func(db *gorm.DB) error {
				if err := uc.insertBlock(ctx, db, block); err != nil {
					uc.in.Logger.Error(ctx, err)
					return err
				}

				for _, txn := range block.Transactions() {
					if err := uc.insertTransaction(ctx, db, txn, block.Number().Uint64()); err != nil {
						uc.in.Logger.Error(ctx, err)
						return err
					}

					rc, err := uc.in.EthApiCli.GetTransactionReceipt(ctx, txn.Hash())
					if err != nil {
						uc.in.Logger.Error(ctx, err)
						return err
					}

					for _, txnLog := range rc.Logs {
						logID := uc.in.IdGenerator.GenerateInt64()
						if err := uc.insertTransactionLog(ctx, db, txnLog, logID, txn.Hash().Hex()); err != nil {
							uc.in.Logger.Error(ctx, err)
							return err
						}
					}
				}

				return nil
			}

			if err := db.Transaction(tx); err != nil {
				uc.in.Logger.Error(ctx, err)
				continue
			}
		}
	}

}

func (uc *blockInsert) insertBlock(ctx context.Context, db *gorm.DB, block *types.Block) error {
	blockToDB := &model.Block{
		Number:     block.Number().Uint64(),
		Hash:       block.Hash().Hex(),
		Time:       block.Time(),
		ParentHash: block.ParentHash().Hex(),
	}

	return uc.in.BlockRepository.Insert(ctx, db, blockToDB)
}

func (uc *blockInsert) insertTransaction(ctx context.Context, db *gorm.DB, txn *types.Transaction, blockID uint64) error {
	txnToDB := &model.Transaction{
		BlockNumber: blockID,
		Hash:        txn.Hash().Hex(),
		Nonce:       txn.Nonce(),
		Data:        hex.EncodeToString(txn.Data()),
	}

	if txn.To() != nil {
		txnToDB.To = txn.To().Hex()
	}

	return uc.in.TransactionRepository.Insert(ctx, db, txnToDB)
}

func (uc *blockInsert) insertTransactionLog(ctx context.Context, db *gorm.DB, txnLog *types.Log, logID int64, txnHash string) error {
	logToDB := &model.TransactionLog{
		ID:              logID,
		TransactionHash: txnHash,
		Index:           txnLog.Index,
		Data:            hex.EncodeToString(txnLog.Data),
	}

	return uc.in.TransactionLogRepository.Insert(ctx, db, logToDB)
}
