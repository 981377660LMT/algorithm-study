a = 1


def fun(a):
    a = 2


fun(a)
print(a)  # 1
# 如果在a=2之前加上global a，声明为全局变量，那么结果为2。
