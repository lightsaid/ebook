# SQL 迁移文件
MIGRATE_SQL=./migrations

# MySQL链接
DB_SOURCE=mysql://root:root.cc@tcp(localhost:3306)/db_ebook

seed:
	go run ./cmd/seed/*.go

## migrate/new name=$1: 创建迁移SQL文件，如 make migrate/new name=create_user_table
migrate/new:
	migrate create -seq -ext=.sql -dir=${MIGRATE_SQL} ${name}

## migrate/up: 向上迁移所有
migrate/up:
	migrate -path=${MIGRATE_SQL} -database="${DB_SOURCE}" -verbose up

## migrate/down: 向下迁移所有
migrate/down:
	migrate -path=${MIGRATE_SQL} -database="${DB_SOURCE}" -verbose down

## migrate/up1: 向上迁移一次
migrate/up1: 
	migrate -path=${MIGRATE_SQL} -database="${DB_SOURCE}" -verbose up 1

## migrate/down1: 向下迁移一次
migrate/down1:
	migrate -path=${MIGRATE_SQL} -database="${DB_SOURCE}" -verbose down 1

## migrate/force version=$1: 强制迁移到指定版本
migrate/force:
	migrate -path=${MIGRATE_SQL} -database="${DB_SOURCE}" -verbose force ${version}

.PHONY: migrate/new migrate/up migrate/down migrate/up1 migrate/down1 migrate/force