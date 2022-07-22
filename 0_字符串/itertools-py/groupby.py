from itertools import groupby

# !groupby的作用是将相同的一段元素合并到一组
# !注意groupby只会线性遍历一次 相同的值要先排在一起
# !一般的分组直接用哈希表adjMap即可

iterable = [[0, "a"], [1, "a"], [2, "b"], [3, "c"], [4, "c"], [5, "c"]]


for key, group in groupby(iterable, key=lambda x: x[1]):
    print(key, *group)
