import React, { Component, Fragment } from "react";

import Router from "next/router";

import Link from "next/link";

import axios from "axios";

import { withStyles } from "@material-ui/core/styles";

import { Button, Typography } from "@material-ui/core";

import TextField from "@material-ui/core/TextField";

import cx from "classnames";

const issueStyles = theme => ({
  body: {
    padding: theme.spacing(1),
    margin: theme.spacing(1)
  },
  title: {
    margin: theme.spacing(1)
  },
  status: {
    backgroundColor: "green"
  },
  image: {
    height: theme.spacing(6),
    width: theme.spacing(6)
  },
  issueBody: {
    border: "1px solid #ddd",
    padding: theme.spacing(1),
    margin: theme.spacing(1)
  },
  username: {
    borderBottom: "1px solid #ddd",
    backgroundColor: "#c0d3eb"
  }
});

class _Issue extends Component {
  static getInitialProps({ query }) {
    return { query };
  }
  constructor(props) {
    super(props);
    this.state = {
      issue: undefined,
      user: undefined
    };
  }

  componentDidMount() {
    const { username, repo, issueNumber } = Router.router.query;
    console.log(Router.router);
    axios
      .get(`/repos/${username}/${repo}/issues/${issueNumber}`, {
        headers: {
          Authorization: `Bearer ${localStorage.getItem("secret")}`
        }
      })
      .then(response =>
        this.setState({
          issue: response.data
        })
      )
      .catch(error => {
        console.log(error);
      });
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
  }

  render() {
    const { issue, user } = this.state;
    const { classes } = this.props;
    if (issue == undefined) {
      return null;
    }
    if (user == undefined) {
      return null;
    }
    return (
      <div className="container">
        <div className="row">
          <div className="col-12">
            <Typography variant="h6" className={classes.title}>
              {issue.title} #{issue.issue_number}
            </Typography>
            <div>
              <Button variant="contained">Edit</Button>
            </div>
            <div>
              <Button variant="contained" color="primary">
                New issue
              </Button>
            </div>
            <div>
              <Button className={classes.status}>{issue.status}</Button>
              <Typography variant="body2">
                opened this issue 9 days ago
              </Typography>
            </div>
            <div>
              <img src={user.image} className={classes.image} />
              <div className={classes.issueBody}>
                <div className={classes.username}>
                  {user.username}commented 9 days ago
                </div>
                <Typography variant="h6" className={classes.body}>
                  {issue.body}
                </Typography>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }
}

const Issue = withStyles(issueStyles)(_Issue);

export default Issue;
