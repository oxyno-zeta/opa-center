import React from "react";
import { Route, Switch, Redirect } from "react-router-dom";
import { makeStyles } from "@material-ui/core/styles";
import Partitions from "./routes/Partitions";
import Statuses from "./routes/Statuses";
import DecisionLogs from "./routes/DecisionLogs";

// Create styles
const useStyles = makeStyles((theme) => ({
  space: { margin: theme.spacing(4), marginTop: 0 },
}));

function Routes() {
  // Get classes
  const classes = useStyles();

  return (
    <div className={classes.space}>
      <Switch>
        <Route exact path="/" component={Partitions} />
        <Route exact path="/partitions/:id/statuses" component={Statuses} />
        <Route
          exact
          path="/partitions/:id/decision-logs"
          component={DecisionLogs}
        />
        <Route path="*">
          <Redirect to="/"></Redirect>
        </Route>
      </Switch>
    </div>
  );
}

export default Routes;
