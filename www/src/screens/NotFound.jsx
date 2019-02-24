import React from 'react'
import { Box, Heading } from "grommet";
import { StatusWarning } from "grommet-icons";


export default function NotFoundScreen() {
    return (
        <Box fill align="center" margin={{ top: "xlarge" }}>
            <StatusWarning size='xlarge' color='plain'/>
            <Heading size='medium'>
                Not found
            </Heading>
        </Box>
    )
}
