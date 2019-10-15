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
  render() {
    const { classes } = this.props;
    return (
      <div className={cx(classes.mainSection, "container")}>
        <div className="row">
          <div className="col">
            <Typography variant="h6" className={classes.title}>
              Sign in to GitHub
            </Typography>
            <form onSubmit={this.handleSubmit} className={"col-12 col-md-4"}>
              <Typography variant="h6">Username or email address </Typography>
              <TextField
                placeholder="username or email"
                margin="normal"
                variant="outlined"
              />
              <Typography variant="h6">Password </Typography>
              <TextField
                placeholder="password"
                margin="normal"
                variant="outlined"
              />
              <Button variant="contained" color="primary">
                Sign in
              </Button>
            </form>
          </div>
        </div>
      </div>
    );
  }
}

const Login = withStyles(loginStyles)(_Login);

export default Login;
