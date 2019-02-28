import { userConst } from '../constants';

export function auth(state = {}, action) {
  switch (action.type) {
    case userConst.USER_LOGIN:
      return {
        loading: true
      };
    case userConst.USER_LOGIN_SUCCESS:
      return {
        user: action.user
      };
    case userConst.USER_LOGIN_FAILURE:
      return {
        error: action.error
      };
    default:
      return state
  }
}
