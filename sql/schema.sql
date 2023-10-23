CREATE TABLE "public"."draw" (
    "draw_id" int8 NOT NULL,
    "winner_ticket" varchar(50),
    "drew_at" timestamptz NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "update_at" timestamptz,
    PRIMARY KEY ("draw_id")
);

CREATE TABLE "public"."ticket" (
    "ticket_id" varchar(50) NOT NULL,
    "draw_id" int8 NOT NULL,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("ticket_id")
);

ALTER TABLE "public"."draw"
    ADD CONSTRAINT "draw_winner_ticket_fkey"
        FOREIGN KEY ("winner_ticket")
        REFERENCES "public"."ticket"("ticket_id") ON DELETE SET NULL ON UPDATE RESTRICT;

ALTER TABLE "public"."ticket"
    ADD CONSTRAINT "ticket_draw_id_fkey"
        FOREIGN KEY ("draw_id")
        REFERENCES "public"."draw"("draw_id") ON DELETE RESTRICT ON UPDATE RESTRICT;