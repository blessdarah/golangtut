-- Modify "tickets" table
ALTER TABLE "public"."tickets" DROP CONSTRAINT "tickets_pkey", DROP COLUMN "ticket_id", ADD COLUMN "id" text NOT NULL, ADD PRIMARY KEY ("id");
-- Modify "payments" table
ALTER TABLE "public"."payments" ADD CONSTRAINT "fk_payments_ticket" FOREIGN KEY ("ticket_id") REFERENCES "public"."tickets" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
