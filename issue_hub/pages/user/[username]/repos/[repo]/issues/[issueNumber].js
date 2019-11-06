import React, { Component, Fragment } from "react";

import Router from "next/router";

import Link from "next/link";

import axios from "axios";

import { withStyles } from "@material-ui/core/styles";

import { Button, Typography } from "@material-ui/core";

import TextField from "@material-ui/core/TextField";

import cx from "classnames";

const commentStyles = theme => ({
  image: {
    height: theme.spacing(3),
    width: theme.spacing(3),
    marginRight: theme.spacing(1)
  },
  comment: {
    border: "1px solid #ddd",
    marginBottom: theme.spacing(1),
    marginTop: theme.spacing(1)
  },
  commentUser: {
    display: "flex",
    backgroundColor: "#f6f8fa",
    padding: theme.spacing(1),
    borderBottom: "1px solid #ddd"
  },
  body: {
    margin: theme.spacing(1),
    padding: theme.spacing(1)
  },
  username: {
    fontWeight: "bold",
    marginRight: theme.spacing(1)
  }
});

class _Comment extends Component {
  render() {
    const { comments, classes } = this.props;
    return (
      <div>
        {comments.map(comment => (
          <div key={comment.id} className={classes.comment}>
            <div className={classes.commentUser}>
              <img src={comment.user.image} className={cx(classes.image)} />
              <Typography variant="body2" className={classes.username}>
                {comment.user.username}
              </Typography>
              <div>commented 9 days ago</div>
            </div>
            <div className={classes.body}> {comment.body}</div>
          </div>
        ))}
      </div>
    );
  }
}

const Comment = withStyles(commentStyles)(_Comment);

