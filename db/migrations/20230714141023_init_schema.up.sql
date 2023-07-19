CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "full_name" varchar NOT NULL,
  "hashed_password" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "level" boolean NOT NULL DEFAULT false,
  "password_changed_at" timestamptz NOT NULL DEFAULT (now()),
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "bank_accounts" (
  "username" varchar NOT NULL,
  "card_number" varchar UNIQUE NOT NULL,
  "balance" serial NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "update_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "products" (
  "id_product" serial UNIQUE PRIMARY KEY,
  "product_name" varchar NOT NULL DEFAULT '',
  "kind_of_product" varchar NOT NULL DEFAULT '',
  "owner" varchar NOT NULL,
  "currency" varchar NOT NULL,
  "price" serial NOT NULL,
  "quantity" serial NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "update_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "purchase_history" (
  "id_purchase_history" serial PRIMARY KEY,
  "id_product" serial NOT NULL,
  "buyer" varchar NOT NULL,
  "card_number_of_buyer" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "update_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "users" ("username");

CREATE INDEX ON "bank_accounts" ("username");

CREATE INDEX ON "bank_accounts" ("card_number");

CREATE UNIQUE INDEX ON "bank_accounts" ("username", "currency");

CREATE INDEX ON "products" ("id_product");

CREATE INDEX ON "products" ("owner");

CREATE INDEX ON "purchase_history" ("id_purchase_history");

CREATE INDEX ON "purchase_history" ("id_product");

CREATE INDEX ON "purchase_history" ("buyer");

ALTER TABLE "bank_accounts" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "products" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

ALTER TABLE "purchase_history" ADD FOREIGN KEY ("id_product") REFERENCES "products" ("id_product");

ALTER TABLE "purchase_history" ADD FOREIGN KEY ("buyer") REFERENCES "users" ("username");

ALTER TABLE "purchase_history" ADD FOREIGN KEY ("card_number_of_buyer") REFERENCES "bank_accounts" ("card_number");
