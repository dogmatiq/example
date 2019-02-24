import React from 'react';
import { connect } from 'react-redux';
import { createRouteNodeSelector } from 'redux-router5';
import { startsWithSegment } from 'router5-helpers';
import NotFoundScreen from '../screens/NotFound';

function Root({ route }) {
    // const { params, name } = route;
    const testRoute = startsWithSegment(name);

    // if (testRoute('home')) {
    //     return <Home params={ params } />;
    // } else if (testRoute('about')) {
    //     return <About params={ params } />;
    // } else if (testRoute('contact')) {
    //     return <Contact params={ params } />;
    // }

    return (<NotFoundScreen/>)
}

export default connect(createRouteNodeSelector(''))(Root);
