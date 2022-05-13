递增的数组，bisect_left 给出的是`首个大于等于` target 的数所在的下标，bisect_right 给出的是`首个严格大于` target 的数所在的下标

# 插入不要用 bisect_insort 而是直接切片 使用 insort_left 直接超时

Python 的 list 实现基本类似 C++ vector (可以参考 CPython 的源码 listobject.c)
[x:x] = [v] 确实要比 insert 要快，因为 [x:x] = [v] 的底层实现调用 memmove 库函数来搬运插入之后的元素，而 insert 采用 for 循环搬运元素，参考 list 源码
