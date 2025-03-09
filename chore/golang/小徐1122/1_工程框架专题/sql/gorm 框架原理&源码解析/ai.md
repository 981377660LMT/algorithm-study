# gorm 框架讲解

### 极其详细的GORM框架讲解

GORM是Go语言中最流行的ORM（对象关系映射）库之一，支持主流数据库（MySQL、PostgreSQL、SQLite、SQL Server等），提供链式API、事务、钩子、关联关系、迁移等特性。本文将从基础到高级全面讲解GORM，覆盖所有核心功能。

---

#### 一、安装与配置

1. **安装GORM及驱动**

   ```bash
   go get -u gorm.io/gorm
   # 根据数据库选择驱动，例如MySQL：
   go get -u gorm.io/driver/mysql
   ```

2. **连接数据库**

   ```go
   import (
     "gorm.io/driver/mysql"
     "gorm.io/gorm"
   )

   func main() {
     dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
     db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
     if err != nil {
       panic("数据库连接失败")
     }
   }
   ```

---

#### 二、模型定义与约定

1. **模型结构体**
   GORM使用结构体映射数据库表，默认规则：

   - 表名：结构体名的复数形式（如`User` → `users`）。
   - 主键：字段名为`ID`或标记`gorm:"primaryKey"`。
   - 时间戳：`CreatedAt`、`UpdatedAt`、`DeletedAt`（软删除）。

   ```go
   type User struct {
     ID        uint      `gorm:"primaryKey"`
     Name      string    `gorm:"size:255;not null"`
     Age       int       `gorm:"default:18"`
     CreatedAt time.Time
     UpdatedAt time.Time
     DeletedAt gorm.DeletedAt `gorm:"index"` // 软删除支持
   }
   ```

2. **自定义表名**
   实现`TableName()`方法或使用`Table()`方法：
   ```go
   func (User) TableName() string { return "profiles" }
   // 或者
   db.Table("profiles").Create(&user)
   ```

---

#### 三、自动迁移（AutoMigrate）

根据模型自动创建/更新表结构：

```go
db.AutoMigrate(&User{}, &Product{}) // 迁移多个模型
```

**注意**：AutoMigrate不会删除列或修改列类型，复杂迁移需手动处理。

---

#### 四、CRUD操作

1. **创建（Create）**

   ```go
   user := User{Name: "Jinzhu", Age: 18}
   result := db.Create(&user)       // 插入记录
   result.Error                     // 错误信息
   result.RowsAffected              // 影响行数

   // 批量插入
   users := []User{{Name: "A"}, {Name: "B"}}
   db.CreateInBatches(users, 100)  // 每批100条

   // 选择/忽略字段
   db.Select("Name", "Age").Create(&user)
   db.Omit("Age").Create(&user)
   ```

2. **查询（Read）**

   - **单条查询**
     ```go
     db.First(&user)          // 按主键升序第一条
     db.Last(&user)           // 按主键降序最后一条
     db.Take(&user)           // 无排序第一条
     db.First(&user, 10)      // 主键=10的记录
     db.First(&user, "name = ?", "jinzhu")
     ```
   - **条件查询**
     ```go
     db.Where("name = ?", "jinzhu").First(&user)
     db.Where(&User{Name: "jinzhu", Age: 20}).First(&user) // 结构体条件（忽略零值）
     db.Where(map[string]interface{}{"name": "jinzhu", "age": 0}).Find(&users) // map包含零值
     ```
   - **高级查询**
     ```go
     db.Not("name = ?", "jinzhu").Find(&users)
     db.Or("name = ?", "jinzhu").Or("age >= ?", 20).Find(&users)
     db.Order("age desc, name").Find(&users)
     db.Limit(10).Offset(5).Find(&users) // 分页
     db.Select("name", "age").Find(&users)
     db.Distinct("name").Find(&users)
     ```

3. **更新（Update）**

   ```go
   user.Name = "New Name"
   db.Save(&user) // 更新所有字段

   db.Model(&user).Update("name", "hello") // 更新单个字段
   db.Model(&user).Updates(User{Name: "hello", Age: 20}) // 结构体（忽略零值）
   db.Model(&user).Updates(map[string]interface{}{"name": "hello", "age": 0})

   db.Where("age > ?", 18).Updates(User{Age: 20}) // 批量更新
   ```

4. **删除（Delete）**
   ```go
   db.Delete(&user)          // 软删除（需模型包含DeletedAt）
   db.Unscoped().Delete(&user) // 硬删除
   db.Where("age < ?", 18).Delete(&User{}) // 条件删除
   ```

---

#### 五、高级查询

1. **预加载关联（Preload）**

   ```go
   type User struct {
     gorm.Model
     Orders []Order
   }
   type Order struct {
     gorm.Model
     UserID uint
   }

   db.Preload("Orders").Find(&users) // 加载用户的所有订单
   db.Preload("Orders", "state = ?", "paid").Find(&users) // 带条件的预加载
   ```

2. **连接查询（Joins）**

   ```go
   db.Joins("JOIN orders ON users.id = orders.user_id").Find(&users)
   ```

