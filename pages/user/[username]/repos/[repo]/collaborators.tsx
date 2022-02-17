import React, { useState, useEffect } from "react";
import Router from "next/router";
import Link from "next/link";
import axios from "axios";

import { withStyles } from "@material-ui/core/styles";
import { TextField, Button, Typography, makeStyles } from "@material-ui/core";
import cx from "classnames";
import CloseIcon from "@material-ui/icons/Close";

import { TextInput, ButtonPrimary } from "@primer/react";

import { Navbar, authHeaders } from "../../../../../utils/utils";

const useCollaboratorsStyles = makeStyles((theme) => ({
	container: {
		margin: "50px auto",
	},
	mainSection: {
		border: "1px solid #ddd",
		borderRadius: "3px",
	},
	name: {
		padding: theme.spacing(1),
		fontWeight: "bold",
		borderBottom: "1px solid #ddd",
	},
	descripition: {
		padding: theme.spacing(4),
		borderBottom: "1px solid #ddd",
	},
}));

const useSidebarStyles = makeStyles((theme) => ({
	sideSection: {
		display: "flex",
		flexDirection: "column",
		border: "1px solid #ddd",
		padding: theme.spacing(1),
	},
}));

function Sidebar() {
	const classes = useSidebarStyles();
	return (
		<div className={classes.sideSection}>
			<a>Options</a>
			<a>Collaborators</a>
			<a>Options</a>
			<a>Options</a>
		</div>
	);
}

const useCollaboratorStyles = makeStyles((theme) => ({
	collaboratorDetails: {
		display: "flex",
		margin: theme.spacing(1),
		justifyContent: "space-between",
	},
	collaboratorImage: {
		height: theme.spacing(5),
		width: theme.spacing(5),
		marginRight: theme.spacing(2),
	},
	closeButton: {
		border: 0,
		backgroundColor: "#fff",
	},
}));

function Collaborator({ collaborator, removeCollaborator }) {
	const classes = useCollaboratorStyles();
	return (
		<div className={classes.collaboratorDetails}>
			<div>
				<img
					src={collaborator.user.image}
					className={classes.collaboratorImage}
				/>
				<div>{collaborator.user.username}</div>
			</div>
			<button
				className={classes.closeButton}
				onClick={() => removeCollaborator(collaborator)}
			>
				<CloseIcon />
			</button>
		</div>
	);
}

const useAddCollaboratorStyles = makeStyles((theme) => ({
	search: {
		padding: theme.spacing(1),
		fontWeight: "bold",
	},
	searchText: {
		padding: theme.spacing(1),
	},
	addCollaborator: {
		display: "flex",
	},
}));

const AddCollaborator = ({ fetchCollaborators }) => {
	const classes = useAddCollaboratorStyles();
	const [collaboratorName, setCollaboratorName] = useState("");
	// constructor(props) {
	// 	super(props);
	// 	this.state = {
	// 		collaboratorName: "",
	// 	};
	// }

	const handleSubmit = async (event) => {
		const { username, repo } = Router.router.query;
		event.preventDefault();
		try {
			const response = await axios.post(
				`/repos/${username}/${repo}/collaborators/${collaboratorName}`,
				{},
				authHeaders()
			);
			if (response.status === 201) {
				setCollaboratorName("");
				fetchCollaborators();
			}
		} catch (error) {
			console.log(error);
		}
	};

	return (
		<div>
			<Typography variant="body2" className={classes.search}>
				Search by username, full name or email address
			</Typography>
			<Typography variant="body2" className={classes.searchText}>
				You’ll only be able to find a GitHub user by their email address if
				they’ve chosen to list it publicly. Otherwise, use their username instea
			</Typography>
			<form onSubmit={handleSubmit}>
				<div className={classes.addCollaborator}>
					<TextInput
						name="collaboratorName"
						value={collaboratorName}
						onChange={(e) => setCollaboratorName(e.target.value)}
					/>
					<Button type="submit">Add collaborator</Button>
				</div>
			</form>
		</div>
	);
};

const Collaborators = ({ username, repo }) => {
	// static getInitialProps({ query }) {
	// 	return { query, username: query.username, repo: query.repo };
	// }
	const classes = useCollaboratorsStyles();
	const [collaborators, setCollaborators] = useState([]);
	// constructor(props) {
	// 	super(props);
	// 	this.state = {
	// 		collaborators: [],
	// 	};
	// }

	const fetchCollaborators = async () => {
		try {
			const response = await axios.get(
				`/repos/${username}/${repo}/collaborators`,
				authHeaders()
			);
			if (response.status === 200) {
				setCollaborators(response.data);
			}
		} catch (error) {
			console.log(error);
		}
	};

	useEffect(() => {
		fetchCollaborators();
	}, []);

	const removeCollaborator = async (collaborator) => {
		event.preventDefault();
		try {
			const response = await axios.delete(
				`/repos/${username}/${repo}/collaborators/${collaborator.user.username}`,
				authHeaders()
			);
			if (response.status === 204) {
				fetchCollaborators();
			}
		} catch (error) {
			console.log(error);
		}
	};

	const collaboratorDetails =
		collaborators === [] ? (
			<Typography variant="body2" className={classes.descripition}>
				This repository doesn’t have any collaborators yet. Use the form below
				to add a collaborator.
			</Typography>
		) : (
			collaborators.map((collaborator) => (
				<Collaborator
					key={collaborator.user.id}
					collaborator={collaborator}
					removeCollaborator={removeCollaborator}
				/>
			))
		);
	return (
		<>
			<Navbar />
			<div className={cx(classes.container, "container")}>
				<div className="row">
					<div className="col-3">
						<Sidebar />
					</div>
					<div className="col-9">
						<div className={classes.mainSection}>
							<Typography variant="body2" className={classes.name}>
								Collaborators
							</Typography>
							<div className={classes.descripition}>{collaboratorDetails}</div>
							<AddCollaborator fetchCollaborators={fetchCollaborators} />
						</div>
					</div>
				</div>
			</div>
		</>
	);
};

export default Collaborators;
