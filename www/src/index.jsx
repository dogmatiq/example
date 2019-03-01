import React from 'react';
import { render } from 'react-dom';
import { Provider } from 'react-redux';
import { RouterProvider } from 'react-router5'
import configureStore  from './store';
import createRouter from './router'
import App from './components/App';

const router = createRouter()
const store = configureStore(
    router,
    window.__REDUX_DEVTOOLS_EXTENSION__ && window.__REDUX_DEVTOOLS_EXTENSION__()
)
const wrappedApp = (
    <Provider store={store}>
        <RouterProvider router={router}>
            <App/>
        </RouterProvider>
    </Provider>
)

router.start((err, state) => {
    render(wrappedApp, document.getElementById('app'))
})
