const {OpenAccountRequest, OpenAccountResponse} = require('./pb/account_pb.js');
const {AccountClient} = require('./pb/account_grpc_web_pb');


var client = new AccountClient('http://localhost:9900');

var request = new OpenAccountRequest();
request.setAccountId("fake-account-id");
request.setName("fake-account-name");

client.openAccount(request, {}, (err, response) => {
    console.log(response)
});
