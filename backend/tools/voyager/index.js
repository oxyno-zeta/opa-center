const server = require("express");
const { express: voyagerMiddleware } = require("graphql-voyager/middleware");

const app = server();
const port = 3000;
let graphqlEndpoint = process.env.GRAPHQL_ENDPOINT;

if (graphqlEndpoint == "") {
  graphqlEndpoint = "http://localhost:8080/api/graphql";
}

app.use("/", voyagerMiddleware({ endpointUrl: graphqlEndpoint }));

app.listen(port, () =>
  console.log(`voyager listening: http://localhost:${port}`)
);
