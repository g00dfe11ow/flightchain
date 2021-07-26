#!/bin/bash

# To create Org applicant
function createApplicant() {
  infoln "Enrolling the CA admin"
  mkdir -p organizations/peerOrganizations/applicant.flight.com/

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/applicant.flight.com/

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:7054 --caname ca-applicant --tls.certfiles "${PWD}/organizations/fabric-ca/applicant/tls-cert.pem"
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-applicant.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-applicant.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-applicant.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-applicant.pem
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/peerOrganizations/applicant.flight.com/msp/config.yaml"

  infoln "Registering peer0"
  set -x
  fabric-ca-client register --caname ca-applicant --id.name peer0 --id.secret peer0pw --id.type peer --tls.certfiles "${PWD}/organizations/fabric-ca/applicant/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering user"
  set -x
  fabric-ca-client register --caname ca-applicant --id.name user1 --id.secret user1pw --id.type client --tls.certfiles "${PWD}/organizations/fabric-ca/applicant/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering the org admin"
  set -x
  fabric-ca-client register --caname ca-applicant --id.name applicantadmin --id.secret applicantadminpw --id.type admin --tls.certfiles "${PWD}/organizations/fabric-ca/applicant/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Generating the peer0 msp"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:7054 --caname ca-applicant -M "${PWD}/organizations/peerOrganizations/applicant.flight.com/peers/peer0.applicant.flight.com/msp" --csr.hosts peer0.applicant.flight.com --tls.certfiles "${PWD}/organizations/fabric-ca/applicant/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/applicant.flight.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/applicant.flight.com/peers/peer0.applicant.flight.com/msp/config.yaml"

  infoln "Generating the peer0-tls certificates"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:7054 --caname ca-applicant -M "${PWD}/organizations/peerOrganizations/applicant.flight.com/peers/peer0.applicant.flight.com/tls" --enrollment.profile tls --csr.hosts peer0.applicant.flight.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/applicant/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/applicant.flight.com/peers/peer0.applicant.flight.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/applicant.flight.com/peers/peer0.applicant.flight.com/tls/ca.crt"
  cp "${PWD}/organizations/peerOrganizations/applicant.flight.com/peers/peer0.applicant.flight.com/tls/signcerts/"* "${PWD}/organizations/peerOrganizations/applicant.flight.com/peers/peer0.applicant.flight.com/tls/server.crt"
  cp "${PWD}/organizations/peerOrganizations/applicant.flight.com/peers/peer0.applicant.flight.com/tls/keystore/"* "${PWD}/organizations/peerOrganizations/applicant.flight.com/peers/peer0.applicant.flight.com/tls/server.key"

  mkdir -p "${PWD}/organizations/peerOrganizations/applicant.flight.com/msp/tlscacerts"
  cp "${PWD}/organizations/peerOrganizations/applicant.flight.com/peers/peer0.applicant.flight.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/applicant.flight.com/msp/tlscacerts/ca.crt"

  mkdir -p "${PWD}/organizations/peerOrganizations/applicant.flight.com/tlsca"
  cp "${PWD}/organizations/peerOrganizations/applicant.flight.com/peers/peer0.applicant.flight.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/applicant.flight.com/tlsca/tlsca.applicant.flight.com-cert.pem"

  mkdir -p "${PWD}/organizations/peerOrganizations/applicant.flight.com/ca"
  cp "${PWD}/organizations/peerOrganizations/applicant.flight.com/peers/peer0.applicant.flight.com/msp/cacerts/"* "${PWD}/organizations/peerOrganizations/applicant.flight.com/ca/ca.applicant.flight.com-cert.pem"

  infoln "Generating the user msp"
  set -x
  fabric-ca-client enroll -u https://user1:user1pw@localhost:7054 --caname ca-applicant -M "${PWD}/organizations/peerOrganizations/applicant.flight.com/users/User1@applicant.flight.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/applicant/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/applicant.flight.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/applicant.flight.com/users/User1@applicant.flight.com/msp/config.yaml"

  infoln "Generating the org admin msp"
  set -x
  fabric-ca-client enroll -u https://applicantadmin:applicantadminpw@localhost:7054 --caname ca-applicant -M "${PWD}/organizations/peerOrganizations/applicant.flight.com/users/Admin@applicant.flight.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/applicant/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/applicant.flight.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/applicant.flight.com/users/Admin@applicant.flight.com/msp/config.yaml"
}

