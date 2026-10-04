# Effect guidance

This app uses Effect 4.

Before writing Effect code, read `node_modules/effect/AGENTS.md` completely from
this directory and follow its relevant links. Resolve API questions against the
installed source in `node_modules/effect/src`.

Keep API definitions in `src/api.ts` separate from server implementations in
`src/handlers.ts`. HTTP modules are imported from `effect/http` and
`effect/http-api`. Keep `effect` and `@effect/platform-node` on the same version.
