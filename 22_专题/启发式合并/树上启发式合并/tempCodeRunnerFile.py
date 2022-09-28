class Solution:
    def numberOfGoodPaths(self, vals: List[int], edges: List[List[int]]) -> int:
        def dfs(cur: int, pre: int) -> "SortedDict":
            """后序dfs返回子树内的每种结点个数 在当前结点处统计经过当前结点能产生多少条新路径"""
            self.res += 1
            curRes = SortedDict({vals[cur]: 1})

            for next in adjList[cur]:
                if next == pre:
                    continue
                nextRes = dfs(next, cur)
                it = iter(nextRes)
                for key in it:
                    # 我们枚举到的每个元素都会被删掉，所以总共最多只有 n 次删除，复杂度还是正确的。
                    if key < vals[cur]:
                        nextRes.pop(key)
                    else:
                        break
                # !merge
                if len(curRes) < len(nextRes):
                    curRes, nextRes = nextRes, curRes
                for key in nextRes:
                    self.res += curRes.get(key, 0) * nextRes[key]
                    curRes[key] = curRes.get(key, 0) + nextRes[key]

            return curRes

        def merge(curRes: "SortedDict", nextRes: "SortedDict") -> None:
            """启发式合并"""

        n = len(vals)
        adjList = [[] for _ in range(n)]
        for u, v in edges:
            adjList[u].append(v)
            adjList[v].append(u)

        self.res = 0
        dfs(0, -1)
        return self.res