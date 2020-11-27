import React from "react";
import { Link, withRouter } from "react-router-dom";
import { makeStyles } from "@material-ui/core/styles";
import Grid from "@material-ui/core/Grid";
import Typography from "@material-ui/core/Typography";
import Toolbar from "@material-ui/core/Toolbar";
import AppBar from "@material-ui/core/AppBar";

// Create styles
const useStyles = makeStyles((theme) => ({
  appBar: {
    position: "relative",
    boxShadow: "none",
    borderBottom: `1px solid ${theme.palette.grey["300"]}`,
  },
  inline: {
    display: "inline",
  },
  flex: {
    display: "flex",
  },
  link: {
    textDecoration: "none",
    color: "inherit",
  },
}));

function TopBar() {
  // Get classes
  const classes = useStyles();

  // Return
  return (
    <AppBar color="default" className={classes.appBar}>
      <Toolbar>
        <Grid container spacing={0} alignItems="baseline">
          <Grid item xs={12} className={classes.flex}>
            <div className={classes.inline}>
              <Typography variant="h6" color="inherit" noWrap>
                <Link to="/" className={classes.link}>
                  <span className={classes.tagline}>OPA Center</span>
                </Link>
              </Typography>
            </div>
          </Grid>
        </Grid>
      </Toolbar>
    </AppBar>
  );
}

export default withRouter(TopBar);
