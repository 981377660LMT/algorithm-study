# 达达帮翰翰给女生送礼物，翰翰一共准备了 N 个礼物，其中第 i 个礼物的重量是 G[i]。
# 达达的力气很大，他一次可以搬动重量之和不超过 W 的任意多个物品。
# 达达希望一次搬掉尽量重的一些物品，请你告诉达达在他的力气范围内一次性能搬动的最大重量是多少。

# 1≤N≤46,
# 如果直接枚举的话, 最多有2^46种(如果不考虑W的话)
# 如果对半枚举, 每一半最多2^23=8388608. 可以接受.
# 分成了left和right集合.
# 将right进行排序.
# 对于left的每个元素, 用二分搜索找到right中可以配对的最大元素.

# !折半枚举/折半搜索

from bisect import bisect_right


limit, n = map(int, input().split())
left = set([0])
for i in range(n // 2):
    cur = int(input())
    for pre in list(left):
        if pre + cur <= limit:
            left.add(pre + cur)

right = set([0])
for i in range(n - n // 2):
    cur = int(input())
    for pre in list(right):
        if pre + cur <= limit:
            right.add(pre + cur)
right = sorted(right)

res = 0
for num in left:
    curLimit = limit - num
    index = bisect_right(right, curLimit) - 1
    if index >= 0:
        res = max(res, num + right[index])
print(res)
