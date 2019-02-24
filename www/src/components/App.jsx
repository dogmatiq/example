import React from 'react';
import { connect } from 'react-redux';
import { Box, Button, Heading, Grommet } from 'grommet';
import { Notification } from 'grommet-icons';
import Root from './Root';


const AppBar = (props) => (
    <Box
        tag='header'
        direction='row'
        align='center'
        justify='between'
        background='brand'
        pad={{ left: 'medium', right: 'small', vertical: 'small' }}
        elevation='medium'
        style={{ zIndex: '1' }}
        {...props}
    />
);

function App(props) {
    const { router, theme } = props
    return (
        <Grommet theme={theme}>
            <Box fill>
                <AppBar>
                    <Heading level='3' margin='none'>Dogma Banking Example</Heading>
                    <Button icon={<Notification />} onClick={() => { console.log(arguments)}} />
                </AppBar>
                <Root router={router}></Root>
            </Box>
        </Grommet>
    );
};

export default connect(state => ({
    store: state.store,
    router: state.router,
    theme: {
        global: {
            font: {
                family: 'Roboto',
                size: '14px',
                height: '20px',
            },
        },
    }
}))(App)
