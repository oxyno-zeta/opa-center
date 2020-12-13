import React from "react";
import { Link as RouterLink } from "react-router-dom";
import { makeStyles } from "@material-ui/core/styles";
import Card from "@material-ui/core/Card";
import CardActions from "@material-ui/core/CardActions";
import CardContent from "@material-ui/core/CardContent";
import Button from "@material-ui/core/Button";
import Typography from "@material-ui/core/Typography";
import Collapse from "@material-ui/core/Collapse";
import Editor from "@monaco-editor/react";
import Divider from "@material-ui/core/Divider";
import * as dayjs from "dayjs";

const useStyles = makeStyles(() => ({
  opaCodeTitle: {
    margin: "15px 0",
  },
}));

function PartitionCard({ partition, openEdit }) {
  const classes = useStyles();
  const [expanded, setExpanded] = React.useState(false);

  const handleExpandClick = () => {
    setExpanded(!expanded);
  };

  return (
    <Card>
      <CardContent>
        <Typography variant="body2" component="p">
          Name: <b>{partition.name}</b>
        </Typography>
        <Typography variant="body2" component="p">
          Created At: <b>{dayjs(partition.createdAt).format("LLL")}</b>
        </Typography>
        <Typography variant="body2" component="p">
          Updated At: <b>{dayjs(partition.updatedAt).format("LLL")}</b>
        </Typography>
        <Typography variant="body2" component="p">
          Decision logs retention:{" "}
          <b>{partition.decisionLogRetention || "Unlimited"}</b>
        </Typography>
        <Typography variant="body2" component="p">
          Status data retention:{" "}
          <b>{partition.statusDataRetention || "Unlimited"}</b>
        </Typography>
      </CardContent>
      <CardActions>
        <Button
          variant="outlined"
          size="small"
          onClick={handleExpandClick}
          classes={{ root: classes.button }}
        >
          {expanded ? "Hide" : "See"} OPA configuration
        </Button>
        <Button
          variant="outlined"
          size="small"
          classes={{ root: classes.button }}
          component={RouterLink}
          to={`/partitions/${partition.id}/decision-logs/`}
        >
          Go to Decision Logs
        </Button>
        <Button
          variant="outlined"
          size="small"
          classes={{ root: classes.button }}
          component={RouterLink}
          to={`/partitions/${partition.id}/statuses/`}
        >
          Go to Status data
        </Button>
        <Button
          variant="outlined"
          size="small"
          onClick={openEdit}
          classes={{ root: classes.button }}
        >
          Edit
        </Button>
      </CardActions>
      <Collapse in={expanded} timeout="auto" unmountOnExit>
        <CardContent>
          <Divider />
          <Typography
            variant="body"
            component="p"
            className={classes.opaCodeTitle}
          >
            OPA YAML configuration content:
          </Typography>
          <Editor
            height="250px"
            language="yaml"
            theme="light"
            value={partition.opaConfiguration}
            options={{
              readOnly: true,
              scrollBeyondLastLine: false,
            }}
          />
        </CardContent>
      </Collapse>
    </Card>
  );
}

export default PartitionCard;
