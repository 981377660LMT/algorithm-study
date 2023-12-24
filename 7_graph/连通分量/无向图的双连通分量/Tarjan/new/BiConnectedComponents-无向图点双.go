package main

func main() {

}

type Neighbor = struct{ to, id int }

// def findVBCC(n: int, graph: List[List[int]]) -> Tuple[List[List[int]], List[int], List[bool]]:
//     """
//     !Tarjan 算法求无向图的 v-BCC

//     Args:
//         n (int): 图的顶点数
//         graph (List[List[int]]):  邻接表

//     Returns:
//         Tuple[List[List[int]], List[int], List[bool]]:
//         每个 v-BCC 组里包含哪些点，每个点所在 v-BCC 的编号(从0开始)，每个顶点是否为割点(便于缩点成树)

//     Notes:
//         - 原图的割点`至少`在两个不同的 v-BCC 中
//         - 原图不是割点的点都`只存在`于一个 v-BCC 中
//         - v-BCC 形成的子图内没有割点
//     """

//     def dfs(cur: int, pre: int) -> int:
//         nonlocal dfsId, idCount
//         dfsId += 1
//         dfsOrder[cur] = dfsId
//         curLow = dfsId
//         childCount = 0
//         for _, next in enumerate(graph[cur]):
//             # edge = (cur, next, ei)
//             edge = (cur, next)
//             if dfsOrder[next] == 0:
//                 stack.append(edge)
//                 childCount += 1
//                 nextLow = dfs(next, cur)
//                 if nextLow >= dfsOrder[cur]:
//                     isCut[cur] = True
//                     idCount += 1
//                     group = []
//                     # eids = []
//                     while True:
//                         topEdge = stack.pop()
//                         v1, v2 = topEdge[0], topEdge[1]
//                         if vbccId[v1] != idCount:
//                             vbccId[v1] = idCount
//                             group.append(v1)
//                         if vbccId[v2] != idCount:
//                             vbccId[v2] = idCount
//                             group.append(v2)
//                         # eids.append(topEdge[2])
//                         if v1 == cur and v2 == next:
//                             break
//                     # 点数和边数相同，说明该 v-BCC 是一个简单环，且环上所有的边只属于一个简单环
//                     # if len(comp) == len(eids):
//                     #     for eid in eids:
//                     #         onSimpleCycle[eid] = True
//                     groups.append(group)
//                 if nextLow < curLow:
//                     curLow = nextLow
//             elif next != pre and dfsOrder[next] < dfsOrder[cur]:
//                 stack.append(edge)
//                 if dfsOrder[next] < curLow:
//                     curLow = dfsOrder[next]
//         if pre == -1 and childCount == 1:
//             isCut[cur] = False
//         return curLow

//     dfsId = 0
//     dfsOrder = [0] * n
//     vbccId = [0] * n
//     idCount = 0
//     isCut = [False] * n
//     stack = []  # (u, v, eid)
//     groups = []

//     for i, order in enumerate(dfsOrder):
//         if order == 0:
//             if len(graph[i]) == 0:  # 零度，即孤立点（isolated vertex）
//                 idCount += 1
//                 vbccId[i] = idCount
//                 groups.append([i])
//                 continue
//             dfs(i, -1)

//     return groups, [v - 1 for v in vbccId], isCut

// 求无向图的点双连通分量.
// 返回值为 (每个双连通分量包含哪些点, 每个点所属的双连通分量编号, 每个点是否为割点)
func BiConnectedComponent(n int, graph [][]int) (groups [][]int, belong []int, isCut []bool) {
	dfsId := 0
	dfsOrder := make([]int, n)
	vbccId := make([]int, n)
	idCount := 0
	isCut = make([]bool, n)
	stack := [][3]int{}

}

// 边双缩点成树.
// !BCC 和割点作为新图中的节点，并在每个割点与包含它的所有 BCC 之间连边
// !bcc1 - 割点1 - bcc2 - 割点2 - ...
func ToTree(graph [][]int, groups [][]int, isCut []bool) (tree [][]int) {

}
