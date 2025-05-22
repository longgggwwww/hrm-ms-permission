-- Create "perm_groups" table
CREATE TABLE "public"."perm_groups" ("id" uuid NOT NULL, "code" character varying NOT NULL, "name" character varying NOT NULL, PRIMARY KEY ("id"));
-- Create index "perm_groups_code_key" to table: "perm_groups"
CREATE UNIQUE INDEX "perm_groups_code_key" ON "public"."perm_groups" ("code");
-- Create "perms" table
CREATE TABLE "public"."perms" ("id" uuid NOT NULL, "code" character varying NOT NULL, "name" character varying NOT NULL, "description" character varying NULL, "perm_group" uuid NULL, PRIMARY KEY ("id"), CONSTRAINT "perms_perm_groups_group" FOREIGN KEY ("perm_group") REFERENCES "public"."perm_groups" ("id") ON UPDATE NO ACTION ON DELETE SET NULL);
-- Create index "perms_code_key" to table: "perms"
CREATE UNIQUE INDEX "perms_code_key" ON "public"."perms" ("code");
-- Create "roles" table
CREATE TABLE "public"."roles" ("id" uuid NOT NULL, "code" character varying NOT NULL, "name" character varying NOT NULL, "color" character varying NULL, "description" character varying NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, PRIMARY KEY ("id"));
-- Create index "roles_code_key" to table: "roles"
CREATE UNIQUE INDEX "roles_code_key" ON "public"."roles" ("code");
-- Create "role_perms" table
CREATE TABLE "public"."role_perms" ("role_id" uuid NOT NULL, "perm_id" uuid NOT NULL, PRIMARY KEY ("role_id", "perm_id"), CONSTRAINT "role_perms_perm_id" FOREIGN KEY ("perm_id") REFERENCES "public"."perms" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "role_perms_role_id" FOREIGN KEY ("role_id") REFERENCES "public"."roles" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create "user_perms" table
CREATE TABLE "public"."user_perms" ("id" uuid NOT NULL, "user_id" character varying NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, "perm_id" uuid NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "user_perms_perms_user_perms" FOREIGN KEY ("perm_id") REFERENCES "public"."perms" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION);
-- Create index "userperm_perm_id_user_id" to table: "user_perms"
CREATE UNIQUE INDEX "userperm_perm_id_user_id" ON "public"."user_perms" ("perm_id", "user_id");
-- Create "user_roles" table
CREATE TABLE "public"."user_roles" ("id" uuid NOT NULL, "user_id" character varying NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, "role_id" uuid NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "user_roles_roles_user_roles" FOREIGN KEY ("role_id") REFERENCES "public"."roles" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION);
-- Create index "userrole_role_id_user_id" to table: "user_roles"
CREATE UNIQUE INDEX "userrole_role_id_user_id" ON "public"."user_roles" ("role_id", "user_id");
