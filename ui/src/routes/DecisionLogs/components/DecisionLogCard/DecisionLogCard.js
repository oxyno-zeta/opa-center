import React from "react";
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

const useStyles = makeStyles((theme) => ({
  opaCodeTitle: {
    margin: "15px 0",
  },
}));

function DecisionLogCard({ decisionLog }) {
  const classes = useStyles();
  const [expanded, setExpanded] = React.useState(false);

  const handleExpandClick = () => {
    setExpanded(!expanded);
  };

  return (
    <Card>
      <CardContent>
        <Typography variant="body2" component="p">
          Decision ID: <b>{decisionLog.decisionId}</b>
        </Typography>
        <Typography variant="body2" component="p">
          Created At: <b>{dayjs(decisionLog.createdAt).format("LLL")}</b>
        </Typography>
        <Typography variant="body2" component="p">
          Updated At: <b>{dayjs(decisionLog.updatedAt).format("LLL")}</b>
        </Typography>
        <Typography variant="body2" component="p">
          Path: <b>{decisionLog.path}</b>
        </Typography>
        <Typography variant="body2" component="p">
          Requested By: <b>{decisionLog.requestedBy}</b>
        </Typography>
        <Typography variant="body2" component="p">
          Timestamp: <b>{dayjs(decisionLog.timestamp).format("LLL")}</b>
        </Typography>
      </CardContent>
      <CardActions>
        <Button
          variant="outlined"
          size="small"
          onClick={handleExpandClick}
          classes={{ root: classes.button }}
        >
          {expanded ? "Hide" : "See"} original message
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
            Original JSON message content:
          </Typography>
          <Editor
            height="250px"
            language="json"
            theme="light"
            value={JSON.stringify(
              JSON.parse(decisionLog.originalMessage),
              null,
              2
            )}
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

export default DecisionLogCard;
