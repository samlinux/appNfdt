/**
 * Hyperledger Fabric REST API
 * @rbole 
 */

'use strict';
module.exports = async function (req, gateway, config) {
  // Get the key from the GET request and set it to lowercase, because of the chaincode.
  let txId = req.params.txid;

  try {
    
    // get a client object from the gateway
    let client = gateway.getClient();

    // get the channel from the client
    let channel = client.getChannel(config.channel);

    // get tx payload
    let response_payload = await channel.queryTransaction(txId);
    let result = response_payload;
    
    // get the rwset from this transaction
    // this is the path to the rwset of this transaction
    // is there a better way to do that ???
    let rwset = result.transactionEnvelope.payload.data.actions[0].payload.action.proposal_response_payload.extension.results.ns_rwset[1].rwset.writes[0];

    // parse the rwset of the given transaction to json
    rwset.value = JSON.parse(rwset.value)

    // Construct the finale return object.
    let r = {
      key: txId,
      value: rwset
    };
    return r;
  } catch(err){
    //console.log('Failed to evaluate transaction:', err)
    let r = {result:'Failed to evaluate transaction: '+err};
    return r; 
  }
}
