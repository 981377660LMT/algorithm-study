0. Go 语言包源代码不仅是核心库，还是如何使用该语言的示例
   https://go.dev/src/
1. 包命名

   - 包名应该是小写的，没有下划线或者驼峰式的命名，例如 "fmt"、"encoding/json"。
     包名应该短小、简洁、富有表现力
   - 一个简短的例子是 once.Do；once.Do(setup) 是一个好的名称，没有必要将其改为 once.DoOrWaitUntilDone(setup)。
     **长名称并不一定能使事物更易读**。有时，有用的文档注释比额外的长名称更有价值。
   - 导入包的代码将使用该名称引用其内容，因此，包中导出的名称可以利用这一点避免重复。例如，bufio 包中的缓冲读取器类型被称为 Reader，**而不是 BufReader**，因为用户将其视为 bufio.Reader，这是一个明确、简洁的名称

2. getter、setter
   在 getter 方法的名字中加上 Get 前缀既不符合惯用法，也不必要。如果你有一个名为 owner（小写、未导出）的字段，**getter 方法应该被称为 Owner（大写、导出）**，而不是 GetOwner。使用大写字母命名导出函数提供了区分字段和方法的钩子。如果需要 setter 函数，**则可能称为 SetOwner。**

   ```go
   owner := obj.Owner()
    if owner != user {
        obj.SetOwner(user)
    }
   ```

3. 接口命名
   单方法接口的命名方式是将方法名**加上-er 后缀或类似的修改**来构建代理名词：Reader、Writer、Formatter、CloseNotifier 等。
4. 对于字符串，range 更多地为您做了工作，通过解析 UTF-8，将单个 Unicode 代码点拆分出来。错误的编码会消耗一个字节，并生成替换符 U+FFFD。
5. defer 懒执行
   被推迟函数的参数（如果函数是方法，则包括接收器）是在 **defer 执行时计算的**，而不是在调用执行时计算的。
