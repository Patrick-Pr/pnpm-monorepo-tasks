import { Schema } from "effect";
import { HttpApi, HttpApiEndpoint, HttpApiGroup } from "effect/http-api";

export const Api = HttpApi.make("server").add(
  HttpApiGroup.make("system").add(
    HttpApiEndpoint.get("hello", "/", {
      success: Schema.Struct({ message: Schema.String }),
    }),
    HttpApiEndpoint.get("health", "/health", {
      success: Schema.Struct({ status: Schema.Literal("ok") }),
    }),
  ),
);
