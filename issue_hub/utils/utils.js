import React, { Component } from "react";
import Link from "next/link";
import Router from "next/router";
import cx from "classnames";
import axios from "axios";

import { Button, TextField, Typography } from "@material-ui/core";
import { withStyles } from "@material-ui/core/styles";
import MenuIcon from "@material-ui/icons/Menu";
import AppBar from "@material-ui/core/AppBar";
import Toolbar from "@material-ui/core/Toolbar";
import IconButton from "@material-ui/core/IconButton";
import InputBase from "@material-ui/core/InputBase";
import GitHubIcon from "@material-ui/icons/GitHub";
import NotificationsIcon from "@material-ui/icons/Notifications";

const menuStyles = theme => ({
  dashboard: {
    padding: theme.spacing(1),
    margin: theme.spacing(1),
    borderBottom: "1px solid #ddd"
  },
  issues: {
    padding: theme.spacing(1),
    margin: theme.spacing(1),
    borderBottom: "1px solid #ddd"
  },
  userDetails: {
    display: "flex",
    padding: theme.spacing(1),
    margin: theme.spacing(1),
    borderBottom: "1px solid #ddd"
  },
  userImage: {
    heigth: theme.spacing(3),
    width: theme.spacing(3),
    margin: theme.spacing(0, 1, 0, 0)
  },
  signout: {
    color: "#fff",
    padding: theme.spacing(1),
    margin: theme.spacing(1)
  }
});

function _Menu({ classes, user }) {
  const handleSignout = () => {
    localStorage.removeItem("secret");
    Router.push("/login");
  };
  return (
    <div>
      <div className={classes.dashboard}>Dashboard</div>
      <div className={classes.issues}>Issues</div>
      <div className={classes.userDetails}>
        <img src={user.image} className={classes.userImage} />
        <div>{user.username}</div>
      </div>
      <div>
        <Link href="/login">
          <Button className={classes.signout} onClick={handleSignout}>
            Signout
          </Button>
        </Link>
      </div>
    </div>
  );
}

const Menu = withStyles(menuStyles)(_Menu);

const navbarStyles = theme => ({
  container: {
    backgroundColor: "black",
    padding: theme.spacing(1)
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
      navbar: false,
      user: undefined
    };
  }

  fetchUser = async () => {
    try {
      const response = await axios.get("/user", authHeaders());
      if (response.status === 200) {
        this.setState({ user: response.data });
      }
    } catch (error) {
      console.log(error);
    }
  };

  componentDidMount() {
    this.fetchUser();
  }

  handleClick = event => {
    this.setState({
      navbar: !this.state.navbar
    });
  };

  render() {
    const { classes } = this.props;
    const { navbar, user } = this.state;
    const menu = navbar ? <Menu user={user} /> : <div></div>;

    return (
      <AppBar position="static" className={classes.container}>
        <div className={cx(classes.container, "container")}>
          <div className="row">
            <div className="col-12">
              <div className={classes.navbar}>
                <Button
                  onClick={this.handleClick}
                  className={classes.Button}
                  type="button"
                >
                  <MenuIcon />
                </Button>
                <GitHubIcon />
                <NotificationsIcon />
              </div>
              {menu}
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

export function authHeaders() {
  return {
    headers: {
      Authorization: `Bearer ${localStorage.getItem("secret")}`
    }
  };
}
