``` mysql
-- 查询某个库所有表
SELECT
    *
FROM
    information_schema. TABLES
WHERE
    table_schema = 'cc_center';

-- 查询指定库所有列
SELECT
    *
FROM
    information_schema. COLUMNS
WHERE
    table_schema = 'cc_center';

-- 查询所有表的注释和字段注释
SELECT
    a.table_name 表名,
    a.table_comment 表说明,
    b.COLUMN_NAME 字段名,
    b.column_comment 字段说明,
    b.column_type 字段类型,
    b.column_key 约束
FROM
    information_schema. TABLES a
LEFT JOIN information_schema. COLUMNS b ON a.table_name = b.TABLE_NAME
WHERE
    a.table_schema = 'cc_center'
ORDER BY
    a.table_name
```