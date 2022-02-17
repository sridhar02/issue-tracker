import React, { useState, useEffect } from "react";
import Router from "next/router";
import Link from "next/link";
import axios from "axios";
import cx from "classnames";
import { formatDistance, parseISO } from "date-fns";

import SettingsIcon from "@material-ui/icons/Settings";
import { withStyles } from "@material-ui/core/styles";
import {
	Button,
	Typography,
	TextField,
	Popover,
	Paper,
	makeStyles,
} from "@material-ui/core";
import LockIcon from "@material-ui/icons/Lock";
import DeleteIcon from "@material-ui/icons/Delete";
import ErrorOutlineIcon from "@material-ui/icons/ErrorOutline";
import CheckIcon from "@material-ui/icons/Check";
import MoreHorizIcon from "@material-ui/icons/MoreHoriz";

import { TextInput } from "@primer/react";

import { Navbar, authHeaders } from "../../../../../../utils/utils";

// const commentStyles = makeStyles((theme) => ({ => ({}));

function Comment({ comments, collaborators, repoUserName }) {
	return (
		<div>
			{comments.map((comment) => (
				<div key={comment.id}>
					<Body
						collaborators={collaborators}
						user={comment.user}
						body={comment.body}
						updated_at={comment.created_at}
						repoUserName={repoUserName}
					/>
				</div>
			))}
		</div>
	);
}

const usePostCommentStyles = makeStyles((theme) => ({
	commentSection: {
		display: "block",
		width: "100%",
		minHeight: "100px",
		maxHeight: "500px",
		marginBottom: theme.spacing(1),
		padding: theme.spacing(1),
		resize: "vertical",
	},
	commentButton: {
		color: "#fff",
		backgroundColor: "#2cbe4e",
		"&:hover": {
			backgroundColor: "green",
		},
	},
	commentClose: {
		marginRight: theme.spacing(1),
		backgroundColor: "#eff3f6",
	},
	commentOpen: {
		marginRight: theme.spacing(1),
		backgroundColor: "#eff3f6",
	},
}));

const PostComments = ({ issue, fetchComments, fetchIssue }) => {
	const classes = usePostCommentStyles();
	const [body, setBody] = useState("");

	const postCommentHandler = async (event) => {
		const { username, repo, issueNumber } = Router.router.query;
		event.preventDefault();
		// const { user } = this.state;
		try {
			const response = await axios.post(
				`/repos/${username}/${repo}/issues/${issueNumber}/comments`,
				{
					body: body,
				},
				authHeaders()
			);
			if (response.status === 201) {
				fetchComments();

				setBody("");
			}
		} catch (error) {
			console.log(error);
		}
	};

	const toggleIssueStatus = async (event) => {
		const { username, repo, issueNumber } = Router.router.query;
		const newStatus =
			issue.status === STATUS_OPEN ? STATUS_CLOSED : STATUS_OPEN;
		try {
			const response = await axios.put(
				`/repos/${username}/${repo}/issues/${issueNumber}`,
				{
					status: newStatus,
				},
				authHeaders()
			);
			if (response.status === 204) {
				fetchIssue();
			}
		} catch (error) {
			console.log(error);
		}
	};

	const closeButton = (
		<Button
			variant="contained"
			onClick={toggleIssueStatus}
			className={
				issue.status === STATUS_OPEN
					? classes.commentOpen
					: classes.commentClose
			}
		>
			{issue.status === STATUS_OPEN && <ErrorOutlineIcon />}
			{issue.status === STATUS_OPEN ? "Close" : "Reopen"} issue
		</Button>
	);

	return (
		<form onSubmit={postCommentHandler}>
			<div>
				<textarea
					className={cx(classes.commentSection, "form-control")}
					name="body"
					value={body}
					onChange={(e) => setBody(e.target.value)}
				/>
				{closeButton}
				<Button
					variant="contained"
					className={classes.commentButton}
					type="submit"
				>
					Comment
				</Button>
			</div>
		</form>
	);
};

