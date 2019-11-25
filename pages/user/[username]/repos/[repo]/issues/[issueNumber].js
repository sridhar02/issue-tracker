import React, { Component, Fragment } from 'react';
import Router from 'next/router';
import Link from 'next/link';
import axios from 'axios';
import cx from 'classnames';
import { formatDistance, parseISO } from 'date-fns';

import SettingsIcon from '@material-ui/icons/Settings';
import { withStyles } from '@material-ui/core/styles';
import {
  Button,
  Typography,
  TextField,
  Popover,
  Paper
} from '@material-ui/core';
import LockIcon from '@material-ui/icons/Lock';
import DeleteIcon from '@material-ui/icons/Delete';
import ErrorOutlineIcon from '@material-ui/icons/ErrorOutline';
import CheckIcon from '@material-ui/icons/Check';

import { TextInput } from '@primer/components';

import { Navbar, authHeaders } from '../../../../../../utils/utils.js';

const commentStyles = theme => ({});

function _Comment({ comments, classes }) {
  return (
    <div>
      {comments.map(comment => (
        <div key={comment.id}>
          <Body
            user={comment.user}
            body={comment.body}
            updated_at={comment.created_at}
          />
        </div>
      ))}
    </div>
  );
}

const Comment = withStyles(commentStyles)(_Comment);

const postCommentStyles = theme => ({
  commentSection: {
    display: 'block',
    width: '100%',
    minHeight: '100px',
    maxHeight: '500px',
    marginBottom: theme.spacing(1),
    padding: theme.spacing(1),
    resize: 'vertical'
  },
  commentButton: {
    color: '#fff',
    backgroundColor: '#2cbe4e',
    '&:hover': {
      backgroundColor: 'green'
    }
  },
  commentClose: {
    marginRight: theme.spacing(1),
    backgroundColor: '#eff3f6'
  },
  commentOpen: {
    marginRight: theme.spacing(1),
    backgroundColor: '#eff3f6'
  }
});

class _PostComments extends Component {
  constructor(props) {
    super(props);
    this.state = {
      body: ''
    };
  }
  handleChange = event => {
    this.setState({
      [event.target.name]: event.target.value
    });
  };

  postCommenthandler = async event => {
    const { username, repo, issueNumber } = Router.router.query;
    const { fetchComments } = this.props;
    event.preventDefault();
    const { user } = this.state;
    try {
      const response = await axios.post(
        `/repos/${username}/${repo}/issues/${issueNumber}/comments`,
        {
          body: this.state.body
        },
        authHeaders()
      );
      if (response.status === 201) {
        fetchComments();
        this.setState({
          body: ''
        });
      }
    } catch (error) {
      console.log(error);
    }
  };

  toggleIssueStatus = async event => {
    const { username, repo, issueNumber } = Router.router.query;
    const { issue, fetchIssue } = this.props;
    const newStatus =
      issue.status === STATUS_OPEN ? STATUS_CLOSED : STATUS_OPEN;
    try {
      const response = await axios.put(
        `/repos/${username}/${repo}/issues/${issueNumber}`,
        {
          status: newStatus
        },
        authHeaders()
      );
      if (response.status === 204) {
        fetchIssue();
      }
    } catch (error) {
      console.log(error);
    }
  };

  render() {
    const { classes, issue } = this.props;
    const closeButton = (
      <Button
        variant="contained"
        onClick={this.toggleIssueStatus}
        className={
          issue.status === STATUS_OPEN
            ? classes.commentOpen
            : classes.commentClose
        }
      >
        {issue.status === STATUS_OPEN && (
          <ErrorOutlineIcon className={classes.statusCloseIcon} />
        )}
        {issue.status === STATUS_OPEN ? 'Close' : 'Reopen'} issue
      </Button>
    );

    return (
      <form onSubmit={this.postCommenthandler}>
        <div>
          <textarea
            className={cx(classes.commentSection, 'form-control')}
            name="body"
            value={this.state.body}
            onChange={this.handleChange}
          />
          {closeButton}
          <Button
            variant="contained"
            className={classes.commentButton}
            type="submit"
          >
            Comment
          </Button>
        </div>
      </form>
    );
  }
}
const PostComments = withStyles(postCommentStyles)(_PostComments);

