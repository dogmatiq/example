import React from 'react';
import { connect } from 'react-redux';
import { Box, Grommet } from 'grommet';
import { grommet } from "grommet/themes";
import Header from './Header';
import Root from './Root';


function App(props) {
    const { router } = props
    return (
        <Grommet theme={grommet}>
            <Box fill>
                <Header/>
                <Root router={router}></Root>
            </Box>
        </Grommet>
    );
};

export default connect(state => ({
    store: state.store,
    router: state.router,
}))(App)
