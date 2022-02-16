import React, { Component, Fragment, useState, useEffect } from "react";
import Router, { useRouter } from "next/router";
import Link from "next/link";
import axios from "axios";
import cx from "classnames";

import { Button, Typography, makeStyles } from "@material-ui/core";
import TextField from "@material-ui/core/TextField";
import GitHubIcon from "@material-ui/icons/GitHub";

import { Navbar, authHeaders } from "../../../../../../utils/utils";

const useLabelStyles = makeStyles((theme) => ({
	newLabelButton: {
		margin: theme.spacing(2),
		marginTop: theme.spacing(3),
		display: "flex",
		backgroundColor: "#28a745",
		color: "white",
	},
	newLabel: {
		display: "flex",
		justifyContent: "space-between",
	},
	createLabelLayOut: {
		border: "1px solid #ddd",
		margin: theme.spacing(1),
		padding: theme.spacing(2),
	},
}));

const Labels = () => {
	const classes = useLabelStyles();

	const [name, setName] = useState("");
	const [description, setDescription] = useState("");
	const [color, setColor] = useState("");
	const [addLabel, setAddLabel] = useState(false);
	const [labels, setLabels] = useState([]);

	const handleSubmit = async (event) => {
		const { username, repo } = Router.router.query;
		event.preventDefault();
		try {
			const response = await axios.post(
				`/repos/${username}/${repo}/labels`,
				{
					name: name,
					description: description,
					color: color,
				},
				authHeaders()
			);
			if (response.status === 201) {
				setAddLabel(false);
			}
		} catch (error) {
			console.log(error);
		}
	};

	const fetchLabels = async () => {
		const { username, repo } = Router.router.query;
		try {
			const response = await axios.get(
				`/repos/${username}/${repo}/labels`,
				authHeaders()
			);
			if (response.status === 200) {
				setLabels(response.data);
			}
		} catch (error) {
			console.log(error);
		}
	};

	useEffect(() => {
		fetchLabels();
	}, []);

	const addLabelInput = addLabel ? (
		<div className={classes.createLabelLayOut}>
			<div className={classes.newLabel}>
				<label>Label Name</label>
				<label>Description</label>
				<label>Color</label>
			</div>
			<form onSubmit={handleSubmit} className={classes.newLabel}>
				<TextField
					variant="outlined"
					name="name"
					value={name}
					onChange={(e) => setName(e.target.value)}
				/>
				<TextField
					variant="outlined"
					name="description"
					value={description}
					onChange={(e) => setDescription(e.target.value)}
				/>
				<TextField
					variant="outlined"
					name="color"
					value={color}
					onChange={(e) => setColor(e.target.value)}
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
							onClick={() => setAddLabel(true)}
						>
							New Label
						</Button>
						{addLabelInput}
					</div>
				</div>
			</div>
		</Fragment>
	);
};

export default Labels;
