import { CustomerPromiseClient } from '../pb/customer_grpc_web_pb';
import { LoginRequest } from '../pb/customer_pb';

const client = new CustomerPromiseClient(GPRCWEB_SERVER);

export function login(name, password) {
    let request = new LoginRequest();
    request.setCustomerName(name);
    request.setPassword(password);

    return client.login(request, {})
}

