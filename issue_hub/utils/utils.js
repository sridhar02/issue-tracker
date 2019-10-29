import React, { Component } from "react";

import Link from "next/link";

import { Button, TextField, Typography } from "@material-ui/core";

import { withStyles } from "@material-ui/styles";

import Router from "next/router";
import { makeStyles } from "@material-ui/core/styles";
import AppBar from "@material-ui/core/AppBar";
import Toolbar from "@material-ui/core/Toolbar";
import IconButton from "@material-ui/core/IconButton";
import InputBase from "@material-ui/core/InputBase";

const navbarStyles = theme => ({});

class _Navbar extends Component {
  constructor(props) {
    super(props);
  }
  onSignout = () => {
    localStorage.removeItem("secret");
    Router.push("/login");
  };
  render() {
    const { classes } = this.props;
    return (
      <AppBar position="static">
        <div className="container">
          <div className="row">
            <div className="col-12">
              <Link href="/user/[username]">
                <Button color="inherit">USER</Button>
              </Link>
              <Link href="/account">
                <Button color="inherit" onClick={this.onSignout}>
                  Signout
                </Button>
              </Link>
            </div>
          </div>
        </div>
      </AppBar>
    );
  }
}

export const Navbar = withStyles(navbarStyles)(_Navbar);
