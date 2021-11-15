https://github.com/jackfrued/Python-Interview-Bible/blob/master/Python%E9%9D%A2%E8%AF%95%E5%AE%9D%E5%85%B8-%E5%9F%BA%E7%A1%80%E7%AF%87-2020.md
https://github.com/taizilongxu/interview_python
https://github.com/lengyue1024/BAT_interviews/blob/master/

1. python dict 原理
   Python 字典实现为散列表。
   Python dict 使用开放寻址来解决散列冲突
   Python 3.6 之后是 OrderedDict
2. python list 原理
   在 CPython 中， Python 中的列表基于 **PyListObject** 实现, 中， 是指针的数组（arrays of pointers）
   元组的标准实现也只是一个数组
3. python 循环引用会不会内存泄漏
4. `对__if__name__ == 'main'的理解陈述`
   **name**是当前模块名，当模块被直接运行时模块名为*main*，也就是当前的模块，当模块被导入时，模块名就不是**main**，即代码将不会执行。
5. 介绍一下 except 的用法和作用？
   try…except…except…else…
   try 下的语句正常执行，则执行 else 块代码。如果发生异常，就不会执行如果存在 finally 语句，最后总是会执行。

6. ` Python中__new__与__init方法的区别`
   `__new__`:它是`创建对象时`调用，会返回当前对象的一个实例，可以用*new*来实现单例
   `__init__`:它是创建对象后调用，对当前对象的一些实例初始化，无返回值

   1、`__new__`**至少要有一个参数 cls，代表当前类**，此参数在实例化时由 Python 解释器自动识别
   2、`__new__`**必须要有返回值**，返回实例化出来的实例，这点在自己实现`__new__`时要特别注意，可以 return 父类（通过 super(当前类名, cls)）`__new__`出来的实例，或者直接是 object 的`__new__`出来的实例
   3、`__init__`有一个参数 self，就是这个`__new__`返回的实例，`__init__`在**new**的基础上可以完成一些其它初始化的动作，`__init__`**不需要返回值**

   执行`__new__`方法获得保存对象所需的**内存空间**，再通过`__init__`执行对内存空间数据的**填充**（对象属性的初始化）。`__new__`方法的返回值是创建好的 Python 对象（的引用），而`__init__`方法的第一个参数就是这个对象（的引用），所以在`__init__`中可以完成对对象的初始化操作。`__new__`是类方法，它的第一个参数是类，`__init__`是对象方法，它的第一个参数是对象。

7. 什么是栈溢出?怎么解决?
   因为栈一般默认为 1-2m，一旦出现死循环或者是大量的递归调用，在不断的压栈过程中，造成栈容量超过 1m 而导致溢出。
   栈溢出的情况有两种：1)`局部数组过大`。当函数内部数组过大时，有可能导致堆栈溢出。2)`递归调用层次太多`。递归函数在运行时会执行压栈操作，当压栈次数太多时，也会导致堆栈溢出。
   解决方法：1)用栈把递归转换成非递归。2)增大栈空间。
8. 建立一个简单 tcp 服务器须要的流程？
   1.socket 建立一个套接字

   2.bind 绑定 ip 和 port

   3.listen 使套接字变为能够被动连接

   4.accept `等待客户端的连接`

   5.recv/send 接收发送数据

9. TTL，MSL，RTT？
   （1）MSL：报文最大生存时间”，他是任何报文在网络上存在的最长时间，超过这个时间报文将被丢弃。
   （2）TTL：TTL 是 time to live 的缩写，中文能够译为“生存时间”，
   （3）RTT： RTT 是客户到服务器往返所花时间（round-trip time，简称 RTT），TCP 含有动态估算 RTT 的算法。TCP 还持续估算一个给定链接的 RTT，这是由于 `RTT 受网络传输拥塞程序的变化而变化。`

