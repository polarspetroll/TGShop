CREATE DATABASE TGShop;
USE TGShop;

-- products table
CREATE TABLE products (
  id INT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(100),
  price VARCHAR(10),
  stat BOOLEAN,
  filename VARCHAR(30) UNIQUE
);


-- login table
CREATE TABLE login (
  username VARCHAR(30) UNIQUE,
  password VARCHAR(64)
);
