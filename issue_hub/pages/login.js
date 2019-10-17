import React, { Component, Fragment } from "react";

import Router from "next/router";

import Link from "next/link";

import axios from "axios";

import { withStyles } from "@material-ui/core/styles";

import { Button, Typography } from "@material-ui/core";

import TextField from "@material-ui/core/TextField";

import cx from "classnames";

const loginStyles = theme => ({
  mainSection: {
    margin: theme.spacing(1)
  },
  title: {
    paddingBottom: theme.spacing(2),
    paddingLeft: theme.spacing(2)
  }
});

class _Login extends Component {
  constructor(props) {
    super(props);
    this.state = {
      username: "",
      password: ""
    };
  }

  handleChange = event => {
    this.setState({
      [event.target.name]: event.target.value
    });
  };

  handleSubmit = event => {
    event.preventDefault();
    axios
      .post("/signin", {
        username: this.state.username,
        password: this.state.password
      })
      .then(response => {
        localStorage.setItem("secret", response.data.secret);
        Router.push(`/user/${this.state.username}`);
      });
  };
  render() {
    const { classes } = this.props;
    return (
      <div className={cx(classes.mainSection, "container")}>
        <form onSubmit={this.handleSubmit}>
          <div className="row">
            <div className="col-12">
              <Typography variant="h6" className={classes.title}>
                Sign in to GitHub
              </Typography>
              <Typography variant="h6">Username or email address </Typography>
              <TextField
                placeholder="username"
                name="username"
                margin="normal"
                variant="outlined"
                value={this.state.name}
                onChange={this.handleChange}
              />
              <Typography variant="h6">Password </Typography>
              <TextField
                name="password"
                placeholder="password"
                margin="normal"
                variant="outlined"
                value={this.state.name}
                onChange={this.handleChange}
              />
              <div>
                <Button variant="contained" color="primary" type="submit">
                  Sign in
                </Button>
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
