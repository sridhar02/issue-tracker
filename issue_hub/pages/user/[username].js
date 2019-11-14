import React, { Component, Fragment } from "react";
import Router from "next/router";
import Link from "next/link";
import axios from "axios";
import cx from "classnames";

import { withStyles } from "@material-ui/core/styles";
import { Button, Typography } from "@material-ui/core";
import TextField from "@material-ui/core/TextField";

import { Navbar } from "../../utils/utils.js";

const userStyles = theme => ({
  image: {
    heigth: theme.spacing(13.75),
    width: theme.spacing(13.75),
    margin: theme.spacing(0, 1.875, 0, 0),
    [theme.breakpoints.up("md")]: {
      heigth: theme.spacing(30.9975),
      width: theme.spacing(30.8725)
    }
  },
  container: {
    [theme.breakpoints.up("md")]: {
      margin: theme.spacing(3, 4.4375, 0)
    },
    marginBottom: 50
  },
  repoWrapper: {
    // padding: theme.spacing(2),
    fontWeight: "bold",
    width: "100%",
    "&:last-child $repo": {
      borderBottom: "1px solid #ddd"
    }
  },
  repo: {
    borderTop: "1px solid #ddd",
    borderLeft: "1px solid #ddd",
    borderRight: "1px solid #ddd",

    padding: theme.spacing(2),
    [theme.breakpoints.up("md")]: {
      border: "1px solid #ddd",
      marginBottom: theme.spacing(2)
    }
  },
  navbar: {
    marginBottom: theme.spacing(1)
  },
  userDetails: {
    display: "flex",
    padding: theme.spacing(1.875, 1.875, 1.25),
    [theme.breakpoints.up("md")]: {
      display: "block"
    }
  },
  userName: {
    [theme.breakpoints.up("md")]: {
      padding: theme.spacing(2),
      fontWeight: "bold",
      fontSize: theme.spacing(2.5),
      marginTop: theme.spacing(1)
    }
  },
  userProfile: {
    padding: "0 15px 15px",
    margin: theme.spacing(1, 0, 0),
    borderBottom: "1px solid #ddd",
    whiteSpace: "pre-wrap",
    [theme.breakpoints.up("md")]: {
      borderBottom: "0px"
    }
  },
  popular: {
    padding: theme.spacing(2, 0),
    [theme.breakpoints.up("md")]: {
      borderTop: "1px solid #ddd",
      marginTop: theme.spacing(2),
      display: "flex",
      justifyContent: "space-between"
    }
  },
  newRepo: {
    backgroundColor: "#2cbe4e",
    color: "#fff",
    fontWeight: "bold",
    "&:hover": {
      backgroundColor: "green"
    }
  }
});

class _User extends Component {
  static getInitialProps({ query }) {
    return { query };
  }
  constructor(props) {
    super(props);
    this.state = {
      user: undefined,
      repos: []
    };
  }

  fetchUser = async () => {
    const { username } = Router.router.query;
    try {
      const response = await axios.get(`/users/${username}`);
      if (response.status === 200) {
        this.setState({ user: response.data });
      }
    } catch (error) {
      console.log(error);
    }
  };

  fetchUserRepos = async () => {
    const { username } = Router.router.query;
    try {
      const response = await axios.get(`/users/${username}/repos`);
      if (response.status === 200) {
        this.setState({ repos: response.data });
      }
    } catch (error) {
      console.log(error);
    }
  };

  componentDidMount() {
    this.fetchUser();
    this.fetchUserRepos();
  }

  render() {
    const { user, repos } = this.state;
    const { classes } = this.props;
    if (user === undefined) {
      return null;
    }

    return (
      <Fragment>
        <Navbar className={classes.navbar} user={user} />
        <div className={cx(classes.container, "container")}>
          <div className="row">
            <div className="col-lg-3">
              <div className={classes.userDetails}>
                <img src={user.image} className={classes.image} />
                <div></div>
                <Typography variant="body2" className={classes.userName}>
                  {user.username}
                </Typography>
              </div>
              <div className={classes.userProfile}></div>
            </div>
            <div className="col-lg-9">
              <div className={cx(classes.popular)}>
                <div>Popular respositories</div>
                <div className="d-none d-md-block">
                  <Link href="/new">
                    <Button className={classes.newRepo} variant="contained">
                      <a>New Respository</a>
                    </Button>
                  </Link>
                </div>
              </div>
              <div className={"row"}>
                {repos.map(repo => (
                  <div
                    key={repo.id}
                    className={cx(classes.repoWrapper, "col-md-6")}
                  >
                    <div className={classes.repo}>
                      <Link
                        href={`/user/${user.username}/repos/${repo.name}/issues`}
                      >
                        <a>
                          <span className="d-md-none">
                            {repo.user.username}/
                          </span>
                          {repo.name}
                        </a>
                      </Link>
                    </div>
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
