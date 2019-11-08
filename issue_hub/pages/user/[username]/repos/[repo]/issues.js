import React, { Component, Fragment } from "react";

import Router from "next/router";

import Link from "next/link";

import axios from "axios";

import { withStyles } from "@material-ui/core/styles";

import { Button, Typography } from "@material-ui/core";

import TextField from "@material-ui/core/TextField";

import cx from "classnames";

import { Navbar } from "../../../../../utils/utils.js";

const issueStyles = theme => ({
  issue: {
    borderTop: "1px solid #ddd",
    // borderBottom: "1px solid #ddd",
    padding: theme.spacing(2),
    marginBottom: theme.spacing(1)
  },
  mainSection: {
    marginTop: theme.spacing(4),
    padding: "0px"
  },
  button: {
    marginBottom: theme.spacing(2)
  },
  issueDisplay: {
    display: "flex",
    justifyContent: "space-between"
  },
  title: {
    fontWeight: "bold"
  },
  issueStatus: {
    display: "flex"
  },
  status: {
    padding: theme.spacing(2)
  }
});

class _Issues extends Component {
  static getInitialProps({ query }) {
    return { query, username: query.username, repo: query.repo };
  }

  constructor(props) {
    super(props);
    this.state = {
      issues: []
    };
  }

  componentDidMount() {
    const { username, repo } = Router.router.query;
    axios
      .get(`/repos/${username}/${repo}/issues`, {
        headers: {
          Authorization: `Bearer ${localStorage.getItem("secret")}`
        }
      })
      .then(response =>
        this.setState({
          issues: response.data
        })
      )
      .catch(error => {
        console.log(error);
      });
  }

  render() {
    const { classes, username, repo } = this.props;
    const { issues } = this.state;
    if (issues === undefined) {
      return null;
    }
    return (
      <Fragment>
        <Navbar user={issues.user} />
        <div className={cx(classes.mainSection, "container")}>
          <div className="row">
            <div className="col-12">
              <div className="d-none">
                <Link href={`/user/${username}/repos/${repo}/issues/new`}>
                  <Button color="primary" variant="contained">
                    <a>New issue</a>
                  </Button>
                </Link>
              </div>
              <div className={classes.issueStatus}>
                <a className={classes.status}>Open</a>
                <a className={classes.status}>Closed</a>
                <a className={classes.status}>Yours</a>
              </div>
              {issues.map(issue => (
                <div key={issue.id} className={classes.issue}>
                  <div className={classes.issueDisplay}>
                    <Link
                      href={`/user/${username}/repos/${repo}/issues/${issue.issue_number}`}
                    >
                      <a className={classes.title}>{issue.title}</a>
                    </Link>
                    <div>#{issue.issue_number}</div>
                  </div>
                  <div> opened 3 days ago by sridhar02</div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </Fragment>
    );
  }
}

const Issues = withStyles(issueStyles)(_Issues);

export default Issues;
