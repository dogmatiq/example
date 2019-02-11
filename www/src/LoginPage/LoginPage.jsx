import React from 'react';
import { connect } from 'react-redux';
import { Box, Button, FormField, TextInput } from 'grommet';
import { FormNext } from "grommet-icons";

import { userActions } from '../_actions';

class LoginPage extends React.Component {
    constructor(props) {
        super(props);

        // reset login status
        this.props.dispatch(userActions.logout());

        this.state = {
            username: '',
            password: '',
            submitted: false
        };

        this.handleChange = this.handleChange.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);
    }

    handleChange(e) {
        const { name, value } = e.target;
        this.setState({ [name]: value });
    }

    handleSubmit(e) {
        e.preventDefault();

        this.setState({ submitted: true });
        const { username, password } = this.state;
        const { dispatch } = this.props;
        if (username && password) {
            dispatch(userActions.login(username, password));
        }
    }

    render() {
        const { loggingIn } = this.props;
        const { username, password, submitted } = this.state;
        return (
            <Box pad="large" align="center">
                <form onSubmit={event => event.preventDefault()}>
                    <Box>
                        <FormField label="Name">
                            <TextInput />
                        </FormField>
                        <FormField label="Password">
                            <TextInput type="password" />
                        </FormField>
                        <Button
                            icon={<FormNext />}
                            type="submit"
                            reverse="true"
                            label="Login"
                            primary={true}
                        />
                    </Box>
                </form>
            </Box>
        );
    }
}

function mapStateToProps(state) {
    const { loggingIn } = state.authentication;
    return {
        loggingIn
    };
}

const connectedLoginPage = connect(mapStateToProps)(LoginPage);
export { connectedLoginPage as LoginPage };
