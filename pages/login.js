import React, { Component, Fragment } from 'react';
import Router from 'next/router';
import Link from 'next/link';
import axios from 'axios';
import cx from 'classnames';

import { withStyles } from '@material-ui/core/styles';
import { Button, Typography } from '@material-ui/core';
import TextField from '@material-ui/core/TextField';
import GitHubIcon from '@material-ui/icons/GitHub';

import { TextInput, ButtonPrimary } from '@primer/components';

const loginStyles = theme => ({
  mainSection: {
    margin: theme.spacing(1),
    [theme.breakpoints.up('md')]: {
      margin: '80px auto',
      display: 'flex',
      justifyContent: 'center'
    }
  },
  title: {
    padding: theme.spacing(2),
    marginLeft: theme.spacing(6)
  },
  gitIcon: {
    padding: theme.spacing(0.5),
    fontSize: '55px',
    marginLeft: theme.spacing(14)
  },
  inputFields: {
    border: '1px solid #ddd',
    padding: theme.spacing(2),
    margin: theme.spacing(2)
  },
  username: {
    fontWeight: 'bold'
  },
  password: {
    fontWeight: 'bold'
  },
  signinButton: {
    backgroundColor: '#28A745',
    margin: '20px 0px 0px',
    padding: '6px 12px',
    color: '#fff',
    width: '100%'
  },
  newAccount: {
    border: '1px solid #ddd',
    padding: theme.spacing(2),
    margin: theme.spacing(2)
  },
  topHeaders: {
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'center'
  },
  Input: {
    marginTop: theme.spacing(1)
  }
});

class _Login extends Component {
  constructor(props) {
    super(props);
    this.state = {
      username: '',
      password: ''
    };
  }

  handleChange = event => {
    this.setState({
      [event.target.name]: event.target.value
    });
  };

  handleSubmit = async event => {
    event.preventDefault();
    try {
      const response = await axios.post('/signin', {
        username: this.state.username,
        password: this.state.password
      });
      if (response.status === 201) {
        localStorage.setItem('secret', response.data.secret);
        Router.push(`/user/${this.state.username}`);
      }
    } catch (error) {
      console.log(error);
    }
  };
  render() {
    const { classes } = this.props;
    return (
      <div className={cx(classes.mainSection, 'container')}>
        <form onSubmit={this.handleSubmit}>
          <div className="row">
            <div className="col-12">
              <div className={classes.topHeaders}>
                <div>
                  <GitHubIcon className={classes.gitIcon} />
                </div>
                <div>
                  <Typography variant="h6" className={classes.title}>
                    Sign in to GitHub
                  </Typography>
                </div>
              </div>
              <div className={classes.inputFields}>
                <Typography variant="body2" className={classes.username}>
                  Username or email address
                </Typography>
                <TextInput
                  placeholder="username"
                  name="username"
                  margin="normal"
                  variant="outlined"
                  value={this.state.name}
                  onChange={this.handleChange}
                  className={classes.Input}
                />
                <Typography variant="body2" className={classes.password}>
                  Password
                </Typography>
                <TextInput
                  name="password"
                  placeholder="password"
                  type="password"
                  margin="normal"
                  variant="outlined"
                  value={this.state.name}
                  className={classes.Input}
                  onChange={this.handleChange}
                />
                <div>
                  <Button
                    variant="contained"
                    className={classes.signinButton}
                    type="submit"
                  >
                    Sign in
                  </Button>
                </div>
              </div>
              <div className={classes.newAccount}>
                <Typography variant="body2">
                  New to GitHub?{' '}
                  <Link href="/join">
                    <a>Create an account</a>
                  </Link>
                </Typography>
              </div>
            </div>
          </div>
        </form>
      </div>
    );
  }
}

const Login = withStyles(loginStyles)(_Login);

export default Login;
