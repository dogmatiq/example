import { customerConst } from '../constants';

export function auth(state = {}, action) {
  switch (action.type) {
    case customerConst.CUSTOMER_LOGIN:
      return {
        loading: true
      };
    case customerConst.CUSTOMER_LOGIN_SUCCESS:
      return {
        name: action.name
      };
    case customerConst.CUSTOMER_LOGIN_FAILURE:
      return {
        error: action.error
      };
    default:
      return state
  }
}
