package libwallet

import (
	"fmt"

	"code.cryptopower.dev/group/cryptopower/libwallet/utils"
	"decred.org/dcrwallet/v2/errors"
	"github.com/asdine/storm"

	sharedW "code.cryptopower.dev/group/cryptopower/libwallet/assets/wallet"

	"golang.org/x/crypto/bcrypt"
)

const (
	logFileName   = "libwallet.log"
	walletsDbName = "wallets.db"

	Mainnet  = utils.Mainnet
	Testnet3 = utils.Testnet

	walletsMetadataBucketName    = "metadata"
	walletstartupPassphraseField = "startup-passphrase"
)

// setDBInterface extract the assets manager db interface that is available
// in each wallet by default from one of the validly created wallets.
func (mgr *AssetsManager) setDBInterface(db sharedW.AssetsManagerDB) {
	if db != nil {
		mgr.db = db
	}
}

func (mgr *AssetsManager) SetStartupPassphrase(passphrase string, passphraseType int32) error {
	return mgr.ChangeStartupPassphrase("", passphrase, passphraseType)
}

func (mgr *AssetsManager) VerifyStartupPassphrase(startupPassphrase string) error {
	var startupPassphraseHash []byte
	err := mgr.db.ReadWalletConfigValue(walletstartupPassphraseField, &startupPassphraseHash)
	if err != nil && err != storm.ErrNotFound {
		return err
	}

	if startupPassphraseHash == nil {
		// startup passphrase was not previously set
		if len(startupPassphrase) > 0 {
			return errors.E(utils.ErrInvalidPassphrase)
		}
		return nil
	}

	// startup passphrase was set, verify
	err = bcrypt.CompareHashAndPassword(startupPassphraseHash, []byte(startupPassphrase))
	if err != nil {
		return errors.E(utils.ErrInvalidPassphrase)
	}

	return nil
}

