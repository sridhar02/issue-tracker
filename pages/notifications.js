import React, { Component, Fragment } from 'react';
import Router from 'next/router';
import Link from 'next/link';
import axios from 'axios';
import cx from 'classnames';

import { withStyles } from '@material-ui/core/styles';
import { Button, Typography } from '@material-ui/core';
import TextField from '@material-ui/core/TextField';
import GitHubIcon from '@material-ui/icons/GitHub';

import { TabNav, UnderlineNav } from '@primer/components';

import { Navbar, authHeaders } from '../utils/utils.js';

const sidebarStyles = theme => ({
  space: {
    padding: theme.spacing(1)
  }
});

function _Sidebar({ classes }) {
  return (
    <div>
      <div className={classes.space}>Unread</div>
      <div className={classes.space}>Read</div>
      <div className={classes.space}>Participating</div>
      <div className={classes.space}>Save for later</div>
      <div className={classes.space}>All notifications</div>
    </div>
  );
}

const Sidebar = withStyles(sidebarStyles)(_Sidebar);

const notificationStyles = theme => ({
  container: {
    margin: '50px auto'
  },
  tabs: {
    margin: ' 20px 400px '
  }
});

class _Notifications extends Component {
  render() {
    const { classes } = this.props;
    return (
      <Fragment>
        <Navbar />
        <UnderlineNav className={classes.tabs}>
          <UnderlineNav.Link href="/notifications">
            Notifications
          </UnderlineNav.Link>
          <UnderlineNav.Link href="/watching">Watching</UnderlineNav.Link>
          <UnderlineNav.Link href="/subscriptions">
            Subscripitions
          </UnderlineNav.Link>
        </UnderlineNav>

        <div className={cx(classes.container, 'container')}>
          <div className="row">
            <div className="col-12"></div>
          </div>
          <div className="row">
            <div className="col-3">
              <Sidebar />
            </div>
            <div className="col-9"></div>
          </div>
        </div>
      </Fragment>
    );
  }
}

const Notifications = withStyles(notificationStyles)(_Notifications);

export default Notifications;
