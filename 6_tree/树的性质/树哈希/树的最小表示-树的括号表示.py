# AcWing 157. 树形地铁系统（树的最小表示）
# 有n组数据，每个数据给两个树的dfs序，0表示向下走，1表示向上走，问这两个树是不是同构的。
# O(n^2)


# 如果两个字符串描述的探索路线可以视为同一个地铁系统的两种探索路线，则输出 same
from collections import deque


def check(s1: str, s2: str) -> bool:
    """两棵有根树是否同构(对应子树位置可以不同)

    len(s1) == len(s2) <= 3000

    如果两颗树是同构的，当且仅当 这两颗树的最小表示是相同的。
    可以对树的每个子树，都进行字典序的排序，逐层递归。
    最后如果两棵树同构的话，那么只有可能是排列上的不同，即可以通过最小表示，对每棵树生成唯一的最小表示。
    """

    def dfs(s: str, index: int) -> str:
        subtree = []
        depthDiff = 0
        for next in range(index + 1, len(s)):  # 找子树结点
            depthDiff += 1 if s[next] == "0" else -1
            if depthDiff == 1 and s[next] == "0":
                subtree.append(dfs(s, next))
            if depthDiff < 0:  # 回上面去了
                break

        subtree = deque(sorted(subtree))  # 最小表示
        subtree.appendleft("(")
        subtree.append(")")
        subtree.appendleft("#")  # 当前元素
        return "".join(subtree)

    return dfs(s1, 0) == dfs(s2, 0)


n = int(input())
for _ in range(n):
    # 加一个虚拟的根节点，从这个点出发dfs，比较的是`子树`的最小表示
    s1 = "0" + input() + "1"
    s2 = "0" + input() + "1"
    if check(s1, s2):
        print("same")
    else:
        print("different")
