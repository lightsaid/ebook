CREATE TABLE IF NOT EXISTS `shopping_carts` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户id',
  `book_id` BIGINT UNSIGNED NOT NULL COMMENT '图书id',
  `quantity` INT UNSIGNED NOT NULL DEFAULT 1 COMMENT '数量',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_user_book` (`user_id`, `book_id`),
  INDEX `idx_user_id` (`user_id`),
  INDEX `idx_book_id` (`book_id`),
  INDEX `idx_created_at` (`created_at`),
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
