import { createStore, applyMiddleware, combineReducers } from 'redux';
import thunk from 'redux-thunk';
import { router5Middleware, router5Reducer } from 'redux-router5';
import { createLogger } from 'redux-logger';
import * as reducers from './reducers';
import services from './services';

export default function configureStore(router, initialState = {}) {
    const createStoreWithMiddleware = applyMiddleware(
        router5Middleware(router),
        createLogger(),
        thunk.withExtraArgument(services)
    )(createStore);
    const store = createStoreWithMiddleware(
        combineReducers({
            ...reducers,
            router: router5Reducer,
        }), initialState);

    window.store = store;
    return store;
}
