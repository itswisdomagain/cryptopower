package assets

import (
	"path/filepath"

	"decred.org/dcrwallet/v3/errors"
	"github.com/asdine/storm"
	bolt "go.etcd.io/bbolt"
)

const (
	walletsDbName          = "wallets.db"
	ErrWalletDatabaseInUse = "wallet_db_in_use"
)

func (mgr *AssetsManager) initializeDatabase() (err error) {
	// Attempt to acquire lock on the wallets.db file.
	mgr.db, err = storm.Open(filepath.Join(mgr.rootDir, walletsDbName))
	if err != nil {
		if err == bolt.ErrTimeout {
			return errors.New(ErrWalletDatabaseInUse)
		}
		return errors.Errorf("error opening wallets database: %s", err.Error())
	}

	// init database for persistence of wallet objects
	if err = mgr.db.Init(&ManagedWallet{}); err != nil {
		return errors.E("error initializing wallets database: %v", err)
	}

	return nil
}

func (mgr *AssetsManager) batchDbTransaction(dbOp func(node storm.Node) error) (err error) {
	dbTx, err := mgr.db.Begin(true)
	if err != nil {
		return err
	}

	// Commit or rollback the transaction after f returns or panics.  Do not
	// recover from the panic to keep the original stack trace intact.
	panicked := true
	defer func() {
		if panicked || err != nil {
			dbTx.Rollback()
			return
		}

		err = dbTx.Commit()
	}()

	err = dbOp(dbTx)
	panicked = false
	return err
}
