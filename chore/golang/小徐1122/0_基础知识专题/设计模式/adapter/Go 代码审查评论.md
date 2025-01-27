## Comment Sentences

用于文档声明的注释应为完整句子，即使这看起来有些多余。这种方法在提取到 godoc 文档时格式良好。注释应以所描述事物的名称开头，并以句号结束：

```go
// Request represents a request to run a command.
type Request struct { ...

// Encode writes the JSON encoding of req to w.
func Encode(w io.Writer, req *Request) { ...
```

## Contexts

不要向结构体类型添加 Context 成员；而是为该类型中需要传递它的每个方法添加一个 ctx 参数。
唯一的例外是那些签名必须与标准库或第三方库中的接口匹配的方法。

## Declaring Empty Slices

- 在声明一个空切片时，优先使用`var`关键字

```go
var users []string
users2 := []string{}
```

`前者声明了一个 nil 切片值，而后者是非 nil 但长度为零。`
它们在功能上是等价的——它们的 len 和 cap 都是零——但 nil 切片是首选风格。
请注意，在某些有限情况下，非空但长度为零的切片是首选的，`例如在编码 JSON 对象时（ nil 切片编码为 null ，而 []string{} 编码为 JSON 数组 [] ）。`

## Doc Comments

所有顶级导出名称都应该有文档注释，非平凡的未导出类型或函数声明也应如此。

## Examples

在添加新包时，请包含预期用法的示例：一个可运行的示例或一个演示完整调用序列的简单测试。

## Interfaces

在 Go 语言中，接口的设计哲学与传统面向对象语言有显著差异，主要体现在 **接口的归属**、**模拟测试方式** 和 **定义时机** 三个方面。以下通过具体场景和对比分析，帮助你理解这些核心原则：

---

### 一、接口的归属：由消费者定义，而非生产者

#### ❌ 错误做法（生产者定义接口）

```go
// 生产者包 producer.go
type Thinger interface { Thing() bool }  // 生产者定义接口

type defaultThinger struct{}
func (t defaultThinger) Thing() bool { return true }

func NewThinger() Thinger { return defaultThinger{} }  // 返回接口类型
```

**问题**：

1. 生产者绑定了接口，后续新增方法需修改接口，导致所有实现被迫更新
2. 消费者被迫依赖生产者的接口，失去灵活性

#### ✅ 正确做法（消费者定义接口）

```go
// 生产者包 producer.go
type Thinger struct{}  // 返回具体类型
func (t Thinger) Thing() bool { return true }

func NewThinger() Thinger { return Thinger{} }
```

```go
// 消费者包 consumer.go
type Thinger interface { Thing() bool }  // 消费者定义所需接口

func Foo(t Thinger) string { ... }  // 依赖自己定义的接口
```

**优势**：

1. 生产者可自由扩展结构体方法，无需担心破坏接口契约
2. 消费者只声明自己关心的接口方法，实现解耦

---

### 二、模拟测试：消费者自行模拟，不依赖生产者

#### ❌ 错误做法（绑定生产者接口）

```go
// 测试代码被迫依赖生产者接口
import "producer"

type fakeThinger struct{}
func (t fakeThinger) Thing() bool { ... }  // 必须实现生产者包的接口

func TestFoo(t *testing.T) {
    result := Foo(fakeThinger{})
    // ...
}
```

**问题**：

1. 测试代码与生产者实现耦合
2. 生产者接口变更会导致测试代码被迫修改

#### ✅ 正确做法（消费者自定义接口）

```go
// 测试代码 consumer_test.go
type mockThinger struct{}  // 定义自己的测试接口
func (t mockThinger) Thing() bool { return false }

func TestFoo(t *testing.T) {
    result := Foo(mockThinger{})
    // ...
}
```

**优势**：

1. 测试完全独立，不依赖生产者包的接口定义
2. 可精准控制模拟行为，只需实现测试关心的接口方法

---

### 三、接口定义时机：按需定义，而非提前设计

#### ❌ 错误做法（未经验证的抽象）

```go
// 设计阶段过早定义接口
type Storage interface {
    Save(data []byte) error
    Load(id string) ([]byte, error)
    Delete(id string) error
}

// 实际使用时发现部分场景只需要 Save
func process(s Storage) { ... }  // 被迫依赖不需要的方法
```

**问题**：

1. 接口包含冗余方法，违反接口隔离原则
2. 增加不必要的实现复杂度

#### ✅ 正确做法（使用时提炼接口）

```go
// 先写具体实现
type DiskStorage struct{}

func (s DiskStorage) Save(data []byte) error { ... }
func (s DiskStorage) Load(id string) ([]byte, error) { ... }

// 使用时按需定义最小接口
type Saver interface { Save([]byte) error }

func process(s Saver) { ... }  // 只依赖实际需要的方法
```

**优势**：

1. 接口最小化，精准描述依赖
2. 更易维护，新增方法不影响历史代码

---

### 关键原则总结

| 维度             | 正确做法                          | 错误做法               |
| ---------------- | --------------------------------- | ---------------------- |
| **接口定义位置** | 由消费者包定义所需接口            | 由生产者包定义接口     |
| **返回类型**     | 生产者返回具体类型（结构体/指针） | 生产者返回接口类型     |
| **测试模拟**     | 消费者自定义模拟接口              | 依赖生产者接口进行模拟 |
| **接口设计**     | 按需定义小接口，延迟抽象          | 过早定义大而全的接口   |

---

### 典型案例：数据库操作适配

#### 生产者包实现

```go
// producer/mysql.go
type MySQL struct{}  // 返回具体类型

func (m MySQL) Query(query string) ([]byte, error) { ... }  // 方法1
func (m MySQL) Exec(cmd string) error { ... }               // 方法2
```

#### 消费者按需定义接口

```go
// consumer/service.go
type QueryExecutor interface {  // 只定义需要的方法
    Query(string) ([]byte, error)
}

func GetData(qe QueryExecutor) { ... }  // 依赖精准
```

#### 测试模拟

```go
// consumer/service_test.go
type mockExecutor struct{}
func (m mockExecutor) Query(q string) ([]byte, error) {
    return []byte("test"), nil
}

func TestGetData(t *testing.T) {
    GetData(mockExecutor{})  // 无需依赖真实数据库
}
```

---

### 这些设计带来的长期优势

1. **降低耦合度**：生产者与消费者通过具体类型和隐式接口交互
2. **提升可维护性**：生产者可自由扩展结构体方法，不影响消费者
3. **增强测试灵活性**：消费者可自主定义测试专用的最小接口
4. **减少过度设计**：避免创建从未使用过的复杂接口层次

遵循这些原则的 Go 代码，能更好地适应需求变化，在大型项目中显著降低模块间的依赖成本，是符合 Go 语言设计哲学的最佳实践。
