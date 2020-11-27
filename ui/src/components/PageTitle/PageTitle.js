import React from "react";
import { makeStyles } from "@material-ui/core/styles";
import Toolbar from "@material-ui/core/Toolbar";
import Typography from "@material-ui/core/Typography";
import Divider from "@material-ui/core/Divider";

// Create styles
const useStyles = makeStyles({
  root: {
    flexGrow: 1,
  },
  flex: {
    flex: 1,
  },
  flexDisplay: {
    display: "flex",
  },
});

function PageTitle({ rightElement, title, titleElement, leftElement }) {
  // Get classes
  const classes = useStyles();

  return (
    <div className={classes.root}>
      <Toolbar>
        <div>{leftElement}</div>
        <div className={classes.flex}>
          <div className={classes.flexDisplay}>
            {title && (
              <Typography variant="h6" color="inherit">
                {title}
              </Typography>
            )}
            {titleElement}
          </div>
        </div>
        <div>{rightElement}</div>
      </Toolbar>
      <Divider />
    </div>
  );
}

export default PageTitle;
