## 经典建造者模式

包含了两个关键角色：

• 待建造的实例 Instance：这个类中包含了大量需要在构造过程中完成初始化的成员属性，其中一部分字段是选填的，因此导致构造器方法的入参组合非常灵活
• 建造者 InstanceBuilder：建造者暴露出一系列链式调用风格的方法，让用户能够灵活地进行成员属性的设置，最后再调用 Build 方法一次性完成实例的构造

如果使用 Golang 中比较推崇的 NewXXX 风格进行 Food 构造器函数的定义，此时构造器就会存在很多种不同的入参组合数，最终我们不得不重复声明多个构造器函数，以应对不同的入参组合.

## Options 模式 (推荐)

1. 创建一个 defaultOption
2. 用函数修改 defaultOption 结构体
3. 兜底修复方法 repair，完成构造 option 实例过程中一些缺省值的设置

## Options 场景扩展

数据库查询条件的组合可以是多种多样的，假如我们为了支持每一种查询条件组合，都声明一个查询方法，那样的查询方法的数量会急剧膨胀，大量冗余的代码由此滋生
为了规避上面的问题，我们可以选择使用 Options 模式进行优化改造.

```go
func WithID(id int64) Option {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("id = ?", id)
    }
}


func WithName(name string) Option {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("name = ?", name)
    }
}


func WithAge(age int64) Option {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("age = ?", age)
    }
}


func WithCityID(cityID int) Option {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("city_id = ?", cityID)
    }
}


func WithPhone(phone string) Option {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("phone = ?", phone)
    }
}


func WithCreateTime(begin, end time.Time) Option {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("create_time >= ? and create_time <= end", begin, end)
    }
}
```

这种存在多种可选条件或属性组合的场景，都可以使用 Options 建造者模式，来帮助我们实现代码的优化设计.
