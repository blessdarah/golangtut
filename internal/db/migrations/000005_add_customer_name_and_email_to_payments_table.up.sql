-- Modify "payments" table
ALTER TABLE "public"."payments" ADD COLUMN "customer_name" character varying(80) NOT NULL, ADD COLUMN "customer_email" character varying(255) NOT NULL;
