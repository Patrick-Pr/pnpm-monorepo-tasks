import { Console, Effect } from "effect";

const program = Console.log("Hello, World!");

const result = Effect.runSync(program); // => undefined
