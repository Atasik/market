DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS products_users;
DROP TABLE IF EXISTS products_orders;
DROP TABLE IF EXISTS orders_users;

CREATE TABLE users 
(
  id        serial       not null unique,
  user_mode varchar(255) not null,
  username  varchar(255) not null unique,
  password  varchar(255) not null
);
INSERT INTO users (user_mode, username, password) VALUES
('admin',	'admin',	'admin');

CREATE TABLE products
(
  id            serial       not null unique,
  title         varchar(255) not null,
  price         numeric      not null,
  tag           varchar(255) not null,
  type          varchar(255) not null,
  description   varchar(255) not null,
  count         int          not null,
  creation_date date         not null,
  views         int          not null,
  image_url     varchar(255) not null
);

CREATE TABLE orders
(
  id            serial not null unique,
  creation_date date   not null,
  delivery_date date   not null
);

CREATE TABLE reviews
(
  id              serial                                         not null unique,
  creation_date   date                                           not null,
  product_id      int references products (id) on delete cascade not null,
  user_id         int references users    (id)                   not null,
  username        varchar(255)                                   not null,
  review_text     varchar(255)                                   not null,
  rating          int                                            not null
);

CREATE TABLE products_users
(
  id              serial                                         not null unique,
  product_id      int references products (id) on delete cascade not null,
  user_id         int references users (id) on delete cascade    not null,
  purchased_count int                                            not null
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

