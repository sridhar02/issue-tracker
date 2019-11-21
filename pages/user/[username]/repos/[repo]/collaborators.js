import React, { Component, Fragment } from 'react';
import Router from 'next/router';
import Link from 'next/link';
import axios from 'axios';

import { withStyles } from '@material-ui/core/styles';
import { TextField, Button, Typography } from '@material-ui/core';
import cx from 'classnames';
import CloseIcon from '@material-ui/icons/Close';

import { TextInput, ButtonPrimary } from '@primer/components';

import { Navbar, authHeaders } from '../../../../../utils/utils.js';

const collaboratorsStyles = theme => ({
  container: {
    margin: '50px auto'
  },
  mainSection: {
    border: '1px solid #ddd',
    borderRadius: '3px'
  },
  name: {
    padding: theme.spacing(1),
    fontWeight: 'bold',
    borderBottom: '1px solid #ddd'
  },
  descripition: {
    padding: theme.spacing(4),
    borderBottom: '1px solid #ddd'
  }
});

const sidebarStyles = theme => ({
  sideSection: {
    display: 'flex',
    flexDirection: 'column',
    border: '1px solid #ddd',
    padding: theme.spacing(1)
  }
});

function _Sidebar({ classes }) {
  return (
    <div className={classes.sideSection}>
      <a>Options</a>
      <a>Collaborators</a>
      <a>Options</a>
      <a>Options</a>
    </div>
  );
}
const Sidebar = withStyles(sidebarStyles)(_Sidebar);

const listCollaboratorStyles = theme => ({
  collaboratorDetails: {
    display: 'flex',
    margin: theme.spacing(1),
    justifyContent: 'space-between'
  },
  collaboratorImage: {
    height: theme.spacing(5),
    width: theme.spacing(5),
    marginRight: theme.spacing(2)
  },
  closeButton: {
    border: 0,
    backgroundColor: '#fff'
  }
});

function _ListOfCollaborator({ collaborator, removeCollaborator, classes }) {
  return (
    <div key={collaborator.name} className={classes.collaboratorDetails}>
      <div key={collaborator.userImage}>
        <img
          src={collaborator.userImage}
          className={classes.collaboratorImage}
        />
        <div>{collaborator.username}</div>
      </div>
      <button
        key={collaborator.username}
        className={classes.closeButton}
        onClick={() => removeCollaborator(collaborator)}
      >
        <CloseIcon />
      </button>
    </div>
  );
}

const ListOfCollaborator = withStyles(listCollaboratorStyles)(
  _ListOfCollaborator
);

const addCollaboratorStyles = theme => ({
  search: {
    padding: theme.spacing(1),
    fontWeight: 'bold'
  },
  searchText: {
    padding: theme.spacing(1)
  },
  addCollaborator: {
    display: 'flex'
  }
});

class _AddCollaborator extends Component {
  constructor(props) {
    super(props);
    this.state = {
      collaboratorName: ''
    };
  }
  handleChange = event => {
    this.setState({
      [event.target.name]: event.target.value
    });
  };

  handleSubmit = async event => {
    const { username, repo } = Router.router.query;
    const { collaboratorName } = this.state;
    const { fetchCollaborators } = this.props;
    event.preventDefault();
    try {
      const response = await axios.post(
        `/repos/${username}/${repo}/collaborators/${collaboratorName}`,
        {},
        authHeaders()
      );
      if (response.status === 201) {
        this.setState({ collaboratorName: '' });
        fetchCollaborators();
      }
    } catch (error) {
      console.log(error);
    }
  };

  render() {
    const { classes } = this.props;
    const { collaboratorName } = this.props;
    return (
      <div>
        <Typography variant="body2" className={classes.search}>
          Search by username, full name or email address
        </Typography>
        <Typography variant="body2" className={classes.searchText}>
          You’ll only be able to find a GitHub user by their email address if
          they’ve chosen to list it publicly. Otherwise, use their username
          instea
        </Typography>
        <form onSubmit={this.handleSubmit}>
          <div className={classes.addCollaborator}>
            <TextInput
              name="collaboratorName"
              value={this.state.collaboratorName}
              onChange={this.handleChange}
              className={classes.textInput}
            />
            <Button type="submit">Add collaborator</Button>
          </div>
        </form>
      </div>
    );
  }
}
const AddCollaborator = withStyles(addCollaboratorStyles)(_AddCollaborator);

class _Collaborators extends Component {
  static getInitialProps({ query }) {
    return { query, username: query.username, repo: query.repo };
  }
  constructor(props) {
    super(props);
    this.state = {
      collaborators: []
    };
  }

  fetchCollaborators = async () => {
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
    this.fetchCollaborators();
  }

  removeCollaborator = async collaborator => {
    const { username, repo } = Router.router.query;
    event.preventDefault();
    try {
      const response = await axios.delete(
        `/repos/${username}/${repo}/collaborators/${collaborator.username}`,
        authHeaders()
      );
      if (response.status === 204) {
        this.fetchCollaborators();
      }
    } catch (error) {
      console.log(error);
    }
  };

  render() {
    const { classes } = this.props;
    const { collaborators } = this.state;
    const collaboratorDetails =
      collaborators === [] ? (
        <Typography variant="body2" className={classes.descripition}>
          This repository doesn’t have any collaborators yet. Use the form below
          to add a collaborator.
        </Typography>
      ) : (
        collaborators.map(collaborator => (
          <ListOfCollaborator
            collaborator={collaborator}
            removeCollaborator={this.removeCollaborator}
          />
        ))
      );
    return (
      <Fragment>
        <Navbar />
        <div className={cx(classes.container, 'container')}>
          <div className="row">
            <div className="col-3">
              <Sidebar />
            </div>
            <div className="col-9">
              <div className={classes.mainSection}>
                <Typography variant="body2" className={classes.name}>
                  Collaborators
                </Typography>
                <div className={classes.descripition}>
                  {collaboratorDetails}
                </div>
                <AddCollaborator fetchCollaborators={this.fetchCollaborators} />
              </div>
            </div>
          </div>
        </div>
      </Fragment>
    );
  }
}

const Collaborators = withStyles(collaboratorsStyles)(_Collaborators);

export default Collaborators;