const assigneePopperStyles = theme => ({
  popper: {
    padding: theme.spacing(2),
    width: '250px',
    backgroundColor: '#F6F8FA'
  },
  collaboratorDetails: {
    display: 'flex',
    justifyContent: 'flex-start',
    width: '100%',
    marginTop: theme.spacing(1)
  },
  collaboratorImage: {
    height: theme.spacing(3),
    width: theme.spacing(3),
    marginRight: theme.spacing(1),
    marginLeft: theme.spacing(1)
  },
  paperBottom: {
    backgroundColor: '#fff'
  },
  assigneeText: {
    fontWeight: 'bold'
  },
  assigneeImage: {
    height: theme.spacing(3),
    width: theme.spacing(3),
    marginBottom: theme.spacing(1),
    marginRight: theme.spacing(2)
  },
  assigneeButton: {
    backgroundColor: '#fff',
    border: 0,
    marginBottom: theme.spacing(1),
    width: '100%',
    display: 'flex',
    justifyContent: 'space-between',
    '&:hover': {
      border: 0
    }
  },
  assigneeClick: {
    '&:hover': {
      backgroundColor: 'blue',
      border: 0
    }
  },
  assigneeIcon: {
    width: '20px'
  },
  collaborator: {
    display: 'flex',
    justifyContent: 'flex-end'
  }
});

class _Assignee extends Component {
  constructor(props) {
    super(props);
    this.state = {
      anchorEl: null,
      open: false,
      collaborators: [],
      addAssignees: [],
      removeAssignees: []
    };
  }

  togglePopper = async event => {
    const { open, addAssignees, removeAssignees, assigneed } = this.state;
    const { fetchIssue } = this.props;
    if (open) {
      if (removeAssignees.length !== 0) {
        await this.deleteAssignee(removeAssignees);
      }
      if (addAssignees.length !== 0) {
        await this.postAssignee(addAssignees);
      }
      if (removeAssignees.length !== 0 || addAssignees.length !== 0) {
        fetchIssue();
      }
      this.setState({
        anchorEl: null,
        open: !open
      });
    } else {
      const { currentTarget } = event;
      this.setState({
        anchorEl: currentTarget,
        open: !open
      });
    }
  };

  fetchCollaborators = async () => {
    const { username, repo } = Router.router.query;
    try {
      const response = await axios.get(
        `/repos/${username}/${repo}/collaborators`,
        authHeaders()
      );
      if (response.status === 200) {
        this.setState({
          collaborators: response.data
        });
      }
    } catch (error) {
      console.log(error);
    }
  };

  componentDidMount() {
    this.fetchCollaborators();
  }

  postAssignee = async addAssignees => {
    event.preventDefault();
    const { username, repo, issueNumber } = Router.router.query;
    try {
      const response = await axios.post(
        `/repos/${username}/${repo}/issues/${issueNumber}/assignees`,
        {
          usernames: addAssignees
        },
        authHeaders()
      );
      if (response.status === 201) {
        this.setState({
          addAssignees: []
        });
      }
    } catch (error) {
      console.log(error);
    }
  };

  deleteAssignee = async removeAssignees => {
    const { username, repo, issueNumber } = Router.router.query;
    try {
      const response = await axios({
        method: 'delete',
        url: `/repos/${username}/${repo}/issues/${issueNumber}/assignees`,
        headers: authHeaders().headers,
        data: { usernames: removeAssignees }
      });
      if (response.status === 204) {
        this.setState({ removeAssignees: [] });
      }
    } catch (error) {
      console.log(error);
    }
  };

