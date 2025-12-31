-- Create "documents" table
CREATE TABLE "public"."documents" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "description" text NULL,
  "user_id" text NULL,
  "file_name" text NULL,
  "file_size" bigint NULL,
  "object_key" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_documents_deleted_at" to table: "documents"
CREATE INDEX "idx_documents_deleted_at" ON "public"."documents" ("deleted_at");