# To create Org approver
function createApprover() {
  infoln "Enrolling the CA admin"
  mkdir -p organizations/peerOrganizations/approver.flight.com/

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/approver.flight.com/

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:8054 --caname ca-approver --tls.certfiles "${PWD}/organizations/fabric-ca/approver/tls-cert.pem"
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-8054-ca-approver.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-8054-ca-approver.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-8054-ca-approver.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-8054-ca-approver.pem
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/peerOrganizations/approver.flight.com/msp/config.yaml"

  infoln "Registering peer0"
  set -x
  fabric-ca-client register --caname ca-approver --id.name peer0 --id.secret peer0pw --id.type peer --tls.certfiles "${PWD}/organizations/fabric-ca/approver/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering user"
  set -x
  fabric-ca-client register --caname ca-approver --id.name user1 --id.secret user1pw --id.type client --tls.certfiles "${PWD}/organizations/fabric-ca/approver/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering the org admin"
  set -x
  fabric-ca-client register --caname ca-approver --id.name approveradmin --id.secret approveradminpw --id.type admin --tls.certfiles "${PWD}/organizations/fabric-ca/approver/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Generating the peer0 msp"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:8054 --caname ca-approver -M "${PWD}/organizations/peerOrganizations/approver.flight.com/peers/peer0.approver.flight.com/msp" --csr.hosts peer0.approver.flight.com --tls.certfiles "${PWD}/organizations/fabric-ca/approver/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/approver.flight.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/approver.flight.com/peers/peer0.approver.flight.com/msp/config.yaml"

  infoln "Generating the peer0-tls certificates"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:8054 --caname ca-approver -M "${PWD}/organizations/peerOrganizations/approver.flight.com/peers/peer0.approver.flight.com/tls" --enrollment.profile tls --csr.hosts peer0.approver.flight.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/approver/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/approver.flight.com/peers/peer0.approver.flight.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/approver.flight.com/peers/peer0.approver.flight.com/tls/ca.crt"
  cp "${PWD}/organizations/peerOrganizations/approver.flight.com/peers/peer0.approver.flight.com/tls/signcerts/"* "${PWD}/organizations/peerOrganizations/approver.flight.com/peers/peer0.approver.flight.com/tls/server.crt"
  cp "${PWD}/organizations/peerOrganizations/approver.flight.com/peers/peer0.approver.flight.com/tls/keystore/"* "${PWD}/organizations/peerOrganizations/approver.flight.com/peers/peer0.approver.flight.com/tls/server.key"

  mkdir -p "${PWD}/organizations/peerOrganizations/approver.flight.com/msp/tlscacerts"
  cp "${PWD}/organizations/peerOrganizations/approver.flight.com/peers/peer0.approver.flight.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/approver.flight.com/msp/tlscacerts/ca.crt"

  mkdir -p "${PWD}/organizations/peerOrganizations/approver.flight.com/tlsca"
  cp "${PWD}/organizations/peerOrganizations/approver.flight.com/peers/peer0.approver.flight.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/approver.flight.com/tlsca/tlsca.approver.flight.com-cert.pem"

  mkdir -p "${PWD}/organizations/peerOrganizations/approver.flight.com/ca"
  cp "${PWD}/organizations/peerOrganizations/approver.flight.com/peers/peer0.approver.flight.com/msp/cacerts/"* "${PWD}/organizations/peerOrganizations/approver.flight.com/ca/ca.approver.flight.com-cert.pem"

  infoln "Generating the user msp"
  set -x
  fabric-ca-client enroll -u https://user1:user1pw@localhost:8054 --caname ca-approver -M "${PWD}/organizations/peerOrganizations/approver.flight.com/users/User1@approver.flight.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/approver/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/approver.flight.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/approver.flight.com/users/User1@approver.flight.com/msp/config.yaml"

  infoln "Generating the org admin msp"
  set -x
  fabric-ca-client enroll -u https://approveradmin:approveradminpw@localhost:8054 --caname ca-approver -M "${PWD}/organizations/peerOrganizations/approver.flight.com/users/Admin@approver.flight.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/approver/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/approver.flight.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/approver.flight.com/users/Admin@approver.flight.com/msp/config.yaml"
}