10. python 里的单例模式?
    在程序运行中该类只实例化一次，并提供了一个全局访问点 Python 的模块就是天然的单例模式 当模块在第一次导入时，就会生成 .pyc 文件 当第二次导入时就会直接先加载 .pyc 文件，而不会再次执行模块代码。

11. Python2 和 Python3 的区别，如何实现 python2 代码迁移到 Python3 环境
    代码迁移：python3 有个内部工具叫做 2to3.py，位置在 Python3/tool/script 文件夹。
    python2 和 python3 区别？列举 5 个

    - Python 2 中的 xrange 函数在 Python 3 中被 range 函数取代。
    - python2 中使用 ascii 编码，python3 中使用 utf-8 编码
    - python2 中 unicode 表示字符串序列，str 表示字节序列;python3 中 str 表示字符串序列，byte 表示字节序列
    - python2 range(1,10)返回列表，python3 中返回迭代器，节约内存
    - Python3 使用 print 必须要以小括号包裹打印内容，比如 print('hi'),Python2 既可以使用带小括号的方式，也可以使用一个空格来分隔打印内容，比如 print 'hi'
    - python2 中是 raw_input()函数，python3 中是 input()函数
    - Python 3 中字典的 keys、values、items 方法都不再返回 list 对象，而是返回 view object，内置的 map、filter 等函数也不再返回 list 对象，而是返回迭代器对象。

12. Python2 和 Python3 的编码方式有什么差别
    在 python2 中主要有 str 和 unicode 两种字符串类型，而到 python3 中改为了 bytes 和 str,
    在 python3 中如果是写或者读 bytes 类型就必需带上’b’.
13. 斐波那契生成器
    ```Python
    def fib(num):
        a, b = 0, 1
        for _ in range(num):
            a, b = b, a + b
            yield a
    ```
14. 查找变量的顺序
    LEGB 法则
    L-----LOCAL 局部

    E-------ENCLOSE------嵌套作用域

    G-------GLOBLE-------全局

    B---------BUILT-IN------内置

```Python
a = 100
b = 2
c = 10
def waibu():
    b=200
    c=2
    def neibu():
        c=300
        print(c)#LOCAL局部变量
        print(b)#ENCLOSE嵌套
        print(a)#GLOBAL全局
        print(max)#BUILT-IN内置
    neibu()
waibu()
```

15. 如何在一个 function 里面设置一个全局的变量？
    Global VarName
16. 谈下 python 的 GIL
    GIL 是 python 的全局解释器锁，同一进程中假如有多个线程运行，`一个线程在运行 python 程序的时候会霸占 python 解释器（加了一把锁即 GIL）`，使该进程内的其他线程无法运行，等该线程运行完后其他线程才能运行。如果线程运行过程中遇到耗时操作，则解释器锁解开，使其他线程运行。所以在多线程中，`线程的运行仍是有先后顺序的，并不是同时进行。`

    多进程中因为每个进程都能被系统分配资源，相当于`每个进程有了一个 python 解释器`，所以多进程可以实现多个进程的同时运行，缺点是进程系统资源开销大

17. python2 和 python3 的 range（100）的区别
    python2 返回列表，python3 返回迭代器，节约内存
18. 一句话解释什么样的语言能够用装饰器?
    函数可以作为参数传递的语言，可以使用装饰器
19. 正则 re.complie 作用
    re.compile 是将正则表达式编译成一个对象，加快速度，并重复使用
20. 提高 python 运行效率的方法
    1、使用生成器，因为可以节约大量内存

    2、循环代码优化，避免过多重复代码的执行

    3、核心模块用 Cython PyPy 等，提高效率

    4、多进程、多线程、协程

    5、多个 if elif 条件判断，`可以把最有可能先发生的条件放到前面写`，这样可以减少程序判断的次数，提高效率

21. any()和 all()方法
    any():只要`迭代器中有一个`元素为真就为真
    all():迭代器中所有的判断项返回都是真，结果才为真
    python 中什么元素为假？
    答案：（0，空字符串，空列表、空字典、空元组、None, False）

