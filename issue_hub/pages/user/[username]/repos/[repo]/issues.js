import React, { Component, Fragment } from "react";
import Router from "next/router";
import Link from "next/link";
import axios from "axios";
import { formatDistance, parseISO } from "date-fns";

import { withStyles } from "@material-ui/core/styles";
import { Button, Typography } from "@material-ui/core";
import TextField from "@material-ui/core/TextField";
import cx from "classnames";

import { Navbar, authHeaders } from "../../../../../utils/utils.js";

const issueStyles = theme => ({
  issueDisplay: {
    display: "flex",
    justifyContent: "space-between"
  },
  issue: {
    borderTop: "1px solid #ddd",
    padding: theme.spacing(2),
    [theme.breakpoints.up("md")]: {
      border: "1px solid #ddd",
      marginBottom: 0
    }
  },
  title: {
    fontWeight: "bold"
  }
});

function _Issue(props) {
  const { username, repo } = Router.router.query;
  const { classes, issue } = props;
  return (
    <div key={issue.id} className={classes.issue}>
      <div className={classes.issueDisplay}>
        <Link
          href={`/user/${username}/repos/${repo}/issues/${issue.issue_number}`}
        >
          <a className={classes.title}>{issue.title}</a>
        </Link>
        <div>#{issue.issue_number}</div>
      </div>
      <div>
        opened {formatDistance(Date.now(), parseISO(issue.created_at))} ago by{" "}
        {issue.user.username}
      </div>
    </div>
  );
}

const Issue = withStyles(issueStyles)(_Issue);

const issuesStyles = theme => ({
  mainSection: {
    marginTop: theme.spacing(4),
    padding: "0px",
    [theme.breakpoints.up("md")]: {
      margin: "25px auto"
    }
  },
  button: {
    marginBottom: theme.spacing(2)
  },

  issueStatus: {
    display: "flex"
  },
  status: {
    padding: theme.spacing(2)
  },
  newIssue: {
    [theme.breakpoints.up("md")]: {
      "&:hover": {
        backgroundColor: "red"
      },
      backgroundColor: "#2cbe4e",
      padding: theme.spacing(1)
    }
  },
  newIssueButton: {
    [theme.breakpoints.up("md")]: {
      display: "flex",
      justifyContent: "flex-end",
      marginBottom: theme.spacing(2.5)
    }
  },
  pinnedTitle: {
    [theme.breakpoints.up("md")]: {
      padding: theme.spacing(1.5)
    }
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

  fetchIssues = async () => {
    const { username, repo } = Router.router.query;
    try {
      const response = await axios.get(
        `/repos/${username}/${repo}/issues`,
        authHeaders()
      );
      if (response.status === 200) {
        this.setState({
          issues: response.data
        });
      }
    } catch (error) {
      console.log(error);
    }
  };

  componentDidMount() {
    this.fetchIssues();
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
              <div>
                <div className={classes.pinnedTitle}> Pinned Issues</div>
              </div>
              <div className="d-none d-md-block">
                <Link href={`/user/${username}/repos/${repo}/issues/new`}>
                  <div className={classes.newIssueButton}>
                    <Button
                      color="primary"
                      variant="contained"
                      className={classes.newIssue}
                    >
                      <a>New issue</a>
                    </Button>
                  </div>
                </Link>
              </div>
              <div className={cx(classes.issueStatus, " d-md-none")}>
                <a className={classes.status}>Open</a>
                <a className={classes.status}>Closed</a>
                <a className={classes.status}>Yours</a>
              </div>
              {issues.map(issue => (
                <Issue key={issue.id} issue={issue} />
              ))}
            </div>
          </div>
        </div>
      </Fragment>
    );
  }
}

const Issues = withStyles(issuesStyles)(_Issues);

export default Issues;
