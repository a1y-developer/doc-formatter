-- Modify "documents" table
ALTER TABLE "public"."documents" DROP COLUMN "user_id", DROP COLUMN "file_name", DROP COLUMN "file_size", ALTER COLUMN "object_key" SET NOT NULL, ADD COLUMN "parent_id" uuid NULL, ADD COLUMN "owner_id" uuid NOT NULL, ADD COLUMN "name" text NOT NULL, ADD COLUMN "size" bigint NOT NULL, ADD COLUMN "mime_type" character varying(100) NULL, ADD COLUMN "status" character varying(20) NULL DEFAULT 'PENDING', ADD COLUMN "metadata" jsonb NULL, ADD CONSTRAINT "uni_documents_object_key" UNIQUE ("object_key");
-- Create index "idx_documents_owner_id" to table: "documents"
CREATE INDEX "idx_documents_owner_id" ON "public"."documents" ("owner_id");
-- Create index "idx_documents_parent_id" to table: "documents"
CREATE INDEX "idx_documents_parent_id" ON "public"."documents" ("parent_id");
-- Create index "idx_documents_status" to table: "documents"
CREATE INDEX "idx_documents_status" ON "public"."documents" ("status");
-- Create index "idx_file_parent_name" to table: "documents"
CREATE UNIQUE INDEX "idx_file_parent_name" ON "public"."documents" ("name");
-- Create "folders" table
CREATE TABLE "public"."folders" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "description" text NULL,
  "owner_id" uuid NOT NULL,
  "parent_id" uuid NULL,
  "name" text NOT NULL,
  "path_tokens" text[] NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_documents_parent" FOREIGN KEY ("parent_id") REFERENCES "public"."documents" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_folders_parent" FOREIGN KEY ("parent_id") REFERENCES "public"."folders" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create index "idx_folders_deleted_at" to table: "folders"
CREATE INDEX "idx_folders_deleted_at" ON "public"."folders" ("deleted_at");
-- Create index "idx_folders_parent_id" to table: "folders"
CREATE INDEX "idx_folders_parent_id" ON "public"."folders" ("parent_id");
