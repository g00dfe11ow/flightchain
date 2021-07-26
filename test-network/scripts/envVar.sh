#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#

# This is a collection of bash functions used by different scripts

# imports
. scripts/utils.sh

export CORE_PEER_TLS_ENABLED=true
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/flight.com/orderers/orderer.flight.com/msp/tlscacerts/tlsca.flight.com-cert.pem
export PEER0_APPLICANT_CA=${PWD}/organizations/peerOrganizations/applicant.flight.com/peers/peer0.applicant.flight.com/tls/ca.crt
export PEER0_APPROVER_CA=${PWD}/organizations/peerOrganizations/approver.flight.com/peers/peer0.approver.flight.com/tls/ca.crt
export PEER0_USER_CA=${PWD}/organizations/peerOrganizations/user.flight.com/peers/peer0.user.flight.com/tls/ca.crt
export ORDERER_ADMIN_TLS_SIGN_CERT=${PWD}/organizations/ordererOrganizations/flight.com/orderers/orderer.flight.com/tls/server.crt
export ORDERER_ADMIN_TLS_PRIVATE_KEY=${PWD}/organizations/ordererOrganizations/flight.com/orderers/orderer.flight.com/tls/server.key

# Set environment variables for the peer org
setGlobals() {
  local USING_ORG=""
  if [ -z "$OVERRIDE_ORG" ]; then
    USING_ORG=$1
  else
    USING_ORG="${OVERRIDE_ORG}"
  fi
  infoln "Using organization ${USING_ORG}"
  if [ $USING_ORG = "applicant" ]; then
    export CORE_PEER_LOCALMSPID="ApplicantMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_APPLICANT_CA
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/applicant.flight.com/users/Admin@applicant.flight.com/msp
    export CORE_PEER_ADDRESS=localhost:7051
  elif [ $USING_ORG = "approver" ]; then
    export CORE_PEER_LOCALMSPID="ApproverMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_APPROVER_CA
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/approver.flight.com/users/Admin@approver.flight.com/msp
    export CORE_PEER_ADDRESS=localhost:8051

  elif [ $USING_ORG = "user" ]; then
    export CORE_PEER_LOCALMSPID="UserMSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_USER_CA
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/user.flight.com/users/Admin@user.flight.com/msp
    export CORE_PEER_ADDRESS=localhost:9051
  else
    errorln "ORG Unknown"
  fi

  if [ "$VERBOSE" == "true" ]; then
    env | grep CORE
  fi
}

# Set environment variables for use in the CLI container 
setGlobalsCLI() {
  setGlobals $1

  local USING_ORG=""
  if [ -z "$OVERRIDE_ORG" ]; then
    USING_ORG=$1
  else
    USING_ORG="${OVERRIDE_ORG}"
  fi
  if [ $USING_ORG = "applicant" ]; then
    export CORE_PEER_ADDRESS=peer0.applicant.flight.com:7051
  elif [ $USING_ORG = "approver" ]; then
    export CORE_PEER_ADDRESS=peer0.approver.flight.com:8051
  elif [ $USING_ORG = "user" ]; then
    export CORE_PEER_ADDRESS=peer0.user.flight.com:9051
  else
    errorln "ORG Unknown"
  fi
}

# parsePeerConnectionParameters $@
# Helper function that sets the peer connection parameters for a chaincode
# operation
parsePeerConnectionParameters() {
  PEER_CONN_PARMS=()
  PEERS=""
  while [ "$#" -gt 0 ]; do
    setGlobals $1
    PEER="peer0.$1"
    ## Set peer addresses
    if [ -z "$PEERS" ]
    then
      PEERS="$PEER"
    else
	    PEERS="$PEERS $PEER"
    fi
    PEER_CONN_PARMS=("${PEER_CONN_PARMS[@]}" --peerAddresses $CORE_PEER_ADDRESS)
    ## Set path to TLS certificate
    typeset -u upperOrg
    upperOrg=$1
    CA=PEER0_${upperOrg}_CA
    TLSINFO=(--tlsRootCertFiles "${!CA}")
    PEER_CONN_PARMS=("${PEER_CONN_PARMS[@]}" "${TLSINFO[@]}")
    # shift by one to get to the next organization
    shift
  done
}

verifyResult() {
  if [ $1 -ne 0 ]; then
    fatalln "$2"
  fi
}