22. python 中的异常
    IOError：输入输出异常

    AttributeError：试图访问一个对象没有的属性

    ImportError：无法引入模块或包，基本是路径问题

    IndentationError：语法错误，代码没有正确的对齐

    IndexError：下标索引超出序列边界

    KeyError:试图访问你字典里不存在的键

    SyntaxError:Python 代码逻辑语法出错，不能执行

    NameError:使用一个还未赋予对象的变量

23. 解释 Python 中的 dir()函数
    dir() 函数不带参数时，返回当前范围内的变量、方法和定义的类型列表；带参数时，返回参数的属性、方法列表。

```Python
>>> dir(copy.copy)
[‘__annotations__’, ‘__call__’, ‘__class__’, ‘__closure__’, ‘__code__’, ‘__defaults__’, ‘__delattr__’, ‘__dict__’, ‘__dir__’, ‘__doc__’, ‘__eq__’, ‘__format__’, ‘__ge__’, ‘__get__’, ‘__getattribute__’, ‘__globals__’, ‘__gt__’, ‘__hash__’, ‘__init__’, ‘__init_subclass__’, ‘__kwdefaults__’, ‘__le__’, ‘__lt__’, ‘__module__’, ‘__name__’, ‘__ne__’, ‘__new__’, ‘__qualname__’, ‘__reduce__’, ‘__reduce_ex__’, ‘__repr__’, ‘__setattr__’, ‘__sizeof__’, ‘__str__’, ‘__subclasshook__’]

```

24. 当退出 Python 时，是否释放全部内存？
    答案是 No。循环引用其它对象或引用自全局命名空间的对象的模块，在 Python 退出时并非完全释放。
    另外，也不会释放 C 库保留的内存部分。
25. 下面代码的输出结果将是什么？

```Python
list = ['a','b','c','d','e']
print(list[10:])

仅仅返回一个空列表。这成为特别让人恶心的疑难杂症，因为运行的时候没有错误产生，导致Bug很难被追踪到。
```

26. python 新式类和经典类的区别？
    a. 在 python 里凡是继承了 object 的类，都是新式类
    b. Python3 里只有新式类
    c. Python2 里面继承 object 的是新式类，没有写父类的是经典类
    d. 经典类目前在 Python 里基本没有应用
    e. **新式类继承是根据 C3 算法,旧式类是从左到右深度优先**
27. 什么是僵尸进程和孤儿进程？怎么避免僵尸进程？
    **孤儿进程**： 父进程退出，子进程还在运行的这些子进程都是孤儿进程，孤儿进程将被 init 进程（进程号为 1）所收养，并由 init 进程对他们完成状态收集工作。
    **僵尸进程**： 进程使用 fork 创建子进程，如果子进程退出，而父进程并没有调用 wait 获 waitpid 获取子进程的状态信息，那么子进程的进程描述符仍然保存在系统中的这些进程是僵尸进程。 避免僵尸进程的方法： 1.fork 两次用孙子进程去完成子进程的任务 2.用 wait()函数使父进程阻塞 3.使用信号量，在 signal handler 中调用 waitpid,这样父进程不用阻塞
28. 什么是 cgi,wsgi,uwsgi,uWSGI, FastCGI ?
    **CGI**:是通用网关接口，是连接 web 服务器和应用程序的接口，用户通过 CGI 来获取动态数据或文件等。 CGI 程序是一个独立的程序，它可以用几乎所有语言来写，包括 perl，c，lua，python 等等。
    **FastCGI**: FastCGI 是语言无关的、可伸缩架构的 CGI 开放扩展，其主要行为是将 CGI 解释器进程保持在内存中并因此获得较高的性能。
    **WSGI**: WSGI `专指 Python 应用程序`,python 的`web服务器网关接口，是一套协议`。用于接收用户请求并将请求进行初次封装，然后将请求交给 web 框架。 实现 wsgi 协议的模块：wsgiref,本质上就是编写一 socket 服务端，用于接收用户请求（django) werkzeug,本质上就是编写一个 socket 服务端，用于接收用户请求(flask)
    **uwsgi**: 与 WSGI 一样是一种通信协议，它是 uWSGI 服务器的独占协议，用于定义传输信息的类型。
    **uWSGI**: 是一个 web 服务器，实现了 WSGI 的协议，uWSGI 协议，http 协议
