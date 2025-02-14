// golang sql 标准库
// http://go-database-sql.org/index.html
//
// - Overview
//   sql.DB 抽象旨在让您无需担心如何管理对底层数据存储的并发访问。
//   创建 sql.DB 后，您可以使用它查询所表示的数据库，以及创建语句和事务。
//
// - Importing a Database Driver
//   !标准库提供接口，驱动包具体实现，匿名导入使得驱动自动注册。
//   _ "github.com/go-sql-driver/mysql"
//   只引入mysql驱动，但是不使用它的API，只是为了让它注册自己(sql.Register(name, driver))。
//   这样在代码中使用sql.Open("xxx", ...)时，就能找到对应的驱动实现。
//
// - Accessing the Database
//   sql.Open 是懒惰的，不会立即连接数据库。
//   不要频繁 Open() 和 Close() 数据库。相反，为每个需要访问的不同数据存储创建一个 sql.DB 对象。
//
// - Retrieving Result Sets
//   Query() 返回一个 Rows 对象，它是一个迭代器，可以使用 Next() 和 Scan() 方法遍历结果集。
//   QueryRow() 返回一个 Row 对象，它是一个单行结果。
//   Exec() 用于执行不返回结果集的查询，如 INSERT、UPDATE 和 DELETE。
//   Prepare() 用于准备一个语句，然后可以多次执行它。
//
// - Modifying Data and Using Transactions
//   使用 Exec() ，最好配合预处理语句，以完成 INSERT 、 UPDATE 、 DELETE 或其他不返回行的语句。
//   !执行该语句会生成一个 sql.Result ，该对象提供对语句元数据的访问：最后插入的 ID 和受影响的行数。
//   不要使用Query，因为它会返回一个Rows对象资源，但是不会关闭它，可能会导致连接泄漏。
//
// - Using Prepared Statements
//
// - Handling Errors
//
// - Working with Nulls
//
// - Working with Unknown Columns
//
// - The Connection Pool
//
// - Suprises,Antipatterns and Limitations
//
// - Related Reading and Resources

package main

import (
	"database/sql"
	"log"
	// _ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/hello")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	{
		var (
			id   int
			name string
		)
		rows, err := db.Query("select id, name from users where id = ?", 1)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&id, &name)
			if err != nil {
				log.Fatal(err)
			}
			log.Println(id, name)
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
	}
}
