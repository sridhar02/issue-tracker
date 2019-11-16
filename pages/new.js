import React, { Component, Fragment } from 'react';
import Router from 'next/router';
import Link from 'next/link';
import axios from 'axios';
import cx from 'classnames';

import { withStyles } from '@material-ui/core/styles';
import { Button, Typography } from '@material-ui/core';
import TextField from '@material-ui/core/TextField';
import Input from '@material-ui/core/Input';
import Radio from '@material-ui/core/Radio';
import RadioGroup from '@material-ui/core/RadioGroup';
import FormHelperText from '@material-ui/core/FormHelperText';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import FormControl from '@material-ui/core/FormControl';
import FormLabel from '@material-ui/core/FormLabel';

import { TextInput } from '@primer/components';

import { Navbar, authHeaders } from '../utils/utils.js';

const repoStyles = theme => ({
  container: {
    [theme.breakpoints.up('md')]: {
      margin: '50px auto',
      display: 'flex',
      justifyContent: 'center'
    }
  },
  username: {
    border: '1px solid #ddd',
    padding: theme.spacing(1)
    // margin: theme.spacing(2)
  },
  formControl: {
    margin: theme.spacing(1)
  },
  formGroup: {
    borderBottom: '1px solid #ddd',
    padding: theme.spacing(1),
    marginBottom: theme.spacing(2)
  },
  image: {
    heigth: theme.spacing(3),
    width: theme.spacing(3),
    marginRight: theme.spacing(1)
  },
  generalDescription: {
    paddingBottom: theme.spacing(3),
    paddingTop: theme.spacing(2),
    borderBottom: '1px solid #ddd'
  },
  user: {
    display: 'flex',
    marginTop: theme.spacing(2)
  },
  name: {
    heigth: '50px',
    width: '100%',
    padding: 0,
    border: '1px solid #ddd',
    borderRadius: '4px'
  },
  repositoryName: {
    fontWeight: 'bold',
    padding: '10px'
  },
  owner: {
    fontWeight: 'bold',
    padding: '10px'
  },
  createButton: {
    backgroundColor: '#28A745',
    color: '#fff'
  },
  description: {
    fontWeight: 'bold',
    paddingBottom: theme.spacing(2)
  },
  nameProperty: {
    paddingTop: theme.spacing(2),
    paddingBottom: theme.spacing(2)
  },
  descriptionTextField: {
    borderBottom: '1px solid #ddd',
    marginBottom: theme.spacing(2),
    paddingBottom: theme.spacing(2),
    width: '100%'
  }
});

class _NewRepo extends Component {
  constructor(props) {
    super(props);
    this.state = {
      user: undefined,
      name: '',
      type: '',
      description: ''
    };
  }

  fetchUser = async () => {
    try {
      const response = await axios.get('/user', authHeaders());
      if (response.status === 200) {
        this.setState({
          user: response.data
        });
      }
    } catch (error) {
      console.log(error);
    }
  };
  componentDidMount() {
    this.fetchUser();
  }
  handleChange = event => {
    this.setState({
      [event.target.name]: event.target.value
    });
  };
  handleSubmit = async event => {
    event.preventDefault();
    const { user } = this.state;
    try {
      const response = await axios.post(
        '/user/repos',
        {
          name: this.state.name,
          type: this.state.type,
          description: this.state.description
        },
        authHeaders()
      );
      if (response.status === 201) {
        Router.push(`/user/${user.username}`);
      }
    } catch (error) {
      console.log(error);
    }
  };
  render() {
    const { classes } = this.props;
    const { user } = this.state;

    if (user === undefined) {
      return null;
    }

    return (
      <Fragment>
        <Navbar />
        <form onSubmit={this.handleSubmit}>
          <div className={cx(classes.container, 'container')}>
            <div className="row">
              <div className="col-12">
                <Typography variant="h6">Create a new repository</Typography>
                <Typography
                  variant="body2"
                  className={classes.generalDescription}
                >
                  A repository contains all project files, including the
                  revision history. Already have a project repository elsewhere?
                </Typography>
                <div className={cx(classes.user, 'row')}>
                  <div className="col-lg-3">
                    <Typography variant="body2" className={classes.owner}>
                      Owner
                    </Typography>
                    <Button className={classes.username}>
                      <img src={user.image} className={classes.image} />
                      {user.username}
                    </Button>
                  </div>
                  <div className="col-lg-3">
                    <Typography
                      variant="body2"
                      className={classes.repositoryName}
                    >
                      Repository name
                    </Typography>
                    <TextInput
                      label="name"
                      name="name"
                      value={this.state.name}
                      onChange={this.handleChange}
                      placeholder=" Enter repo name"
                      className={classes.name}
                      InputProps={{ classes: { underline: classes.underline } }}
                    />
                  </div>
                </div>

                <Typography variant="body2" className={classes.nameProperty}>
                  Great repository names are short and memorable. Need
                  inspiration? How about laughing-dollop?
                </Typography>
                <Typography variant="body2" className={classes.description}>
                  Description
                </Typography>
                <TextField
                  variant="outlined"
                  name="description"
                  value={this.state.description}
                  onChange={this.handleChange}
                  className={classes.descriptionTextField}
                />
                <div className={classes.formGroup}>
                  <FormControl
                    component="fieldset"
                    className={classes.formControl}
                  >
                    <FormLabel component="legend">Repo Type</FormLabel>
                    <RadioGroup
                      aria-label="repo type"
                      name="type"
                      value={this.state.type}
                      onChange={this.handleChange}
                    >
                      <FormControlLabel
                        value="public"
                        control={<Radio />}
                        label="Public"
                      />
                      <FormControlLabel
                        value="private"
                        control={<Radio />}
                        label="private"
                      />
                    </RadioGroup>
                  </FormControl>
                </div>
                <div>
                  <Button
                    variant="contained"
                    className={classes.createButton}
                    type="submit"
                  >
                    Create Repository
                  </Button>
                </div>
              </div>
            </div>
          </div>
        </form>
      </Fragment>
    );
  }
}

const NewRepo = withStyles(repoStyles)(_NewRepo);

export default NewRepo;