  checkAssignee = collaborator => {
    const { issue } = this.props;
    const { addAssignees, removeAssignees } = this.state;

    // if collaborator is in issue(assignees) && removeAssignees not includes collaborator then  add collaborator to removeAssignees list.
    // if collaborator is not in issue assignees list && removeAssignees not includes collaborator then add collaborator to addAssignee list.
    // if collaborator is not in issue Assignees list then remove collaborator from addAssignees list && collaborator is in addAssignees list.
    // if collaborator is in issue Assignees list them remove from removeAssignees list && collaborator is in removeAssignees list.

    const alreadyAssigned = issue.assignees.some(
      assignee => assignee.user.username === collaborator.user.username
    );

    const alreadyInAdd = addAssignees.includes(collaborator.user.username);

    const alreadyInRemove = removeAssignees.includes(
      collaborator.user.username
    );

    if (alreadyAssigned && !alreadyInRemove) {
      this.setState({
        removeAssignees: [...removeAssignees, collaborator.user.username]
      });
    } else if (!alreadyAssigned && !alreadyInAdd) {
      this.setState({
        addAssignees: [...addAssignees, collaborator.user.username]
      });
    } else if (!alreadyAssigned && alreadyInAdd) {
      this.setState({
        addAssignees: addAssignees.filter(
          addAssignee => addAssignee !== collaborator.user.username
        )
      });
    } else if (alreadyAssigned && alreadyInRemove) {
      this.setState({
        removeAssignees: removeAssignees.filter(
          removeAssignee => removeAssignee !== collaborator.user.username
        )
      });
    }
  };

  render() {
    const { classes, issue } = this.props;
    const {
      anchorEl,
      open,
      collaborators,
      addAssignees,
      removeAssignees
    } = this.state;
    return (
      <Fragment>
        <button
          type="button"
          onClick={this.togglePopper}
          className={classes.assigneeButton}
        >
          Assignee
          <div>
            <SettingsIcon />
          </div>
        </button>
        <div>
          {issue.assignees.map(assignee => (
            <div key={assignee.user.id}>
              <img
                src={assignee.user.image}
                className={classes.assigneeImage}
              />
              {assignee.user.username}
            </div>
          ))}
        </div>
        <Popover open={open} anchorEl={anchorEl} onClose={this.togglePopper}>
          <Paper className={classes.popper}>
            <div className={classes.paperTop}>
              <Typography variant="body2" className={classes.assigneeText}>
                Assign up to 10 people to this issue
              </Typography>
              <TextInput />
            </div>
            <div className={classes.paperBottom}>
              {collaborators.map(collaborator => {
                const assigned =
                  (issue.assignees.some(
                    assignee =>
                      assignee.user.username === collaborator.user.username
                  ) ||
                    addAssignees.includes(collaborator.user.username)) &&
                  !removeAssignees.includes(collaborator.user.username);
                return (
                  <div className={classes.assigneeClick}>
                    <Button
                      key={collaborator.user.username}
                      className={classes.collaboratorDetails}
                      onClick={() => this.checkAssignee(collaborator)}
                    >
                      <div className={classes.assigneeIcon}>
                        {assigned && <CheckIcon />}
                      </div>
                      <div className={classes.collaborator}>
                        <img
                          key={collaborator.user.image}
                          src={collaborator.user.image}
                          className={classes.collaboratorImage}
                        />
                        <div key={collaborator.user.username}>
                          {collaborator.user.username}
                        </div>
                      </div>
                    </Button>
                  </div>
                );
              })}
            </div>
          </Paper>
        </Popover>
      </Fragment>
    );
  }
}

const Assignee = withStyles(assigneePopperStyles)(_Assignee);

const sidebarStyles = theme => ({
  sidebar: {
    padding: theme.spacing(1),
    margin: theme.spacing(1),
    borderBottom: '1px solid #ddd'
  }
});

