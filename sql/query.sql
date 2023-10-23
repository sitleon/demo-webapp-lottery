-- name: CreateDraw :exec
INSERT INTO "public"."draw" ("draw_id", "drew_at") VALUES ($1, $2);

-- name: CreateTicket :exec
INSERT INTO "public"."ticket" ("ticket_id", "draw_id") VALUES ($1, $2);

-- name: CountTicket :one
SELECT count(*) FROM "public"."ticket" WHERE draw_id = $1;

-- name: GetTicket :one
SELECT
    "ticket"."ticket_id" AS "ticket_id",
    "ticket"."draw_id" AS "draw_id",
    CAST("draw"."winner_ticket" IS NOT NULL AS BOOLEAN) AS "is_draw",
    CAST("draw"."winner_ticket" = "ticket"."ticket_id" AS BOOLEAN) AS "is_winner"
FROM "public"."ticket"
LEFT JOIN "public"."draw" ON "draw"."draw_id" = "ticket"."draw_id"
WHERE "ticket_id" = $1;

-- name: GetDraw :one
SELECT "draw_id", "winner_ticket", "drew_at" FROM "public"."draw" WHERE "draw_id" = $1;

-- name: SetDrawWinner :one
UPDATE "public"."draw"
	SET "winner_ticket" = "winner"."ticket_id",
		"update_at" = NOW()
	FROM (
        SELECT "ticket_id" FROM "public"."ticket"
        	WHERE "draw_id" = $1 OFFSET $2 LIMIT 1
    ) AS "winner"
WHERE "draw"."draw_id" = $1 AND "winner_ticket" IS NULL
RETURNING "winner_ticket";
