/**
 * test for our super fabric REST API
 */
const supertest = require('supertest');
const util = require('util');
const api = supertest('localhost:3000');

describe("Hyperledger Fabric API tests", function() {
  it("checks if api is running", async function() {
    this.skip();
    let result = await api.get('/')
    console.log(result.body)
  }) 
  
  it("add a new nfdt", async function() {
    this.skip();

    // data to store in the blockchain
    let data = {
      "type":"art",
      "name":"Hyperledger Fabric Functionalities",
      "description": "Hyperledger Fabric is an implementation of distributed ledger technology (DLT) that delivers enterprise-ready network security, scalability, confidentiality and performance, in a modular blockchain architecture. Hyperledger Fabric delivers the following blockchain network functionalities:",
      "owner": {
        "firstName": "Roland",
        "lastName": "Bole",
        "departement": "Development"
      }
    }
    
    // payload send to API
    let payload = {
      value: data
    };

    let result = await api.post('/add').send(payload)
    console.log(result.body)
  })

  it("update an existing nfdt", async function() {
    this.skip();

    // data to store in the blockchain
    let data = {
      "name":"Hyperledger Fabric’s first long term support release",
      "description": "Hyperledger Fabric has matured since the initial v1.0 release, and so has the community of Fabric operators. The Fabric developers have been working with network operators to deliver v1.4 with a focus on stability and production operations. As such, v1.4.x will be our first long term support release.Our policy to date has been to provide bug fix (patch) releases for our most recent major or minor release until the next major or minor release has been published. We plan to continue this policy for subsequent releases. However, for Hyperledger Fabric v1.4, the Fabric maintainers are pledging to provide bug fixes for a period of one year from the date of release. This will likely result in a series of patch releases (v1.4.1, v1.4.2, and so on), where multiple fixes are bundled into a patch release.If you are running with Hyperledger Fabric v1.4.x, you can be assured that you will be able to safely upgrade to any of the subsequent patch releases. In the advent that there is need of some upgrade process to remedy a defect, we will provide that process with the patch release."
    }
    
    // payload send to API
    let payload = {
      key:'01a83570-cdbf-11ea-8c3a-2ff57c1be9a3',
      value: data
    };

    let result = await api.post('/update').send(payload)
    console.log(result.body)
  })

  it("query an asset by Id", async function() {
    this.skip();
    let key = '01a83570-cdbf-11ea-8c3a-2ff57c1be9a3';
    let result = await api.get('/queryById/'+key)
    console.log(util.inspect(result.body,false, null, true));
  })

  it("query an asset by owner", async function() {
    this.skip();
    let key = 'Bole';
    let result = await api.get('/queryByOwner/'+key)
    console.log(util.inspect(result.body,false, null, true));
  })

  it("do an adHocQuery", async function() {
    this.skip();
    let query = {
      "selector": {"owner.departement": "Development"}, 
      "use_index":["_design/indexOwnerDoc", "indexOwner"]
    }

    let payload = {
      value: query
    };

    let result = await api.post('/queryAdHoc').send(payload)
    console.log(util.inspect(result.body,false, null, true));
  })

  it("get getAllTxByKey", async function() {
    this.skip();
    let key = '01a83570-cdbf-11ea-8c3a-2ff57c1be9a3';

    let result = await api.get('/getAllTxByKey/'+key)
    let tx = JSON.parse(result.body.value);
    console.log(util.inspect(tx,false, null, true));
    tx.forEach(tx => {
       console.log(util.format('TxId: %s',tx.TxId)) 
    });

  })


  it("get getTxByTxId", async function() {
    //this.skip();
    let txId = '2ccce674dc5ff59d4077385339d67d23a5e5e7799727816c5f6070770935f2c3';

    let result = await api.get('/getTxByTxId/'+txId)
    console.log(util.inspect(result.body.value,false, null, true));
  })


})