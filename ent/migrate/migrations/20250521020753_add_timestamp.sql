-- Modify "user_roles" table
ALTER TABLE "public"."user_roles" ADD COLUMN "created_at" timestamptz NOT NULL, ADD COLUMN "updated_at" timestamptz NOT NULL;
-- Modify "user_perms" table
ALTER TABLE "public"."user_perms" ADD COLUMN "created_at" timestamptz NOT NULL, ADD COLUMN "updated_at" timestamptz NOT NULL, ADD CONSTRAINT "user_perms_perms_user_perms" FOREIGN KEY ("perm_id") REFERENCES "public"."perms" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
