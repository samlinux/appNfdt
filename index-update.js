/**
 * Hyperledger Fabric REST API
 * @rbole 
 */

'use strict';
const util = require('util');

module.exports = async function (req, contract) {

  // Get the keys and value from the POST request.
  let key = req.body.key;
  let value = req.body.value;
  value = JSON.stringify(value)

  try {
     /* 
     Submit the specified transaction.
     Submit a transaction to the ledger. The transaction function name will be evaluated on the endorsing peers and then submitted to the ordering service for committing to the ledger. 
     */
    let result = await contract.submitTransaction('update', key, value);
    
    result = result.toString();
    result = JSON.parse(result);
    
    // Prepare the return value.
    let r = util.format('An existing Asset has successfully updated: Key: %s, TxId: %s ', result.Key, result.TxId);
    return r;
  }
  catch(error){
    let r = {r:'Failed to submit transaction: '+error};
    return r;
  }
}