const useAssigneePopperStyles = makeStyles((theme) => ({
	popper: {
		padding: theme.spacing(2),
		width: "250px",
		backgroundColor: "#F6F8FA",
	},
	collaboratorDetails: {
		display: "flex",
		justifyContent: "flex-start",
		width: "100%",
		marginTop: theme.spacing(1),
	},
	collaboratorImage: {
		height: theme.spacing(3),
		width: theme.spacing(3),
		marginRight: theme.spacing(1),
		marginLeft: theme.spacing(1),
	},
	paperBottom: {
		backgroundColor: "#fff",
	},
	assigneeText: {
		fontWeight: "bold",
	},
	assigneeImage: {
		height: theme.spacing(3),
		width: theme.spacing(3),
		marginBottom: theme.spacing(1),
		marginRight: theme.spacing(2),
	},
	assigneeButton: {
		backgroundColor: "#fff",
		border: 0,
		marginBottom: theme.spacing(1),
		width: "100%",
		display: "flex",
		justifyContent: "space-between",
		"&:hover": {
			border: 0,
		},
	},
	assigneeClick: {
		"&:hover": {
			backgroundColor: "blue",
			border: 0,
		},
	},
	assigneeIcon: {
		width: "20px",
	},
	collaborator: {
		display: "flex",
		justifyContent: "flex-end",
	},
	removeAllAssignees: {
		width: "100%",
		"&:hover": {
			backgroundColor: "blue",
		},
		display: "flex",
		justifyContent: "flex-start",
	},
	postYourSelfButton: {
		border: 0,
		padding: 0,
		margin: 0,
		fontStyle: "normal",
		backgroundColor: "#fff",
		fontSize: "14px",
	},
}));

