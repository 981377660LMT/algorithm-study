"""
换根dp框架

op是相邻结点转移时,fromRes如何变化
merge是如何合并两个子节点的res
e是每个节点res的初始值

框架传入op和merge看似只求根节点0的值,实际上求出了每个点的dp值
"""


from typing import Callable, List

Op = Callable[[int, int, int, int], int]
Merge = Callable[[int, int], int]
E = Callable[[int], int]


class Rerooting:
    """https://atcoder.jp/contests/dp/submissions/22766939"""

    __slots__ = ("adjList", "_n", "_decrement", "_root", "_parent", "_order")

    def __init__(self, n: int, decrement: int = 0):
        """
        n: 頂点数
        decrement: 頂点の値を減らす量 (1-indexedなら1, 0-indexedなら0)
        """
        self.adjList = [[] for _ in range(n)]
        self._n = n
        self._decrement = decrement
        self._root = None  # 一番最初に根とする頂点

    def addEdge(self, u: int, v: int):
        """
        u,v 間に無向辺を張る (u->v, v->u)
        """
        u -= self._decrement
        v -= self._decrement
        self.adjList[u].append(v)
        self.adjList[v].append(u)

    def rerooting(self, op: Op, merge: Merge, e: E, root: int = 0) -> List[int]:
        """
        - op: 頂点の値を更新する関数

          (fromRes,parent,cur,direction) -> newRes
          direction: 0表示用cur更新parent的dp1,1表示用parent更新cur的dp2

          dpをmergeする前段階で実行する演算
          例:最も遠い点までの距離を求める場合 return fromRes+1

        - merge: 子の値を親にマージする関数

          (childRes1,childRes2) -> newRes

          モノイドの性質を満たす演算を定義する それが全方位木DPをする条件
          例:最も遠い点までの距離を求める場合 return max(childRes1,childRes2)

        - e: 単位元
          (root) -> res

          mergeの単位元
          例:最も遠い点までの距離を求める場合 e=0

        - root: 根とする頂点

        <概要>
        1. rootを根としてまず一度木構造をbfsで求める 多くの場合rootは任意 (0)
        2. 自身の部分木のdpの値をdp1に、自身を含まない兄弟のdpの値のmergeをdp2に入れる
          木構造が定まっていることからこれが効率的に求められる。 葉側からボトムアップに実行する
        3. 任意の頂点を新たに根にしたとき、部分木は
          ①元の部分木 ②兄弟を親とした部分木 ③元の親を親とした(元の根の方向に伸びる)部分木の三つに分かれる。
          ①はstep2のdp1であり、かつdp2はstep3において、②から②と③をmergeした値へと更新されているので
          ②も③も分かっている。 根側からトップダウンに実行する(このことが上記の更新において重要)

        計算量 O(|V|) (Vは頂点数)
        参照 https://qiita.com/keymoon/items/2a52f1b0fb7ef67fb89e
        """
        # step1
        root -= self._decrement
        assert 0 <= root < self._n
        self._root = root
        self._parent = [-1] * self._n  # 親の番号を記録
        self._order = [root]  # bfsの訪問順を記録 深さが広義単調増加している
        stack = [root]
        while stack:
            cur = stack.pop()
            for next in self.adjList[cur]:
                if next == self._parent[cur]:
                    continue
                self._parent[next] = cur
                self._order.append(next)
                stack.append(next)

        # step2
        dp1 = [e(i) for i in range(self._n)]  # !子树部分的dp值
        dp2 = [e(i) for i in range(self._n)]  # !非子树部分的dp值
        for cur in self._order[::-1]:  # 从下往上拓扑序dp
            res = e(cur)
            for next in self.adjList[cur]:
                if self._parent[cur] == next:
                    continue
                dp2[next] = res
                res = merge(res, op(dp1[next], cur, next, 0))  # op从下往上更新dp1
            res = e(cur)
            for next in self.adjList[cur][::-1]:
                if self._parent[cur] == next:
                    continue
                dp2[next] = merge(res, dp2[next])
                res = merge(res, op(dp1[next], cur, next, 0))
            dp1[cur] = res

        # step3
        for newRoot in self._order[1:]:  # 元の根に関するdp1は既に求まっている
            parent = self._parent[newRoot]
            dp2[newRoot] = op(merge(dp2[newRoot], dp2[parent]), parent, newRoot, 1)  # op从上往下更新dp2
            dp1[newRoot] = merge(dp1[newRoot], dp2[newRoot])
        return dp1


if __name__ == "__main__":
    # 310-求树上每个节点到其他节点的最远距离

    def op(fromRes: int, parent: int, cur: int, direction: int) -> int:
        # dpをmergeする前段階で実行する演算
        # 例:最も遠い点までの距離を求める場合 return res+1
        return fromRes + 1

    def merge(childRes1: int, childRes2: int) -> int:
        # モノイドの性質を満たす演算を定義する それが全方位木DPをする条件
        # 例:最も遠い点までの距離を求める場合 return max(childRes1,childRes2)
        return max(childRes1, childRes2)

    def e(root: int) -> int:
        # mergeの単位元
        # 例:最も遠い点までの距離を求める場合 e=0
        return 0

    class Solution:
        def findMinHeightTrees(self, n: int, edges: List[List[int]]) -> List[int]:
            R = Rerooting(n)
            for u, v in edges:
                R.addEdge(u, v)
            maxDist = R.rerooting(op=op, merge=merge, e=e, root=0)
            min_ = min(maxDist)
            return [i for i in range(n) if maxDist[i] == min_]
