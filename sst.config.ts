import { SSTConfig } from "sst";
import { Api } from "sst/constructs";

export default {
  config(_input) {
    return {
      name: "currency-exchange-api-go",
      region: "us-east-1",
    };
  },
  stacks(app) {
    app.setDefaultFunctionProps({
      runtime: "go1.x",
    });
    app.stack(function Stack({ stack }) {
      const api = new Api(stack, "api", {
        routes: {
          "GET /": "functions/lambda/exchange.go",
        },
      });
      stack.addOutputs({
        ApiEndpoint: api.url,
      });
    });
  },
} satisfies SSTConfig;
