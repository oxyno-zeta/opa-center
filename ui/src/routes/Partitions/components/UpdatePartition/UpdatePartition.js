import React from "react";
import Button from "@material-ui/core/Button";
import Dialog from "@material-ui/core/Dialog";
import DialogTitle from "@material-ui/core/DialogTitle";
import DialogContent from "@material-ui/core/DialogContent";
import DialogActions from "@material-ui/core/DialogActions";
import Typography from "@material-ui/core/Typography";
import { useForm } from "react-hook-form";
import FormControl from "@material-ui/core/FormControl";
import FormHelperText from "@material-ui/core/FormHelperText";
import Input from "@material-ui/core/Input";
import InputLabel from "@material-ui/core/InputLabel";
import { gql, useMutation } from "@apollo/client";
import CenterLoadingSpinner from "../../../../components/CenterLoadingSpinner";

const UPDATE_PARTITION = gql`
  mutation updatePartition(
    $id: ID!
    $statusDataRetention: String
    $decisionLogRetention: String
  ) {
    updatePartition(
      input: {
        id: $id
        statusDataRetention: $statusDataRetention
        decisionLogRetention: $decisionLogRetention
      }
    ) {
      partition {
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
`;

const durationRegex = /^\d+[smh]+(?:\d+[smh]+)*$/;

function UpdatePartition({ partition, isOpened, handleClose }) {
  // Form hook
  const { register, handleSubmit, errors: formErrors } = useForm({
    defaultValues: {
      statusDataRetention: partition.statusDataRetention,
      decisionLogRetention: partition.decisionLogRetention,
    },
  });
  // Mutation hook
  const [updatePartition, { loading, error }] = useMutation(UPDATE_PARTITION);
  // Callback for form answer
  const onSubmit = async (data) => {
    try {
      await updatePartition({
        variables: {
          id: partition.id,
          statusDataRetention: data.statusDataRetention,
          decisionLogRetention: data.decisionLogRetention,
        },
      });
      // Close modal
      handleClose();
    } catch (e) {}
  };

  return (
    <Dialog
      onClose={handleClose}
      aria-labelledby="dialog-title"
      open={isOpened}
    >
      <form noValidate onSubmit={handleSubmit(onSubmit)}>
        <DialogTitle id="dialog-title">Update Partition</DialogTitle>
        <DialogContent dividers>
          {loading && <CenterLoadingSpinner />}
          {error && (
            <Typography color="error" style={{ marginBottom: "5px" }}>
              Error: {error.message}
            </Typography>
          )}
          <FormControl error={!!formErrors.statusDataRetention} fullWidth>
            <InputLabel htmlFor="statusDataRetention">
              Status Data Retention
            </InputLabel>
            <Input
              inputRef={register({
                pattern: durationRegex,
              })}
              id="statusDataRetention"
              fullWidth
              label="Status Data Retention"
              name="statusDataRetention"
            />
            {formErrors.statusDataRetention && (
              <FormHelperText>
                Duration must contains only digits and seconds (s), minutes (m)
                and hours (h). Example: 1h2m3s.
              </FormHelperText>
            )}
          </FormControl>
          <FormControl
            error={!!formErrors.decisionLogRetention}
            fullWidth
            style={{ marginTop: "10px" }}
          >
            <InputLabel htmlFor="decisionLogRetention">
              Decision Logs Retention
            </InputLabel>
            <Input
              inputRef={register({
                pattern: durationRegex,
              })}
              id="decisionLogRetention"
              fullWidth
              label="Decision Logs Retention"
              name="decisionLogRetention"
            />
            {formErrors.decisionLogRetention && (
              <FormHelperText>
                Duration must contains only digits and seconds (s), minutes (m)
                and hours (h). Example: 1h2m3s.
              </FormHelperText>
            )}
          </FormControl>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose} disabled={loading} color="primary">
            Cancel
          </Button>
          <Button type="submit" disabled={loading} color="primary">
            Submit
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  );
}

export default UpdatePartition;