func (mgr *AssetsManager) ChangeStartupPassphrase(oldPassphrase, newPassphrase string, passphraseType int32) error {
	if len(newPassphrase) == 0 {
		return mgr.RemoveStartupPassphrase(oldPassphrase)
	}

	err := mgr.VerifyStartupPassphrase(oldPassphrase)
	if err != nil {
		return err
	}

	startupPassphraseHash, err := bcrypt.GenerateFromPassword([]byte(newPassphrase), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	mgr.db.SaveWalletConfigValue(walletstartupPassphraseField, startupPassphraseHash)
	mgr.db.SaveWalletConfigValue(sharedW.IsStartupSecuritySetConfigKey, true)
	mgr.db.SaveWalletConfigValue(sharedW.StartupSecurityTypeConfigKey, passphraseType)

	return nil
}

func (mgr *AssetsManager) RemoveStartupPassphrase(oldPassphrase string) error {
	err := mgr.VerifyStartupPassphrase(oldPassphrase)
	if err != nil {
		return err
	}

	mgr.db.DeleteWalletConfigValue(walletstartupPassphraseField)
	mgr.db.SaveWalletConfigValue(sharedW.IsStartupSecuritySetConfigKey, false)
	mgr.db.DeleteWalletConfigValue(sharedW.StartupSecurityTypeConfigKey)

	return nil
}

func (mgr *AssetsManager) IsStartupSecuritySet() bool {
	var data bool
	mgr.db.ReadWalletConfigValue(sharedW.IsStartupSecuritySetConfigKey, &data)
	return data
}

func (mgr *AssetsManager) IsDarkModeOn() bool {
	var data bool
	mgr.db.ReadWalletConfigValue(sharedW.DarkModeConfigKey, &data)
	return data
}

func (mgr *AssetsManager) SetDarkMode(data bool) {
	mgr.db.SaveWalletConfigValue(sharedW.DarkModeConfigKey, data)
}

func (mgr *AssetsManager) GetDexServers() (map[string][]byte, error) {
	var servers = make(map[string][]byte, 0)
	err := mgr.db.ReadWalletConfigValue(sharedW.KnownDexServersConfigKey, &servers)
	return servers, err
}

func (mgr *AssetsManager) SaveDexServers(servers map[string][]byte) {
	mgr.db.SaveWalletConfigValue(sharedW.KnownDexServersConfigKey, servers)
}

func (mgr *AssetsManager) GetCurrencyConversionExchange() string {
	var key string
	mgr.db.ReadWalletConfigValue(sharedW.CurrencyConversionConfigKey, &key)
	if key == "" {
		return "none" // default exchange value
	}
	return key
}

func (mgr *AssetsManager) SetCurrencyConversionExchange(data string) {
	mgr.db.SaveWalletConfigValue(sharedW.CurrencyConversionConfigKey, data)
}

func (mgr *AssetsManager) GetLanguagePreference() string {
	var lang string
	mgr.db.ReadWalletConfigValue(sharedW.LanguagePreferenceKey, &lang)
	return lang
}

func (mgr *AssetsManager) SetLanguagePreference(lang string) {
	mgr.db.SaveWalletConfigValue(sharedW.LanguagePreferenceKey, lang)
}

func (mgr *AssetsManager) GetUserAgent() string {
	var data string
	mgr.db.ReadWalletConfigValue(sharedW.UserAgentConfigKey, data)
	return data
}

func (mgr *AssetsManager) SetUserAgent(data string) {
	mgr.db.SaveWalletConfigValue(sharedW.UserAgentConfigKey, data)
}

func (mgr *AssetsManager) IsTransactionNotificationsOn() bool {
	var data bool
	mgr.db.ReadWalletConfigValue(sharedW.TransactionNotificationConfigKey, &data)
	return data
}

func (mgr *AssetsManager) SetTransactionsNotifications(data bool) {
	mgr.db.SaveWalletConfigValue(sharedW.TransactionNotificationConfigKey, data)
}

func (mgr *AssetsManager) SetPrivacyMode(isActive bool) {
	mgr.db.SaveWalletConfigValue(sharedW.PrivacyModeConfigKey, isActive)
}

func (mgr *AssetsManager) GetPrivacyMode() bool {
	var data bool
	mgr.db.ReadWalletConfigValue(sharedW.PrivacyModeConfigKey, data)
	return data
}

func (mgr *AssetsManager) SetHttpAPIPrivacyMode(apiType utils.HttpAPIType, isActive bool) {
	dataKey := genKey(sharedW.PrivacyModeConfigKey, apiType)
	mgr.db.SaveWalletConfigValue(dataKey, isActive)
}

func (mgr *AssetsManager) GetHttpAPIPrivacyMode(apiType utils.HttpAPIType) bool {
	var data bool
	dataKey := genKey(sharedW.PrivacyModeConfigKey, apiType)
	mgr.db.ReadWalletConfigValue(dataKey, &data)
	return data
}

func (mgr *AssetsManager) GetLogLevels() {
	//TODO: loglevels should have a custom type supported on libwallet.
	// Issue is to be addressed in here: https://code.cryptopower.dev/group/cryptopower/-/issues/965
	var logLevel string
	mgr.db.ReadWalletConfigValue(sharedW.LogLevelConfigKey, &logLevel)
	SetLogLevels(logLevel)
}

func (mgr *AssetsManager) SetLogLevels(logLevel string) {
	mgr.db.SaveWalletConfigValue(sharedW.LogLevelConfigKey, logLevel)
	SetLogLevels(logLevel)
}

// SetExchangeConfig sets the exchnage config for the asset
func (mgr *AssetsManager) SetExchangeConfig(fromCurrency utils.AssetType, sourceWalletID int32, toCurrency utils.AssetType, destinationWalletID, sourceAccountID, DestinationAccountID int32) {
	mgr.db.SaveWalletConfigValue(sharedW.ExchangeSourceAssetTypeConfigKey, fromCurrency)
	mgr.db.SaveWalletConfigValue(sharedW.ExchangeDestinationAssetTypeConfigKey, toCurrency)
	mgr.db.SaveWalletConfigValue(sharedW.ExchangeSourceWalletConfigKey, sourceWalletID)
	mgr.db.SaveWalletConfigValue(sharedW.ExchangeSourceAccountConfigKey, sourceAccountID)
	mgr.db.SaveWalletConfigValue(sharedW.ExchangeDestinationWalletConfigKey, destinationWalletID)
	mgr.db.SaveWalletConfigValue(sharedW.ExchangeDestinationAccountConfigKey, DestinationAccountID)
}

// ExchangeConfig returns the previously set exchange config for
// the asset.
func (mgr *AssetsManager) ExchangeConfig() *sharedW.ExchangeConfig {
	var sourceAsset utils.AssetType
	var destinationAsset utils.AssetType
	var sourceWalletID int32
	var destinationWalletID int32
	var sourceAccoutNumber int32
	var destinationAccountNumber int32

	mgr.db.ReadWalletConfigValue(sharedW.ExchangeSourceWalletConfigKey, &sourceWalletID)
	mgr.db.ReadWalletConfigValue(sharedW.ExchangeSourceAssetTypeConfigKey, &sourceAsset)
	mgr.db.ReadWalletConfigValue(sharedW.ExchangeSourceAccountConfigKey, &sourceAccoutNumber)
	mgr.db.ReadWalletConfigValue(sharedW.ExchangeDestinationAssetTypeConfigKey, &destinationAsset)
	mgr.db.ReadWalletConfigValue(sharedW.ExchangeDestinationWalletConfigKey, &destinationWalletID)
	mgr.db.ReadWalletConfigValue(sharedW.ExchangeDestinationAccountConfigKey, &destinationAccountNumber)

	return &sharedW.ExchangeConfig{
		SourceAsset:      sourceAsset,
		DestinationAsset: destinationAsset,

		SourceWalletID:      sourceWalletID,
		DestinationWalletID: destinationWalletID,

		SourceAccountNumber:      sourceAccoutNumber,
		DestinationAccountNumber: destinationAccountNumber,
	}
}

// ExchangeConfigIsSet checks if exchange config is set for the asset.
func (mgr *AssetsManager) ExchangeConfigIsSet() bool {
	var sourceWalletID int32 = -1

	mgr.db.ReadWalletConfigValue(sharedW.ExchangeSourceWalletConfigKey, &sourceWalletID)

	return sourceWalletID != -1
}

// ClearExchangeConfig clears the wallet's exchange config.
func (mgr *AssetsManager) ClearExchangeConfig() error {
	mgr.db.DeleteWalletConfigValue(sharedW.ExchangeSourceAssetTypeConfigKey)
	mgr.db.DeleteWalletConfigValue(sharedW.ExchangeDestinationAssetTypeConfigKey)
	mgr.db.DeleteWalletConfigValue(sharedW.ExchangeSourceWalletConfigKey)
	mgr.db.DeleteWalletConfigValue(sharedW.ExchangeSourceAccountConfigKey)
	mgr.db.DeleteWalletConfigValue(sharedW.ExchangeDestinationWalletConfigKey)
	mgr.db.DeleteWalletConfigValue(sharedW.ExchangeDestinationAccountConfigKey)

	return nil
}

func genKey(prefix, identifier interface{}) string {
	return fmt.Sprintf("%v-%v", prefix, identifier)
}