29. 假设你使用的是官方的 CPython，说出下面代码的运行结果

```Python
a, b, c, d = 1, 1, 1000, 1000
print(a is b, c is d)

def foo():
    e = 1000
    f = 1000
    print(e is f, e is d)
    g = 1
    print(g is a)

foo()
```

30. Python 计数引用？
    对于 CPython 解释器来说，Python 中的`每一个对象其实就是 PyObject 结构体`，它的内部有一个名为 `ob_refcnt` 的引用计数器成员变量。

```C++
typedef struct _object {
    _PyObject_HEAD_EXTRA
    Py_ssize_t ob_refcnt;
    struct _typeobject *ob_type;
} PyObject
```

程序在运行的过程中 ob_refcnt 的值会被更新并藉此来反映引用有多少个变量引用到该对象。当对象的引用计数值为 0 时，它的内存就会被释放掉。
可以通过 sys 模块的 getrefcount 函数来获得对象的引用计数。

31. 下面这段代码的执行结果是什么。

```Python
def multiply():
    return [lambda x: i * x for i in range(4)]

print([m(100) for m in multiply()])

```

32. Python 中为什么没有函数重载？

    Python 是解释型语言，函数重载现象通常出现在编译型语言中
    其次 Python 是动态类型语言，函数的参数没有类型约束，也就无法根据参数类型来区分重载。
    再者 Python 中函数的参数可以有默认值，可以使用可变参数和关键字参数，因此即便没有函数重载，也要可以让一个函数根据调用者传入的参数产生不同的行为。

33. 说出下面代码的运行结果。

```Python
class Parent:
    x = 1

class Child1(Parent):
    pass

class Child2(Parent):
    pass

print(Parent.x, Child1.x, Child2.x)
Child1.x = 2
print(Parent.x, Child1.x, Child2.x)
Parent.x = 3
print(Parent.x, Child1.x, Child2.x)

1 1 1
1 2 1
3 2 3

```

34. 什么是鸭子类型（duck typing）？
    在 Python 语言中，有很多 bytes-like 对象（如：bytes、bytearray、array.array、memoryview）、file-like 对象（如：StringIO、BytesIO、GzipFile、socket）、path-like 对象（如：str、bytes），其中 file-like 对象都能支持 read 和 write 操作，可以像文件一样读写，这就是所谓的对象有鸭子的行为就可以判定为鸭子的判定方法。再比如 Python 中列表的 extend 方法，它需要的参数并不一定要是列表，只要是可迭代对象就没有问题。
35. 谈谈你对“猴子补丁”（monkey patching）的理解。
    “猴子补丁”是动态类型语言的一个特性，代码运行时在不修改源代码的前提下改变代码中的方法、属性、函数等以达到热补丁（hot patch）的效果。
    `单元测试中的 Mock 技术`也是对猴子补丁的应用，Python 中的 unittest.mock 模块就是解决单元测试中用 Mock 对象替代被测对象所依赖的对象的模块。
36. 何剖析 Python 代码的执行性能？
    剖析代码性能可以使用 Python 标准库中的 **cProfile** 和 **pstats** 模块
37. 解释一下线程池的工作原理。
    `池化技术`就是一种典型空间换时间的策略，我们使用的数据库连接池、线程池等都是池化技术的应用，
    当系统**比较空闲**时，大部分线程长时间处于闲置状态时，线程池可以自动销毁一部分线程，回收系统资源。
    当线程池中所有的**线程都被占用后**，可以选择自动创建一定数量的新线程，用于处理更多的任务，也可以选择让任务排队等待直到有空闲的线程可用。
