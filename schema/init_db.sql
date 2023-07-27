DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS carts;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS products_carts;
DROP TABLE IF EXISTS products_users;
DROP TABLE IF EXISTS products_orders;

CREATE TABLE users 
(
  id        serial       not null unique,
  role      varchar(255) not null,
  username  varchar(255) not null unique,
  password  varchar(255) not null
);
INSERT INTO users (role, username, password) VALUES
('admin',	'admin',	'$argon2id$v=19$m=65536,t=3,p=1$kMwiCJlyCi2xXKy/U1c8hA$FtPqNnpdWNc7cD0hOcbTxMav4s/HyGUhew6bhlWqy5c');

CREATE TABLE carts
(
  id      serial                                     not null unique,
  user_id int references users(id) on delete cascade not null unique
);
INSERT INTO carts (user_id) VALUES
(1);

CREATE TABLE products
(
  id            serial                                                        not null unique,
  user_id       int references users(id) on delete cascade                    not null,
  title         varchar(255)                                                  not null,
  price         numeric                                    check (price > 0)  not null, 
  tag           varchar(255), 
  category      varchar(255)                                                  not null,
  description   varchar(255), 
  amount        int                                        check (amount > 0) not null,
  created_at    timestamp                                                     not null DEFAULT (now() AT TIME ZONE 'utc'),
  updated_at    timestamp                                                     ,
  views         int                                                           not null,
  image_url     varchar(255)                                                  not null,
  image_id      varchar(255)                                                  not null unique
);

CREATE TABLE orders
(
  id            serial                                        not null unique,
  user_id       int references users (id) on delete cascade   not null,
  created_at    timestamp                                     not null DEFAULT (now() AT TIME ZONE 'utc'),
  delivered_at  timestamp                                     not null
);

CREATE TABLE reviews
(
  id              serial                                         not null unique,
  created_at      timestamp                                      not null DEFAULT (now() AT TIME ZONE 'utc'),
  updated_at      timestamp                                      ,
  product_id      int references products (id) on delete cascade not null,
  user_id         int references users (id)                      not null,
  text            varchar(255)                                   not null,
  category        varchar(255)                                   not null
);

CREATE TABLE products_carts
(
  id               serial                                                                      not null unique,
  product_id       int references products (id) on delete cascade                              not null,
  cart_id          int references carts (id) on delete cascade                                 not null,
  purchased_amount int                                            check (purchased_amount > 0) not null
);

CREATE TABLE products_users 
(
  id              serial                                         not null unique,
  product_id      int references products (id) on delete cascade not null,
  user_id         int references users (id) on delete cascade    not null
);

CREATE TABLE products_orders
(
  id               serial                                                                      not null unique,
  product_id       int references products (id) on delete cascade                              not null,
  order_id         int references orders (id) on delete cascade                                not null,
  purchased_amount int                                            check (purchased_amount > 0) not null
);

