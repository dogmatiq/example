import React from 'react';
import { render } from 'react-dom';
import { Provider } from 'react-redux';
import { RouterProvider } from 'react-router5'
import { configureStore } from './store';
import { configureRouter } from './router';
import App from './components/App';

const router = configureRouter()
const store = configureStore(router)
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
