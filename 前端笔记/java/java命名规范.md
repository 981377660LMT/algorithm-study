1. Java 中常用到的命名形式共有三种，既首字母大写的 UpperCamelCase，首字母小写的 lowerCamelCase 以及全部大写的并用下划线分割单词的 UPPERCAMELUNSER_SCORE。通常约定，类一般采用大驼峰命名，方法和局部变量使用小驼峰命名，而大写下划线命名通常是常量和枚举中使用。
   项目:spring-cloud 小写，短横线分割
   包: com.google.foo
   常量:CCCCC_DDDDD
2. 类命通常时名词或名词短语，接口名除了用名词和名词短语以外，还可以使用形容词或形容词短语，如 **Cloneable**，**Callable** 等。表示实现该接口的类有某种功能或能力。对于测试类则以它要测试的类开头，以 Test 结尾，如 **HashMapTest**。
3. 对于一些特殊特有名词缩写也可以使用全大写命名，比如 XMLHttpRequest，不过笔者认为缩写三个字母以内都大写，超过三个字母则按照要给单词算。这个没有标准如阿里巴巴中 fastjson 用 JSONObject 作为类命，而 google 则使用 JsonObjectRequest 命名，对于这种特殊的缩写，原则是统一就好。
   抽象类 Abstract 或者 Base 开头 BaseUserService
   工具类 Utils 作为后缀 StringUtils
   异常类 Exception 结尾 RuntimeException
   接口实现类 接口名+ Impl UserServiceImpl
   领域模型相关 /DO/DTO/VO/DAO 正例：UserDAO 反例： UserDo， UserDao
   设计模式相关类 Builder，Factory 等 当使用到设计模式时，需要使用对应的设计模式作为后缀，如 ThreadFactory
   处理特定功能的 Handler，Predicate, Validator 表示处理器，校验器，断言，这些类工厂还有配套的方法名如 handle，predicate，validate
4. 方法
   **返回真伪值的方法**
   is 对象是否符合期待的状态 isValid
   can 对象能否执行所期待的动作 canRemove
   should 调用方执行某个命令或方法是好还是不好,应不应该，或者说推荐还是不推荐 shouldMigrate
   has 对象是否持有所期待的数据和属性 hasObservers
   needs 调用方是否需要执行某个命令或方法 needsMigrate
   **用来检查的方法**
   ensure 检查是否为期待的状态，不是则抛出异常或返回 error code ensureCapacity
   validate 检查是否为正确的状态，不是则抛出异常或返回 error code validateInputs
   **按需求才执行的方法**
   might 同上 mightCreate
   try 尝试执行，失败时抛出异常或是返回 errorcode tryCreate
   OrDefault 尝试执行，失败时返回默认值 getOrDefault
   OrElse 尝试执行、失败时返回实际参数中指定的值 getOrElse
   force 强制尝试执行。error 抛出异常或是返回值 forceCreate, forceStop
   IfNeeded 需要的时候执行，不需要的时候什么都不做 drawIfNeeded
   **异步相关方法**
   Async 异步方法 sendAsync
   Sync 对应已有异步方法的同步方法 sendSync
   schedule Job 和 Task **放入队列** schedule, scheduleJob
   execute 执行异步方法（注：我一般拿这个做同步方法名） execute, executeTask
   start 同上 start, startJob
   cancel 停止异步方法 cancel, cancelJob
   stop 同上 stop, stopJob
   **回调方法**
   on 事件发生时执行 onCompleted
   before 事件发生前执行 beforeUpdate
   will 同上 willUpdate
   did 同上 didUpdate
   should 确认事件是否可以发生时执行 shouldUpdate
   **操作对象生命周期的方法**
   初始化。也可作为延迟初始化使用 initialize
   销毁的替代 destroy
   **与数据相关的方法**
   新创建 createAccount
   从既有的某物新建，或是从其他的数据新建 fromConfig
   转换 toString
   读取 loadAccount
   更新既有某物 updateAccount
   远程读取 fetchAccount
   保存 commitChange
   保存或应用 applyChange
   清除数据或是恢复到初始状态 clearAll
   **成对出现的动词**
   backup 备份 restore 恢复
   create 创建 destory 移除
   add 增加 remove 删除 (listener...)
   load 载入 save 保存
   split 分割 merge 合并
   bind 绑定 separate 分离
   parse 解析 emit 生成
   update 更新 revert 复原
   expand 展开 collapse 折叠
   collect 收集 aggregate 聚集
