import React from "react";
import App from "next/app";

import { ThemeProvider } from "@material-ui/styles";

import { createMuiTheme } from "@material-ui/core/styles";
import blue from "@material-ui/core/colors/blue";

import axios from "axios";

axios.defaults.baseURL = process.env.API_URL;

export const theme = createMuiTheme({
  palette: {
    primary: {
      main: "#00F"
    }
  }
});

class MyApp extends App {
  // Only uncomment this method if you have blocking data requirements for
  // every single page in your application. This disables the ability to
  // perform automatic static optimization, causing every page in your app to
  // be server-side rendered.
  //
  // static async getInitialProps(appContext) {
  //   // calls page's `getInitialProps` and fills `appProps.pageProps`
  //   const appProps = await App.getInitialProps(appContext);
  //
  //   return { ...appProps }
  // }

  render() {
    const { Component, pageProps } = this.props;
    return (
      <ThemeProvider theme={theme}>
        <Component {...pageProps} />
      </ThemeProvider>
    );
  }
}

export default MyApp;
