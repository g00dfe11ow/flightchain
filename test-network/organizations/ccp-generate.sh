#!/bin/bash

function one_line_pem {
    echo "`awk 'NF {sub(/\\n/, ""); printf "%s\\\\\\\n",$0;}' $1`"
}

function json_ccp {
    local PP=$(one_line_pem $4)
    local CP=$(one_line_pem $5)
    sed -e "s/\${ORG}/$1/" \
        -e "s/\${P0PORT}/$2/" \
        -e "s/\${CAPORT}/$3/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${CAPEM}#$CP#" \
        -e "s#\${upperORG}#$6#" \
        organizations/ccp-template.json
}

function yaml_ccp {
    local PP=$(one_line_pem $4)
    local CP=$(one_line_pem $5)
    sed -e "s/\${ORG}/$1/" \
        -e "s/\${P0PORT}/$2/" \
        -e "s/\${CAPORT}/$3/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${CAPEM}#$CP#" \
        -e "s#\${upperORG}#$6#" \
        organizations/ccp-template.yaml | sed -e $'s/\\\\n/\\\n          /g'
}

ORG=applicant
upperORG=Applicant
P0PORT=7051
CAPORT=7054
PEERPEM=organizations/peerOrganizations/applicant.flight.com/tlsca/tlsca.applicant.flight.com-cert.pem
CAPEM=organizations/peerOrganizations/applicant.flight.com/ca/ca.applicant.flight.com-cert.pem

echo "$(json_ccp $ORG $P0PORT $CAPORT $PEERPEM $CAPEM $upperORG)" > organizations/peerOrganizations/applicant.flight.com/connection-applicant.json
echo "$(yaml_ccp $ORG $P0PORT $CAPORT $PEERPEM $CAPEM $upperORG)" > organizations/peerOrganizations/applicant.flight.com/connection-applicant.yaml

ORG=approver
upperORG=Approver
P0PORT=8051
CAPORT=8054
PEERPEM=organizations/peerOrganizations/approver.flight.com/tlsca/tlsca.approver.flight.com-cert.pem
CAPEM=organizations/peerOrganizations/approver.flight.com/ca/ca.approver.flight.com-cert.pem

echo "$(json_ccp $ORG $P0PORT $CAPORT $PEERPEM $CAPEM $upperORG)" > organizations/peerOrganizations/approver.flight.com/connection-approver.json
echo "$(yaml_ccp $ORG $P0PORT $CAPORT $PEERPEM $CAPEM $upperORG)" > organizations/peerOrganizations/approver.flight.com/connection-approver.yaml

ORG=user
upperORG=User
P0PORT=9051
CAPORT=9054
PEERPEM=organizations/peerOrganizations/user.flight.com/tlsca/tlsca.user.flight.com-cert.pem
CAPEM=organizations/peerOrganizations/user.flight.com/ca/ca.user.flight.com-cert.pem

echo "$(json_ccp $ORG $P0PORT $CAPORT $PEERPEM $CAPEM $upperORG)" > organizations/peerOrganizations/user.flight.com/connection-user.json
echo "$(yaml_ccp $ORG $P0PORT $CAPORT $PEERPEM $CAPEM $upperORG)" > organizations/peerOrganizations/user.flight.com/connection-user.yaml