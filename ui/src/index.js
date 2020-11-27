import React from "react";
import ReactDOM from "react-dom";
import { BrowserRouter as Router } from "react-router-dom";
import { ApolloProvider } from "@apollo/client";
import * as dayjs from "dayjs";
import localizedFormat from "dayjs/plugin/localizedFormat";
import "./index.css";
import TopBar from "./components/TopBar";
import Routes from "./routes";
import reportWebVitals from "./reportWebVitals";
import client from "./client";

// Extend dayjs
dayjs.extend(localizedFormat);

ReactDOM.render(
  <React.StrictMode>
    <ApolloProvider client={client}>
      <Router>
        <TopBar />
        <Routes />
      </Router>
    </ApolloProvider>
  </React.StrictMode>,
  document.getElementById("root")
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
