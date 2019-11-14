import React, { Component, Fragment } from 'react';
import Router from 'next/router';
import Link from 'next/link';
import cx from 'classnames';
import axios from 'axios';

import TextField from '@material-ui/core/TextField';
import { withStyles } from '@material-ui/core/styles';
import { Button, Typography } from '@material-ui/core';
import GitHubIcon from '@material-ui/icons/GitHub';

import { TextInput, ButtonPrimary } from '@primer/components';

const navbarStyles = theme => ({
  navbar: {
    display: 'flex',
    backgroundColor: 'black',
    justifyContent: 'space-between',
    color: '#fff',
    padding: theme.spacing(2)
  },
  signup: {
    color: '#fff',
    border: '1px solid #fff'
  },
  signin: {
    color: '#fff'
  },
  searchbar: {
    backgroundColor: '#fff',
    borderRadius: '4px',
    margin: theme.spacing(1),
    height: '30px'
  }
});

function _Navbar({ classes }) {
  return (
    <div className={cx(classes.navbar, 'd-none d-md-flex')}>
      <GitHubIcon />
      <div>Why Github</div>
      <div>Enterprise</div>
      <div>Explore</div>
      <div>Marketplace</div>
      <div>Pricing</div>
      <div></div>
      <TextField variant="outlined" className={classes.searchbar} />
      <Link href="/login">
        <Button className={classes.signin}>Sign in</Button>
      </Link>
      <Button className={classes.signup}>Sign up</Button>
    </div>
  );
}
const Navbar = withStyles(navbarStyles)(_Navbar);

const signupStyles = theme => ({
  container: {
    backgroundColor: '#2b3137',
    color: '#fff',
    padding: '20px',
    margin: 0,
    [theme.breakpoints.up('md')]: {
      padding: '100px',
      paddingRight: '150px',
      maxWidth: '1368px',
      margin: 0,
      display: 'flex'
    }
  },
  signup: {
    backgroundColor: '#28A745',
    margin: '20px 0px 0px',
    padding: '25px',
    color: '#fff',
    width: '100%'
  },
  sidebar: {
    borderRadius: '4px',
    padding: theme.spacing(2),
    backgroundColor: '#fff',
    color: 'black',
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'center'
  },
  Input: {
    marginTop: theme.spacing(1)
  }
});

class _Signup extends Component {
  constructor(props) {
    super(props);
    this.state = {
      name: '',
      username: '',
      email: '',
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
      const response = await axios.post('/signup', {
        name: this.state.name,
        username: this.state.username,
        email: this.state.email,
        password: this.state.password
      });
      if (response.status === 201) {
        Router.push('/login');
      }
    } catch (error) {
      console.log(error);
    }
  };

  render() {
    const { classes } = this.props;
    return (
      <Fragment>
        <Navbar />
        <div className={cx(classes.container, 'container')}>
          <form onSubmit={this.handleSubmit}>
            <div className="row">
              <div className="col-lg-7 d-none d-md-block">
                <Typography variant="h2">Built for developers</Typography>
                <Typography variant="h6">
                  GitHub is a development platform inspired by the way you work.
                  From open source to business, you can host and review code,
                  manage projects, and build software alongside 40 million
                  developers.
                </Typography>
              </div>
              <div className={cx(classes.sidebar, 'col-lg-5')}>
                <div>
                  <Typography variant="h5">Create your account</Typography>
                </div>
                <Typography variant="body2">Name</Typography>
                <TextInput
                  type="text"
                  name="name"
                  margin="normal"
                  variant="outlined"
                  value={this.state.name}
                  onChange={this.handleChange}
                  className={classes.Input}
                />
                <Typography variant="body2">Username</Typography>
                <TextInput
                  type="text"
                  name="username"
                  margin="normal"
                  variant="outlined"
                  className={classes.Input}
                  value={this.state.username}
                  onChange={this.handleChange}
                />
                <Typography variant="body2">Email</Typography>
                <TextInput
                  type="text"
                  name="email"
                  margin="normal"
                  variant="outlined"
                  className={classes.Input}
                  value={this.state.email}
                  onChange={this.handleChange}
                />
                <Typography variant="body2">Password</Typography>
                <TextInput
                  type="password"
                  name="password"
                  margin="normal"
                  variant="outlined"
                  className={classes.Input}
                  value={this.state.password}
                  onChange={this.handleChange}
                />
                <div>
                  <Button
                    variant="contained"
                    className={classes.signup}
                    type="submit"
                  >
                    Sign up for Github
                  </Button>
                </div>
              </div>
            </div>
          </form>
        </div>
      </Fragment>
    );
  }
}

const Signup = withStyles(signupStyles)(_Signup);

export default Signup;
