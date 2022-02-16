import React, { useState, useEffect } from "react";
// import { useRouter } from "next/router";
import Link from "next/link";
import axios from "axios";
import cx from "classnames";

import { makeStyles } from "@material-ui/core";
import { Button, Typography, TextField } from "@material-ui/core";

import { TabNav, UnderlineNav } from "@primer/react";

import { Navbar } from "../../utils/utils";

import { IUser } from "../types";

const useUserStyles = makeStyles((theme) => ({
	image: {
		heigth: theme.spacing(13.75),
		width: theme.spacing(13.75),
		margin: theme.spacing(0, 1.875, 0, 0),
		borderRadius: "4px",
		[theme.breakpoints.up("md")]: {
			heigth: theme.spacing(30.9975),
			width: theme.spacing(30.8725),
		},
	},
	container: {
		[theme.breakpoints.up("md")]: {
			margin: "50px auto",
		},
		marginBottom: 50,
	},
	repoWrapper: {
		// padding: theme.spacing(2),
		fontWeight: "bold",
		width: "100%",
		"&:last-child $repo": {
			borderBottom: "1px solid #ddd",
		},
	},
	repo: {
		borderTop: "1px solid #ddd",
		borderLeft: "1px solid #ddd",
		borderRight: "1px solid #ddd",

		padding: theme.spacing(2),
		[theme.breakpoints.up("md")]: {
			border: "1px solid #ddd",
			marginBottom: theme.spacing(2),
		},
	},
	navbar: {
		marginBottom: theme.spacing(1),
	},
	userDetails: {
		display: "flex",
		padding: theme.spacing(1.875, 1.875, 1.25),
		[theme.breakpoints.up("md")]: {
			display: "block",
		},
	},
	userName: {
		[theme.breakpoints.up("md")]: {
			padding: theme.spacing(2),
			fontWeight: "bold",
			fontSize: theme.spacing(2.5),
			marginTop: theme.spacing(1),
		},
	},
	userProfile: {
		padding: "0 15px 15px",
		margin: theme.spacing(1, 0, 0),
		borderBottom: "1px solid #ddd",
		whiteSpace: "pre-wrap",
		[theme.breakpoints.up("md")]: {
			borderBottom: "0px",
		},
	},
	popular: {
		padding: theme.spacing(2, 0),
		[theme.breakpoints.up("md")]: {
			marginTop: theme.spacing(2),
			display: "flex",
			justifyContent: "space-between",
		},
	},
	newRepo: {
		backgroundColor: "#2cbe4e",
		color: "#fff",
		fontWeight: "bold",
		"&:hover": {
			backgroundColor: "green",
		},
	},
}));

const User = ({ username }) => {
	const [user, setUser] = useState<IUser | undefined>(undefined);
	const [repos, setRepos] = useState([]);

	const classes = useUserStyles();

	const fetchUser = async () => {
		try {
			const response = await axios.get(`/users/${username}`);
			if (response.status === 200) {
				setUser(response.data);
			}
		} catch (error) {
			console.log(error);
		}
	};

	const fetchUserRepos = async () => {
		try {
			const response = await axios.get(`/users/${username}/repos`);
			if (response.status === 200) {
				setRepos(response.data);
			} else setRepos([]);
		} catch (error) {
			console.log(error);
		}
	};

	useEffect(() => {
		fetchUser();
	}, [username]);

	useEffect(() => {
		fetchUserRepos();
	}, [username]);

	if (user === undefined) {
		return null;
	}

	return (
		<>
			<Navbar />
			<div className={cx(classes.container, "container")}>
				<div className="row">
					<div className="col-lg-3">
						<div className={classes.userDetails}>
							<img src={user.image} className={classes.image} />
							<div></div>
							<Typography variant="body2" className={classes.userName}>
								{user.username}
							</Typography>
						</div>
						<div className={classes.userProfile}></div>
					</div>
					<div className="col-lg-9">
						<div className="d-none d-md-block">
							<UnderlineNav arial-label="main">
								<UnderlineNav.Link href="#home">Overview</UnderlineNav.Link>
								<UnderlineNav.Link href="/new">Respositories</UnderlineNav.Link>
								<UnderlineNav.Link href="#Projects">Projects</UnderlineNav.Link>
							</UnderlineNav>
						</div>
						<div className={cx(classes.popular)}>
							<div>Popular respositories</div>
							<div className="d-none d-md-block">
								<Link href="/new">
									<Button className={classes.newRepo} variant="contained">
										<a>New Respository</a>
									</Button>
								</Link>
							</div>
						</div>
						<div className={"row"}>
							{repos.map((repo) => (
								<div
									key={repo.id}
									className={cx(classes.repoWrapper, "col-md-6")}
								>
									<div className={classes.repo}>
										<Link
											href={`/user/${user.username}/repos/${repo.name}/issues`}
										>
											<a>
												<span className="d-md-none">{repo.user.username}/</span>
												{repo.name}
											</a>
										</Link>
									</div>
								</div>
							))}
						</div>
					</div>
				</div>
			</div>
		</>
	);
};

export async function getServerSideProps(context) {
	return {
		props: {
			username: context?.query.username,
		}, // will be passed to the page component as props
	};
}

export default User;