const Assignee = ({
	repo,
	authorizedUser,
	fetchIssue,
	issue,
	collaborators,
}) => {
	const classes = useAssigneePopperStyles();
	const [anchorEl, setAnchorEl] = useState(null);
	const [open, setOpen] = useState(false);
	const [addAssignees, setAddAssignees] = useState([]);
	const [removeAssignees, setRemoveAssignees] = useState([]);

	// constructor(props) {
	// 	super(props);
	// 	this.state = {
	// 		anchorEl: null,
	// 		open: false,
	// 		addAssignees: [],
	// 		removeAssignees: [],
	// 	};
	// }

	const togglePopper = async (event) => {
		// const { open, addAssignees, removeAssignees, assigneed } = this.state;
		if (open) {
			if (removeAssignees.length !== 0) {
				await deleteAssignee(removeAssignees);
			}
			if (addAssignees.length !== 0) {
				await postAssignee(addAssignees);
			}
			if (removeAssignees.length !== 0 || addAssignees.length !== 0) {
				fetchIssue();
			}

			setAnchorEl(null);
			setOpen(!open);
		} else {
			const { currentTarget } = event;

			setAnchorEl(currentTarget);
			setOpen(!open);
		}
	};

	const postAssignee = async (addAssignees) => {
		console.log(addAssignees);
		event.preventDefault();
		const { username, repo, issueNumber } = Router.router.query;
		try {
			const response = await axios.post(
				`/repos/${username}/${repo}/issues/${issueNumber}/assignees`,
				{
					usernames: addAssignees,
				},
				authHeaders()
			);
			if (response.status === 201) {
				setAddAssignees([]);
			}
		} catch (error) {
			console.log(error);
		}
	};

	const deleteAssignee = async (removeAssignees) => {
		const { username, repo, issueNumber } = Router.router.query;
		try {
			const response = await axios({
				method: "delete",
				url: `/repos/${username}/${repo}/issues/${issueNumber}/assignees`,
				headers: authHeaders().headers,
				data: { usernames: removeAssignees },
			});
			if (response.status === 204) {
				// this.setState({ removeAssignees: [] });
				setRemoveAssignees([]);
			}
		} catch (error) {
			console.log(error);
		}
	};

	const removeAllAssignees = () => {
		const allAssignees = issue.assignees.map(
			(assignee) => assignee.user.username
		);
		// this.setState(
		// 	{
		// 		removeAssignees: allAssignees,
		// 	},
		// 	() => this.togglePopper()
		// );
		setRemoveAssignees(allAssignees);
		// togglePopper();
	};

	const postYourSelfAsAssignee = async () => {
		let newAssignee = [authorizedUser.username];
		await postAssignee(newAssignee);
		await fetchIssue();
	};

	const checkAssignee = (collaborator) => {
		// if collaborator is in issue(assignees) && removeAssignees not includes collaborator then  add collaborator to removeAssignees list.
		// if collaborator is not in issue assignees list && removeAssignees not includes collaborator then add collaborator to addAssignee list.
		// if collaborator is not in issue Assignees list then remove collaborator from addAssignees list && collaborator is in addAssignees list.
		// if collaborator is in issue Assignees list them remove from removeAssignees list && collaborator is in removeAssignees list.

		const alreadyAssigned = issue.assignees.some(
			(assignee) => assignee.user.username === collaborator.user.username
		);

		const alreadyInAdd = addAssignees.includes(collaborator.user.username);

		const alreadyInRemove = removeAssignees.includes(
			collaborator.user.username
		);

		if (alreadyAssigned && !alreadyInRemove) {
			// this.setState({
			// 	removeAssignees: [...removeAssignees, collaborator.user.username],
			// });
			setRemoveAssignees([...removeAssignees, collaborator.user.username]);
		} else if (!alreadyAssigned && !alreadyInAdd) {
			setAddAssignees([...addAssignees, collaborator.user.username]);
		} else if (!alreadyAssigned && alreadyInAdd) {
			setAddAssignees(
				addAssignees.filter(
					(addAssignee) => addAssignee !== collaborator.user.username
				)
			);
		} else if (alreadyAssigned && alreadyInRemove) {
			setRemoveAssignees(
				removeAssignees.filter(
					(removeAssignee) => removeAssignee !== collaborator.user.username
				)
			);
		}
	};

	// const { classes, issue, collaborators } = this.props;
	// const { anchorEl, open, addAssignees, removeAssignees } = this.state;
	const assigneeList =
		issue.assignees.length !== 0 ? (
			issue.assignees.map((assignee) => (
				<div key={assignee.user.id}>
					<img src={assignee.user.image} className={classes.assigneeImage} />
					{assignee.user.username}
				</div>
			))
		) : (
			<div>
				No one —{" "}
				<button
					onClick={postYourSelfAsAssignee}
					className={classes.postYourSelfButton}
				>
					assignee yourself
				</button>
			</div>
		);
	return (
		<>
			<button
				type="button"
				onClick={togglePopper}
				className={classes.assigneeButton}
			>
				Assignee
				<div>
					<SettingsIcon />
				</div>
			</button>
			<div>{assigneeList}</div>
			<Popover
				open={open}
				anchorEl={anchorEl}
				onClose={togglePopper}
				anchorOrigin={{
					vertical: "bottom",
					horizontal: "right",
				}}
				transformOrigin={{
					vertical: "top",
					horizontal: "right",
				}}
			>
				<Paper className={classes.popper}>
					<div>
						<Typography variant="body2" className={classes.assigneeText}>
							Assign up to 10 people to this issue
						</Typography>
						<TextInput />
					</div>
					<div>
						<Button
							onClick={removeAllAssignees}
							className={classes.removeAllAssignees}
						>
							clear all assignees
						</Button>
					</div>
					<div className={classes.paperBottom}>
						{collaborators.map((collaborator) => {
							const assigned =
								(issue.assignees.some(
									(assignee) =>
										assignee.user.username === collaborator.user.username
								) ||
									addAssignees.includes(collaborator.user.username)) &&
								!removeAssignees.includes(collaborator.user.username);
							return (
								<div
									className={classes.assigneeClick}
									key={collaborator.user.username}
								>
									<Button
										className={classes.collaboratorDetails}
										onClick={() => checkAssignee(collaborator)}
									>
										<div className={classes.assigneeIcon}>
											{assigned && <CheckIcon />}
										</div>
										<div className={classes.collaborator}>
											<img
												src={collaborator.user.image}
												className={classes.collaboratorImage}
											/>
											<div>{collaborator.user.username}</div>
										</div>
									</Button>
								</div>
							);
						})}
					</div>
				</Paper>
			</Popover>
		</>
	);
};

const useLabelStyles = makeStyles((theme) => ({
	labelButton: {
		backgroundColor: "#fff",
		border: 0,
		marginBottom: theme.spacing(1),
		width: "100%",
		display: "flex",
		justifyContent: "space-between",
		"&:hover": {
			border: 0,
		},
	},
}));

const Labels = () => {
	const classes = useLabelStyles();
	const [anchorEl, setAnchorEl] = useState(null);
	const [open, setOpen] = useState(false);

	const togglePopper = async (event) => {
		if (open) {
			setAnchorEl(null);
			setOpen(!open);
		} else {
			const { currentTarget } = event;

			setAnchorEl(currentTarget);
			setOpen(!open);
		}
	};

	return (
		<>
			<button className={classes.labelButton} onClick={togglePopper}>
				labels {""}
				<SettingsIcon />
			</button>
			<Popover
				open={open}
				anchorEl={anchorEl}
				onClose={togglePopper}
				anchorOrigin={{
					vertical: "bottom",
					horizontal: "right",
				}}
				transformOrigin={{
					vertical: "top",
					horizontal: "right",
				}}
			>
				<Paper></Paper>
			</Popover>
		</>
	);
};