38. 举例说明什么情况下会出现 KeyError、TypeError、ValueError。
    变量 a 是一个字典，执行 int(a['x'])这个操作就有可能引发上述三种类型的异常。
    如果字典中**没有键 x**，会引发 KeyError；
    如果键 x 对应的值不是 str、float、int、bool 以 bytes-like 类型，在**调用 int 函数构造 int 类型的对象时，会引发 TypeError；**
    如果 a[x]是一个字符串或者字节串，而对应的内容又**无法处理成 int 时**，将引发 ValueError。
39. 运行结果

        ```Python
        def extend_list(val, items=[]):
        items.append(val)
        return items

        list1 = extend_list(10)
        list2 = extend_list(123, [])
        list3 = extend_list('a')
        print(list1)
        print(list2)
        print(list3)

        [10, 'a']
        [123]
        [10, 'a']
        ```

    Python 函数在定义的时候，默认参数 items 的值就被计算出来了，即[]

40. 如何读取大文件，例如内存只有 4G，如何读取一个大小为 8G 的文件？
    用内置函数 iter 将文件对象处理成迭代器对象，每次只读取少量的数据进行处理，
    或者切分

    ```BASH
    <!-- 将名为filename的文件切割为10个文件 -->
    split -n 10 filename
    ```

41. 说一下你知道的 Python 编码规范。
    1. 变量、函数和属性应该使用小写字母来拼写，如果有多个单词就使用下划线进行连接。
    2. 类中受保护的实例属性，应该以一个下划线开头。
    3. 类中私有的实例属性，应该以两个下划线开头。
    4. 模块级别的常量，应该采用全大写字母，如果有多个单词就用下划线进行连接。
    5. 类的实例方法，应该把第一个参数命名为 self 以表示对象自身。
    6. 类的类方法，应该把第一个参数命名为 cls 以表示该类自身。
    7. 采用内联形式的否定词，而不要把否定词放在整个表达式的前面。例如：if a is not b 就比 if not a is b 更容易让人理解。
    8. 不要用检查长度的方式来判断字符串、列表等是否为 None 或者没有元素，应该用 if not x 这样的写法来检查它。
    9. 引入模块的时候，from math import sqrt 比 import math 更好。
    10. 如果有多个 import 语句，应该将其分为三部分，从上到下分别是 **Python 标准模块、第三方模块和自定义模块**，**每个部分内部应该按照模块名称的字母表顺序来排列。**
42. 运行下面的代码是否会报错，如果报错请说明哪里有什么样的错，如果不报错请说出代码的执行结果。

```Python
class A:
    def __init__(self, value):
        self.__value = value

    @property
    def value(self):
        return self.__value

obj = A(1)
obj.__value = 2
print(obj.value)
print(obj.__value)

1
2

如果不希望代码运行时动态的给对象添加新属性，可以在定义类时使用__slots__魔法。例如，我们可以在上面的A中添加一行__slots__ = ('__value', )，再次运行上面的代码，将会在原来的第10行处产生AttributeError错误。
```

43. 对下面给出的字典按值从大到小对键进行排序。

```Python
prices = {
    'AAPL': 191.88,
    'GOOG': 1186.96,
    'IBM': 149.24,
    'ORCL': 48.44,
    'ACN': 166.89,
    'FB': 208.09,
    'SYMC': 21.29
}

sorted(prices, key=lambda x: prices[x], reverse=True)
```

44. 说一下 namedtuple 的用法和作用。
    命名元组与普通元组一样是不可变容器
    命名元组的本质就是一个类，所以它还可以作为父类创建子类。
45. 返回该列表最大的嵌套深度
    例如：列表[1, 2, 3]的嵌套深度为 1，列表[[1], [2, [3]]]的嵌套深度为 3。

