import { createStore, applyMiddleware, combineReducers } from 'redux';
import { router5Middleware, router5Reducer } from 'redux-router5';
import { createLogger } from 'redux-logger';

export function configureStore (router, initialState = {}) {
    const createStoreWithMiddleware = applyMiddleware(
        router5Middleware(router),
        createLogger()
    )(createStore);
    const store = createStoreWithMiddleware(
        combineReducers({
        router: router5Reducer,
    }), initialState);

    window.store = store;
    return store;
}
