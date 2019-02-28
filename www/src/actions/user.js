import { userConst }  from '../constants';

export const userActions = {
    login
};

function login(username, password) {

    return dispatch => {
        dispatch(request({ username, password }));

        fetch("http://someurl.com")
            .then(
                user => {
                    dispatch(success(user));
                },
                error => {
                    dispatch(failure(error));
                }
            );
    };

    function request(user, password) { return { type: userConst.USER_LOGIN, user: user, password:  password} }
    function success(user) { return { type: userConst.USER_LOGIN_SUCCESS, user: user } }
    function failure(error) { return { type: userConst.USER_LOGIN_FAILURE, error: error } }
}

