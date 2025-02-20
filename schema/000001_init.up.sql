CREATE TABLE users
(
    id serial NOT NULL UNIQUE,
    name varchar(255) not null,
    username varchar(255) not null unique,
    password_hash varchar(255) not null
);

CREATE TABLE todo_lists
(
    id serial NOT NULL UNIQUE,
    title varchar(255) NOT NULL,
    description varchar(255)
);

CREATE TABLE users_lists
(
    id serial NOT NULL UNIQUE,
    user_id int REFERENCES users (id) on delete CASCADE NOT NULL,
    list_id int REFERENCES todo_lists (id) on delete CASCADE NOT NULL
);

CREATE TABLE todo_items
(
    id serial NOT NULL UNIQUE,
    title varchar(255) NOT NULL,
    description varchar(255),
    done BOOLEAN DEFAULT false
);

CREATE TABLE lists_items
(
    id serial NOT NULL UNIQUE,
    item_id int REFERENCES todo_items (id) on delete CASCADE NOT NULL,
    list_id int REFERENCES todo_lists (id) on delete CASCADE NOT NULL
);