# 把 A 队的人员的平均实力值和 B 队人员的平均实力值相加，从而得到一个参赛方的综合实力评估
# 综合实力评估尽可能高，请你帮助他完成分队

# 公式变形
x, y = map(int, input().split())
nums = list(map(int, input().split()))

# 按数字大小进行索引降序排序
indexes = sorted(range(x + y), key=nums.__getitem__, reverse=True)

if x == y:
    res = ['A'] * x + ['B'] * y
elif x < y:
    res = ['B'] * (x + y)
    # 选 x 个最大值
    for i in range(x):
        res[indexes[i]] = 'A'
else:
    res = ['B'] * (x + y)
    # 选 x 个最小值
    for i in range(x):
        res[indexes[~i]] = 'A'

print(''.join(res))
