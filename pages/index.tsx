import React, { useEffect } from "react";
import Router from "next/router";

const Index = () => {
	useEffect(() => {
		Router.push("/join");
	}, []);

	return <></>;
};

export default Index;
