# 计算两点的距离
import math


# 暴力算法主体函数
def calDirect(seq):
    minDis = float('inf')
    pair = []
    for i in range(len(seq)):
        for j in range(i + 1, len(seq)):
            dis = math.dist(seq[i], seq[j])
            if dis < minDis:
                minDis = dis
                pair = [seq[i], seq[j]]
    return [pair, minDis]


# 递归的对两个子问题进行求解

# ①先求出左半部分的最短距离（第一对点）
# ②再求出右半部分的最短距离（第二对点）
# ③求得中间部分的最短距离（第三对点）

# 合并问题的解
# 比较三个值中最小的一个，将其返回
print(calDirect([[1, 2], [3, 4], [2, 3]]))