# To create Org user
function createUser() {
  infoln "Enrolling the CA admin"
  mkdir -p organizations/peerOrganizations/user.flight.com/

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/user.flight.com/

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:9054 --caname ca-user --tls.certfiles "${PWD}/organizations/fabric-ca/user/tls-cert.pem"
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-user.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-user.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-user.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-user.pem
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/peerOrganizations/user.flight.com/msp/config.yaml"

  infoln "Registering peer0"
  set -x
  fabric-ca-client register --caname ca-user --id.name peer0 --id.secret peer0pw --id.type peer --tls.certfiles "${PWD}/organizations/fabric-ca/user/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering user"
  set -x
  fabric-ca-client register --caname ca-user --id.name user1 --id.secret user1pw --id.type client --tls.certfiles "${PWD}/organizations/fabric-ca/user/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering the org admin"
  set -x
  fabric-ca-client register --caname ca-user --id.name useradmin --id.secret useradminpw --id.type admin --tls.certfiles "${PWD}/organizations/fabric-ca/user/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Generating the peer0 msp"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:9054 --caname ca-user -M "${PWD}/organizations/peerOrganizations/user.flight.com/peers/peer0.user.flight.com/msp" --csr.hosts peer0.user.flight.com --tls.certfiles "${PWD}/organizations/fabric-ca/user/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/user.flight.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/user.flight.com/peers/peer0.user.flight.com/msp/config.yaml"

  infoln "Generating the peer0-tls certificates"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:9054 --caname ca-user -M "${PWD}/organizations/peerOrganizations/user.flight.com/peers/peer0.user.flight.com/tls" --enrollment.profile tls --csr.hosts peer0.user.flight.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/user/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/user.flight.com/peers/peer0.user.flight.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/user.flight.com/peers/peer0.user.flight.com/tls/ca.crt"
  cp "${PWD}/organizations/peerOrganizations/user.flight.com/peers/peer0.user.flight.com/tls/signcerts/"* "${PWD}/organizations/peerOrganizations/user.flight.com/peers/peer0.user.flight.com/tls/server.crt"
  cp "${PWD}/organizations/peerOrganizations/user.flight.com/peers/peer0.user.flight.com/tls/keystore/"* "${PWD}/organizations/peerOrganizations/user.flight.com/peers/peer0.user.flight.com/tls/server.key"

  mkdir -p "${PWD}/organizations/peerOrganizations/user.flight.com/msp/tlscacerts"
  cp "${PWD}/organizations/peerOrganizations/user.flight.com/peers/peer0.user.flight.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/user.flight.com/msp/tlscacerts/ca.crt"

  mkdir -p "${PWD}/organizations/peerOrganizations/user.flight.com/tlsca"
  cp "${PWD}/organizations/peerOrganizations/user.flight.com/peers/peer0.user.flight.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/user.flight.com/tlsca/tlsca.user.flight.com-cert.pem"

  mkdir -p "${PWD}/organizations/peerOrganizations/user.flight.com/ca"
  cp "${PWD}/organizations/peerOrganizations/user.flight.com/peers/peer0.user.flight.com/msp/cacerts/"* "${PWD}/organizations/peerOrganizations/user.flight.com/ca/ca.user.flight.com-cert.pem"

  infoln "Generating the user msp"
  set -x
  fabric-ca-client enroll -u https://user1:user1pw@localhost:9054 --caname ca-user -M "${PWD}/organizations/peerOrganizations/user.flight.com/users/User1@user.flight.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/user/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/user.flight.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/user.flight.com/users/User1@user.flight.com/msp/config.yaml"

  infoln "Generating the org admin msp"
  set -x
  fabric-ca-client enroll -u https://useradmin:useradminpw@localhost:9054 --caname ca-user -M "${PWD}/organizations/peerOrganizations/user.flight.com/users/Admin@user.flight.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/user/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/peerOrganizations/user.flight.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/user.flight.com/users/Admin@user.flight.com/msp/config.yaml"
}

