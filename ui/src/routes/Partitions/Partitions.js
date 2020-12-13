import React from "react";
import { gql, useQuery } from "@apollo/client";
import Grid from "@material-ui/core/Grid";
import Typography from "@material-ui/core/Typography";
import { makeStyles } from "@material-ui/core/styles";
import IconButton from "@material-ui/core/IconButton";
import AddIcon from "@material-ui/icons/Add";
import CenterLoadingSpinner from "../../components/CenterLoadingSpinner";
import GraphqlErrors from "../../components/GraphqlErrors";
import PageTitle from "../../components/PageTitle";
import PartitionCard from "./components/PartitionCard";
import TablePagination from "../../components/TablePagination";
import CreatePartition from "./components/CreatePartition";
import UpdatePartition from "./components/UpdatePartition";

const GET_PARTITIONS = gql`
  query getPartitions(
    $after: String
    $before: String
    $first: Int
    $last: Int
  ) {
    partitions(after: $after, before: $before, first: $first, last: $last) {
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
          name
          opaConfiguration
          statusDataRetention
          decisionLogRetention
        }
      }
    }
  }
`;

// Create styles
const useStyles = makeStyles((theme) => ({
  gridRoot: { marginTop: theme.spacing(2), marginBottom: theme.spacing(2) },
}));

function Partitions() {
  // Get classes
  const classes = useStyles();
  // MAX PAGINATION
  const MAX_PAGINATION = 10;
  // Query variables
  const qVariables = { first: MAX_PAGINATION };
  // Save variables
  const [variables, setVariables] = React.useState(qVariables);
  // Query data
  const { loading, error, data, fetchMore, refetch } = useQuery(
    GET_PARTITIONS,
    {
      variables,
    }
  );
  // Get is create modal opened
  const [isCreateModalOpened, setCreateModalOpened] = React.useState(false);
  // Get is update modal opened
  const [updateModalData, setUpdateModalData] = React.useState(null);

  const handleNextPage = async () => {
    const variables = {
      ...qVariables,
      after: data.partitions.pageInfo.endCursor,
      first: MAX_PAGINATION,
      last: undefined,
      before: undefined,
    };
    // Fetch more
    await fetchMore({
      variables,
      updateQuery: (previousResult, { fetchMoreResult }) => fetchMoreResult,
    });
    // Save variables for refetch
    setVariables(variables);
  };

  const handlePreviousPage = async () => {
    const variables = {
      ...qVariables,
      before: data.partitions.pageInfo.startCursor,
      last: MAX_PAGINATION,
      first: undefined,
      after: undefined,
    };
    // Fetch more
    await fetchMore({
      variables,
      updateQuery: (previousResult, { fetchMoreResult }) => fetchMoreResult,
    });
    // Save variables for refetch
    setVariables(variables);
  };

  // Check if loading is enabled to display loading
  if (loading) {
    return (
      <>
        <PageTitle title={<>Partitions</>} />
        <CenterLoadingSpinner />
      </>
    );
  }

  // Check if error is raised to display errors
  if (error) return <GraphqlErrors error={error} />;

  let content = null;

  // Check if data exists or not
  if (!data.partitions.edges || data.partitions.edges.length === 0) {
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
          {data.partitions.edges.map(({ node }) => (
            <Grid
              key={node.id}
              item
              xs={12}
              classes={{ root: classes.gridRoot }}
            >
              <PartitionCard
                partition={node}
                openEdit={() => {
                  setUpdateModalData(node);
                }}
              />
            </Grid>
          ))}
        </Grid>
        <TablePagination
          handleNextPage={handleNextPage}
          handlePreviousPage={handlePreviousPage}
          hasNextPage={data.partitions.pageInfo.hasNextPage}
          hasPreviousPage={data.partitions.pageInfo.hasPreviousPage}
        />
      </>
    );
  }

  // Display
  return (
    <>
      <PageTitle
        title={<>Partitions</>}
        rightElement={
          <IconButton
            onClick={() => {
              setCreateModalOpened(true);
            }}
          >
            <AddIcon />
          </IconButton>
        }
      />

      {content}
      <CreatePartition
        isOpened={isCreateModalOpened}
        handleClose={() => {
          setCreateModalOpened(false);
        }}
        refetch={() => {
          refetch(variables);
        }}
      />
      {updateModalData && (
        <UpdatePartition
          partition={updateModalData}
          isOpened={!!updateModalData}
          handleClose={() => {
            setUpdateModalData(null);
          }}
        />
      )}
    </>
  );
}

export default Partitions;
