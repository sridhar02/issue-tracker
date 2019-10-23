import React, { Component, Fragment } from "react";

import Router from "next/router";

import Link from "next/link";

import axios from "axios";

import { withStyles } from "@material-ui/core/styles";

import { Button, Typography } from "@material-ui/core";

import TextField from "@material-ui/core/TextField";

import cx from "classnames";

const issueStyles = theme => ({
  issue: {
    border: "1px solid #ddd",
    padding: theme.spacing(1),
    marginBottom: theme.spacing(1),
    marginRight: theme.spacing(2)
  },
  mainSection: {
    margin: theme.spacing(2)
  },
  button: {
    marginBottom: theme.spacing(2)
  }
});

class _Issues extends Component {
  static getInitialProps({ query }) {
    return { query };
  }

  constructor(props) {
    super(props);
    this.state = {
      issues: []
    };
  }

  componentDidMount() {
    const { username, repo } = Router.router.query;
    console.log(username);
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
    const { classes } = this.props;
    const { issues } = this.state;
    return (
      <div className={cx(classes.mainSection, "container")}>
        <div className="row">
          <div className="col-12">
            <div className={classes.button}>
              <Button color="primary" variant="contained">
                NEW
              </Button>
            </div>
            {issues.map(issue => (
              <div key={issue.id} className={classes.issue}>
                <Link href="/kadssadasd">
                  <a>{issue.title}</a>
                </Link>
                <div>#{issue.issue_number} opened 3 days ago by sridhar02</div>
              </div>
            ))}
          </div>
        </div>
      </div>
    );
  }
}

const Issues = withStyles(issueStyles)(_Issues);

export default Issues;
