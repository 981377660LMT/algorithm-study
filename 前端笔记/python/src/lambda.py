def multiply():
    return [lambda x: i * x for i in range(4)]


print([m(100) for m in multiply()])
# 首先需要注意的是multiply函数用生成式语法返回了一个列表，
# 列表中保存了4个Lambda函数，这4个Lambda函数会返回传入的参数乘以i的结果。需要注意的是这里有闭包（closure）现象，multiply函数中的局部变量i的生命周期被延展了，
# 由于i最终的值是3


# 方法改进：使用生成器，让函数获得i的当前值。
def multiply2():
    for i in range(4):
        yield lambda x: x * i


print([m(100) for m in multiply2()])
