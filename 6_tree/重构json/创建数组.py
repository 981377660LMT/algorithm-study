adjlist1 = [[]] * 10
adjlist1[1].append(2)
print(adjlist1)
# [[2], [2], [2], [2], [2], [2], [2], [2], [2], [2]]
# 原因是浅拷贝，我们以这种方式创建的列表，list_two 里面的三个列表的内存是指向同一块，不管我们修改哪个列表，其他两个列表也会跟着改变。
adjlist2 = [[] for i in range(10)]
adjlist2[1].append(2)
print(adjlist2)
# [[], [2], [], [], [], [], [], [], [], []]
# 我们对 adjlist2 进行更新操作，这次就能正常更新了。
