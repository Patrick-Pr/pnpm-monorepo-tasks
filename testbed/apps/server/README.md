# Server

A basic Effect 4 HTTP API using `effect/http-api`, `effect/http`, and the Node
runtime. `effect` and `@effect/platform-node` are pinned to the same version.

From the `testbed` directory:

```sh
pnpm install
pnpm --filter server dev:watch
```

The server listens on `http://127.0.0.1:3000` by default.

| Endpoint | Response |
| --- | --- |
| `GET /` | `{ "message": "Hello, World!" }` |
| `GET /health` | `{ "status": "ok" }` |
| `GET /docs` | Swagger UI with the generated API specification |

Set `HOST` and `PORT` to override the listen address:

```sh
HOST=0.0.0.0 PORT=4000 pnpm --filter server dev
```

Run `pnpm --filter server typecheck` to check the TypeScript code.

Define endpoint schemas in `src/api.ts` and implement their handlers in
`src/handlers.ts`. `src/main.ts` wires the routes, request logging, documentation,
and HTTP server together. The Node runtime closes the server when it receives
SIGINT or SIGTERM.
