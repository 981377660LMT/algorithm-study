from typing import List
from collections import defaultdict
from LCA import LCA

# n<=1000


class Solution:
    def closestNode(self, n: int, edges: List[List[int]], query: List[List[int]]) -> List[int]:
        """nlogn"""
        adjMap = defaultdict(set)
        for u, v in edges:
            adjMap[u].add(v)
            adjMap[v].add(u)
        lca = LCA(n, adjMap, root=0)

        # 答案是最(靠下)深的LCA
        res = []
        for root1, root2, root3 in query:
            res.append(
                max(
                    lca.queryLCA(root1, root3),
                    lca.queryLCA(root2, root3),
                    lca.queryLCA(root1, root2),
                    key=lambda x: lca.depth[x],  # !按照深度排序
                )
            )
        return res


if __name__ == "__main__":
    # print(
    #     Solution().closestNode(
    #         n=7,
    #         edges=[[0, 1], [0, 2], [0, 3], [1, 4], [2, 5], [2, 6]],
    #         query=[[5, 3, 4], [5, 3, 6]],
    #     )
    # )

    # print(Solution().closestNode(n=1, edges=[], query=[[0, 0, 0]],))
    print(
        Solution().closestNode(
            n=5,
            edges=[[1, 0], [0, 3], [2, 4], [4, 3]],
            query=[[0, 0, 0], [0, 3, 2], [3, 0, 0], [4, 3, 1]],
        )
    )