class _Sidebar extends Component {
  toggleLockIssue = async event => {
    const { username, repo, issueNumber } = Router.router.query;
    const { issue, fetchIssue } = this.props;
    const lockStatus = issue.lock === LOCK_LOCK ? LOCK_UNLOCK : LOCK_LOCK;
    try {
      const response = await axios.put(
        `/repos/${username}/${repo}/issues/${issueNumber}/lock`,
        {
          lock: lockStatus
        },
        authHeaders()
      );
      if (response.status === 204) {
        fetchIssue();
      }
    } catch (error) {
      console.log(error);
    }
  };

  togglePinIssue = async event => {
    const { username, repo, issueNumber } = Router.router.query;
    const { issue, fetchIssue } = this.props;
    const pinStatus = issue.pinned === PIN_UNPIN ? PIN_PIN : PIN_UNPIN;
    try {
      const response = await axios.put(
        `/repos/${username}/${repo}/issues/${issueNumber}/pin`,
        {
          pinned: pinStatus
        },
        authHeaders()
      );
      if (response.status === 204) {
        fetchIssue();
      }
    } catch (error) {
      console.log(error);
    }
  };

  render() {
    const { classes, issue, fetchIssue } = this.props;
    const lockButton = (
      <Button onClick={this.toggleLockIssue} className={classes.lockButton}>
        <LockIcon /> {issue.lock === LOCK_UNLOCK ? 'Lock' : 'Unlock'}
        conversation
      </Button>
    );
    const pinButton = (
      <Button onClick={this.togglePinIssue}>
        {issue.pinned === PIN_PIN ? 'Unpin' : 'Pin'} Issue
      </Button>
    );

    return (
      <Fragment>
        <div className={classes.sidebar}>
          <Assignee issue={issue} fetchIssue={fetchIssue} />
        </div>
        <div className={classes.sidebar}>Labels</div>
        <div className={classes.sidebar}>Projects</div>
        <div className={classes.sidebar}>Milestone</div>
        <div className={classes.sidebar}>{lockButton}</div>
        <div className={classes.sidebar}>{pinButton}</div>
        <div className={classes.sidebar}>
          <Button>
            <DeleteIcon /> Delete Issue
          </Button>
        </div>
      </Fragment>
    );
  }
}

const Sidebar = withStyles(sidebarStyles)(_Sidebar);

const issueHeaderStyles = theme => ({
  issueHeading: {
    display: 'flex'
  },
  title: {
    margin: theme.spacing(1)
  },
  issueStatus: {
    display: 'flex',
    paddingBottom: theme.spacing(2.5),
    marginBottom: theme.spacing(2.5),
    borderBottom: '1px solid #ddd',
    [theme.breakpoints.up('md')]: {
      order: 1
    }
  },
  issueTitle: {
    marginTop: theme.spacing(1)
  },
  titleEditButtons: {
    display: 'flex'
  },
  save: {
    backgroundColor: '#2cbe4e',
    marginBottom: theme.spacing(1),
    marginTop: theme.spacing(1)
  },
  cancel: {
    margin: theme.spacing(1, 0, 1, 1)
  },
  editButton: {
    fontSize: theme.spacing(1.5),
    padding: theme.spacing(0.5),
    backgroundColor: '#eff3f6',
    marginRight: theme.spacing(1)
  },
  newIssue: {
    color: '#fff',
    backgroundColor: '#2cbe4e',
    fontSize: theme.spacing(1.5),
    padding: theme.spacing(0.5),
    fontWeight: 'bold',
    '&:hover': {
      backgroundColor: 'green'
    }
  },
  headerTop: {
    [theme.breakpoints.up('md')]: {
      display: 'flex',
      justifyContent: 'space-between'
    }
  },
  buttonsHeader: {
    [theme.breakpoints.up('md')]: {
      order: 2
    }
  },
  statusOpen: {
    marginRight: theme.spacing(1),
    fontSize: theme.spacing(1.5),
    padding: theme.spacing(0.5),
    backgroundColor: '#2cbe4e',
    color: '#fff',
    fontWeight: 'bold',
    '&:hover': {
      backgroundColor: '#2cbe4e'
    }
  },
  statusClose: {
    marginRight: theme.spacing(1),
    fontSize: theme.spacing(1.5),
    padding: theme.spacing(0.5),
    backgroundColor: '#CB2431',
    color: '#fff',
    fontWeight: 'bold',
    '&:hover': {
      backgroundColor: '#CB2431'
    }
  }
});

