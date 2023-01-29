from typing import List
from collections import defaultdict
from MinCostMaxFlow import MinCostMaxFlowDinic

# 注意网络流的流量守恒
# !如果某个点流入流出的流量不一致，那么就需要从源点加一条虚拟边来平衡流量
# 虚拟边流量为abs(入度-出度) 费用为0
# !出度大于入度的点看做源点，每个入度大于出度的点看做汇点，然后再添加虚拟源点和汇点，是多源多汇的
# !核心是费用流，多源多汇最后也是转换成单源单汇来解的


class Solution:
    def maximumRequests(self, n: int, requests: List[List[int]]) -> int:
        """请你从原请求列表中选出若干个请求，使得它们是一个可行的请求列表，并返回所有可行列表中最大请求数目。"""
        START, END = n + 1, n + 2
        flowDiff = defaultdict(int)
        mcmf = MinCostMaxFlowDinic(n + 3, START, END)
        for u, v in requests:
            mcmf.addEdge(u, v, 1, 1)
            flowDiff[u] -= 1
            flowDiff[v] += 1

        for key, count in flowDiff.items():
            if count > 0:
                mcmf.addEdge(key, END, count, 0)
            elif count < 0:
                mcmf.addEdge(START, key, -count, 0)

        # !因为虚拟源点和汇点是需要去掉的(这两个点流量不守恒) 去掉之后剩下的部分就成了题中需要求的循环自洽的整体
        # !剩下要最大 那么就要在原来图中删去经过边数最少的流 即原图的最小费用流
        return len(requests) - mcmf.work()[1]


print(Solution().maximumRequests(n=5, requests=[[0, 1], [1, 0], [0, 1], [1, 2], [2, 0], [3, 4]]))
print(
    Solution().maximumRequests(
        n=3, requests=[[1, 2], [1, 2], [2, 2], [0, 2], [2, 1], [1, 1], [1, 2]]
    )
)
# 4
