import React from 'react';
import { connect } from 'react-redux';
import { Box, Button, Form, FormField } from 'grommet';
import { Waypoint } from "grommet-icons";
import { customerActions } from '../actions';

class Login extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        const {handleSubmit} = this.props;
        return (
            <Box fill align="center" justify="center" pad="large">
                <Box width="medium">
                    <Form onSubmit={({ value }) => handleSubmit(value)}>
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

export default connect(
    state => ({
        loading: state.customer.loading,
        error: state.customer.error
    }),
    (dispatch) => ({
        handleSubmit: (data) => {
            const { customer, password } = data;
            dispatch(customerActions.login(customer, password));
        }
    })
    )(Login);
