import React from 'react';
import { connect } from 'react-redux';
import { Box, Button, Form, FormField } from 'grommet';
import { Waypoint } from "grommet-icons";
import { customerActions } from '../actions';

class Login extends React.Component {
    constructor(props) {
        super(props);
        this.handleSubmit = this.handleSubmit.bind(this)
    }

    handleSubmit(data) {
        const { dispatch } = this.props;
        const { customer, password } = data;
        console.log(arguments)
        dispatch(customerActions.login(customer, password));
    }

    render() {
        return (
            <Box fill align="center" justify="center" pad="large">
                <Box width="medium">
                    <Form onSubmit={({ value }) => this.handleSubmit(value)}>
                        <FormField
                            label="Name"
                            name="customer"
                            required />
                        <FormField
                            label="Password"
                            name="password"
                            type="password"
                            required />
                        <Box
                            margin={{ top: "large" }}>
                            <Button
                                primary={true}
                                type="submit"
                                icon={<Waypoint />}
                                label="Login"
                                alignSelf="end"
                                color="dark-1">
                            </Button>
                        </Box>
                    </Form>
                </Box>
            </Box>
        );
    }
}

export default connect(state => ({state}))(Login)
