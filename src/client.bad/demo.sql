CREATE DATABASE demo;
USE demo;

CREATE TABLE IF NOT EXISTS user (
    id INT NOT NULL AUTO_INCREMENT COMMENT "自增ID",
    username CHAR(32) NOT NULL COMMENT "用户名",
    password CHAR(32) NOT NULL COMMENT "用户密码",
    email CHAR(32) NOT NULL,
    phone CHAR(16) NOT NULL,
    description VARCHAR(128) NOT NULL,
    PRIMARY KEY(id)
) ENGINE=InnoDB CHARSET=utf8mb4;

INSERT INTO user(username, password, email, phone) VALUES
    ("one", "one", "one@shylinux.com", "110"),
    ("two", "two", "two@shylinux.com", "120");
