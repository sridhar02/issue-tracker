import React, { useState, useEffect } from "react";
import Link from "next/link";
import Router from "next/router";
import cx from "classnames";
import axios from "axios";

import { Button, TextField, Typography, makeStyles } from "@material-ui/core";
import { withStyles } from "@material-ui/core/styles";
import MenuIcon from "@material-ui/icons/Menu";
import AppBar from "@material-ui/core/AppBar";
import Toolbar from "@material-ui/core/Toolbar";
import IconButton from "@material-ui/core/IconButton";
import InputBase from "@material-ui/core/InputBase";
import GitHubIcon from "@material-ui/icons/GitHub";
import NotificationsIcon from "@material-ui/icons/Notifications";
import ArrowDropDownIcon from "@material-ui/icons/ArrowDropDown";
import Popper from "@material-ui/core/Popper";
import Fade from "@material-ui/core/Fade";
import Paper from "@material-ui/core/Paper";
import PropTypes from "prop-types";

import { TextInput } from "@primer/react";

const useMenuStyles = makeStyles((theme) => ({
	dashboard: {
		padding: theme.spacing(1),
		margin: theme.spacing(1),
		borderBottom: "1px solid #ddd",
	},
	issues: {
		padding: theme.spacing(1),
		margin: theme.spacing(1),
		borderBottom: "1px solid #ddd",
	},
	userDetails: {
		display: "flex",
		padding: theme.spacing(1),
		margin: theme.spacing(1),
		borderBottom: "1px solid #ddd",
	},
	userImage: {
		heigth: theme.spacing(3),
		width: theme.spacing(3),
		margin: theme.spacing(0, 1, 0, 0),
	},
	signout: {
		color: "#fff",
		padding: theme.spacing(1),
		margin: theme.spacing(1),
	},
}));

function Menu({ user }) {
	const classes = useMenuStyles();
	const handleSignout = () => {
		localStorage.removeItem("secret");
		Router.push("/login");
	};
	return (
		<div>
			<div className={classes.dashboard}>Dashboard</div>
			<div className={classes.issues}>Issues</div>
			<div className={classes.userDetails}>
				<img src={user.image} className={classes.userImage} />
				<div>{user.username}</div>
			</div>
			<div>
				<Link href="/login">
					<Button className={classes.signout} onClick={handleSignout}>
						Signout
					</Button>
				</Link>
			</div>
		</div>
	);
}

const useNavbarStyles = makeStyles((theme) => ({
	container: {
		backgroundColor: "black",
		padding: theme.spacing(1),
	},
	navbar: {
		display: "flex",
		justifyContent: "space-between",
	},
	overview: {
		display: "flex",
		justifyContent: "space-between",
		marginTop: theme.spacing(3),
		[theme.breakpoints.up("md")]: {
			display: "flex",
			justifyContent: "space-between",
		},
	},
	Button: {
		color: "#fff",
		padding: theme.spacing(0),
		margin: theme.spacing(0),
	},
	textField: {
		paddingTop: 0,
		paddingBottom: 0,
		backgroundColor: "#fff",
		borderRadius: "4px",
		width: "30%",
	},
	space: {
		width: "10%",
	},
	userImage: {
		height: theme.spacing(4),
		width: theme.spacing(4),
		borderRadius: "4px",
	},
	dropDown: {
		backgroundColor: "black",
		color: "#fff",
		padding: 0,
		margin: 0,
		border: 0,
	},
	popper: {
		padding: "10px",
		display: "flex",
		flexDirection: "column",
	},
	usernamePopper: {
		padding: theme.spacing(1),
		borderBottom: "1px solid #ddd",
	},
	typography: {
		padding: theme.spacing(1),
	},
	yourGists: {
		padding: theme.spacing(1),
		borderBottom: "1px solid #ddd",
	},
}));

export const Navbar = () => {
	const classes = useNavbarStyles();

	const [navbar, setNavbar] = useState(false);
	const [user, setUser] = useState(undefined);
	const [anchorEl, setAnchorEl] = useState(null);
	const [open, setOpen] = useState(false);

	const fetchUser = async () => {
		try {
			const response = await axios.get("/user", authHeaders());
			if (response.status === 200) {
				setUser(response.data);
			}
		} catch (error) {
			console.log(error);
		}
	};

	const handleClick = (event) => {
		setNavbar(!navbar);
	};

	const handleClickPoper = (event) => {
		const { currentTarget } = event;
		setAnchorEl(currentTarget);
		setOpen(!open);
	};

	const handleSignout = () => {
		localStorage.removeItem("secret");
		Router.push("/login");
	};

	useEffect(() => {
		fetchUser();
	}, []);

	const menu = navbar ? <Menu user={user} /> : <div></div>;
	const id = open ? "simple-popper" : null;

	return (
		<AppBar position="static" className={classes.container}>
			<div className={cx(classes.container, "container")}>
				<div className="row">
					<div className="col-12">
						<div className={classes.navbar}>
							<Button
								onClick={handleClick}
								className={cx(classes.Button, "d-md-none")}
								type="button"
							>
								<MenuIcon />
							</Button>
							{user && (
								<Link href={`/user/${user.username}`}>
									<GitHubIcon />
								</Link>
							)}
							<TextInput
								className={cx(classes.textField, "d-none d-md-block")}
							/>
							<div className="d-none d-md-block">Pull requests</div>
							<div className="d-none d-md-block">Issues</div>
							<div className="d-none d-md-block">Marketplace</div>
							<div className="d-none d-md-block">Explore</div>
							<div className={cx(classes.space, "d-none d-md-block")}></div>
							<Link href="/notifications">
								<NotificationsIcon />
							</Link>
							{user && (
								<div className="d-none d-md-block">
									<button
										type="button"
										aria-describedby={id}
										onClick={handleClickPoper}
										className={classes.dropDown}
									>
										<img src={user.image} className={classes.userImage} />
										<ArrowDropDownIcon />
									</button>
									<Popper id={id} open={open} anchorEl={anchorEl}>
										<Paper className={classes.popper}>
											<Typography className={classes.usernamePopper}>
												signed as {user.username}
											</Typography>
											<Button className={classes.typography}>
												Your Profile
											</Button>
											<Button className={classes.typography}>
												Your Repositories
											</Button>
											<Button className={classes.typography}>
												Your Projects
											</Button>
											<Button className={classes.typography}>Your stars</Button>
											<Button className={classes.yourGists}>Your gists</Button>
											<Button className={classes.typography}>
												Feature Preview
											</Button>
											<Button className={classes.typography}>Help</Button>
											<Button className={classes.typography}>Settings</Button>
											<Link href="/login">
												<Button
													className={classes.typography}
													onClick={handleSignout}
												>
													sign out
												</Button>
											</Link>
										</Paper>
									</Popper>
								</div>
							)}
						</div>
						{menu}
						<div className={cx(classes.overview, "d-md-none")}>
							<div>Overview</div>
							<div>Repositories</div>
							<div>Projects</div>
						</div>
					</div>
				</div>
			</div>
		</AppBar>
	);
};

export function authHeaders() {
	return {
		headers: {
			Authorization: `Bearer ${localStorage.getItem("secret")}`,
		},
	};
}
