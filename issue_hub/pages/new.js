import React, { Component, Fragment } from "react";

import Router from "next/router";

import Link from "next/link";

import axios from "axios";

import { withStyles } from "@material-ui/core/styles";

import { Button, Typography } from "@material-ui/core";

import TextField from "@material-ui/core/TextField";

import cx from "classnames";

import { Navbar } from "../utils/utils.js";

import Radio from "@material-ui/core/Radio";
import RadioGroup from "@material-ui/core/RadioGroup";
import FormHelperText from "@material-ui/core/FormHelperText";
import FormControlLabel from "@material-ui/core/FormControlLabel";
import FormControl from "@material-ui/core/FormControl";
import FormLabel from "@material-ui/core/FormLabel";

const repoStyles = theme => ({
  username: {
    border: "1px solid #ddd",
    padding: theme.spacing(1),
    margin: theme.spacing(2)
  },
  formControl: {
    margin: theme.spacing(1)
  },
  image: {
    heigth: theme.spacing(4),
    width: theme.spacing(4),
    marginRight: theme.spacing(1)
  }
});

class _NewRepo extends Component {
  constructor(props) {
    super(props);
    this.state = {
      user: undefined,
      name: "",
      type: "",
      description: ""
    };
  }
  componentDidMount() {
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
  }
  handleChange = event => {
    this.setState({
      [event.target.name]: event.target.value
    });
  };
  handleSubmit = event => {
    event.preventDefault();
    const { user } = this.state;
    axios
      .post(
        "/user/repos",
        {
          name: this.state.name,
          type: this.state.type,
          description: this.state.description
        },
        {
          headers: {
            Authorization: `Bearer ${localStorage.getItem("secret")}`
          }
        }
      )
      .then(response => Router.push(`/user/${user.username}`))
      .catch(error => {
        console.log(error);
      });
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
          <div className="container">
            <div className="row">
              <div className="col-12">
                <Typography variant="h6">Create a new repository</Typography>
                <Typography variant="body2">
                  A repository contains all project files, including the
                  revision history. Already have a project repository elsewhere?
                </Typography>
                <Typography variant="h6">Owner</Typography>
                <span className={classes.username}>
                  <img src={user.image} className={classes.image} />
                  {user.username}
                </span>
                <Typography variant="h6">Repository name</Typography>
                <TextField
                  variant="outlined"
                  label="name"
                  name="name"
                  value={this.state.name}
                  onChange={this.handleChange}
                  placeholder=" Enter repo name"
                />
                <Typography variant="body2">
                  Great repository names are short and memorable. Need
                  inspiration? How about laughing-dollop?
                </Typography>
                <Typography variant="h6">Description</Typography>
                <TextField
                  variant="outlined"
                  label="description"
                  name="description"
                  value={this.state.description}
                  onChange={this.handleChange}
                />
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
                <div>
                  <Button variant="contained" color="primary" type="submit">
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
