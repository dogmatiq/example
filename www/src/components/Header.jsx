import React from 'react';
import { Box, Heading, Image } from 'grommet';
import icon from '../images/icon.png';

export default function Header(props){
    return (
        <Box
            tag='header'
            direction='row'
            align='center'
            justify='start'
            background='light-1'
            pad={{ left: 'medium', right: 'small', vertical: 'none' }}
            elevation='medium'
            style={{ zIndex: '1' }}
            {...props}>
            <Image src={icon} width="40px"/>
            <Heading color="dark-1" level='3' margin={{ left: 'medium'}}>
                Dogma Banking Example
            </Heading>
        </Box>
    );
};
