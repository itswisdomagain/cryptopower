package assets

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/crypto-power/cryptopower/libwallet/utils"
)

func walletDataDir(rootDir string, netType utils.NetworkType, assetType AssetType, walletID int) string {
	assetTypeStr := strings.ToLower(string(assetType)) // TODO: chainParams.Name
	return filepath.Join(rootDir, string(netType), assetTypeStr, strconv.Itoa(walletID))
}

func fileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func moveFile(sourcePath, destinationPath string) error {
	if exists, _ := fileExists(sourcePath); exists {
		return os.Rename(sourcePath, destinationPath)
	}
	return nil
}

func backupFile(fileName string, suffix int) (newName string, err error) {
	newName = fileName + ".bak" + strconv.Itoa(suffix)
	exists, err := fileExists(newName)
	if err != nil {
		return "", err
	} else if exists {
		return backupFile(fileName, suffix+1)
	}

	err = moveFile(fileName, newName)
	if err != nil {
		return "", err
	}

	return newName, nil
}
