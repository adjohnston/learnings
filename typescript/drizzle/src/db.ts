import * as schema from "./schema";
import { drizzle } from "drizzle-orm/node-postgres";
import { Client } from "pg";

const client = new Client({
  connectionString: "postgres://dev:password@localhost:5433/drizzle",
});

(async () => {
  await client.connect();
})();

export const db = drizzle(client, { schema });
