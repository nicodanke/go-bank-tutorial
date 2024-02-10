-- SQL dump generated using DBML (dbml-lang.org)
-- Database: PostgreSQL
-- Generated at: 2024-02-10T00:55:44.478Z

CREATE TABLE "accounts" (
  "id" int8 NOT NULL DEFAULT (nextval('accounts_id_seq'::regclass)),
  "owner" varchar NOT NULL,
  "balance" int8 NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  PRIMARY KEY ("id") 
);

CREATE TABLE "entries" (
  "id" int8 NOT NULL DEFAULT (nextval('entries_id_seq'::regclass)),
  "account_id" int8 NOT NULL,
  "amount" int8 NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  PRIMARY KEY ("id")
);

CREATE TABLE "sessions" (
  "id" uuid NOT NULL,
  "username" varchar NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" bool NOT NULL DEFAULT false,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  PRIMARY KEY ("id")
);

CREATE TABLE "transfers" (
  "id" int8 NOT NULL DEFAULT (nextval('transfers_id_seq'::regclass)),
  "from_account_id" int8 NOT NULL,
  "to_account_id" int8 NOT NULL,
  "amount" int8 NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  PRIMARY KEY ("id")
);

CREATE TABLE "users" (
  "username" varchar NOT NULL,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT ('0001-01-01 00:00:00+00'::timestampwithtimezone),
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  PRIMARY KEY ("username")
);

COMMENT ON COLUMN "entries"."amount" IS 'Can be negative or positive';

COMMENT ON COLUMN "transfers"."amount" IS 'Only positive';

ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "sessions" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "transfers" ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");