const useSidebarStyles = makeStyles((theme) => ({
	sidebar: {
		padding: theme.spacing(1),
		margin: theme.spacing(1),
		borderBottom: "1px solid #ddd",
	},
}));

const Sidebar = ({
	issue,
	fetchIssue,
	repo,
	collaborators,
	authorizedUser,
}) => {
	const classes = useSidebarStyles();
	const toggleLockIssue = async (event) => {
		const { username, repo, issueNumber } = Router.router.query;
		const lockStatus = issue.lock === LOCK_LOCK ? LOCK_UNLOCK : LOCK_LOCK;
		try {
			const response = await axios.put(
				`/repos/${username}/${repo}/issues/${issueNumber}/lock`,
				{
					lock: lockStatus,
				},
				authHeaders()
			);
			if (response.status === 204) {
				fetchIssue();
			}
		} catch (error) {
			console.log(error);
		}
	};

	const togglePinIssue = async (event) => {
		const { username, repo, issueNumber } = Router.router.query;
		const pinStatus = issue.pinned === PIN_UNPIN ? PIN_PIN : PIN_UNPIN;
		try {
			const response = await axios.put(
				`/repos/${username}/${repo}/issues/${issueNumber}/pin`,
				{
					pinned: pinStatus,
				},
				authHeaders()
			);
			if (response.status === 204) {
				fetchIssue();
			}
		} catch (error) {
			console.log(error);
		}
	};

	// const { classes, issue, fetchIssue, repo, collaborators, authorizedUser } =
	// 	this.props;
	const lockButton = (
		<Button onClick={toggleLockIssue}>
			<LockIcon /> {issue.lock === LOCK_UNLOCK ? "Lock" : "Unlock"}
			conversation
		</Button>
	);
	const pinButton = (
		<Button onClick={togglePinIssue}>
			{issue.pinned === PIN_PIN ? "Unpin" : "Pin"} Issue
		</Button>
	);

	return (
		<>
			<div className={classes.sidebar}>
				<Assignee
					repo={repo}
					issue={issue}
					fetchIssue={fetchIssue}
					authorizedUser={authorizedUser}
					collaborators={collaborators}
				/>
			</div>
			<div className={classes.sidebar}>
				<Labels />
			</div>
			<div className={classes.sidebar}>Projects</div>
			<div className={classes.sidebar}>Milestone</div>
			<div className={classes.sidebar}>{lockButton}</div>
			<div className={classes.sidebar}>{pinButton}</div>
			<div className={classes.sidebar}>
				<Button>
					<DeleteIcon /> Delete Issue
				</Button>
			</div>
		</>
	);
};

const useIssueHeaderStyles = makeStyles((theme) => ({
	issueHeading: {
		display: "flex",
	},
	title: {
		margin: theme.spacing(1),
	},
	issueStatus: {
		display: "flex",
		paddingBottom: theme.spacing(2.5),
		marginBottom: theme.spacing(2.5),
		borderBottom: "1px solid #ddd",
		[theme.breakpoints.up("md")]: {
			order: 1,
		},
	},
	issueTitle: {
		marginTop: theme.spacing(1),
	},
	titleEditButtons: {
		display: "flex",
	},
	save: {
		backgroundColor: "#2cbe4e",
		marginBottom: theme.spacing(1),
		marginTop: theme.spacing(1),
	},
	cancel: {
		margin: theme.spacing(1, 0, 1, 1),
	},
	editButton: {
		fontSize: theme.spacing(1.5),
		padding: theme.spacing(0.5),
		backgroundColor: "#eff3f6",
		marginRight: theme.spacing(1),
	},
	newIssue: {
		color: "#fff",
		backgroundColor: "#2cbe4e",
		fontSize: theme.spacing(1.5),
		padding: theme.spacing(0.5),
		fontWeight: "bold",
		"&:hover": {
			backgroundColor: "green",
		},
	},
	headerTop: {
		[theme.breakpoints.up("md")]: {
			display: "flex",
			justifyContent: "space-between",
		},
	},
	buttonsHeader: {
		[theme.breakpoints.up("md")]: {
			order: 2,
		},
	},
	statusOpen: {
		marginRight: theme.spacing(1),
		fontSize: theme.spacing(1.5),
		padding: theme.spacing(0.5),
		backgroundColor: "#2cbe4e",
		color: "#fff",
		fontWeight: "bold",
		"&:hover": {
			backgroundColor: "#2cbe4e",
		},
	},
	statusClose: {
		marginRight: theme.spacing(1),
		fontSize: theme.spacing(1.5),
		padding: theme.spacing(0.5),
		backgroundColor: "#CB2431",
		color: "#fff",
		fontWeight: "bold",
		"&:hover": {
			backgroundColor: "#CB2431",
		},
	},
	statusIcon: {
		marginRight: theme.spacing(0.5),
	},
}));

