const {OpenAccountRequest, OpenAccountResponse} = require('./account_pb.js');
const {AccountClient} = require('./account_grpc_web_pb');


var client = new AccountClient('http://localhost:9900');

var request = new OpenAccountRequest();
request.setAccountId("fake-account-id");
request.setName("fake-account-name");

client.openAccount(request, {}, (err, response) => {
    console.error(arguments)
});
