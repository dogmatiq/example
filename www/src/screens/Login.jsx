import React from 'react';
import { connect } from 'react-redux';
import { Box, Button, FormField, TextInput } from 'grommet';
import { LinkNext } from "grommet-icons";
import { customerActions } from '../actions';

class Login extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            customer: '',
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

        this.setState({ submitted: true }, ()=>{
            const { customer, password } = this.state;
            const { dispatch } = this.props;
            if (customer && password) {
                dispatch(customerActions.login(customer, password));
            }
        });
    }

    render() {
        return (
            <Box pad="large" align="center">
                <form onSubmit={this.handleSubmit}>
                    <Box>
                        <FormField label="Name">
                            <TextInput onChange={this.handleChange} name="customer"/>
                        </FormField>
                        <FormField label="Password" required>
                            <TextInput type="password" onChange={this.handleChange} name="password"/>
                        </FormField>
                        <Button
                            icon={<LinkNext />}
                            type="submit"
                            reverse={true}
                            primary={true}
                            label="Login"
                        />
                    </Box>
                </form>
            </Box>
        );
    }
}

function mapStateToProps(state) {
    return state;
}

export default connect(mapStateToProps)(Login);
