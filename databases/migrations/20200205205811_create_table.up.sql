-- Table Definition ----------------------------------------------
CREATE TABLE "users" (
  "id" SERIAL PRIMARY KEY NOT NULL,
  "roleId" int NOT NULL,
  "name" varchar(80) NOT NULL,
  "phone" varchar(80) NOT NULL,
  "password" varchar NOT NULL,
  "token" varchar,
  "tokenExpiredAt" timestamp,
  "createdAt" timestamp NOT NULL DEFAULT (now()),
  "createdBy" varchar(20) DEFAULT 'admin',
  "updatedAt" timestamp NOT NULL DEFAULT (now()),
  "updatedBy" varchar(20) DEFAULT 'admin',
  "deletedAt" timestamp,
  "deletedBy" varchar(20)
);

CREATE TABLE "roles" (
  "id" SERIAL PRIMARY KEY NOT NULL,
  "name" varchar(80) NOT NULL,
  "createdAt" timestamp NOT NULL DEFAULT (now()),
  "createdBy" varchar(20) DEFAULT 'admin',
  "updatedAt" timestamp NOT NULL DEFAULT (now()),
  "updatedBy" varchar(20) DEFAULT 'admin',
  "deletedAt" timestamp,
  "deletedBy" varchar(20)
);

CREATE TABLE "characters" (
  "id" SERIAL PRIMARY KEY NOT NULL,
  "characterTypeID" int NOT NULL,
  "name" varchar(80) NOT NULL,
  "power" int NOT NULL,
  "createdAt" timestamp NOT NULL DEFAULT (now()),
  "createdBy" varchar(20) DEFAULT 'admin',
  "updatedAt" timestamp NOT NULL DEFAULT (now()),
  "updatedBy" varchar(20) DEFAULT 'admin',
  "deletedAt" timestamp,
  "deletedBy" varchar(20)
);

CREATE TABLE "charactersType" (
  "id" SERIAL PRIMARY KEY NOT NULL,
  "name" varchar(80) NOT NULL,
  "code" int NOT NULL,
  "createdAt" timestamp NOT NULL DEFAULT (now()),
  "createdBy" varchar(20) DEFAULT 'admin',
  "updatedAt" timestamp NOT NULL DEFAULT (now()),
  "updatedBy" varchar(20) DEFAULT 'admin',
  "deletedAt" timestamp,
  "deletedBy" varchar(20)
);

ALTER TABLE "users" ADD FOREIGN KEY ("roleId") REFERENCES "roles" ("id");
ALTER TABLE "characters" ADD FOREIGN KEY ("characterTypeID") REFERENCES "charactersType" ("id");
