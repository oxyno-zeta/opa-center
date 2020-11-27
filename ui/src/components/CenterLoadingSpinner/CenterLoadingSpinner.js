import React from "react";
import classnames from "classnames";
import CircularProgress from "@material-ui/core/CircularProgress";
import { makeStyles } from "@material-ui/core/styles";

// Create styles
const useStyles = makeStyles((theme) => ({
  progressContainer: {
    flex: 1,
    display: "flex",
    flexDirection: "column",
    position: "relative",
    alignItems: "center",
  },
  progress: {
    margin: theme.spacing(4),
  },
  spaceTop: { marginTop: theme.spacing(6) },
}));

function CenterLoadingSpinner({ noTopSpaceNeeded }) {
  // Get classes
  const classes = useStyles();

  return (
    <div
      className={classnames(classes.progressContainer, {
        [classes.spaceTop]: !noTopSpaceNeeded,
      })}
    >
      <CircularProgress className={classes.progress} size={50} />
    </div>
  );
}

export default CenterLoadingSpinner;
