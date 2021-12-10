# itertools.islice的基本用法为：
# itertools.islice(iterable, start, stop[, step])
# 可以返回从迭代器中的start位置到stop位置的元素。如果stop为None，则一直迭代到最后位置。
# https://www.jianshu.com/p/4e0344191895
from itertools import islice

print(list(islice('ABCDEFG', 2)))


# 另外，如果在读取文件时也可以使用，比如不想读取文件第一行：
with open('tsconfig.json', 'r') as f:
    for line in islice(f, 1, None):
        print(line)

