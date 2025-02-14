// golang sql 标准库
// http://go-database-sql.org/index.html
//
// - Overview
//   sql.DB 抽象旨在让您无需担心如何管理对底层数据存储的并发访问。
//   创建 sql.DB 后，您可以使用它查询所表示的数据库，以及创建语句和事务。
// - Importing a Database Driver
// - Accessing the Database
// - Retrieving Result Sets
// - Modifying Data and Using Transactions
// - Using Prepared Statements
// - Handling Errors
// - Working with Nulls
// - Working with Unknown Columns
// - The Connection Pool
// - Suprises,Antipatterns and Limitations
// - Related Reading and Resources

package main

import "database/sql"

func main() {
	sql.Open("mysql", "user:password@/dbname")
	sql.Named()
}
