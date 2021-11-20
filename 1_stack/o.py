aa = dict()
aa['b'] = 1
aa['a'] = 1

print(aa)

# 从python 3.6 开始，builtin的dict对象，也能够记录insert order，但还是跟collections.OrderedDict有不一样的细节。
# 两个dict对象比较，无论按什么顺序insert数据，只要数据是一样的，他们就是相等的：
# 而两个OrderedDict对象，要让它们相等，不仅数据要完全一样，插入顺序也要完全一样
