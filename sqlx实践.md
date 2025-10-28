# Go + sqlx 实践与最佳方案

## 1. 为什么要处理 GROUP_CONCAT 错位问题

在使用 `GROUP_CONCAT()` 聚合多个关联表字段（例如分类ID、名称、图标）时，
如果部分字段为 `NULL`，默认行为会导致结果错位：

``` json
"CategoryIDs": "161,162,163",
"CategoryIcons": ",S6y2gGfrxy,tqy25AwOD3"
```

可以看到 ID 有三个，但 ICON 只有两个。这样会造成解析错误或错位。

## 2. 解决方案：使用 `IFNULL` 或 `COALESCE` 统一默认值

通过在 SQL 中为可能为空的字段加默认值，确保每列都返回相同数量的元素：

``` sql
SELECT 
    b.*,
    a.id AS "author.id",
    a.author_name AS "author.author_name",
    p.id AS "publisher.id",
    p.publisher_name AS "publisher.publisher_name",
    GROUP_CONCAT(DISTINCT IFNULL(c.id, '')) AS category_ids,
    GROUP_CONCAT(DISTINCT IFNULL(c.category_name, '')) AS category_names,
    GROUP_CONCAT(DISTINCT IFNULL(c.icon, '')) AS category_icons
FROM books b
LEFT JOIN author a ON a.id = b.author_id
LEFT JOIN publisher p ON p.id = b.publisher_id
LEFT JOIN book_categories bc ON b.id = bc.book_id
LEFT JOIN category c ON bc.category_id = c.id
WHERE b.id = ? AND b.deleted_at IS NULL
GROUP BY b.id;
```

## 3. 更优雅的方式：使用 JSON_OBJECT + JSON_ARRAYAGG

如果数据库版本为 MySQL 5.7+，推荐使用 JSON 聚合函数，避免字符串拆分：

``` sql
SELECT 
    b.*,
    a.id AS "author.id",
    a.author_name AS "author.author_name",
    p.id AS "publisher.id",
    p.publisher_name AS "publisher.publisher_name",
    JSON_ARRAYAGG(
        JSON_OBJECT(
            'id', c.id,
            'name', c.category_name,
            'sort', c.sort,
            'icon', c.icon
        )
    ) AS categories
FROM books b
LEFT JOIN author a ON a.id = b.author_id
LEFT JOIN publisher p ON p.id = b.publisher_id
LEFT JOIN book_categories bc ON b.id = bc.book_id
LEFT JOIN category c ON bc.category_id = c.id
WHERE b.id = ? AND b.deleted_at IS NULL
GROUP BY b.id;
```

这种方式可以直接生成 JSON 数组，不需要在 Go 端拆分字符串。

## 4. Go 模型与反序列化

``` go
type Category struct {
    ID    uint64 `json:"id"`
    Name  string `json:"name"`
    Sort  int    `json:"sort"`
    Icon  string `json:"icon"`
}

type Book struct {
    ID         uint64      `db:"id" json:"id"`
    Title      string      `db:"title" json:"title"`
    Author     *Author     `json:"author"`
    Publisher  *Publisher  `json:"publisher"`
    Categories []*Category `db:"categories" json:"categories"`
}
```

查询后，直接用 `json.Unmarshal` 将 JSON 字段转为结构体：

``` go
var book Book
if err := db.Get(&book, query, id); err != nil {
    return nil, err
}

if err := json.Unmarshal([]byte(book.CategoriesRaw), &book.Categories); err != nil {
    log.Println("解析分类 JSON 出错：", err)
}
```

## 5. 性能建议

-   ✅ 如果数据量大，请限制 `GROUP_CONCAT` 长度：

    ``` sql
    SET SESSION group_concat_max_len = 8192;
    ```

-   ✅ JSON 聚合在 MySQL 8 性能最佳；5.7 虽可用但效率略低。\

-   ✅ 不推荐在 Go 端进行字符串分割，容易错位且难维护。\

-   ✅ sqlx 与命名参数搭配 `Rebind()` 是生产级方案，无性能隐患。

------------------------------------------------------------------------

**总结：**\
- 避免 `GROUP_CONCAT` 错位 → `IFNULL` 兜底。\
- 推荐 `JSON_ARRAYAGG + JSON_OBJECT` → 一步到位 JSON。\
- Go 端解析简单、安全、健壮。
