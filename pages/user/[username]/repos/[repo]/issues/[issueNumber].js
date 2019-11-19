import React, { Component, Fragment } from 'react';
import Router from 'next/router';
import Link from 'next/link';
import axios from 'axios';
import cx from 'classnames';
import { formatDistance, parseISO } from 'date-fns';

import SettingsIcon from '@material-ui/icons/Settings';
import { withStyles } from '@material-ui/core/styles';
import { Button, Typography } from '@material-ui/core';
import TextField from '@material-ui/core/TextField';
import LockIcon from '@material-ui/icons/Lock';
import DeleteIcon from '@material-ui/icons/Delete';
import ErrorOutlineIcon from '@material-ui/icons/ErrorOutline';
import Popper from '@material-ui/core/Popper';
import Paper from '@material-ui/core/Paper';
import CheckIcon from '@material-ui/icons/Check';

import { TextInput } from '@primer/components';

import { Navbar, authHeaders } from '../../../../../../utils/utils.js';

const commentStyles = theme => ({
  image: {
    height: theme.spacing(3),
    width: theme.spacing(3),
    marginRight: theme.spacing(1)
  },
  comment: {
    border: '1px solid #ddd',
    marginBottom: theme.spacing(1),
    marginTop: theme.spacing(1)
  },
  commentUser: {
    display: 'flex',
    backgroundColor: '#f6f8fa',
    borderBottom: '1px solid #ddd'
  },
  body: {
    margin: theme.spacing(1),
    padding: theme.spacing(1)
  },
  username: {
    fontWeight: 'bold',
    marginRight: theme.spacing(1)
  }
});

