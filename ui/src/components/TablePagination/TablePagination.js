import React from "react";
import { makeStyles } from "@material-ui/core/styles";
import Toolbar from "@material-ui/core/Toolbar";
import IconButton from "@material-ui/core/IconButton";
import ChevronLeft from "@material-ui/icons/ChevronLeft";
import ChevronRight from "@material-ui/icons/ChevronRight";

const useStyles = makeStyles((theme) => ({
  toolbar: {
    height: 56,
    minHeight: 56,
    paddingRight: 2,
  },
  spacer: {
    flex: "1 1 100%",
  },
  actions: {
    flexShrink: 0,
    color: theme.palette.text.secondary,
    marginLeft: theme.spacing(2.5),
  },
}));

function TablePagination({
  handleNextPage,
  handlePreviousPage,
  hasNextPage,
  hasPreviousPage,
}) {
  // Get classes
  const classes = useStyles();

  return (
    <Toolbar className={classes.toolbar}>
      <div className={classes.spacer} />
      <div className={classes.actions}>
        <IconButton disabled={!hasPreviousPage} onClick={handlePreviousPage}>
          <ChevronLeft />
        </IconButton>
        <IconButton disabled={!hasNextPage} onClick={handleNextPage}>
          <ChevronRight />
        </IconButton>
      </div>
    </Toolbar>
  );
}

export default TablePagination;
