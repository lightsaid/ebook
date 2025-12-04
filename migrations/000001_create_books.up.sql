
CREATE TABLE IF NOT EXISTS `author` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `author_name` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '作者名称',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` TIMESTAMP NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  INDEX `idx_created_at` (`created_at`),
  INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `publisher` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `publisher_name` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '出版社名称',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` TIMESTAMP NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  INDEX `idx_created_at` (`created_at`),
  INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE IF NOT EXISTS `books` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `isbn` VARCHAR(16) NOT NULL COMMENT 'ISBN号',
  `title` VARCHAR(120) NOT NULL DEFAULT '' COMMENT '书名',
  `subtitle` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '副标题',
  `author_id` BIGINT UNSIGNED NOT NULL COMMENT '作者id',
  `cover_url` VARCHAR(255) NOT NULL COMMENT '封面图片地址',
  `publisher_id` BIGINT UNSIGNED NOT NULL COMMENT '出版社id',
  `pubdate` DATE NOT NULL COMMENT '出版日期',
  `price` INT UNSIGNED NOT null DEFAULT 0 COMMENT '价格,单位分',
  `status` TINYINT NOT NULL DEFAULT 0 COMMENT '0-下架,1-上架',
  `type` TINYINT NOT NULL DEFAULT 1 COMMENT '1-电子书,2-实体,3-电子书+实体',
  `stock` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '库存',
  `source_url` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '电子书资源路径',
  `description` TEXT NOT NULL COMMENT '图书描述',
  `version` INT NOT NULL DEFAULT 1 COMMENT '版本号',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` TIMESTAMP NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  INDEX `idx_title` (`title`),
  INDEX `idx_created_at` (`created_at`),
  INDEX `idx_deleted_at` (`deleted_at`),
  UNIQUE INDEX `idx_isbn` (`isbn`),

  FOREIGN KEY (`author_id`) REFERENCES author(id) ON DELETE RESTRICT,
  FOREIGN KEY (`publisher_id`) REFERENCES publisher(id) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;



