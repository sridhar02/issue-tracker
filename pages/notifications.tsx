import React, { useEffect, useState } from "react";
import Router from "next/router";
import Link from "next/link";
import axios from "axios";
import cx from "classnames";

import { withStyles } from "@material-ui/core/styles";
import { Button, Typography, TextField, makeStyles } from "@material-ui/core";
import GitHubIcon from "@material-ui/icons/GitHub";
import DoneIcon from "@material-ui/icons/Done";

import { TabNav, UnderlineNav } from "@primer/react";

import { Navbar, authHeaders } from "../utils/utils";

const useSidebarStyles = makeStyles((theme) => ({
	space: {
		padding: theme.spacing(1),
	},
}));

function Sidebar() {
	const classes = useSidebarStyles();

	return (
		<div>
			<div className={classes.space}>Unread</div>
			<div className={classes.space}>Read</div>
			<div className={classes.space}>Participating</div>
			<div className={classes.space}>Save for later</div>
			<div className={classes.space}>All notifications</div>
		</div>
	);
}

const useNotificationStyles = makeStyles((theme) => ({
	outline: {
		border: "1px solid #ddd",
		margin: theme.spacing(2),
		backgroundColor: "#f6f8fa",
	},
	notificationHeader: {
		padding: theme.spacing(1),
		borderBottom: "1px solid #ddd",
		fontWeight: "bold",
	},
	notificationTitle: {
		padding: theme.spacing(1),
	},
	body: {
		display: "flex",
		justifyContent: "space-between",
	},
	readButton: {
		border: 0,
		backgroundColor: "white",
		paddingRight: theme.spacing(2),
	},
}));

function Notification({ notification }) {
	const classes = useNotificationStyles();

	return (
		<div className={classes.outline}>
			<Typography variant="body2" className={classes.notificationHeader}>
				{notification.repo.user.username}/{notification.repo.name}
			</Typography>
			<div className={classes.body}>
				<Typography variant="body2" className={classes.notificationTitle}>
					{notification.subject.title}
				</Typography>
				<button className={classes.readButton}>
					<DoneIcon />
				</button>
			</div>
		</div>
	);
}

const useNotificationsStyles = makeStyles((theme) => ({
	container: {
		margin: "10px auto",
	},
	tabs: {
		margin: " 20px 400px ",
	},
}));

const Notifications = () => {
	const classes = useNotificationsStyles();
	const [notifications, setNotifications] = useState([]);

	const fetchNotifications = async () => {
		try {
			const response = await axios.get(`/notifications`, authHeaders());
			if (response.status === 200) {
				setNotifications(response.data);
			}
		} catch (error) {
			console.log(error);
		}
	};

	useEffect(() => {
		fetchNotifications();
	}, []);

	return (
		<>
			<Navbar />
			<UnderlineNav className={classes.tabs}>
				<UnderlineNav.Link href="/notifications">
					Notifications
				</UnderlineNav.Link>
				<UnderlineNav.Link href="/watching">Watching</UnderlineNav.Link>
				<UnderlineNav.Link href="/subscriptions">
					Subscripitions
				</UnderlineNav.Link>
			</UnderlineNav>

			<div className={cx(classes.container, "container")}>
				<div className="row">
					<div className="col-12"></div>
				</div>
				<div className="row">
					<div className="col-3">
						<Sidebar />
					</div>
					<div className="col-9">
						{notifications.map((notification) => (
							<Notification notification={notification} />
						))}
					</div>
				</div>
			</div>
		</>
	);
};

export default Notifications;
