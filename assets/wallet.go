package assets

import "time"

// ManagedWallet contains details of a wallet that is managed by the
// AssetsManager. These details are saved to the AssetsManager database and
// subsequently used to reload the wallet when the AssetsManager starts.
type ManagedWallet struct {
	ID        int       `storm:"id,increment"`
	Name      string    `storm:"unique"`
	CreatedAt time.Time `storm:"index"` // TODO: Why is this indexed??
	Type      AssetType

	// wallet is the main wallet and is assigned a value after the actual wallet
	// is loaded or opened.
	wallet Wallet
}

type Wallet interface {
}
