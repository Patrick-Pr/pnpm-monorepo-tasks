import { Effect, Layer } from "effect";
import { HttpApiBuilder } from "effect/http-api";
import { Api } from "./api.js";

const SystemLive = HttpApiBuilder.group(Api, "system", (handlers) =>
  handlers
    .handle("hello", () => Effect.succeed({ message: "Hello, World!" }))
    .handle("health", () => Effect.succeed({ status: "ok" })),
);

export const ApiLive = HttpApiBuilder.layer(Api).pipe(Layer.provide(SystemLive));