5. 妙用介词，如 for(可以用同音的 4 代替), to(可用同音的 2 代替), from, with，of 等。 如类名采用 UserService4MySqlDAO，方法名 getUserInfoFromRedis，convertJson2Map 等。
6. 注解的原则
   Nothing is strange 没有注解的代码对于阅读者非常不友好，哪怕代码写的在清除，阅读者至少从心理上会有抵触，更何况代码中往往有许多复杂的逻辑，所以`一定要写注解`，不仅要记录代码的逻辑，还有说清楚修改的逻辑。
   Less is more 从代码维护角度来讲，代码中的注解一定是`精华`中的精华。合理清晰的命名能让代码易于理解，对于逻辑简单且命名规范，能够清楚表达代码功能的代码不需要注解。滥用注解会增加额外的负担，更何况大部分都是废话。
   Advance with the time 注解应该`随着代码的变动而改变`，注解表达的信息要与代码中完全一致。通常情况下修改代码后一定要修改注解。

7. 常用命名的建议
   我们在程序中见到的很多名字都很模糊，例如 tmp。就算是看上去合理的词，如 size 或者 get，也都没有装入很多信息。
   **找到更有表现力的词**
   **避免像 tmp 和 retval 这样泛泛的名字**
   **循环迭代器**
   如果不把循环索引命名为（i、j、k），另一个选择可以是（club_i、members_i、user_i）或者，更简化一点（ci、mi、ui）。

```JS
if (clubs[ci].members[ui] == users[mi]) #缺陷！第一个字母不匹配。

```

**带单位的值**
通过给变量结尾追加\_ms，我们可以让所有的地方更明确：

```JS
var start_ms = (new Date()).getTime(); // top of the page
...
var elapsed_ms = (new Date()).getTime() - start_ms; // bottom of the page
document.writeln("Load time was: " + elapsed_ms / 1000 + " seconds");

```

一、问题
由于单词量的匮乏，发现写代码的时候，似乎只会使用 get，set,send,find,start,make 等简单单词。这些单词很多时候能表达的意义非常有限。
比如，从一个栈的数据结构中，取出一个元素。`取出之后，这个元素就不在这个数据结构中了`。此时如果使用 get_element，能表达这个意思吗？可能使用 `pop_element` 更合适。
比如，我想获取栈中的第 N 个元素的值，我要继续使用 get 吗？写成 get_element_value?获取值是可以使用 get 的。
同样都是获取，怎么写能够让代码更可读？我觉得需要积累一些常用的在编程中的单词的用法。

空洞四剑客：get, make, find, stop

```JS
/// 获取当前图片的大小
- (CGSize)getCurrentLayerSize;
```

e.g. "get"这个词过于笼统，蕴含的信息量太少了，获取一个 size 有`很多种途径可以拿`：

直接取：- (CGSize)currentLayerSize;
计算：- (CGSize)`calculate`CurrentLayerSize;
请求：- (CGSize)`fetch`CurrentLayerSize;
生成：- (CGSize)`generate`CurrentLayerSize;

e.g. "make"这个词太业余了，可以有更加专业的动词替换：
生成：- (NSArray \*)generate/create/calculateListViewModels;
有副作用：- (void)`setup`/`update`ListViewModels;

e.g. "find" 似乎用于 js 内置 api find/findIndex 接受一个 predicate 函数
属于`查询方法前缀`： get、find。
`搜索性查询方法`前缀应该用：query、search。
`统计类前缀：count。`
`操作类前缀：insert、add、create、update、delete。`

e.g. "stop"本身没啥问题，只是可能存在更优的固定搭配：
kill + thread
terminate + app / runloop
pause / remove + animation 尽量避免给变量取名为 ret, tmp, case, num 等空泛的词，请绞尽脑汁取一个语义化的名字。
前后缀可以丰富含义，rawContent, escapedString, stringifiedText, stayTimeMs

也可以模拟内置 api 命名风格

## 类定义的顺序

l 静态成员变量 / Static Fields
l 静态初始化块 / Static Initializers
l 成员变量 / Fields
l 初始化块 / Initializers
l 构造器 / Constructors
l 静态成员方法 / Static Methods
l 成员方法 / Methods
l `重载自 Object 的方法如 toString(), hashCode() 和 main 方法`
l 类型(内部类) / Types(Inner Classes)

同等的顺序下，再按 `public, protected,default, private` 的顺序排列。
