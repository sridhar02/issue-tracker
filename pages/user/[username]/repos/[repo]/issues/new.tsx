import React, { Fragment, useState, useEffect } from "react";
import Router from "next/router";
import Link from "next/link";
import cx from "classnames";
import axios from "axios";

import { withStyles } from "@material-ui/core/styles";
import { Button, Typography, makeStyles } from "@material-ui/core";
import TextField from "@material-ui/core/TextField";

import { Navbar, authHeaders } from "../../../../../../utils/utils";

const useIssueStyles = makeStyles((theme) => ({
	container: {
		margin: theme.spacing(1),
		[theme.breakpoints.up("md")]: {
			margin: "25px auto",
			padding: "100px auto",
		},
	},
	image: {
		heigth: theme.spacing(4),
		width: theme.spacing(4),
	},
	issue: {
		border: "1px solid #ddd",
		heigth: theme.spacing(20),
		margin: theme.spacing(1),
		padding: theme.spacing(1),
	},
	write: { paddingRight: theme.spacing(1) },
	options: {
		display: "flex",
		border: "1px solid #ddd",
		borderTop: 0,
		marginBottom: theme.spacing(1),
		padding: theme.spacing(1),
		borderRadius: "3px",
		heigth: theme.spacing(4),
		width: "100%",
	},
	body: {
		marginTop: theme.spacing(1),
		padding: theme.spacing(1),
		width: "100%",
		borderRadius: "3px",
		height: "200px",
		minHeight: "200px",
	},
	issueOptions: {
		padding: theme.spacing(2),
		borderBottom: "1px solid #ddd",
	},
	button: {
		marginTop: theme.spacing(1),
		marginLeft: theme.spacing(2),
		backgroundColor: "#2cbe4e",
	},
	title: {
		width: "100%",
	},
	Headers: {
		borderBottom: "1px solid #ddd",
		margin: theme.spacing(1),
		padding: theme.spacing(1),
	},
	buttonPosition: {
		// [themetheme => ({.breakpoints.up('md')]: {
		//   display: 'flex',
		//   justifyContent: 'flex-end'
		// }
	},
}));

const Issue = ({ username, repo }) => {
	const classes = useIssueStyles();

	const [user, setUser] = useState(undefined);
	const [title, setTitle] = useState("");
	const [body, setBody] = useState("");

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
	//   componentDidMount() {
	//     this.fetchUser();
	//   }

	useEffect(() => {
		fetchUser();
	}, []);

	const handleSubmit = async (event) => {
		const { username, repo } = Router.router.query;
		event.preventDefault();
		try {
			const response = await axios.post(
				`/repos/${username}/${repo}/issues`,
				{
					title: title,
					body: body,
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

	// const { classes, username, repo } = this.props;
	// const { user } = this.state;

	return (
		<Fragment>
			<Navbar />
			<div className={cx(classes.container, "container")}>
				<div className="row">
					<div className="col-12">
						<div className={cx(classes.Headers, "d-none d-md-block")}>
							<Typography variant="h6">
								<Link href={`/user/${username}`}>
									<a>
										{username}/{repo}
									</a>
								</Link>
							</Typography>
						</div>
					</div>
				</div>
				<div className="row">
					<div className="col-lg-9">
						<div>
							<div className={classes.issue}>
								<form onSubmit={handleSubmit}>
									<div className={classes.title}>
										<TextField
											variant="outlined"
											placeholder="Title"
											name="title"
											value={title}
											onChange={(e) => setTitle(e.target.value)}
											className={classes.title}
										/>
										<div className={classes.options}>
											<div className={classes.write}>
												<Typography variant="body2">write</Typography>
											</div>
											<Typography variant="body2">preview</Typography>
										</div>
										<textarea
											className={classes.body}
											// variant="outlined"/
											placeholder="body"
											name="body"
											value={body}
											onChange={(e) => setBody(e.target.value)}
										/>
										<div className={classes.buttonPosition}>
											<Button
												className={classes.button}
												color="primary"
												variant="contained"
												type="submit"
											>
												submit new issue
											</Button>
										</div>
									</div>
								</form>
							</div>
						</div>
					</div>
					<div className="col-lg-3">
						<div className={classes.issueOptions}>
							<Typography variant="body2">Assignees</Typography>
						</div>
						<div className={classes.issueOptions}>
							<Typography variant="body2">Lables</Typography>
						</div>
						<div className={classes.issueOptions}>
							<Typography variant="body2">Projects</Typography>
						</div>
						<div className={classes.issueOptions}>
							<Typography variant="body2">MileStone</Typography>
						</div>
					</div>
				</div>
			</div>
		</Fragment>
	);
};

export async function getServerSideProps(context) {
	return {
		props: {
			username: context?.query?.username,
			repo: context?.query?.repo,
		}, // will be passed to the page component as props
	};
}

export default Issue;
