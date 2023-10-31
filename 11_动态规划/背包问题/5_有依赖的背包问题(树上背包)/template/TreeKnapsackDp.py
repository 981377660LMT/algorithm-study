from typing import List, Tuple


# TODO:优化：上下界优化，二乘木dp
# https://ouuan.github.io/post/%E6%A0%91%E4%B8%8A%E8%83%8C%E5%8C%85%E7%9A%84%E4%B8%8A%E4%B8%8B%E7%95%8C%E4%BC%98%E5%8C%96/
def treeKnapsackDpNaive(
    tree: List[List[int]], items: List[Tuple[int, int]], maxCapacity: int, root=0
) -> int:
    def dfs(cur: int, pre: int) -> List[int]:
        curValue, curWeight = items[cur]
        dp = [0] * (maxCapacity + 1)
        for i in range(curWeight, maxCapacity + 1):
            dp[i] = curValue  # 根节点必须选
        for next_ in tree[cur]:
            if next_ == pre:
                continue
            ndp = dfs(next_, cur)
            # 类似分组背包，枚举分给子树 to 的容量 w，对应的子树的最大价值为 dt[w]
            # w 不可超过 j-it.w，否则无法选择根节点
            for j in range(maxCapacity, curWeight - 1, -1):
                for w in range(j - curWeight + 1):
                    dp[j] = max(dp[j], dp[j - w] + ndp[w])
        return dp

    res = dfs(root, -1)
    return res[maxCapacity]
