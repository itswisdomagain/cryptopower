package assets

import "github.com/crypto-power/cryptopower/libwallet/utils"

type AssetType string

const (
	AssetTypeBTC AssetType = "BTC"
	AssetTypeDCR AssetType = "DCR"
	AssetTypeLTC AssetType = "LTC"
)

type WalletLoader struct {
	// CreateNew generates a new seed or uses the provided seed to create a new
	// wallet for an asset.
	CreateNew       func(netType utils.NetworkType, path, seed string) (Wallet, error)
	CreateWatchOnly func(netType utils.NetworkType, path, xpub string) (Wallet, error)
	OpenExisting    func(netType utils.NetworkType, path string) (Wallet, error)
}

// Asset is a digital currency. An asset can have multiple wallets but only one
// loader.
type Asset struct {
	loader  *WalletLoader
	wallets map[int]*ManagedWallet
	// badWallets map[int]Wallet
}