const IssueHeader = ({ fetchIssue, issue }) => {
	const classes = useIssueHeaderStyles();
	const [editTitle, setEditTitle] = useState(false);
	const [title, setTitle] = useState("");

	// constructor(props) {
	// 	super(props);
	// 	this.state = {
	// 		editTitle: false,
	// 		title: "",
	// 	};
	// }

	const toggleTitle = (event) => {
		event.preventDefault();

		setEditTitle(!editTitle);
	};

	const cancelIssueEdit = (event) => {
		setEditTitle(false);
	};

	// handleChange = (event) => {
	// 	this.setState({
	// 		[event.target.name]: event.target.value,
	// 	});
	// };

	const updateIssueTitle = async (event) => {
		event.preventDefault();
		const { username, repo, issueNumber } = Router.router.query;
		try {
			const response = await axios.put(
				`/repos/${username}/${repo}/issues/${issueNumber}`,
				{
					title: title,
				},
				authHeaders()
			);
			if (response.status === 204) {
				fetchIssue();

				setEditTitle(false);
			}
		} catch (error) {
			console.log(error);
		}
	};

	// const { classes, issue, fetchIssue } = this.props;
	const { username, repo, issueNumber } = Router.router.query;
	// const { editTitle } = this.state;
	const Title =
		editTitle === true ? (
			<div>
				<form onSubmit={updateIssueTitle}>
					<TextField
						name="title"
						variant="outlined"
						placeholder={issue.title}
						value={title}
						onChange={(e) => setTitle(e.target.value)}
						className={classes.issueTitle}
					/>
					<div className={classes.titleEditButtons}>
						<Button className={classes.save} variant="contained" type="submit">
							Save
						</Button>
						<Button
							className={classes.cancel}
							variant="contained"
							onClick={cancelIssueEdit}
						>
							Cancel
						</Button>
					</div>
				</form>
			</div>
		) : (
			<Typography variant="h5" className={classes.title}>
				{issue.title}
			</Typography>
		);
	const status = (
		<Button
			variant="contained"
			className={
				issue.status === STATUS_OPEN ? classes.statusOpen : classes.statusClose
			}
		>
			<ErrorOutlineIcon className={classes.statusIcon} /> {issue.status}
		</Button>
	);

	return (
		<>
			<div className={classes.headerTop}>
				<div className={classes.buttonsHeader}>
					<Button
						variant="contained"
						className={classes.editButton}
						onClick={toggleTitle}
					>
						Edit
					</Button>
					<Link href={`/user/${username}/repos/${repo}/issues/new`}>
						<Button
							variant="contained"
							color="primary"
							className={classes.newIssue}
						>
							New issue
						</Button>
					</Link>
				</div>
				<div className={classes.issueHeading}>
					{Title}
					<Typography variant="h5" className={classes.title}>
						#{issue.issue_number}
					</Typography>
				</div>
			</div>
			<div className={classes.issueStatus}>
				{status}
				<Typography variant="body2">
					opened this issue{" "}
					{formatDistance(Date.now(), parseISO(issue.created_at))} ago
				</Typography>
			</div>
		</>
	);
};

