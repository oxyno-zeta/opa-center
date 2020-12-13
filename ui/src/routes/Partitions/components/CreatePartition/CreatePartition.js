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

const CREATE_PARTITION = gql`
  mutation createPartition(
    $name: String!
    $statusDataRetention: String
    $decisionLogRetention: String
  ) {
    createPartition(
      input: {
        name: $name
        statusDataRetention: $statusDataRetention
        decisionLogRetention: $decisionLogRetention
      }
    ) {
      partition {
        id
      }
    }
  }
`;

function getNameDisplayErrorMessage(type) {
  switch (type) {
    case "maxLength":
      return "Name must have a maximum length of 255 characters.";
    case "required":
      return "Name is required.";
    case "pattern":
      return "Name must be lowercase and contains only letters, digits and dash (-) and cannot starts and ends with dash.";
    default:
      return "Unknown error.";
  }
}

const durationRegex = /^\d+[smh]+(?:\d+[smh]+)*$/;

function CreatePartition({ isOpened, handleClose, refetch }) {
  // Form hook
  const { register, handleSubmit, errors: formErrors } = useForm();
  // Mutation hook
  const [createPartition, { loading, error }] = useMutation(CREATE_PARTITION);
  // Callback for form answer
  const onSubmit = async (data) => {
    try {
      await createPartition({
        variables: {
          name: data.name,
          statusDataRetention: data.statusDataRetention,
          decisionLogRetention: data.decisionLogRetention,
        },
      });
      // Refetch data
      refetch();
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
        <DialogTitle id="dialog-title">Create Partition</DialogTitle>
        <DialogContent dividers>
          {loading && <CenterLoadingSpinner />}
          {error && (
            <Typography color="error" style={{ marginBottom: "5px" }}>
              Error: {error.message}
            </Typography>
          )}
          <FormControl error={!!formErrors.name} fullWidth>
            <InputLabel htmlFor="name">Name</InputLabel>
            <Input
              inputRef={register({
                maxLength: 255,
                required: true,
                pattern: /^[a-z0-9][a-z0-9]*(?:-+[a-z0-9]+)*$/,
              })}
              id="name"
              required
              fullWidth
              label="Name"
              name="name"
              autoFocus
            />
            {formErrors.name && (
              <FormHelperText>
                {getNameDisplayErrorMessage(formErrors.name.type)}
              </FormHelperText>
            )}
          </FormControl>
          <FormControl
            error={!!formErrors.statusDataRetention}
            fullWidth
            style={{ marginTop: "10px" }}
          >
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

export default CreatePartition;