class _IssueHeader extends Component {
  constructor(props) {
    super(props);
    this.state = {
      editTitle: false,
      title: ''
    };
  }

  toggleTitle = event => {
    event.preventDefault();
    this.setState({
      editTitle: !this.state.editTitle
    });
  };

  cancelIssueEdit = event => {
    this.setState({
      editTitle: false
    });
  };

  handleChange = event => {
    this.setState({
      [event.target.name]: event.target.value
    });
  };

  updateIssueTitle = async event => {
    const { fetchIssue } = this.props;
    event.preventDefault();
    console.log(this.state.title);
    const { username, repo, issueNumber } = Router.router.query;
    try {
      const response = await axios.put(
        `/repos/${username}/${repo}/issues/${issueNumber}`,
        {
          title: this.state.title
        },
        authHeaders()
      );
      if (response.status === 204) {
        fetchIssue();
        this.setState({
          editTitle: false
        });
      }
    } catch (error) {
      console.log(error);
    }
  };

  render() {
    const { classes, issue, fetchIssue } = this.props;
    const { username, repo, issueNumber } = Router.router.query;
    const { editTitle } = this.state;
    const title =
      editTitle === true ? (
        <div>
          <form onSubmit={this.updateIssueTitle}>
            <TextField
              name="title"
              variant="outlined"
              placeholder={issue.title}
              value={this.state.title}
              onChange={this.handleChange}
              className={classes.issueTitle}
            />
            <div className={classes.titleEditButtons}>
              <Button
                className={classes.save}
                variant="contained"
                type="submit"
              >
                Save
              </Button>
              <Button
                className={classes.cancel}
                variant="contained"
                onClick={this.cancelIssueEdit}
              >
                Cancel
              </Button>
            </div>
          </form>
        </div>
      ) : (
        <Typography variant="h5" className={classes.title}>
          {issue.title}
        </Typography>
      );
    const status = (
      <Button
        variant="contained"
        className={
          issue.status === STATUS_OPEN
            ? classes.statusOpen
            : classes.statusClose
        }
      >
        <ErrorOutlineIcon className={classes.statusIcon} /> {issue.status}
      </Button>
    );

    return (
      <Fragment>
        <div className={classes.headerTop}>
          <div className={classes.buttonsHeader}>
            <Button
              variant="contained"
              className={classes.editButton}
              onClick={this.toggleTitle}
            >
              Edit
            </Button>
            <Link href={`/user/${username}/repos/${repo}/issues/new`}>
              <Button
                variant="contained"
                color="primary"
                className={classes.newIssue}
              >
                New issue
              </Button>
            </Link>
          </div>
          <div className={classes.issueHeading}>
            {title}
            <Typography variant="h5" className={classes.title}>
              #{issue.issue_number}
            </Typography>
          </div>
        </div>
        <div className={classes.issueStatus}>
          {status}
          <Typography variant="body2">
            opened this issue{' '}
            {formatDistance(Date.now(), parseISO(issue.created_at))} ago
          </Typography>
        </div>
      </Fragment>
    );
  }
}
const IssueHeader = withStyles(issueHeaderStyles)(_IssueHeader);

const bodyStyles = theme => ({
  image: {
    height: theme.spacing(3),
    width: theme.spacing(3),
    marginRight: theme.spacing(1)
  },
  userSection: {
    display: 'flex',
    borderBottom: '1px solid #ddd',
    padding: theme.spacing(1),
    backgroundColor: '#f6f8fa'
  },
  content: {
    border: '1px solid #ddd',
    marginBottom: theme.spacing(1),
    marginTop: theme.spacing(1)
  },
  body: {
    padding: theme.spacing(1),
    margin: theme.spacing(1)
  },
  time: {},
  username: {
    paddingRight: theme.spacing(2)
  }
});

