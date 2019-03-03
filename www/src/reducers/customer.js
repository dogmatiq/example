import { customerConst } from '../constants';

export function customer(state = {}, action) {
  switch (action.type) {
    case customerConst.CUSTOMER_LOGIN:
      return {
        loading: true,
        authenticated: false
      };
    case customerConst.CUSTOMER_LOGIN_SUCCESS:
      return {
        loading: false,
        id: action.id,
        authenticated: true,
      };
    case customerConst.CUSTOMER_LOGIN_FAILURE:
      return {
        loading: false,
        error: action.error,
        authenticated: false
      };
    default:
      return state
  }
}
