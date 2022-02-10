import React, { Component, Fragment, useState } from "react";
import Router from "next/router";
import Link from "next/link";
import cx from "classnames";
import axios from "axios";

import { Button, Typography } from "@material-ui/core";
import GitHubIcon from "@material-ui/icons/GitHub";

import { TextInput, ButtonPrimary } from "@primer/react";
import { makeStyles } from "@material-ui/core";

const useNavbarStyles = makeStyles((theme) => ({
	navbar: {
		display: "flex",
		backgroundColor: "black",
		justifyContent: "space-between",
		alignItems: "center",
		color: "#fff",
		padding: theme.spacing(2),
	},
	signup: {
		color: "#fff",
		border: "1px solid #fff",
	},
	signin: {
		color: "#fff",
	},
	searchbar: {
		backgroundColor: "#fff",
		borderRadius: "4px",
		margin: theme.spacing(1),
		height: "30px",
		width: "30%",
	},
}));

export const Navbar = () => {
	const classes = useNavbarStyles();
	return (
		<div className={cx(classes.navbar, "d-none d-md-flex")}>
			<GitHubIcon />
			<div>Why Github</div>
			<div>Enterprise</div>
			<div>Explore</div>
			<div>Marketplace</div>
			<div>Pricing</div>
			<div></div>
			<TextInput className={classes.searchbar} />
			<Link href="/login">
				<Button className={classes.signin}>Sign in</Button>
			</Link>
			<Button className={classes.signup}>Sign up</Button>
		</div>
	);
};

const useSignupStyles = makeStyles((theme) => ({
	container: {
		backgroundColor: "#2b3137",
		color: "#fff",
		padding: "20px",
		margin: 0,
		[theme.breakpoints.up("md")]: {
			padding: "100px",
			paddingRight: "150px",
			maxWidth: "2500px",
			height: "1080px",
			margin: 0,
			display: "flex",
		},
	},
	signup: {
		backgroundColor: "#28A745",
		margin: "20px 0px 0px",
		padding: "25px",
		color: "#fff",
		width: "100%",
	},
	sidebar: {
		borderRadius: "4px",
		padding: theme.spacing(3),
		backgroundColor: "#fff",
		color: "black",
		display: "flex",
		flexDirection: "column",
		justifyContent: "center",
	},
	Input: {
		width: "100%",
	},
}));

const Signup = () => {
	const classes = useSignupStyles();

	const [name, setName] = useState("");
	const [username, setUsername] = useState("");
	const [email, setEmail] = useState("");
	const [password, setPassword] = useState("");

	const handleSubmit = async (event) => {
		event.preventDefault();
		try {
			const response = await axios.post("/signup", {
				name: name,
				username: username,
				email: email,
				password: password,
			});
			if (response.status === 201) {
				Router.push("/login");
			}
		} catch (error) {
			console.log(error);
		}
	};

	return (
		<Fragment>
			<Navbar />
			<div className={cx(classes.container, "container")}>
				<form onSubmit={handleSubmit}>
					<div className="row">
						<div className="col-lg-8 d-none d-md-block">
							<Typography variant="h2">Built for developers</Typography>
							<Typography variant="h6">
								GitHub is a development platform inspired by the way you work.
								From open source to business, you can host and review code,
								manage projects, and build software alongside 40 million
								developers.
							</Typography>
						</div>
						<div className={cx(classes.sidebar, "col-lg-3")}>
							<div>
								<Typography variant="h5">Create your account</Typography>
							</div>
							<Typography variant="body2">Name</Typography>
							<TextInput
								type="text"
								name="name"
								margin="normal"
								variant="outlined"
								value={name}
								onChange={(event) => setName(event.target.value)}
								className={classes.Input}
							/>
							<Typography variant="body2">Username</Typography>
							<TextInput
								type="text"
								name="username"
								margin="normal"
								variant="outlined"
								className={classes.Input}
								value={username}
								onChange={(event) => setUsername(event.target.value)}
							/>
							<Typography variant="body2">Email</Typography>
							<TextInput
								type="text"
								name="email"
								margin="normal"
								variant="outlined"
								className={classes.Input}
								value={email}
								onChange={(event) => setEmail(event.target.value)}
							/>
							<Typography variant="body2">Password</Typography>
							<TextInput
								type="password"
								name="password"
								margin="normal"
								variant="outlined"
								className={classes.Input}
								value={password}
								onChange={(event) => setPassword(event.target.value)}
							/>
							<div>
								<Button
									variant="contained"
									className={classes.signup}
									type="submit"
								>
									Sign up for Github
								</Button>
							</div>
						</div>
					</div>
				</form>
			</div>
		</Fragment>
	);
};

export default Signup;
