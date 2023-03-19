
-- drop tables
DROP TABLE IF EXISTS "outatime"."subscription";
DROP TABLE IF EXISTS "outatime"."subscriber";
DROP TABLE IF EXISTS "outatime"."event";
DROP TABLE IF EXISTS "outatime"."user_page";
DROP TABLE IF EXISTS "outatime"."link";
DROP TABLE IF EXISTS "outatime"."user";

-- create schema
DROP SCHEMA IF EXISTS "outatime" CASCADE;
CREATE SCHEMA IF NOT EXISTS "outatime";

-- create tables

-- -----------------------------------------------
-- USER TABLE
-- -----------------------------------------------

CREATE TABLE IF NOT EXISTS "outatime"."user" (
    "user_id"       SERIAL          PRIMARY KEY,
    "user_name"     VARCHAR(50)     NOT NULL,
    "user_url"      VARCHAR(30)     NOT NULL    UNIQUE,
    "user_email"    VARCHAR(100)    NOT NULL    UNIQUE,
    "user_password" VARCHAR(100)    NOT NULL,
    "auth_code"     VARCHAR(100)    UNIQUE,
    "user_avatar"   VARCHAR(1024)   NOT NULL,
    "validate_code" VARCHAR(100)    NOT NULL    UNIQUE,
    "isValid"       BOOLEAN         NOT NULL    DEFAULT FALSE,
    "isActive"      BOOLEAN         NOT NULL    DEFAULT FALSE,
    "isAdmin"       BOOLEAN         NOT NULL    DEFAULT FALSE,
    "strikes"       INT             NOT NULL    DEFAULT 0,
    "subscribers"   INT             NOT NULL    DEFAULT 0
);



-- -----------------------------------------------
-- LINK TABLE
-- -----------------------------------------------

CREATE TABLE IF NOT EXISTS "outatime"."link" (
    "link_id"       SERIAL          PRIMARY KEY,
    "user_id"       INT             NOT NULL,
    FOREIGN KEY     ("user_id")     REFERENCES "outatime"."user"("user_id"),
    "link_url"      VARCHAR(1024)   NOT NULL,
    "type"          VARCHAR(100)    NOT NULL
);



-- -----------------------------------------------
-- USER_PAGE TABLE
-- -----------------------------------------------

CREATE TABLE IF NOT EXISTS "outatime"."user_page" (
    "page_id"       SERIAL          PRIMARY KEY,
    "user_id"       INT             NOT NULL,
    FOREIGN KEY     ("user_id")     REFERENCES "outatime"."user"("user_id"),
    "aboutUs"       VARCHAR(5000),
    "banner"        VARCHAR(1024),
    "isPublic"      BOOLEAN         NOT NULL    DEFAULT FALSE
);



-- -----------------------------------------------
-- EVENT TABLE
-- -----------------------------------------------

CREATE TABLE IF NOT EXISTS "outatime"."event" (
    "event_id"      SERIAL          PRIMARY KEY,
    "user_id"       INT             NOT NULL,
    FOREIGN KEY     ("user_id")     REFERENCES "outatime"."user"("user_id"),
    "event_name"    VARCHAR(50)     NOT NULL,
    "isPublic"      BOOLEAN         NOT NULL    DEFAULT FALSE,
    "event_tldr"    VARCHAR(200)    NOT NULL,
    "event_descr"   VARCHAR(2048)   NOT NULL,
    "event_start"   TIMESTAMP       NOT NULL,
    "event_end"     TIMESTAMP       NOT NULL,
    "event_location"VARCHAR(100)    NOT NULL,
    "event_img"     VARCHAR(1024),
    "event_code"    VARCHAR(100)    NOT NULL    UNIQUE,
    "subscribers"   INT             NOT NULL    DEFAULT 0
);



-- -----------------------------------------------
-- SUBSCRIBER TABLE
-- -----------------------------------------------

CREATE TABLE IF NOT EXISTS "outatime"."subscriber" (
    "sub_id"        SERIAL          PRIMARY KEY,
    "sub_code"      VARCHAR(50)     NOT NULL    UNIQUE,
    "email"         VARCHAR(100)    NOT NULL    UNIQUE,
    "validate_code" VARCHAR(100)    NOT NULL    UNIQUE,
    "isValid"       BOOLEAN         NOT NULL    DEFAULT FALSE
);



-- -----------------------------------------------
-- SUBSCRIPTION TABLE
-- -----------------------------------------------

CREATE TABLE IF NOT EXISTS "outatime"."subscription" (
    "subscription_id"   SERIAL          PRIMARY KEY,
    "sub_id"            INT             NOT NULL,
    "user_id"           INT,
    "event_id"          INT,
    FOREIGN KEY     ("sub_id")     REFERENCES "outatime"."subscriber"("sub_id"),
    FOREIGN KEY     ("user_id")     REFERENCES "outatime"."user"("user_id"),
    FOREIGN KEY     ("event_id")     REFERENCES "outatime"."event"("event_id")
);