```Python
def list_depth(items):
    if isinstance(items, list):
        max_depth = 1
        for item in items:
            max_depth = max(list_depth(item) + 1, max_depth)
        return max_depth
    return 0
```

46. 一个通过网络获取数据的函数（可能会因为网络原因出现异常），写一个装饰器让这个函数在出现指定异常时`可以重试指定的次数`，并在`每次重试之前随机延迟一段时间`，最长延迟时间可以`通过参数进行控制`。

```Python
from functools import wraps
from random import random
from time import sleep


def retry(*, retry_times=3, max_wait_secs=5, errors=(Exception, )):

    def decorate(func):

        @wraps(func)
        def wrapper(*args, **kwargs):
            for _ in range(retry_times):
                try:
                    return func(*args, **kwargs)
                except errors:
                    sleep(random() * max_wait_secs)
            return None

        return wrapper

    return decorate

```

47. 列表中有 1000000 个元素，取值范围是[1000, 10000)，设计一个函数找出列表中的重复元素。

```Python
def find_dup(items: list):
    dups = [0] * 9000
    for item in items:
        dups[item - 1000] += 1
    for idx, val in enumerate(dups):
        if val > 1:
            yield idx + 1000

```

48. python/js 函数传参都是传对象的“引用”;也可以看似 c 中`void*`的感觉。

```Python
a = 1
def fun(a):
    print "func_in",id(a)   # func_in 41322472
    a = 2
    print "re-point",id(a), id(2)   # re-point 41322448 41322448
print "func_out",id(a), id(1)  # func_out 41322472 41322472
fun(a)
print a  # 1

类型是属于对象的，而不是变量。而对象有两种,“可更改”（mutable）与“不可更改”（immutable）对象。在python中，strings, tuples, 和numbers是不可更改的对象，而 list, dict, set 等则是可以修改的对象。(这就是这个问题的重点)

当一个引用传递给函数的时候,函数自动复制一份引用,这个函数里的引用和外边的引用没有半毛关系了.所以第一个例子里函数把引用指向了一个不可变对象,当函数返回的时候,外面的引用没半毛感觉.而第二个例子就不一样了,函数内的引用指向的是可变对象,对它的操作就和定位了指针地址一样,在内存里进行修改.
```

49. Python 中单下划线和双下划线
    `__foo__`:一种约定,Python 内部的名字,用来区别其他用户自定义的命名,以防冲突，就是例如`__init__`(),`__del__`(),`__call__`()这些特殊方法

    `_foo`:一种约定,用来指定变量私有.程序员用来指定私有变量的一种方式.不能用`from module import * `导入，其他方面和公有一样访问；

    `__foo`:这个有真正的意义:解析器用`_classname__foo`来代替这个名字,以区别和其他类相同的命名,它无法直接像公有成员一样随便访问,通过对象名`._类名__xxx`这样的方式可以访问.

50. 协程
    协程是进程和线程的升级版,进程和线程都面临着内核态和用户态的切换问题而耗费许多切换时间,而协程就是用户自己控制切换的时机,不再需要陷入系统的内核态.
    Python 里最常见的 yield 就是无栈协程
51. 程序编译与链接
    Bulid 过程可以分解为 4 个步骤:

    预处理(Prepressing):处理以“#”开始的预编译指令,删除所有注释
    编译(Compilation):词法分析、语法分析、AST、优化、生成代码
    汇编(Assembly):将汇编代码转化成机器可以执行的指令
    链接(Linking):把各个模块之间相互引用的部分处理好

52. 静态链接和动态链接
    静态链接方法：静态链接的时候，载入代码就会把程序会用到的动态代码或动态代码的地址确定下来(类似**import**)
    动态链接方法：使用这种方式的程序并不在一开始就完成动态链接，而是直到真正调用动态库代码时，载入程序才计算(被调用的那部分)动态代码的逻辑地址，(类似**require**)
