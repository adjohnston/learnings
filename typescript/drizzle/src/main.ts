import { eq, sql } from "drizzle-orm";
import { db } from "./db.js";
import { user } from "./schema.js";

(async () => {
  // await db
  //   .insert(user)
  //   .values({
  //     firstName: "Dave",
  //     lastName: "Smith",
  //     email: "example@example.com",
  //   })
  //   .returning({ id: user.id });

  await db
    .selectDistinct()
    .from(user)
    .where(eq(user.email, "bob@example.com"))
    .then(console.log);

  process.exit();
})();