const useBodyStyles = makeStyles((theme) => ({
	image: {
		height: theme.spacing(3),
		width: theme.spacing(3),
		marginRight: theme.spacing(1),
		borderRadius: "4px",
	},
	userSection: {
		display: "flex",
	},
	content: {
		border: "1px solid #ddd",
		marginBottom: theme.spacing(1),
		marginTop: theme.spacing(1),
	},
	body: {
		padding: theme.spacing(1),
		margin: theme.spacing(1),
	},
	issueHead: {
		display: "flex",
		justifyContent: " space-between",
		borderBottom: "1px solid #ddd",
		padding: theme.spacing(1),
		backgroundColor: "#f6f8fa",
	},
	username: {
		paddingRight: theme.spacing(2),
		marginTop: theme.spacing(0),
		fontWeight: "bold",
		fontSize: "16px",
	},
	time: {
		marginTop: theme.spacing(0),
	},
	popover: {
		width: "100%",
		"&:hover": {
			backgroundColor: "blue",
		},
	},
	optionsButton: {
		border: 0,
		backgroundColor: "#f6f8fa",
		marginLeft: theme.spacing(2),
	},
	userStatus: {
		border: "1px solid #ddd",
		padding: theme.spacing(0.5),
	},
}));

const Body = ({ user, body, updated_at, repoUserName, collaborators }) => {
	const classes = useBodyStyles();

	const [anchorEl, setAnchorEl] = useState(null);
	const [open, setOpen] = useState(false);

	// constructor(props) {
	// 	super(props);
	// 	this.state = {
	// 		anchorEl: null,
	// 		open: false,
	// 	};
	// }
	const togglePopper = (event) => {
		if (open) {
			setAnchorEl(null);
			setOpen(!open);
		} else {
			const { currentTarget } = event;

			setAnchorEl(currentTarget);
			setOpen(!open);
		}
	};
	// if user.username is equal to repoUserName then we show owner .
	// if user.username is in collaborators we show him as collaborator.
	//
	// render() {
	// const { classes, user, body, updated_at, repoUserName, collaborators } =
	// 	this.props;
	const collaboratorNames = collaborators.map(
		(collaborator) => collaborator.user.username
	);
	const userStatus = collaboratorNames.includes(user.username) ? (
		<div className={classes.userStatus}>Collaborator</div>
	) : (
		<div></div>
	);

	const userInfo =
		repoUserName === user.username ? (
			<div className={classes.userStatus}> Owner</div>
		) : (
			<div></div>
		);

	const collaboratorPopoverMenu = collaboratorNames.includes(user.username) ? (
		<div>
			<Button className={classes.popover}>Report Issue Content</Button>
		</div>
	) : (
		<div></div>
	);

	const ownerPopoverMenu =
		repoUserName === user.username ? (
			<div>
				<Button className={classes.popover}>Edit</Button>
				<Button className={classes.popover}>Delete</Button>
				<Button className={classes.popover}>Report</Button>
			</div>
		) : (
			<div></div>
		);

	return (
		<div className={classes.content}>
			<div className={classes.issueHead}>
				<div className={classes.userSection}>
					<img src={user.image} className={classes.image} />
					<Typography variant="body2" className={classes.username}>
						{user.username}
					</Typography>
					<div className={classes.time}>
						commented {formatDistance(Date.now(), parseISO(updated_at))} ago
					</div>
				</div>
				<div className={classes.userSection}>
					{userStatus}
					{userInfo}
					<button onClick={togglePopper} className={classes.optionsButton}>
						<MoreHorizIcon />
					</button>
					<Popover
						open={open}
						anchorEl={anchorEl}
						onClose={togglePopper}
						anchorOrigin={{
							vertical: "bottom",
							horizontal: "right",
						}}
						transformOrigin={{
							vertical: "top",
							horizontal: "right",
						}}
					>
						<div>
							<Button className={classes.popover}>Copy Link</Button>
						</div>
						<div>
							<Button className={classes.popover}>Quote reply</Button>
						</div>
						{collaboratorPopoverMenu}
						{ownerPopoverMenu}
					</Popover>
				</div>
			</div>

			<div>
				<Typography variant="body2" className={classes.body}>
					{body}
				</Typography>
			</div>
		</div>
	);
};

const STATUS_OPEN = "Open";
const STATUS_CLOSED = "Closed";

const LOCK_LOCK = "Locked";
const LOCK_UNLOCK = "Unlocked";

const PIN_PIN = "Pinned";
const PIN_UNPIN = "Unpinned";

