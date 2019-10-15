import React, { Component, Fragment } from "react";

import Router from "next/router";

import Link from "next/link";

import Button from "@material-ui/core/Button";

import TextField from "@material-ui/core/TextField";

class Signup extends Component {
  render() {
    return (
      <div>
        <div>
          <label>Username </label>
          <TextField
            placeholder="Username "
            margin="normal"
            variant="outlined"
          />
        </div>
        <div>
          <label>email address</label>
          <TextField margin="normal" variant="outlined" />
        </div>
        <div>
          <label>Password</label>
          <TextField margin="normal" variant="outlined" />
        </div>
        <Button variant="contained" color="primary">
          Signup
        </Button>
      </div>
    );
  }
}

export default Signup;
