import React, { Component, Fragment } from 'react';
import Router from 'next/router';
import Link from 'next/link';
import axios from 'axios';
import { formatDistance, parseISO } from 'date-fns';

import { withStyles } from '@material-ui/core/styles';
import { Button, Typography } from '@material-ui/core';
import TextField from '@material-ui/core/TextField';
import cx from 'classnames';

import { TextInput, ButtonPrimary } from '@primer/components';

import { Navbar, authHeaders } from '../../../../../utils/utils.js';

const collaboratorsStyles = theme => ({
  container: {
    margin: '50px auto'
  },
  addCollaborator: {
    display: 'flex'
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
  },
  search: {
    padding: theme.spacing(1),
    fontWeight: 'bold'
  },
  searchText: {
    padding: theme.spacing(1)
  },
  sideSection: {
    display: 'flex',
    flexDirection: 'column',
    border: '1px solid #ddd',
    padding:theme.spacing(1)
  }
});

function _Collaborators({ classes }) {
  return (
    <Fragment>
      <Navbar />
      <div className={cx(classes.container, 'container')}>
        <div className="row">
          <div className="col-3">
            <div className={classes.sideSection}>
              <a>Options</a>
              <a>Collaborators</a>
              <a>Options</a>
              <a>Options</a>
            </div>
          </div>
          <div className="col-9">
            <div className={classes.mainSection}>
              <Typography variant="body2" className={classes.name}>
                Collaborators
              </Typography>
              <Typography variant="body2" className={classes.descripition}>
                This repository doesn’t have any collaborators yet. Use the form
                below to add a collaborator.
              </Typography>
              <div>
                <Typography variant="body2" className={classes.search}>
                  Search by username, full name or email address
                </Typography>
                <Typography variant="body2" className={classes.searchText}>
                  You’ll only be able to find a GitHub user by their email
                  address if they’ve chosen to list it publicly. Otherwise, use
                  their username instea
                </Typography>
                <div className={classes.addCollaborator}>
                  <TextInput />
                  <Button>Add collaborator</Button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Fragment>
  );
}

const Collaborators = withStyles(collaboratorsStyles)(_Collaborators);

export default Collaborators;
