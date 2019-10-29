import React, { Component, Fragment } from "react";

import Router from "next/router";

import Link from "next/link";

import axios from "axios";

import { withStyles } from "@material-ui/core/styles";

import { Button, Typography } from "@material-ui/core";

import TextField from "@material-ui/core/TextField";

import cx from "classnames";

import { Navbar } from "../../utils/utils.js";

const userStyles = theme => ({
  image: {
    heigth: theme.spacing(8),
    width: theme.spacing(8)
  },
  container: {
    padding: theme.spacing(2)
  },
  repo: {
    border: "1px solid #ddd",
    padding: theme.spacing(1),
    margin: theme.spacing(1)
  },
  navbar: {
    marginBottom: theme.spacing(1)
  }
});

class _User extends Component {
  constructor(props) {
    super(props);
    this.state = {
      user: undefined,
      repos: []
    };
  }

  componentDidMount() {
    axios
      .get("/user", {
        headers: {
          Authorization: `Bearer ${localStorage.getItem("secret")}`
        }
      })
      .then(response =>
        this.setState({
          user: response.data
        })
      )
      .catch(error => {
        console.log(error);
      });
    axios
      .get("/user/repos", {
        headers: {
          Authorization: `Bearer ${localStorage.getItem("secret")}`
        }
      })
      .then(response =>
        this.setState({
          repos: response.data
        })
      )
      .catch(error => {
        console.log(error);
      });
  }

  render() {
    const { user, repos } = this.state;
    const { classes } = this.props;

    if (user === undefined) {
      return null;
    }
    return (
      <Fragment>
        <Navbar className={classes.navbar} />
        <div className={cx(classes.container, "container")}>
          <div className="row">
            <div className="col-12">
              <div>
                <img src={user.image} className={classes.image} />
                <Typography variant="body2">{user.name}</Typography>
                <Typography variant="body2">{user.username}</Typography>
              </div>
              <div>
                <Link href="/new">
                  <Button color="primary" variant="contained">
                    <a>New</a>
                  </Button>
                </Link>
              </div>
              <div>
                {repos.map(repo => (
                  <div key={repo.id} className={classes.repo}>
                    <Link
                      href={`/user/${user.username}/repos/${repo.name}/issues`}
                    >
                      <a>{repo.name}</a>
                    </Link>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      </Fragment>
    );
  }
}

const User = withStyles(userStyles)(_User);

export default User;
