# 序列生成需要一个基础序列A，
# 这个序列由n个 不大于100的数字组成，同时给定镜像复制次数m。

# 然后对于A进行m次镜像复制，例如序列A={1，2，3}，
# 则一次镜像复制后得到的序列是{1，2，3，3，2，1}，
# 两次镜像复制得到的序列是B={1，2，3，3，2，1，1，2，3，3，2，1} 。
# 现在给出你生成一个序列所需要的参数，请你计算该序列的第k位是多少。
# 镜像序列的第k位
# !1<=n<=100,1<=m<=60,1<=k<=1e18

# !求出mid后线段树
n, m, k = map(int, input().split())
nums = list(map(int, input().split()))


def dfs(depth: int, cur: int) -> str:
    """复制depth次后的第cur位是多少"""
    if depth == 0:
        return str(nums[cur - 1])
    length = n * 2**depth
    mid = length // 2
    if cur <= mid:
        return dfs(depth - 1, cur)
    return dfs(depth - 1, length - cur + 1)


print(dfs(m, k))
