DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS baskets;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS products_baskets;
DROP TABLE IF EXISTS products_users;
DROP TABLE IF EXISTS products_orders;
DROP TABLE IF EXISTS orders_users;

CREATE TABLE users 
(
  id        serial       not null unique,
  role      varchar(255) not null,
  username  varchar(255) not null unique,
  password  varchar(255) not null
);
INSERT INTO users (role, username, password) VALUES
('admin',	'admin',	'$argon2id$v=19$m=65536,t=3,p=1$kMwiCJlyCi2xXKy/U1c8hA$FtPqNnpdWNc7cD0hOcbTxMav4s/HyGUhew6bhlWqy5c');

CREATE TABLE baskets
(
  id      serial                                     not null unique,
  user_id int references users(id) on delete cascade not null unique
);
INSERT INTO baskets (user_id) VALUES
(1);

CREATE TABLE products
(
  id            serial       not null unique,
  title         varchar(255) not null,
  price         numeric      not null,
  tag           varchar(255) not null,
  type          varchar(255) not null,
  description   varchar(255) not null,
  count         int          not null,
  creation_date timestamp    not null,
  views         int          not null,
  image_url     varchar(255) not null,
  image_id      varchar(255) not null unique
);

CREATE TABLE orders
(
  id            serial      not null unique,
  creation_date timestamp   not null,
  delivery_date timestamp   not null
);

CREATE TABLE reviews
(
  id              serial                                         not null unique,
  creation_date   timestamp                                      not null,
  product_id      int references products (id) on delete cascade not null,
  user_id         int references users    (id)                   not null,
  username        varchar(255)                                   not null,
  review_text     varchar(255)                                   not null,
  rating          int                                            not null
);

CREATE TABLE products_baskets
(
  id              serial                                         not null unique,
  product_id      int references products (id) on delete cascade not null,
  basket_id       int references baskets (id) on delete cascade  not null,
  purchased_count int                                            not null
);

CREATE TABLE products_users 
(
  id              serial                                         not null unique,
  product_id      int references products (id) on delete cascade not null,
  user_id         int references users (id) on delete cascade    not null
);

CREATE TABLE products_orders
(
  id         serial                                         not null unique,
  product_id int references products (id) on delete cascade not null,
  order_id   int references orders (id) on delete cascade   not null
);

CREATE TABLE orders_users
(
  id        serial                                        not null unique,
  order_id  int references orders (id) on delete cascade  not null,
  user_id   int references users (id) on delete cascade   not null
);