function _Comment({ comments, classes }) {
  return (
    <div>
      {comments.map(comment => (
        <div key={comment.id} className={classes.comment}>
          <div className={classes.commentUser}>
            <img src={comment.user.image} className={cx(classes.image)} />
            <Typography variant="body2" className={classes.username}>
              {comment.user.username}
            </Typography>
            <div>
              commented{' '}
              {formatDistance(Date.now(), parseISO(comment.created_at))} ago
            </div>
          </div>
          <div className={classes.body}> {comment.body}</div>
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
  }
});

function _PostComments({ classes, handleSubmit, body, handleChange, button }) {
  return (
    <form onSubmit={handleSubmit}>
      <div>
        <textarea
          className={cx(classes.commentSection, 'form-control')}
          name="body"
          value={body}
          onChange={handleChange}
        />
        {button}
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

const PostComments = withStyles(postCommentStyles)(_PostComments);

const assigneePopperStyles = theme => ({
  popper: {
    padding: theme.spacing(2),
    width: '250px',
    backgroundColor: '#F6F8FA'
  },
  collaboratorDetails: {
    display: 'flex',
    margin: theme.spacing(0.5)
  },
  collaboratorImage: {
    height: theme.spacing(3),
    width: theme.spacing(3),
    marginRight: theme.spacing(1)
  },
  paperBottom: {
    backgroundColor: '#fff'
  },
  assigneeText: {
    fontWeight: 'bold'
  }
});

function _AssigneePopper({
  classes,
  id,
  open,
  anchorEl,
  collaborators,
  issue,
  postAssignee
}) {
  return (
    <Popper id={id} open={open} anchorEl={anchorEl}>
      <Paper className={classes.popper}>
        <div className={classes.paperTop}>
          <Typography variant="body2" className={classes.assigneeText}>
            Assign up to 10 people to this issue
          </Typography>
          <TextInput />
        </div>
        <div className={classes.paperBottom}>
          {collaborators.map(collaborator => (
            <Button
              key={collaborator.username}
              className={classes.collaboratorDetails}
              onClick={() => postAssignee(collaborator)}
            >
              {issue.assignees.some(
                assignee => assignee.user.username === collaborator.username
              ) && <CheckIcon />}
              <img
                key={collaborator.userImage}
                src={collaborator.userImage}
                className={classes.collaboratorImage}
              />
              <div key={collaborator.username}>{collaborator.username}</div>
            </Button>
          ))}
        </div>
      </Paper>
    </Popper>
  );
}

const AssigneePopper = withStyles(assigneePopperStyles)(_AssigneePopper);

const sidebarStyles = theme => ({
  sidebar: {
    padding: theme.spacing(1),
    margin: theme.spacing(1),
    borderBottom: '1px solid #ddd'
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
    marginBottom: theme.spacing(1)
  }
});

class _Sidebar extends Component {
  constructor(props) {
    super(props);
    this.state = {
      anchorEl: null,
      open: false,
      collaborators: [],
      assigneed: false
    };
  }

  handleClickPoper = event => {
    const { currentTarget } = event;
    this.setState(state => ({
      anchorEl: currentTarget,
      open: !state.open
    }));
  };

  fetchCollaborator = async () => {
    const { username, repo } = Router.router.query;
    try {
      const response = await axios.get(
        `/repos/${username}/${repo}/collaborators`,
        authHeaders()
      );
      if (response.status === 200) {
        this.setState({ collaborators: response.data });
      }
    } catch (error) {
      console.log(error);
    }
  };

  componentDidMount() {
    this.fetchCollaborator();
  }

  postAssignee = async collaborator => {
    const { fetchIssue } = this.props;
    event.preventDefault();
    const { username, repo, issueNumber } = Router.router.query;
    try {
      const response = await axios.post(
        `/repos/${username}/${repo}/issues/${issueNumber}/assignees`,
        {
          username: collaborator.username
        },
        authHeaders()
      );
      if (response.status === 201) {
        fetchIssue();
        this.fetchCollaborator();
      }
    } catch (error) {
      console.log(error);
    }
  };

  render() {
    const { classes, lockButton, pinButton, issue } = this.props;
    const { anchorEl, open, collaborators, assigneed } = this.state;
    const id = open ? 'simple-popper' : null;
    return (
      <Fragment>
        <div className={classes.sidebar}>
          <button
            type="button"
            aria-describedby={id}
            onClick={this.handleClickPoper}
            className={classes.assigneeButton}
          >
            Assignee <SettingsIcon />
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
          <AssigneePopper
            id={id}
            open={open}
            anchorEl={anchorEl}
            collaborators={collaborators}
            issue={issue}
            postAssignee={this.postAssignee}
          />
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
    const { classes, issue, status, fetchIssue } = this.props;
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
            {formatDistance(Date.now(), parseISO(issue.created_at))}
            ago
          </Typography>
        </div>
      </Fragment>
    );
  }
}
const IssueHeader = withStyles(issueHeaderStyles)(_IssueHeader);

const STATUS_OPEN = 'Open';
const STATUS_CLOSED = 'Closed';

const LOCK_LOCK = 'Locked';
const LOCK_UNLOCK = 'Unlocked';

const PIN_PIN = 'Pinned';
const PIN_UNPIN = 'Unpinned';

const issueStyles = theme => ({
  body: {
    padding: theme.spacing(1),
    margin: theme.spacing(1)
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
  },
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
  commentClose: {
    marginRight: theme.spacing(1),
    backgroundColor: '#eff3f6'
  },
  commentOpen: {
    marginRight: theme.spacing(1),
    backgroundColor: '#eff3f6'
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

  handleChange = event => {
    this.setState({
      [event.target.name]: event.target.value
    });
  };

  handleSubmit = async event => {
    const { username, repo, issueNumber } = Router.router.query;
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
        this.fetchComments();
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
    const { issue } = this.state;
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
        this.fetchIssue();
      }
    } catch (error) {
      console.log(error);
    }
  };

  toggleLockIssue = async event => {
    const { username, repo, issueNumber } = Router.router.query;
    const { issue } = this.state;
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
        this.fetchIssue();
      }
    } catch (error) {
      console.log(error);
    }
  };

  togglePinIssue = async event => {
    const { username, repo, issueNumber } = Router.router.query;
    const { issue } = this.state;
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
        this.fetchIssue();
      }
    } catch (error) {
      console.log(error);
    }
  };

  render() {
    const { issue, user, editTitle } = this.state;
    const { classes } = this.props;

    if (issue === undefined) {
      return null;
    }

    const button = (
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
        <Navbar />
        <div className={cx(classes.container, 'container')}>
          <div className="row">
            <div className="col-12">
              <div className="row">
                <div className="col-12">
                  <IssueHeader
                    issue={issue}
                    status={status}
                    updateIssueTitle={this.updateIssueTitle}
                    fetchIssue={this.fetchIssue}
                  />
                </div>
              </div>
              <div className="row">
                <div className="col-lg-10">
                  <div className={classes.issueBody}>
                    <div className={classes.username}>
                      <img
                        src={issue.user.image}
                        className={cx(classes.image)}
                      />
                      <Typography variant="body2" className={classes.name}>
                        {issue.user.username}
                      </Typography>
                      <div>
                        commented{' '}
                        {formatDistance(Date.now(), parseISO(issue.updated_at))}{' '}
                        ago
                      </div>
                    </div>
                    <Typography variant="body2" className={classes.body}>
                      {issue.body}
                    </Typography>
                  </div>
                  <Comment comments={this.state.comments} />
                  <PostComments
                    handleSubmit={this.handleSubmit}
                    body={this.state.body}
                    handleChange={this.handleChange}
                    button={button}
                  />
                </div>
                <div className="col-lg-2">
                  <Sidebar
                    issue={issue}
                    lockButton={lockButton}
                    pinButton={pinButton}
                    fetchIssue={this.fetchIssue}
                  />
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
