import { customerConst }  from '../constants';
import { actions as routerActions } from 'redux-router5'

const request = (name)=> ({ type: customerConst.CUSTOMER_LOGIN, name})
const success = (id)=> ({ type: customerConst.CUSTOMER_LOGIN_SUCCESS, id})
const failure = (error)=> ({ type: customerConst.CUSTOMER_LOGIN_FAILURE, error })


function login(name, password) {
    return (dispatch, getState, services) => {
        dispatch(request(name));
        services.customer.login(name, password)
            .then(resp => {
                dispatch(success(
                    resp.getCustomerId(),
                ));
                dispatch(routerActions.navigateTo("home"))
            })
            .catch(error => {
                dispatch(failure(error));
            })
    };
}

export const customerActions = {
    login
};
