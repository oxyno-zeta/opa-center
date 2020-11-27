import React from "react";
import { gql, useQuery } from "@apollo/client";
import { useParams, Link } from "react-router-dom";
import Grid from "@material-ui/core/Grid";
import Typography from "@material-ui/core/Typography";
import { makeStyles } from "@material-ui/core/styles";
import IconButton from "@material-ui/core/IconButton";
import ChevronLeft from "@material-ui/icons/ChevronLeft";
import CenterLoadingSpinner from "../../components/CenterLoadingSpinner";
import GraphqlErrors from "../../components/GraphqlErrors";
import PageTitle from "../../components/PageTitle";
import DecisionLogCard from "./components/DecisionLogCard";
import TablePagination from "../../components/TablePagination";

const GET_DECISION_LOGS = gql`
  query getDecisionLogs(
    $partitionId: ID!
    $after: String
    $before: String
    $first: Int
    $last: Int
  ) {
    partition(id: $partitionId) {
      id
      decisionLogs(after: $after, before: $before, first: $first, last: $last) {
        pageInfo {
          hasNextPage
          hasPreviousPage
          startCursor
          endCursor
        }
        edges {
          cursor
          node {
            id
            createdAt
            updatedAt
            decisionId
            path
            requestedBy
            timestamp
            originalMessage
          }
        }
      }
    }
  }
`;

// Create styles
const useStyles = makeStyles((theme) => ({
  gridRoot: { marginTop: theme.spacing(2), marginBottom: theme.spacing(2) },
  center: {
    textAlign: "center",
    display: "flex",
    justifyContent: "center",
    alignContent: "center",
    flexDirection: "column",
  },
}));

function DecisionLogs() {
  // Get classes
  const classes = useStyles();
  // Get params
  const { id } = useParams();
  // MAX PAGINATION
  const MAX_PAGINATION = 10;
  // Query variables
  const qVariables = { partitionId: id, first: MAX_PAGINATION };
  // Query data
  const { loading, error, data, fetchMore } = useQuery(GET_DECISION_LOGS, {
    variables: qVariables,
  });

  const handleNextPage = async () => {
    const variables = {
      ...qVariables,
      after: data.partition.decisionLogs.pageInfo.endCursor,
      first: MAX_PAGINATION,
    };
    await fetchMore({
      variables,
      updateQuery: (previousResult, { fetchMoreResult }) => fetchMoreResult,
    });
  };

  const handlePreviousPage = async () => {
    const variables = {
      ...qVariables,
      before: data.partition.decisionLogs.pageInfo.startCursor,
      last: MAX_PAGINATION,
      first: undefined,
    };
    await fetchMore({
      variables,
      updateQuery: (previousResult, { fetchMoreResult }) => fetchMoreResult,
    });
  };

  // Check if loading is enabled to display loading
  if (loading) {
    return (
      <>
        <PageTitle title={<>Decision Logs</>} />
        <CenterLoadingSpinner />
      </>
    );
  }

  // Check if error is raised to display errors
  if (error) return <GraphqlErrors error={error} />;

  let content = null;

  // Check if data exists or not
  if (
    !data.partition.decisionLogs.edges ||
    data.partition.decisionLogs.edges.length === 0
  ) {
    content = (
      <div className={classes.center}>
        <Typography component="h1" variant="h5">
          No data available
        </Typography>
      </div>
    );
  } else {
    content = (
      <>
        <Grid container spacing={3}>
          {data.partition.decisionLogs.edges.map(({ node }) => (
            <Grid
              key={node.id}
              item
              xs={12}
              classes={{ root: classes.gridRoot }}
            >
              <DecisionLogCard decisionLog={node} />
            </Grid>
          ))}
        </Grid>
        <TablePagination
          handleNextPage={handleNextPage}
          handlePreviousPage={handlePreviousPage}
          hasNextPage={data.partition.decisionLogs.pageInfo.hasNextPage}
          hasPreviousPage={data.partition.decisionLogs.pageInfo.hasPreviousPage}
        />
      </>
    );
  }

  // Display
  return (
    <>
      <PageTitle
        title={<>Decision Logs</>}
        leftElement={
          <IconButton variant="contained" component={Link} to={"/"}>
            <ChevronLeft />
          </IconButton>
        }
      />

      {content}
    </>
  );
}

export default DecisionLogs;