const useIssueStyles = makeStyles((theme) => ({
	image: {
		height: theme.spacing(3),
		width: theme.spacing(3),
		marginRight: theme.spacing(1),
	},
	issueBody: {
		border: "1px solid #ddd",
		marginBottom: theme.spacing(1),
		marginTop: theme.spacing(1),
	},
	username: {
		display: "flex",
		borderBottom: "1px solid #ddd",
		padding: theme.spacing(1),
		backgroundColor: "#f6f8fa",
	},
	container: {
		margin: theme.spacing(1),
		[theme.breakpoints.up("md")]: {
			margin: "25px auto",
		},
	},
	name: {
		fontWeight: "bold",
		marginRight: theme.spacing(1),
	},
	statusIcon: {
		marginRight: theme.spacing(0.5),
	},
	lockButton: {
		fontSize: "",
	},
	statusCloseIcon: {
		color: "#CB2431",
	},
	issueTitle: {
		marginTop: theme.spacing(1),
	},
}));

const Issue = ({ username, repoName, issueNumber }) => {
	const classes = useIssueStyles();
	// static getInitialProps({ query }) {
	// 	return { query };
	// }
	// const { username, repo } = Router.router.query;

	const [issue, setIssue] = useState(undefined);
	const [repo, setRepo] = useState(undefined);
	const [title, setTitle] = useState("");
	const [body, setBody] = useState("");
	const [editTitle, setEditTitle] = useState(false);
	const [comments, setComments] = useState([]);
	const [collaborators, setCollaborators] = useState([]);
	const [authorizedUser, setAuthorizedUser] = useState(undefined);

	// constructor(props) {
	// 	super(props);
	// 	this.state = {
	// 		issue: undefined,
	// 		repo: undefined,
	// 		title: "",
	// 		body: "",
	// 		editTitle: false,
	// 		comments: [],
	// 		collaborators: [],
	// 		authorizedUser: undefined,
	// 	};
	// }

	const fetchRepo = async () => {
		try {
			const response = await axios.get(`repos/${username}/${repoName}`);
			if (response.status === 200) {
				setRepo(response.data);
			}
		} catch (error) {
			console.log(error);
		}
	};

	const fetchIssue = async () => {
		try {
			const response = await axios.get(
				`/repos/${username}/${repoName}/issues/${issueNumber}`
			);
			if (response.status === 200) {
				setIssue(response.data);
				setTitle(response.data.title);
			}
		} catch (error) {
			console.log(error);
		}
	};

	const fetchComments = async () => {
		const { username, repo, issueNumber } = Router.router.query;
		try {
			const response = await axios.get(
				`/repos/${username}/${repo}/issues/${issueNumber}/comments`
			);
			if (response.status === 200) {
				setComments(response.data);
			}
		} catch (error) {
			console.log(error);
		}
	};

	const fetchCollaborators = async () => {
		const { username, repo } = Router.router.query;
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

	const fetchUser = async () => {
		const { username } = Router.router.query;
		try {
			const response = await axios.get(`/users/${username}`);
			if (response.status === 200) {
				setAuthorizedUser(response.data);
			}
		} catch (error) {
			console.log(error);
		}
	};

	useEffect(() => {
		fetchIssue();
		fetchComments();
		fetchRepo();
		fetchUser();
		fetchCollaborators();
	}, []);

	if (
		issue === undefined ||
		repo === undefined ||
		authorizedUser === undefined
	) {
		return null;
	}

	// console.log(collaborators);
	return (
		<>
			<Navbar />
			<div className={cx(classes.container, "container")}>
				<div className="row">
					<div className="col-12">
						<div className="row">
							<div className="col-12">
								<IssueHeader issue={issue} fetchIssue={fetchIssue} />
							</div>
						</div>
						<div className="row">
							<div className="col-lg-10">
								<Body
									collaborators={collaborators}
									user={issue.user}
									body={issue.body}
									updated_at={issue.updated_at}
									repoUserName={repo.user.username}
								/>
								<Comment
									comments={comments}
									collaborators={collaborators}
									repoUserName={repo.user.username}
								/>
								<PostComments
									issue={issue}
									fetchComments={fetchComments}
									fetchIssue={fetchIssue}
								/>
							</div>
							<div className="col-lg-2">
								<Sidebar
									authorizedUser={authorizedUser}
									collaborators={collaborators}
									repo={repo}
									issue={issue}
									fetchIssue={fetchIssue}
								/>
							</div>
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
			username: context?.query?.username,
			repoName: context?.query?.repo,
			issueNumber: context?.query?.issueNumber,
		}, // will be passed to the page component as props
	};
}

export default Issue;
