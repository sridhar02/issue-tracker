import React, { Component } from "react";

import Link from "next/link";

import { Button, TextField, Typography } from "@material-ui/core";

import { withStyles } from "@material-ui/styles";

import Router from "next/router";

import MenuIcon from "@material-ui/icons/Menu";

import cx from "classnames";

import { makeStyles } from "@material-ui/core/styles";

import AppBar from "@material-ui/core/AppBar";

import Toolbar from "@material-ui/core/Toolbar";

import IconButton from "@material-ui/core/IconButton";

import InputBase from "@material-ui/core/InputBase";

import GitHubIcon from "@material-ui/icons/GitHub";

import NotificationsIcon from "@material-ui/icons/Notifications";

const navbarStyles = theme => ({
  container: {
    backgroundColor: "black",
    padding: theme.spacing(2)
  },
  navbar: {
    display: "flex",
    justifyContent: "space-between"
  },
  overview: {
    display: "flex",
    justifyContent: "space-between",
    marginTop: theme.spacing(3)
  },
  Button: {
    color: "#fff",
    padding: theme.spacing(0),
    margin: theme.spacing(0)
  }
});

class _Navbar extends Component {
  constructor(props) {
    super(props);
    this.state = {
      navbar: "false"
    };
  }
  onSignout = () => {
    localStorage.removeItem("secret");
    Router.push("/login");
  };

  handleClick = event => {
    event.preventDefault();
    this.setState({
      navbar: "true"
    });
  };

  render() {
    const { classes } = this.props;
    return (
      <AppBar position="static">
        <div className={cx(classes.container, "container")}>
          <div className="row">
            <div className="col-12">
              <div className={classes.navbar}>
                <Button onClick={this.handleClick} className={classes.Button}>
                  <MenuIcon />
                </Button>
                <GitHubIcon />
                <NotificationsIcon />
              </div>
              <div className={classes.overview}>
                <div>Overview</div>
                <div>Repositories</div>
                <div>Projects</div>
              </div>
            </div>
          </div>
        </div>
      </AppBar>
    );
  }
}

export const Navbar = withStyles(navbarStyles)(_Navbar);
