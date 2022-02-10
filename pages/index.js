import React, { Component, Fragment } from "react";
import Router from "next/router";

const Index = () => {
	useEffect(() => {
		Router.push("/join");
	}, []);

	return <div></div>;
};

export default Index;
