package di

import tx_manager "github.com/kirillVladov/account-service/pkg/tx"

func (di *DI) TxManager() *tx_manager.TxManager {
	return tx_manager.NewTxManager(di.Database())
}
