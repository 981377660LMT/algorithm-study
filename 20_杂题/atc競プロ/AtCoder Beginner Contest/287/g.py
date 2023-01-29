import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# 将离散化后的所有卡片的分数按照从小到大排序记为 d1 d2 ... dk
# 线段树1维护每种卡片的枚数
# 线段树2维护每种卡片数与对应分数的积
# 第3种查询, 先在线段树1上minLeft二分出left 使得 [left,k)内的卡片数之和恰好不超过x
# 加上这些卡片的分数,剩下 x-sum[left,k] 张卡片全部取 dk-1 分数的卡片
if __name__ == "__main__":
    ...