# To create Orderer
function createOrderer() {
  infoln "Enrolling the CA admin"
  mkdir -p organizations/ordererOrganizations/flight.com

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/ordererOrganizations/flight.com

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:10054 --caname ca-orderer --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/tls-cert.pem"
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-10054-ca-orderer.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-10054-ca-orderer.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-10054-ca-orderer.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-10054-ca-orderer.pem
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/ordererOrganizations/flight.com/msp/config.yaml"

  infoln "Registering orderer"
  set -x
  fabric-ca-client register --caname ca-orderer --id.name orderer --id.secret ordererpw --id.type orderer --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Registering the orderer admin"
  set -x
  fabric-ca-client register --caname ca-orderer --id.name ordererAdmin --id.secret ordererAdminpw --id.type admin --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/tls-cert.pem"
  { set +x; } 2>/dev/null

  infoln "Generating the orderer msp"
  set -x
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:10054 --caname ca-orderer -M "${PWD}/organizations/ordererOrganizations/flight.com/orderers/orderer.flight.com/msp" --csr.hosts orderer.flight.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/ordererOrganizations/flight.com/msp/config.yaml" "${PWD}/organizations/ordererOrganizations/flight.com/orderers/orderer.flight.com/msp/config.yaml"

  infoln "Generating the orderer-tls certificates"
  set -x
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:10054 --caname ca-orderer -M "${PWD}/organizations/ordererOrganizations/flight.com/orderers/orderer.flight.com/tls" --enrollment.profile tls --csr.hosts orderer.flight.com --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/ordererOrganizations/flight.com/orderers/orderer.flight.com/tls/tlscacerts/"* "${PWD}/organizations/ordererOrganizations/flight.com/orderers/orderer.flight.com/tls/ca.crt"
  cp "${PWD}/organizations/ordererOrganizations/flight.com/orderers/orderer.flight.com/tls/signcerts/"* "${PWD}/organizations/ordererOrganizations/flight.com/orderers/orderer.flight.com/tls/server.crt"
  cp "${PWD}/organizations/ordererOrganizations/flight.com/orderers/orderer.flight.com/tls/keystore/"* "${PWD}/organizations/ordererOrganizations/flight.com/orderers/orderer.flight.com/tls/server.key"

  mkdir -p "${PWD}/organizations/ordererOrganizations/flight.com/orderers/orderer.flight.com/msp/tlscacerts"
  cp "${PWD}/organizations/ordererOrganizations/flight.com/orderers/orderer.flight.com/tls/tlscacerts/"* "${PWD}/organizations/ordererOrganizations/flight.com/orderers/orderer.flight.com/msp/tlscacerts/tlsca.flight.com-cert.pem"

  mkdir -p "${PWD}/organizations/ordererOrganizations/flight.com/msp/tlscacerts"
  cp "${PWD}/organizations/ordererOrganizations/flight.com/orderers/orderer.flight.com/tls/tlscacerts/"* "${PWD}/organizations/ordererOrganizations/flight.com/msp/tlscacerts/tlsca.flight.com-cert.pem"

  infoln "Generating the admin msp"
  set -x
  fabric-ca-client enroll -u https://ordererAdmin:ordererAdminpw@localhost:10054 --caname ca-orderer -M "${PWD}/organizations/ordererOrganizations/flight.com/users/Admin@flight.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/tls-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/ordererOrganizations/flight.com/msp/config.yaml" "${PWD}/organizations/ordererOrganizations/flight.com/users/Admin@flight.com/msp/config.yaml"
}
