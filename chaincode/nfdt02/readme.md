# Additional Presentation Material 
Use the chaincode-docker-devmode to build and test the chaincode.

Start the first panel session.
```bash 
# start the session and the first panel
tmux new -s fabric

# switch to the chaincode-docker-devmode folder
cd fabric-samples/chaincode-docker-devmode

# start the dev network
docker-compose -f docker-compose-simple.yaml up
```

Start the second panel session.
```bash 
# create a new panel
CTRL + b \" (one double quote)

# switch to the chaincode-docker-devmode folder
cd fabric-samples/chaincode-docker-devmode

# switch into the chaincode container/folder
docker exec -it chaincode bash

# switch into the chaincode folder
cd nfdt01

# build the chaincode
go build

# run the chaincode
CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=mycc:0 ./nfdt02
```

Start the third panel session.
```bash
# switch into the chaincode container/folder
docker exec -it cli bash
cd /opt/gopath/src

# Install and instantiate the chaincode
peer chaincode install -p chaincodedev/chaincode/nfdt02 -n mycc -v 0
peer chaincode instantiate -n mycc -v 0 -c '{"Args":[]}' -C myc

# Invoke the chaincode
peer chaincode invoke -n mycc -c '{"Args":["add","02b073f0-cd8e-11ea-a8fc-11e4ed7bdc5d","{\"type\":\"art\",\"name\":\"Art 1\",\"description\":\"Art 1\",\"owner\":{\"firstname\":\"Roland\",\"lastname\":\"Bole\",\"departement\":\"Development\"}}"]}' -C myc

peer chaincode invoke -n mycc -c '{"Args":["add","e168f070-cd9a-11ea-8e94-798da3d2ddc5","{\"type\":\"art\",\"name\":\"Art 2\",\"description\":\"Art 2\",\"owner\":{\"firstname\":\"Roland\",\"lastname\":\"Bole\",\"departement\":\"Development\"}}"]}' -C myc


## check the new asset key and replace it for the update !!

peer chaincode invoke -n mycc -c '{"Args":["update","02b073f0-cd8e-11ea-a8fc-11e4ed7bdc5d","{\"name\":\"Slices\"}"]}' -C myc

peer chaincode invoke -n mycc -c '{"Args":["update","02b073f0-cd8e-11ea-a8fc-11e4ed7bdc5d","{\"description\":\"An array has a fixed size. A slice, on the other hand, is a dynamically-sized, flexible view into the elements of an array. In practice, slices are much more common than arrays. The type []T is a slice with elements of type T. A slice is formed by specifying two indices, a low and high bound, separated by a colon:\"}"]}' -C myc


# Query the chaincode

peer chaincode query -n mycc -c '{"Args":["queryById","02b073f0-cd8e-11ea-a8fc-11e4ed7bdc5d"]}' -C myc |jq '.'

peer chaincode query -n mycc -c '{"Args":["queryAdHoc","{\"selector\": {\"owner.departement\": \"Development\"}, \"use_index\":[\"_design/indexOwnerDoc\", \"indexOwner\"]}"]}' -C myc |jq '.'

peer chaincode query -n mycc -c '{"Args":["queryAdHoc","{\"selector\": {\"owner.firstname\": \"Roland\"}, \"use_index\":[\"_design/indexOwnerDoc\", \"indexOwner\"]}"]}' -C myc |jq '.'

peer chaincode query -n mycc -c '{"Args":["getAllTxByKey","02b073f0-cd8e-11ea-a8fc-11e4ed7bdc5d"]}' -C myc |jq '.'


# Init invoke
peer chaincode invoke -n mycc -c '{"Args":["add","{}"]}' -C myc
```


