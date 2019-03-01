import { userConst }  from '../constants';

const request = (name)=> { return { type: userConst.USER_LOGIN, name} }
const success = (name, id)=> { return { type: userConst.USER_LOGIN_SUCCESS, name, id} }
const failure = (error)=> { return { type: userConst.USER_LOGIN_FAILURE, error } }


function login(name, password) {
    return (dispatch, getState, services) => {
        dispatch(request(name));
        services.customer.login(name, password)
            .then(resp => {
                dispatch(success(
                    resp.getCustomerName(),
                    resp.getCustomerId(),
                ));
            })
            .catch(error => {
                dispatch(failure(error));
            })
    };
}

export const customerActions = {
    login
};
