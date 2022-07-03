import sys


a = [0] * 3  # !最好的
b = [0, 0, 0]
c = [0 for _ in range(3)]

# [0, 0, 0] 80
# [0, 0, 0] 120
# [0, 0, 0] 88

for arr in (a, b, c):
    print(arr, sys.getsizeof(arr))

# 1.空的list 是56byte 加上三个指针 3*8=24byte 一共80byte
# 2.含有对齐操作
# 3.listappend操作 类似vector扩容
