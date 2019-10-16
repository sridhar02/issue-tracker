import React, { Component, Fragment } from "react";

import Router from "next/router";

import Link from "next/link";

import TextField from "@material-ui/core/TextField";

import { withStyles } from "@material-ui/core/styles";

import { Button, Typography } from "@material-ui/core";

import cx from "classnames";

import axios from "axios";

class Signup extends Component {
  constructor(props) {
    super(props);
    this.state = {
      name: "",
      username: "",
      email: "",
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
      .post("/signup", {
        name: this.state.name,
        username: this.state.username,
        email: this.state.email,
        password: this.state.password
      })
      .then(response => Router.push("/login"))
      .catch(error => {
        console.log(error);
      });
  };

  render() {
    return (
      <div className="container">
        <form onSubmit={this.handleSubmit}>
          <div className="row">
            <div className="col-12">
              <Typography variant="h6">Name</Typography>
              <TextField
                type="text"
                name="name"
                margin="normal"
                variant="outlined"
                value={this.state.name}
                onChange={this.handleChange}
              />
              <Typography variant="h6">Username</Typography>
              <TextField
                type="text"
                name="username"
                margin="normal"
                variant="outlined"
                value={this.state.username}
                onChange={this.handleChange}
              />
              <Typography variant="h6">Email</Typography>
              <TextField
                type="text"
                name="email"
                margin="normal"
                variant="outlined"
                value={this.state.email}
                onChange={this.handleChange}
              />
              <Typography variant="h6">Password</Typography>
              <TextField
                type="password"
                name="password"
                margin="normal"
                variant="outlined"
                value={this.state.password}
                onChange={this.handleChange}
              />
              <div>
                <Button variant="contained" color="primary" type="submit">
                  Signup
                </Button>
              </div>
            </div>
          </div>
        </form>
      </div>
    );
  }
}

export default Signup;
