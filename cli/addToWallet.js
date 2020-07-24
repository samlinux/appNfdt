/*
 *  SPDX-License-Identifier: Apache-2.0
 */

'use strict';

// Bring key classes into scope, most importantly Fabric SDK network class
const fs = require('fs');
const path = require('path');

const { FileSystemWallet, X509WalletMixin } = require('fabric-network');

const fixtures = path.resolve(__dirname, '../../fabric-samples/rb-first-network');

// A wallet stores a collection of identities
const wallet = new FileSystemWallet('../wallet');

let pkOfTheUser = '0f0c7ee588337ca59eaf89d410d8382b25b26832eea6434218fac7dcdb3a7402_sk';

// config
let config = {
  pathToUser:'/crypto-config/peerOrganizations/org1.example.com/users/User1@org1.example.com',
  pathToUserSignCert: '/msp/signcerts/User1@org1.example.com-cert.pem',
  pathToUserPrivKey: '/msp/keystore/'+pkOfTheUser,
  identityLabel: 'User1@org1.example.com'
}

async function main() {

  // Main try/catch block
  try {

    // Identity to credentials to be stored in the wallet
    const credPath = path.join(fixtures, config.pathToUser);
    const cert = fs.readFileSync(path.join(credPath, config.pathToUserSignCert)).toString();
    const key = fs.readFileSync(path.join(credPath, config.pathToUserPrivKey)).toString();

    // Load credentials into wallet
    const identityLabel = config.identityLabel;
    const identity = X509WalletMixin.createIdentity('Org1MSP', cert, key);

    await wallet.import(identityLabel, identity);

  } catch (error) {
    console.log(`Error adding to wallet. ${error}`);
    console.log(error.stack);
  }
}

main().then(() => {
    console.log('done');
}).catch((e) => {
    console.log(e);
    console.log(e.stack);
    process.exit(-1);
});