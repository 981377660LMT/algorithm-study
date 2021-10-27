# Variables with annotations do not need to be assigned a value.
# So by convention, we omit them in the stub file.
x: int

# Function bodies cannot be completely removed. By convention,
# we replace them with `...` instead of the `pass` statement.

def func_1(code: str) -> int: ...

# We can do the same with default arguments.
def func_2(a: int, b: int = ...) -> int: ...

# 存根文件是用普通的 Python 3语法编写的，
# 但通常省略了运行时逻辑，如变量初始化器、函数体和缺省参数。
# 如果不能完全省略某些运行时逻辑，建议使用省略号表达式(...)替换或省略它们
