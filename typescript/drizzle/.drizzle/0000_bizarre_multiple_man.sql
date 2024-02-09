CREATE TABLE IF NOT EXISTS "org" (
	"id" serial PRIMARY KEY NOT NULL,
	"name" text NOT NULL,
	"users" text DEFAULT '[]' NOT NULL,
	CONSTRAINT "org_name_unique" UNIQUE("name")
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "user" (
	"id" serial PRIMARY KEY NOT NULL,
	"email" text NOT NULL,
	"first_name" text NOT NULL,
	"last_name" text NOT NULL,
	CONSTRAINT "user_email_unique" UNIQUE("email")
);
