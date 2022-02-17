`__foo__`: this is just a convention, a way for the Python system to use names that won't conflict with user names.
这只是一个约定，是 Python 系统使用`不会与用户名冲突的名称的一种方式。`；或者是魔法方法

`_foo`: this is just a convention, a way for the programmer to indicate that the variable is private (whatever that means in Python).
这只是一个约定，是程序员表明变量是私有的(不管在 Python 中是什么意思)的一种方式。

`__foo`: this has real meaning: the interpreter replaces this name with `_classname__foo` as a way to ensure that the name will not overlap with a similar name in another class.
这有真正的意义: 解释器用`_classname__foo `替换这个名称，以确保这个名称不会与另一个类中的类似名称重叠。

`max_`:解决冲突的命名方式
