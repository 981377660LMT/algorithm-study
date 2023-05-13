"""破解保险箱"""


class Solution:
    def crackSafe(self, n: int, k: int) -> str:
        """密码是 n 位数, 密码的每一位是 k 位序列 0, 1, ..., k-1 中的一个 。

        求出一个最短的字符串，使其包含 0~k^n (k进制)中的所有数字
        将所有的 n-1 位数作为节点。每个节点有 k 条边

        如果我们从任一节点出发，能够找出一条路径，
        经过图中的所有边且只经过一次，然后把边上的数字写入字符串（还需加入起始节点的数字），那么这个字符串显然符合要求
        """

        def dfs(cur: int) -> None:
            """非常像求树的欧拉路径"""
            for i in range(k):
                next = cur * 10 + i
                if next not in visited:
                    visited.add(next)
                    dfs(next % max_)
                    res.append(str(i))

        visited = set()
        res = []
        max_ = 10 ** (n - 1)

        dfs(0)
        return "".join(res) + "0" * (n - 1)
