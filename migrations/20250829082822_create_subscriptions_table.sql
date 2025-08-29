-- +goose Up
-- +goose StatementBegin
CREATE TABLE "subscriptions"(
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "user_id" UUID NOT NULL,
    "service_name" VARCHAR(255) NOT NULL,
    "price" INT NOT NULL,
    "start_date" DATE NOT NULL,
    "end_date" DATE NULL
);
CREATE INDEX "subscriptions_user_id_index" ON "subscriptions"("user_id");
CREATE INDEX "subscriptions_service_name_index" ON "subscriptions"("service_name");
CREATE INDEX "subscriptions_start_date_index" ON "subscriptions"("start_date");

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS subscriptions;
DROP INDEX IF EXISTS subscriptions_user_id_index;
DROP INDEX IF EXISTS subscriptions_service_name_index;
DROP INDEX IF EXISTS subscriptions_start_date_index;
-- +goose StatementEnd