3. **子查询**

   ```go
   subQuery := db.Select("AVG(age)").Where("name LIKE ?", "%jin%").Table("users")
   db.Select("AVG(age) as avg_age").Group("name").Having("AVG(age) > (?)", subQuery).Find(&results)
   ```

4. **作用域（Scopes）**
   ```go
   func ActiveUser(db *gorm.DB) *gorm.DB {
     return db.Where("active = ?", true)
   }
   db.Scopes(ActiveUser).Find(&users)
   ```

---

#### 六、事务处理

1. **自动事务**

   ```go
   err := db.Transaction(func(tx *gorm.DB) error {
     if err := tx.Create(&user).Error; err != nil {
       return err
     }
     if err := tx.Update("age", 30).Error; err != nil {
       tx.Rollback() // 可选，Transaction函数会自动处理
       return err
     }
     return nil
   })
   ```

2. **手动事务**
   ```go
   tx := db.Begin()
   if err := tx.Create(&user).Error; err != nil {
     tx.Rollback()
     return
   }
   tx.Commit()
   ```

---

#### 七、关联关系

1. **一对一（Has One）**

   ```go
   type User struct {
     gorm.Model
     Profile Profile
   }
   type Profile struct {
     gorm.Model
     UserID uint
   }

   // 查询时预加载
   db.Preload("Profile").Find(&users)
   ```

2. **一对多（Has Many）**

   ```go
   type User struct {
     gorm.Model
     Orders []Order
   }
   type Order struct {
     gorm.Model
     UserID uint
   }

   // 添加关联
   user := User{ID: 1}
   db.Model(&user).Association("Orders").Append(&Order{})
   ```

3. **多对多（Many2Many）**

   ```go
   type User struct {
     gorm.Model
     Languages []Language `gorm:"many2many:user_languages;"`
   }
   type Language struct {
     gorm.Model
     Name string
   }

   // 查询关联
   db.Preload("Languages").Find(&users)
   // 替换关联
   db.Model(&user).Association("Languages").Replace([]Language{languageZH, languageEN})
   ```

---

#### 八、钩子（Hooks）

在模型生命周期中插入逻辑：

```go
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
  u.UUID = uuid.New() // 生成UUID
  return
}
func (u *User) AfterDelete(tx *gorm.DB) (err error) {
  fmt.Println("用户已删除")
  return
}
```

支持钩子：`BeforeSave`, `AfterFind`, `BeforeUpdate`, 等。

---

#### 九、性能调优

1. **禁用默认事务**
   GORM默认每个操作在事务中执行，可通过配置关闭：

   ```go
   db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
     SkipDefaultTransaction: true,
   })
   ```

2. **预处理语句缓存**

   ```go
   db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
     PrepareStmt: true,
   })
   ```

3. **连接池配置**
   ```go
   sqlDB, _ := db.DB()
   sqlDB.SetMaxIdleConns(10)        // 最大空闲连接数
   sqlDB.SetMaxOpenConns(100)       // 最大打开连接数
   sqlDB.SetConnMaxLifetime(time.Hour)
   ```

---

#### 十、日志与错误处理

1. **日志配置**

   ```go
   newLogger := logger.New(
     log.New(os.Stdout, "\r\n", log.LstdFlags),
     logger.Config{
       SlowThreshold: time.Second,
       LogLevel:      logger.Info,
       Colorful:      true,
     },
   )
   db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
   ```

2. **错误处理**
   ```go
   if result := db.First(&user); result.Error != nil {
     if errors.Is(result.Error, gorm.ErrRecordNotFound) {
       fmt.Println("记录不存在")
     } else {
       panic(result.Error)
     }
   }
   ```

---

#### 十一、高级特性

1. **自定义数据类型**
   使用`datatypes`包处理JSON、时间等：

   ```go
   import "gorm.io/datatypes"

   type User struct {
     Settings datatypes.JSON `gorm:"type:json"`
   }

   user.Settings = datatypes.JSON(`{"color": "blue"}`)
   ```

2. **原生SQL与命名参数**

   ```go
   db.Raw("SELECT * FROM users WHERE age > @age", sql.Named("age", 18)).Scan(&users)
   ```

3. **复合主键**
   ```go
   type Product struct {
     ID   string `gorm:"primaryKey"`
     Code string `gorm:"primaryKey"`
   }
   ```

---

#### 十二、最佳实践

- **避免全局变量**：传递`*gorm.DB`实例而非使用全局变量。
- **合理使用Select/Omit**：减少不必要的数据传输。
- **处理零值**：使用`map`或`Select`更新零值字段。
- **预加载优化**：避免N+1查询问题。
- **索引优化**：为高频查询字段添加索引。

---

#### 十三、常见问题

- **软删除与唯一索引冲突**：在唯一索引字段加上`DeletedAt`。
- **零值更新问题**：使用`map`而非结构体更新零值。
- **连接泄露**：确保调用`sqlDB.Close()`释放连接。

---

通过以上内容，您应能全面掌握GORM的核心功能与高级用法。建议结合官方文档（[https://gorm.io](https://gorm.io) 实践以加深理解。
