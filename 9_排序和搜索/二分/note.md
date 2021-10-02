`1649. 通过指令创建有序数组`

# 插入不要用 bisect_insort 而是直接切片

# 使用 insort_left 直接超时

Python 的 list 实现基本类似 C++ vector (可以参考 CPython 的源码 listobject.c)

Python 不公平的点在于, 虽然它语言解释运行很慢, 但一些底层库(list 等)却基本都是用 c 或 c++实现的, 所以快一些(相对于不用库,直接类 c 数组操作而言),但是 Python 时限又高(8-10s),

所以总体来说相当于用 c++写了个暴力又放宽了时间限制,成就了暴力的奇迹

P.S. [x:x] = [v] 确实要比 insert 要快，因为 [x:x] = [v] 的底层实现调用 memmove 库函数来搬运插入之后的元素，而 insert 采用 for 循环搬运元素，参考 list 源码

维护 sortedList 更快
