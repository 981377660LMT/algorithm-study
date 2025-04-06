# Modern SQL

## 1. 现代 SQL 的背景

- **历史回顾**：传统 SQL（如 SQL-92 标准）主要关注基本的 SELECT、INSERT、UPDATE、DELETE 等操作，及简单的连接和聚合。但随着应用场景越来越复杂，对数据处理的要求不断提高，SQL 标准不断演进，引入了大量新特性。
- **现代 SQL**：指的是 SQL 语言在新标准（如 SQL:1999、SQL:2003、SQL:2008、SQL:2011 及后续版本）中新增的特性，这些特性使 SQL 更加灵活、强大，并能更好地支持大数据和复杂查询需求。

---

## 2. 现代 SQL 的主要特性

### 2.1 窗口函数（Window Functions）

- **概念**：窗口函数允许在不改变结果集行数的情况下，对数据进行排名、累计求和、移动平均等计算。
- **用法**：使用 `OVER()` 子句指定窗口范围。  
  例如，计算某字段的累计和：
  ```sql
  SELECT
      id,
      value,
      SUM(value) OVER (ORDER BY id) AS cumulative_sum
  FROM my_table;
  ```
- **优势**：无需 GROUP BY 即可在保持原有行的同时完成聚合计算，特别适用于需要计算排名、百分比、移动平均等场景。

### 2.2 公共表表达式（CTE，Common Table Expressions）

- **概念**：CTE 通过 `WITH` 子句定义一个临时命名结果集，可在后续查询中重复引用。
- **递归 CTE**：特别适合处理层次结构数据（如组织结构、图数据），通过递归查询实现层级遍历。
- **示例**：
  ```sql
  WITH RECURSIVE Subordinates AS (
      SELECT employee_id, manager_id, name
      FROM employees
      WHERE manager_id IS NULL
      UNION ALL
      SELECT e.employee_id, e.manager_id, e.name
      FROM employees e
      JOIN Subordinates s ON e.manager_id = s.employee_id
  )
  SELECT * FROM Subordinates;
  ```
- **优势**：提高查询的可读性和可维护性，使复杂查询更易于理解和调试。

### 2.3 JSON 与 XML 支持

- **背景**：随着半结构化数据的流行，传统的关系模型面临存储和查询 JSON、XML 等数据的挑战。
- **现代 SQL 特性**：新增 JSON 数据类型和操作函数，使得数据库能直接存储、解析和查询 JSON 数据。例如，MySQL 和 PostgreSQL 都支持 JSONB 类型和相关操作。
- **优势**：允许关系数据库兼容处理半结构化数据，提供了更大的灵活性。

### 2.4 MERGE 语句（UPSERT）

- **功能**：MERGE 语句可以根据匹配条件，自动判断是执行插入、更新还是删除操作。
- **示例**：
  ```sql
  MERGE INTO target_table AS t
  USING source_table AS s
  ON t.id = s.id
  WHEN MATCHED THEN
      UPDATE SET t.value = s.value
  WHEN NOT MATCHED THEN
      INSERT (id, value) VALUES (s.id, s.value);
  ```
- **优势**：大大简化了数据同步、更新合并等操作的逻辑，减少了冗长的条件判断和多条语句执行。

### 2.5 扩展的数据类型与模式匹配

- **扩展数据类型**：现代 SQL 支持更多数据类型，如地理空间数据（GIS）、XML、数组等。
- **模式匹配**：SQL:2016 及后续标准引入了更强大的模式匹配能力，如行模式识别，用于处理复杂的序列模式匹配问题。

---

## 3. 现代 SQL 的意义

- **增强表达力**：现代 SQL 新增的功能，如窗口函数、递归 CTE 等，大大增强了 SQL 表达复杂业务逻辑的能力，使得写出高效、可维护的查询变得更容易。
- **兼容半结构化数据**：支持 JSON、XML 等数据类型，使传统关系数据库能处理更多类型的数据场景。
- **简化操作**：MERGE 语句等特性简化了业务逻辑，使得开发者无需编写大量条件分支，直接通过一条语句完成复杂操作。
- **优化查询性能**：一些新特性在底层也能帮助优化查询性能，例如窗口函数能避免过度分组，提高计算效率。

---

## 4. 总结

CMU-15-445 的“02 Modern SQL”讲解主要涵盖了现代 SQL 新标准中引入的各种扩展特性，帮助学生理解 SQL 语言是如何从基础的关系查询发展到支持复杂分析和半结构化数据处理的。主要内容包括：

- **窗口函数**：实现分组、排名和累计计算等功能。
- **CTE 和递归查询**：支持复杂层次数据的查询和处理。
- **JSON/XML 支持**：兼容半结构化数据。
- **MERGE/UPSERT**：简化数据更新逻辑。
- **扩展数据类型与模式匹配**：提高数据存储与查询的灵活性。

---

- sql (Structured Query Language) 是一种用于管理关系数据库的标准语言，广泛应用于数据查询、更新和管理。
  Data Manipulation Language (DML) -> 怎删改查
  Data Definition Language (DDL) -> 怎建表、操作元数据
  Data Control Language (DCL) -> 权限控制

  `Important: SQL is based on bags (duplicates) not sets (no duplicates).`

- Aggregations + Group By  
  利用聚合函数（如 COUNT、SUM、AVG、MIN、MAX）对表中的数据进行求和、计数、平均等运算，使用 GROUP BY 将数据按照一个或多个字段分组，再分别对各组应用聚合计算，从而按组统计汇总数据。

- String / Date / Time Operations  
  对字符串、日期和时间进行操作，包括拼接、截取、格式化、日期差值计算等。例如，可以用 SUBSTRING、CONCAT 操作字符串，用 DATEADD、DATEDIFF 处理日期与时间，为数据格式和数据清洗提供支持。

- Output Control + Redirection  
  控制查询结果的输出格式和位置，可以通过指定排序、限制返回行数，以及将查询结果重定向输出到文件或其他介质，实现结果展示的定制化和数据传输。

- Nested Queries  
  嵌套查询指在一个 SQL 查询中嵌入另一个查询（子查询），常用于数据过滤、数据匹配等。子查询可以出现在 SELECT、FROM 或 WHERE 子句中，增强查询表达能力和灵活性。

- Common Table Expressions (CTE，内部临时视图)
  **通用表表达式（CTE）允许将一个复杂查询的部分结果暂时命名**，使后续查询引用，提高可读性和维护性。CTE 可以简化嵌套查询的层级结构，也支持递归查询场景。

- Window Functions  
  窗口函数在 SQL 中对一组相关行（窗口）执行计算，如排名、累计和移动平均等，而不需要将行汇聚成单个结果。它们`使用 OVER 子句定义窗口`，提供高效数据分析能力。
