package libwallet

import (
	"code.cryptopower.dev/group/cryptopower/libwallet/assets/ltc"
	sharedW "code.cryptopower.dev/group/cryptopower/libwallet/assets/wallet"
	"code.cryptopower.dev/group/cryptopower/libwallet/utils"
	"github.com/ltcsuite/ltcd/chaincfg"
)

// initializeLTCWalletParameters initializes the fields each LTC wallet is going to need to be setup
func initializeLTCWalletParameters(netType utils.NetworkType) (*chaincfg.Params, error) {
	chainParams, err := utils.LTCChainParams(netType)
	if err != nil {
		return chainParams, err
	}
	return chainParams, nil
}

// CreateNewLTCWallet creates a new LTC wallet and returns it.
func (mgr *AssetsManager) CreateNewLTCWallet(walletName, privatePassphrase string, privatePassphraseType int32) (sharedW.Asset, error) {
	pass := &sharedW.WalletAuthInfo{
		Name:            walletName,
		PrivatePass:     privatePassphrase,
		PrivatePassType: privatePassphraseType,
	}

	wallet, err := ltc.CreateNewWallet(pass, mgr.params)
	if err != nil {
		return nil, err
	}

	mgr.Assets.LTC.Wallets[wallet.GetWalletID()] = wallet

	// extract the db interface if it hasn't been set already.
	if mgr.db == nil && wallet != nil {
		mgr.setDBInterface(wallet.(sharedW.AssetsManagerDB))
	}

	return wallet, nil
}
