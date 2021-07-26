package blockchain

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

var (
	contract *gateway.Contract
)

func Init() {

	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
	}

	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		log.Fatalf("Failed to create wallet: %v", err)
	}

	if !wallet.Exists("applicant") {
		err = populateWallet(wallet)
		if err != nil {
			log.Fatalf("Failed to populate wallet contents: %v", err)
		}
	}

	ccpPath := filepath.Join(
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"applicant.flight.com",
		"connection-applicant.yaml",
	)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "applicant"),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
	}

	contract = network.GetContract("plan")
}

func populateWallet(wallet *gateway.Wallet) error {

	credPath := filepath.Join(
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"applicant.flight.com",
		"users",
		"User1@applicant.flight.com",
		"msp",
	)

	certPath := filepath.Join(credPath, "signcerts", "cert.pem")
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return fmt.Errorf("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity("ApplicantMSP", string(cert), string(key))

	return wallet.Put("applicant", identity)
}

// ChannelExecute 区块链交互，会改变账本的操作
func ChannelExecute(fcn string, args []string) (string, error) {
	result, err := contract.SubmitTransaction(fcn, args...)
	if err != nil {
		log.Fatalf("Failed to Submit transaction: %v", err)
		return err.Error(), err
	}
	return string(result), nil
}

// ChannelQuery 区块链交互，不会改变账本的操作
func ChannelQuery(fcn string, args []string) (string, error) {
	result, err := contract.EvaluateTransaction(fcn, args...)
	if err != nil {
		log.Fatalf("Failed to Evaluate transaction: %v", err)
		return err.Error(), err
	}
	return string(result), nil
}

// func GetAllPlans() (string, error) {
// 	result, err := contract.EvaluateTransaction("GetAllPlans")
// 	if err != nil {
// 		log.Fatalf("Failed to Evaluate transaction: %v", err)
// 		return err.Error(), err
// 	}
// 	return string(result), nil
// }
