# Support Material

Create a new base network folder.
```bash
cd fabric/fabric-samples
cp -r first-network rb-first-network
cd rb-first-network
```

Create a file to start the network.

```bash
vi start.sh
#!/bin/bash
# Exit on first error
set -e

CC_CHANNEL_NAME=myc
CC_NAME=mycc
CC_SRC_PATH=github.com/chaincode/nfdt02

startNetworkWithChaincode() {
  cd ../rb-first-network
  echo y | ./byfn.sh down
  echo y | ./byfn.sh up -a -n -s couchdb -c $CC_CHANNEL_NAME

  CONFIG_ROOT=/opt/gopath/src/github.com/hyperledger/fabric/peer
  Org1Path=$CONFIG_ROOT/crypto/peerOrganizations/org1.example.com
  Org2Path=$CONFIG_ROOT/crypto/peerOrganizations/org2.example.com
  TlsPath=$CONFIG_ROOT/crypto/ordererOrganizations/example.com/orderers/orderer.example.com

  ORG1_MSPCONFIGPATH=${Org1Path}/users/Admin@org1.example.com/msp
  ORG1_TLS_ROOTCERT_FILE=${Org1Path}/peers/peer0.org1.example.com/tls/ca.crt
  ORG2_MSPCONFIGPATH=${Org2Path}/users/Admin@org2.example.com/msp
  ORG2_TLS_ROOTCERT_FILE=${Org2Path}/peers/peer0.org2.example.com/tls/ca.crt
  ORDERER_TLS_ROOTCERT_FILE=${TlsPath}/msp/tlscacerts/tlsca.example.com-cert.pem

  set -x

  echo "Installing smart contract: $CC_NAME on peer0.org1.example.com"
  docker exec \
    -e CORE_PEER_LOCALMSPID=Org1MSP \
    -e CORE_PEER_ADDRESS=peer0.org1.example.com:7051 \
    -e CORE_PEER_MSPCONFIGPATH=${ORG1_MSPCONFIGPATH} \
    -e CORE_PEER_TLS_ROOTCERT_FILE=${ORG1_TLS_ROOTCERT_FILE} \
    cli \
    peer chaincode install \
      -n "$CC_NAME" \
      -v 1.0 \
      -p "$CC_SRC_PATH" 

  echo "Installing smart contract: $CC_NAME on peer1.org1.example.com"
  docker exec \
    -e CORE_PEER_LOCALMSPID=Org1MSP \
    -e CORE_PEER_ADDRESS=peer1.org1.example.com:8051 \
    -e CORE_PEER_MSPCONFIGPATH=${ORG1_MSPCONFIGPATH} \
    -e CORE_PEER_TLS_ROOTCERT_FILE=${ORG1_TLS_ROOTCERT_FILE} \
    cli \
    peer chaincode install \
      -n "$CC_NAME" \
      -v 1.0 \
      -p "$CC_SRC_PATH" 

  echo "Installing smart contract: $CC_NAME on peer0.org2.example.com"
  docker exec \
    -e CORE_PEER_LOCALMSPID=Org2MSP \
    -e CORE_PEER_ADDRESS=peer0.org2.example.com:9051 \
    -e CORE_PEER_MSPCONFIGPATH=${ORG2_MSPCONFIGPATH} \
    -e CORE_PEER_TLS_ROOTCERT_FILE=${ORG2_TLS_ROOTCERT_FILE} \
    cli \
    peer chaincode install \
      -n "$CC_NAME" \
      -v 1.0 \
      -p "$CC_SRC_PATH" 

  echo "Installing smart contract: $CC_NAME on peer1.org2.example.com"
  docker exec \
    -e CORE_PEER_LOCALMSPID=Org2MSP \
    -e CORE_PEER_ADDRESS=peer1.org2.example.com:10051 \
    -e CORE_PEER_MSPCONFIGPATH=${ORG2_MSPCONFIGPATH} \
    -e CORE_PEER_TLS_ROOTCERT_FILE=${ORG2_TLS_ROOTCERT_FILE} \
    cli \
    peer chaincode install \
      -n "$CC_NAME" \
      -v 1.0 \
      -p "$CC_SRC_PATH"     

  echo "Instantiating smart contract: $CC_NAME on $CC_CHANNEL_NAME"
  docker exec \
    -e CORE_PEER_LOCALMSPID=Org1MSP \
    -e CORE_PEER_MSPCONFIGPATH=${ORG1_MSPCONFIGPATH} \
    cli \
    peer chaincode instantiate \
      -o orderer.example.com:7050 \
      -C $CC_CHANNEL_NAME \
      -n $CC_NAME \
      -v 1.0 \
      -c '{"Args":[]}' \
      -P "AND('Org1MSP.peer','Org2MSP.peer')" \
      --tls \
      --cafile ${ORDERER_TLS_ROOTCERT_FILE} \
      --tlsRootCertFiles ${ORG1_TLS_ROOTCERT_FILE} #\
      #--peerAddresses peer0.org1.example.com:7051 peer0.org2.example.com:9051
  
  echo "Waiting for instantiation request to be committed ..."
  sleep 10
  echo "Ready to use the network ..."
}

# start the network with a custom chaincode
startNetworkWithChaincode
```

Create a script to stop the network.
```bash
vi stop.sh

#!/bin/bash
# Exit on first error
set -e

cd ../first-network
echo y | ./byfn.sh down

```