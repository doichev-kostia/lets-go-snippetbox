-- driver SQLite

create table "snippets" (
    "id" text primary key,
    "title" text not null,
    "content" text not null,
    "create_time" timestamp not null default current_timestamp,
    "expire_time" timestamp not null
);

create index "idx_snippets_create_time" on "snippets" ("create_time");

-- For the github.com/alexedwards/scs/v2
create table "sessions" (
    "token"  text primary key,
    "data"   BLOB NOT NULL,
    "expiry" REAL NOT NULL
);

create index "idx_session_expiry" on "sessions" ("expiry");

create table "users" (
    "id" text primary key,
    "name" text not null unique,
    "email" text not null,
    "hashed_password" text not null,
    "create_time" timestamp not null default current_timestamp
);

create index "idx_users_email" on "users" ("email");
