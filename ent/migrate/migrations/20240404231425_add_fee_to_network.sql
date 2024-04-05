-- Modify "networks" table
ALTER TABLE "networks" ADD COLUMN "fee" double precision NOT NULL DEFAULT 0.0;
