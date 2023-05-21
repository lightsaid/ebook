CREATE DATABASE db_ebook; 
USE db_ebook;

CREATE TABLE book (
    id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
    isbn VARCHAR(13) NOT NULL COMMENT '国际标准书号',
    title VARCHAR(255) NOT NULL COMMENT '图书名字',
    poster  VARCHAR(255) NOT NULL COMMENT '图书封面图地址',
    pages INT UNSIGNED NOT NULL COMMENT '总页数',
    price DECIMAL(6, 2) UNSIGNED COMMENT '图书单价',
    published_at DATE NOT NULL COMMENT '发售日期',
    created_at TIMESTAMP NOT NULL default NOW() COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL default NOW() COMMENT '更新时间',
    unique index unq_isbn(`isbn`)
) ENGINE=InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT '图书表';
