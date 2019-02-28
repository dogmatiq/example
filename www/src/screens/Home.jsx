import React from 'react';
import { connect } from 'react-redux';
import { Box, Heading } from "grommet";
import { Home as HomeIcon } from "grommet-icons";

class Home extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        return (
            <Box fill align="center" margin={{ top: "xlarge" }}>
            <HomeIcon size='xlarge' color='plain'/>
            <Heading size='medium'>
               Home
            </Heading>
        </Box>
        );
    }
}

function mapStateToProps(state) {
   return state
}

export default connect(mapStateToProps)(Home);
