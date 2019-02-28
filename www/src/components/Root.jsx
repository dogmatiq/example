import React from 'react';
import { connect } from 'react-redux';
import { createRouteNodeSelector } from 'redux-router5';
import { startsWithSegment } from 'router5-helpers';
import Login from '../screens/Login';
import Home from '../screens/Home';
import NotFoundScreen from '../screens/NotFound';

function Root({ route }) {
    const isAuthenticated = false
    const { params, name } = route;

    const testRoute = startsWithSegment(name);

    if (params.requireAuth && !isAuthenticated ) {
        return <Login params={ params } />;
    }

    if (testRoute('login')) {
        return <Login params={ params } />;
    }

    if (testRoute('home')) {
        return <Home params={ params } />;
    }

    return (<NotFoundScreen/>)
}

export default connect(createRouteNodeSelector(''))(Root);
