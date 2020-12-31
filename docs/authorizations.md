# Authorizations

Here is the list of authorizations available in the application per domain and per action.

This will be used in OPA servers with using this [format](opa-formats.md).

## Partitions

| Action                     | OPA Action                            | OPA Resource                   | GraphQL field                                                                                                            |
| -------------------------- | ------------------------------------- | ------------------------------ | ------------------------------------------------------------------------------------------------------------------------ |
| Get All                    | `partitions:List`                     | `""`                           | Object: Query / Field: `partitions`                                                                                      |
| Create                     | `partitions:Create`                   | `partitions:${partition-name}` | Object: Mutation / Field: `createPartition`                                                                              |
| Update                     | `partitions:Update`                   | `partitions:${partition-name}` | Object: Mutation / Field: `updatePartition`                                                                              |
| Find By ID                 | `partitions:FindByID`                 | `partitions:${id}`             | Object: Query -> Field: `partition` // Object: DecisionLog -> Field: `partition` // Object: Status -> Field: `partition` |
| Generate OPA Configuration | `partitions:GenerateOPAConfiguration` | `partitions:${id}`             | Object: Partition / Field: `opaConfiguration`                                                                            |

## Decisions

| Action              | OPA Action              | OPA Resource         | GraphQL field                              |
| ------------------- | ----------------------- | -------------------- | ------------------------------------------ |
| Find By Decision ID | `decisionlogs:FindByID` | `decisionlogs:${id}` | Object: Query / Field: `decisionLog`       |
| Find By ID          | `decisionlogs:FindByID` | `decisionlogs:${id}` | Object: Query / Field: `decisionLog`       |
| Get All             | `decisionlogs:List`     | `""`                 | Object: Partition / Field: `decisionLogs`  |

## Statuses

| Action     | OPA Action          | OPA Resource     | GraphQL field                          |
| ---------- | ------------------- | ---------------- | -------------------------------------- |
| Find By ID | `statuses:FindByID` | `statuses:${id}` | Object: Query / Field: `status`        |
| Get All    | `statuses:List`     | `""`             | Object: Partition / Field: `statuses`  |
