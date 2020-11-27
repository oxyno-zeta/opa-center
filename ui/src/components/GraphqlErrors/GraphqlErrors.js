import React from "react";
import Typography from "@material-ui/core/Typography";
import { makeStyles } from "@material-ui/core/styles";

// Create styles
const useStyles = makeStyles((theme) => ({
  content: { margin: 8, color: theme.palette.error.main },
}));

function GraphqlErrors({ error }) {
  // Get classes
  const classes = useStyles();

  return (
    <div className={classes.content}>
      <Typography color="error">Errors:</Typography>
      <ul>
        {error.graphQLErrors.map(({ message }, i) => (
          <li key={`graphQLErrors-${i}`}>
            <Typography color="error">{message}</Typography>
          </li>
        ))}
        {error.networkError && (
          <li key="networkError">
            <Typography color="error">
              {error.networkError.toString()}
            </Typography>
          </li>
        )}
      </ul>
    </div>
  );
}

export default GraphqlErrors;
