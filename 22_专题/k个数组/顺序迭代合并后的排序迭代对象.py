from heapq import merge

a = [1, 4, 7, 10]
b = [2, 5, 6, 11]
for c in merge(a, b):
    print(c)

# heapq.merge 返回一个迭代器
#  这就意味着你可以在非常长的序列中使用它，而不会有太大的开销
# 有一点要强调的是 heapq.merge() 需要所有输入序列必须是排过序的