function _Body({ classes, user, body, updated_at }) {
  return (
    <div className={classes.content}>
      <div className={classes.userSection}>
        <img src={user.image} className={classes.image} />
        <Typography variant="body2" className={classes.username}>
          {user.username}
        </Typography>
        <div className={classes.time}>
          commented {formatDistance(Date.now(), parseISO(updated_at))} ago
        </div>
      </div>
      <div>
        <Typography variant="body2" className={classes.body}>
          {body}
        </Typography>
      </div>
    </div>
  );
}

const Body = withStyles(bodyStyles)(_Body);

const STATUS_OPEN = 'Open';
const STATUS_CLOSED = 'Closed';

const LOCK_LOCK = 'Locked';
const LOCK_UNLOCK = 'Unlocked';

const PIN_PIN = 'Pinned';
const PIN_UNPIN = 'Unpinned';

const issueStyles = theme => ({
  image: {
    height: theme.spacing(3),
    width: theme.spacing(3),
    marginRight: theme.spacing(1)
  },
  issueBody: {
    border: '1px solid #ddd',
    marginBottom: theme.spacing(1),
    marginTop: theme.spacing(1)
  },
  username: {
    display: 'flex',
    borderBottom: '1px solid #ddd',
    padding: theme.spacing(1),
    backgroundColor: '#f6f8fa'
  },
  container: {
    margin: theme.spacing(1),
    [theme.breakpoints.up('md')]: {
      margin: '25px auto'
    }
  },
  name: {
    fontWeight: 'bold',
    marginRight: theme.spacing(1)
  },
  statusIcon: {
    marginRight: theme.spacing(0.5)
  },
  lockButton: {
    fontSize: ''
  },
  statusCloseIcon: {
    color: '#CB2431'
  },
  issueTitle: {
    marginTop: theme.spacing(1)
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
      title: '',
      body: '',
      editTitle: false,
      comments: []
    };
  }

  fetchIssue = async () => {
    const { username, repo, issueNumber } = Router.router.query;
    try {
      const response = await axios.get(
        `/repos/${username}/${repo}/issues/${issueNumber}`
      );
      if (response.status === 200) {
        this.setState({
          issue: response.data,
          title: response.data.title
        });
      }
    } catch (error) {
      console.log(error);
    }
  };

  fetchComments = async () => {
    const { username, repo, issueNumber } = Router.router.query;
    try {
      const response = await axios.get(
        `/repos/${username}/${repo}/issues/${issueNumber}/comments`
      );
      if (response.status === 200) {
        this.setState({ comments: response.data });
      }
    } catch (error) {
      console.log(error);
    }
  };

  componentDidMount() {
    const { username, repo, issueNumber } = Router.router.query;
    this.fetchIssue();
    this.fetchComments();
  }

  render() {
    const { issue, user, editTitle } = this.state;
    const { classes } = this.props;

    if (issue === undefined) {
      return null;
    }

    return (
      <Fragment>
        <Navbar />
        <div className={cx(classes.container, 'container')}>
          <div className="row">
            <div className="col-12">
              <div className="row">
                <div className="col-12">
                  <IssueHeader issue={issue} fetchIssue={this.fetchIssue} />
                </div>
              </div>
              <div className="row">
                <div className="col-lg-10">
                  <Body
                    user={issue.user}
                    body={issue.body}
                    updated_at={issue.updated_at}
                  />
                  <Comment comments={this.state.comments} />
                  <PostComments
                    issue={issue}
                    fetchComments={this.fetchComments}
                    fetchIssue={this.fetchIssue}
                  />
                </div>
                <div className="col-lg-2">
                  <Sidebar issue={issue} fetchIssue={this.fetchIssue} />
                </div>
              </div>
            </div>
          </div>
        </div>
      </Fragment>
    );
  }
}

const Issue = withStyles(issueStyles)(_Issue);

export default Issue;
