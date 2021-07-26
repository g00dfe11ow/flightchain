/*
Copyright 2020 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"

	"github.com/fatih/color"
)

func main() {
	log.Println("============ application-golang starts ============")

	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
	}

	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		log.Fatalf("Failed to create wallet: %v", err)
	}

	if !wallet.Exists("approver") {
		err = populateWallet(wallet)
		if err != nil {
			log.Fatalf("Failed to populate wallet contents: %v", err)
		}
	}

	ccpPath := filepath.Join(
		"..",
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"approver.flight.com",
		"connection-approver.yaml",
	)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "approver"),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
	}

	contract := network.GetContract("plan")

	magenta := color.New(color.FgMagenta)

	magenta.Println("--> Submit Transaction: InitLedger, function creates the initial set of plans on the ledger")
	result, err := contract.SubmitTransaction("InitLedger")
	if err != nil {
		log.Fatalf("Failed to Submit transaction: %v", err)
	}
	log.Println(string(result))

	magenta.Println("--> Evaluate Transaction: GetAllPlans, function returns all the current plans on the ledger")
	result, err = contract.EvaluateTransaction("GetAllPlans")
	if err != nil {
		log.Fatalf("Failed to evaluate transaction: %v", err)
	}
	log.Println(string(result))

	magenta.Println("--> Submit Transaction: CreatePlan, creates new plan")
	result, err = contract.SubmitTransaction("CreatePlan", "plan5", "旅游", "18018595217", "{\"nationality\": \"China\", \"type\": \"波音747\", \"number\": 1, \"callsign\": 115, \"registration\": 5}")
	if err != nil {
		log.Fatalf("Failed to Submit transaction: %v", err)
	}
	log.Println(string(result))

	magenta.Println("--> Evaluate Transaction: ReadPlan, function returns an plan with a given planID")
	result, err = contract.EvaluateTransaction("ReadPlan", "plan5")
	if err != nil {
		log.Fatalf("Failed to evaluate transaction: %v\n", err)
	}
	log.Println(string(result))

	magenta.Println("--> Evaluate Transaction: PlanExists, function returns 'true' if an plan with given planID exist")
	result, err = contract.EvaluateTransaction("PlanExists", "plan1")
	if err != nil {
		log.Fatalf("Failed to evaluate transaction: %v\n", err)
	}
	log.Println(string(result))

	magenta.Println("--> Submit Transaction: ApprovalPlan plan1")
	_, err = contract.SubmitTransaction("ApprovalPlan", "plan1", "approved")
	if err != nil {
		log.Fatalf("Failed to Submit transaction: %v", err)
	}

	magenta.Println("--> Evaluate Transaction: ReadPlan, function returns 'plan1' attributes")
	result, err = contract.EvaluateTransaction("ReadPlan", "plan1")
	if err != nil {
		log.Fatalf("Failed to evaluate transaction: %v", err)
	}
	log.Println(string(result))
	log.Println("============ application-golang ends ============")
}

func populateWallet(wallet *gateway.Wallet) error {
	log.Println("============ Populating wallet ============")
	credPath := filepath.Join(
		"..",
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"approver.flight.com",
		"users",
		"User1@approver.flight.com",
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

	identity := gateway.NewX509Identity("ApproverMSP", string(cert), string(key))

	return wallet.Put("approver", identity)
}