const issueStyles = theme => ({
  body: {
    padding: theme.spacing(1),
    margin: theme.spacing(1)
  },
  title: {
    margin: theme.spacing(1)
  },
  status: {
    marginRight: theme.spacing(1),
    fontSize: theme.spacing(1.5),
    padding: theme.spacing(0.5),
    backgroundColor: "#2cbe4e",
    color: "#fff"
  },
  image: {
    height: theme.spacing(3),
    width: theme.spacing(3),
    marginRight: theme.spacing(1)
  },
  issueBody: {
    border: "1px solid #ddd",
    marginBottom: theme.spacing(1),
    marginTop: theme.spacing(1)
  },
  username: {
    display: "flex",
    borderBottom: "1px solid #ddd",
    padding: theme.spacing(1),
    backgroundColor: "#f6f8fa"
  },
  commentSection: {
    display: "block",
    width: "100%",
    minHeight: "100px",
    maxHeight: "500px",
    marginBottom: theme.spacing(1),
    padding: theme.spacing(1),
    resize: "vertical"
  },
  editButton: {
    fontSize: theme.spacing(1.5),
    padding: theme.spacing(0.5),
    backgroundColor: "#eff3f6",
    marginRight: theme.spacing(1)
  },
  newIssue: {
    color: "#fff",
    backgroundColor: "#2cbe4e",
    fontSize: theme.spacing(1.5),
    padding: theme.spacing(0.5)
  },
  container: {
    margin: theme.spacing(1)
  },
  name: {
    fontWeight: "bold",
    marginRight: theme.spacing(1)
  },
  sidebar: {
    padding: theme.spacing(1),
    margin: theme.spacing(1),
    borderBottom: "1px solid #ddd"
  },
  issueStatus: {
    display: "flex",
    paddingBottom: theme.spacing(2.5),
    marginBottom: theme.spacing(2.5),
    borderBottom: "1px solid #ddd"
  },
  commentButton: {
    color: "#fff",
    backgroundColor: "#2cbe4e"
  },
  commentClose: {
    marginRight: theme.spacing(1),
    backgroundColor: "#eff3f6"
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
      user: undefined,
      body: "",
      status: "",
      comments: []
    };
  }

  componentDidMount() {
    const { username, repo, issueNumber } = Router.router.query;
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
    axios
      .get(`/repos/${username}/${repo}/issues/${issueNumber}/comments`)
      .then(response =>
        this.setState({
          comments: response.data
        })
      )
      .catch(error => {
        console.log(error);
      });
  }

  handleChange = event => {
    this.setState({
      [event.target.name]: event.target.value
    });
  };

  handleSubmit = event => {
    const { username, repo, issueNumber } = Router.router.query;
    event.preventDefault();
    const { user } = this.state;
    axios
      .post(
        `/repos/${username}/${repo}/issues/${issueNumber}/comments`,
        {
          body: this.state.body
        },
        {
          headers: {
            Authorization: `Bearer ${localStorage.getItem("secret")}`
          }
        }
      )
      .then(response => {
        if (response.status === 201) {
          axios
            .get(`/repos/${username}/${repo}/issues/${issueNumber}/comments`)
            .then(response =>
              this.setState({
                comments: response.data
              })
            )
            .catch(error => {
              console.log(error);
            });
          this.setState({
            body: ""
          });
        }
      })
      .catch(error => {
        console.log(error);
      });
  };

  handleClick = event => {
    const { username, repo, issueNumber } = Router.router.query;
    event.preventDefault();
    const { issue } = this.state;
    if (issue.status === "Open") {
      this.setState({
        status: "Closed"
      });
    } else {
      this.setState({
        status: "Open"
      });
    }

    console.log(this.state.status);
    axios
      .put(
        `/repos/${username}/${repo}/issues/${issueNumber}`,
        {
          status: this.state.status
        },
        {
          headers: {
            Authorization: `Bearer ${localStorage.getItem("secret")}`
          }
        }
      )
      .then(response => {
        if (response.status === 204) {
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
        }
      })
      .catch(error => {
        console.log(error);
      });
  };

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
      <div className={cx(classes.container, "container")}>
        <div className="row">
          <div className="col-lg-10">
            <div>
              <Button variant="contained" className={classes.editButton}>
                Edit
              </Button>
              <Button
                variant="contained"
                color="primary"
                className={classes.newIssue}
              >
                New issue
              </Button>
            </div>
            <Typography variant="h5" className={classes.title}>
              {issue.title} #{issue.issue_number}
            </Typography>
            <div className={classes.issueStatus}>
              <Button variant="contained" className={classes.status}>
                {issue.status}
              </Button>
              <Typography variant="body2">
                opened this issue 9 days ago
              </Typography>
            </div>
            <div>
              <div className={classes.issueBody}>
                <div className={classes.username}>
                  <img src={issue.user.image} className={cx(classes.image)} />
                  <Typography variant="body2" className={classes.name}>
                    {issue.user.username}
                  </Typography>
                  <div>commented 9 days ago</div>
                </div>
                <Typography variant="body2" className={classes.body}>
                  {issue.body}
                </Typography>
              </div>
              <Comment comments={this.state.comments} />
              <form onSubmit={this.handleSubmit}>
                <div>
                  <textarea
                    className={cx(classes.commentSection, "form-control")}
                    name="body"
                    value={this.state.body}
                    onChange={this.handleChange}
                  />
                  <Button
                    variant="contained"
                    onClick={this.handleClick}
                    className={classes.commentClose}
                  >
                    close issue
                  </Button>
                  <Button
                    variant="contained"
                    className={classes.commentButton}
                    type="submit"
                  >
                    Comment
                  </Button>
                </div>
              </form>
            </div>
          </div>
          <div className="col-lg-2">
            <div className={classes.sidebar}>Assignee</div>
            <div className={classes.sidebar}>Labels</div>
            <div className={classes.sidebar}>Projects</div>
            <div className={classes.sidebar}>Milestone</div>
            <div className={classes.sidebar}>Lock conversation</div>
            <div className={classes.sidebar}>Pin</div>
            <div className={classes.sidebar}>Delete Issue</div>
          </div>
        </div>
      </div>
    );
  }
}

const Issue = withStyles(issueStyles)(_Issue);

export default Issue;
