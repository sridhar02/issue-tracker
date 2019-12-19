import React, { Component, Fragment } from 'react';
import Router from 'next/router';
import Link from 'next/link';
import axios from 'axios';
import cx from 'classnames';

import { withStyles } from '@material-ui/core/styles';
import { Button, Typography } from '@material-ui/core';
import TextField from '@material-ui/core/TextField';
import GitHubIcon from '@material-ui/icons/GitHub';

import { Navbar, authHeaders } from '../../../../../../utils/utils.js';

const labelStyles = theme => ({
  newLabelButton: {
    margin: theme.spacing(2),
    marginTop: theme.spacing(3),
    display: 'flex',
    backgroundColor: '#28a745',
    color: 'white'
  },
  newLabel: {
    display: 'flex',
    justifyContent: 'space-between'
  },
  createLabelLayOut: {
    border: '1px solid #ddd',
    margin: theme.spacing(1),
    padding: theme.spacing(2)
  }
});

class _Labels extends Component {
  static getInitialProps({ query }) {
    return { query };
  }
  constructor(props) {
    super(props);
    this.state = {
      name: '',
      description: '',
      color: '',
      addLabel: false,
      labels: []
    };
  }

  handleChange = event => {
    this.setState({
      [event.target.name]: event.target.value
    });
  };

  handleSubmit = async event => {
    const { username, repo } = Router.router.query;
    event.preventDefault();
    try {
      const response = await axios.post(
        `/repos/${username}/${repo}/labels`,
        {
          name: this.state.name,
          description: this.state.description,
          color: this.state.color
        },
        authHeaders()
      );
      if (response.status === 201) {
        this.setState({
          addLabel: false
        });
      }
    } catch (error) {
      console.log(error);
    }
  };

  fetchLabels = async () => {
    const { username, repo } = Router.router.query;
    try {
      const response = await axios.get(
        `/repos/${username}/${repo}/labels`,
        authHeaders()
      );
      if (response.status === 200) {
        this.setState({
          labels: response.data
        });
      }
    } catch (error) {
      console.log(error);
    }
  };
  componentDidMount() {
    this.fetchLabels();
  }

  addLabel = () => {
    this.setState({
      addLabel: true
    });
  };

  render() {
    const { classes } = this.props;
    console.log(this.state.labels);
    const { addLabel, labels } = this.state;

    const addLabelInput = addLabel ? (
      <div className={classes.createLabelLayOut}>
        <div className={classes.newLabel}>
          <label>Label Name</label>
          <label>Description</label>
          <label>Color</label>
        </div>
        <form onSubmit={this.handleSubmit} className={classes.newLabel}>
          <TextField
            variant="outlined"
            name="name"
            value={this.state.name}
            onChange={this.handleChange}
          />
          <TextField
            variant="outlined"
            name="description"
            value={this.state.description}
            onChange={this.handleChange}
          />
          <TextField
            variant="outlined"
            name="color"
            value={this.state.color}
            onChange={this.handleChange}
          />
          <Button>Cancel</Button>
          <Button type="submit">Create Label</Button>
        </form>
      </div>
    ) : (
      <div></div>
    );

    return (
      <Fragment>
        <Navbar />
        <div className="container">
          <div className="row">
            <div className="col-12">
              <Button
                className={classes.newLabelButton}
                onClick={this.addLabel}
              >
                New Label
              </Button>
              {addLabelInput}
            </div>
          </div>
        </div>
      </Fragment>
    );
  }
}

const Labels = withStyles(labelStyles)(_Labels);

export default Labels;
