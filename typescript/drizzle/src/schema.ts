import { relations } from "drizzle-orm";
import { pgTable, serial, text } from "drizzle-orm/pg-core";

export const org = pgTable("org", {
  id: serial("id").primaryKey(),
  name: text("name").notNull().unique(),
  users: text("users").notNull().default("[]"),
});

export const orgRelations = relations(org, ({ many }) => ({
  users: many(user),
}));

export const user = pgTable("user", {
  id: serial("id").primaryKey(),
  email: text("email").notNull().unique(),
  firstName: text("first_name").notNull(),
  lastName: text("last_name").notNull(),
});
