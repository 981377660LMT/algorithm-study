import sys

print(2 ** 30000)

print(sys.getsizeof(2 ** 30))
print(sys.getsizeof(2 ** 30 - 1))

# int对象的真实值存储在动态数组里
# Include/longintrepr.h
# struct _longobject {
#     PyObject_VAR_HEAD  PyObject_VAR_HEAD 是所有python对象共有的头部，
#     digit ob_digit[1];  digit 是32位无符号整型（和编译环境有关），ob_digit是一个动态数组，他存储了int对象的真实值
# };


# python的int是不会溢出的
print((2147483647 + 2147483648) >> 31)
