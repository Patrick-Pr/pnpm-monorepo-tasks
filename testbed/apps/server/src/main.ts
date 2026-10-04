import { NodeHttpServer, NodeRuntime } from "@effect/platform-node";
import { Config, Layer } from "effect";
import { HttpRouter } from "effect/http";
import { HttpApiSwagger } from "effect/http-api";
import { createServer } from "node:http";
import { Api } from "./api.js";
import { ApiLive } from "./handlers.js";

const ServerLive = NodeHttpServer.layerConfig(createServer, {
  host: Config.String("HOST").pipe(Config.withDefault("127.0.0.1")),
  port: Config.Port("PORT").pipe(Config.withDefault(3000)),
});

const HttpLive = HttpRouter.serve(
  Layer.mergeAll(ApiLive, HttpApiSwagger.layer(Api, { path: "/docs" })),
).pipe(
  Layer.provide(ServerLive),
);

NodeRuntime.runMain(Layer.launch(HttpLive));
