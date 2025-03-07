CREATE TABLE if not exists users (
   id bigserial PRIMARY KEY ,
   firstname text,
   lastname text,
   username text,
   email text unique not null,
   phone text,
   gender text,
   password TEXT NOT NULL,
   role text,
   created_at timestamp default now(),
   created_by bigint references users(id),
   updated_at timestamp,
   updated_by bigint references users(id),
   deleted_at timestamp,
   deleted_by bigint references users(id)
);

CREATE TABLE if not exists posts (
     id bigserial PRIMARY KEY ,
     user_id bigint REFERENCES users(id),
     title text NOT NULL,
     content TEXT NOT NULL,
     created_at timestamp default now(),
     created_by bigint references users(id),
     updated_at timestamp,
     updated_by bigint references users(id),
     deleted_at timestamp,
     deleted_by bigint references users(id)
);

CREATE TABLE if not exists logs (
     id bigserial PRIMARY KEY,
     level VARCHAR(10) CHECK (level IN ('INFO', 'WARNING', 'ERROR')),
     message TEXT NOT NULL,
     service_name text,
     created_at timestamp default now(),
     created_by bigint references users(id),
     updated_at timestamp,
     updated_by bigint references users(id),
     deleted_at timestamp,
     deleted_by bigint references users(id)
);

