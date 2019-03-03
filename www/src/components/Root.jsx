import React from 'react';
import { connect } from 'react-redux';
import { createRouteNodeSelector } from 'redux-router5';
import { startsWithSegment } from 'router5-helpers';
import { actions as routerActions } from 'redux-router5'
import Login from '../screens/Login';
import Home from '../screens/Home';
import NotFoundScreen from '../screens/NotFound';

function Root(props) {
    const { authenticated, route, toLoginPage } = props;
    const { params, name } = route;
    const testRoute = startsWithSegment(name);

    if (params.requireAuth && !authenticated) {
        toLoginPage()
        return <Login params={params} />;
    }

    if (testRoute('login')) {
        return <Login params={params} />;
    }

    if (testRoute('home')) {
        return <Home params={params} />;
    }

    return (<NotFoundScreen />)
}

export default connect(
    state => {
        const selector = createRouteNodeSelector('');
        return (state) => ({
            authenticated: state.customer.authenticated,
            ...selector(state),
        })
    },
    (dispatch) => ({
      toLoginPage: () => dispatch(routerActions.navigateTo("login"))
    })

)(Root);
