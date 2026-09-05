-- Modify "payments" table
ALTER TABLE "public"."payments" DROP CONSTRAINT "fk_payments_ticket";
-- Modify "tickets" table
ALTER TABLE "public"."tickets" DROP CONSTRAINT "tickets_pkey", DROP COLUMN "id", ADD COLUMN "ticket_id" text NOT NULL, ADD PRIMARY KEY ("ticket_id"), ADD CONSTRAINT "fk_payments_ticket" FOREIGN KEY ("ticket_id") REFERENCES "public"."payments" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
