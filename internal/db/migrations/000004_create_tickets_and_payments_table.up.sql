-- Create "tickets" table
CREATE TABLE "public"."tickets" (
  "id" text NOT NULL,
  "type" character varying(30) NOT NULL,
  "price" numeric NOT NULL,
  "event_id" text NOT NULL,
  "description" text NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_tickets_event" FOREIGN KEY ("event_id") REFERENCES "public"."events" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "tickets_event_id_idx" to table: "tickets"
CREATE INDEX "tickets_event_id_idx" ON "public"."tickets" ("event_id");
-- Create "payments" table
CREATE TABLE "public"."payments" (
  "id" text NOT NULL,
  "event_id" text NOT NULL,
  "ticket_id" text NOT NULL,
  "amount" numeric NOT NULL,
  "quantity" integer NOT NULL,
  "total" numeric NOT NULL,
  "payment_provider" character varying(30) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_events_payments" FOREIGN KEY ("event_id") REFERENCES "public"."events" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_payments_ticket" FOREIGN KEY ("ticket_id") REFERENCES "public"."tickets" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "payments_event_id_idx" to table: "payments"
CREATE INDEX "payments_event_id_idx" ON "public"."payments" ("event_id");
-- Create index "payments_ticket_id_idx" to table: "payments"
CREATE INDEX "payments_ticket_id_idx" ON "public"."payments" ("ticket_id");
