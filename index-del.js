/**
 * Hyperledger Fabric REST API
 * @rbole 
 */

'use strict';
module.exports = async function (req, contract) {

  const util = require('util');

  // Get the key from the GET request and set it to lowercase, because of the chaincode.
  let queryKey = req.params.key;
  queryKey = queryKey.toLowerCase();

  try {
    
     /* 
     Submit the specified transaction.
     Submit a transaction to the ledger. The transaction function name will be evaluated on the endorsing peers and then submitted to the ordering service for committing to the ledger. 
     */
   await contract.submitTransaction('delete', queryKey);
    
    // Construct the finale return object.
    let r = {
      result: util.format('Delete the key: %s from the state in ledger' , queryKey)
    };
    return r;
  } catch(err){
    let r = {result:'Failed to delete key state: '+err};
    return r; 
  }
}
