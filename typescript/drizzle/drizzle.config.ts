import { Config } from "drizzle-kit";

export default {
  schema: "src/schema.ts",
  out: ".drizzle",
  driver: "pg",
  dbCredentials: {
    connectionString: "postgres://dev:password@localhost:5433/drizzle",
  },
} satisfies Config;
