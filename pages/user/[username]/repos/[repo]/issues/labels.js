import React, { Component, Fragment } from 'react';
import Router from 'next/router';
import Link from 'next/link';
import axios from 'axios';
import cx from 'classnames';

import { withStyles } from '@material-ui/core/styles';
import { Button, Typography } from '@material-ui/core';
import TextField from '@material-ui/core/TextField';
import GitHubIcon from '@material-ui/icons/GitHub';

const labelStyles = theme => ({});

class _Label extends Component {
  reder() {
    return (
      <div className="container">
        <div className="row">
          <div className="col-12"></div>
        </div>
      </div>
    );
  }
}

const Label = withStyles(labelStyles)(_Label);

export default Label;
