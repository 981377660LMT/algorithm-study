# 在树上选取若干个点，这些点之间都是连通的
# 对v从0到n-1 输出`选取顶点v的组合`有多少种
# n<=1e5


"""メモ 20210521
参照
https://algo-logic.info/tree-dp/
https://ei1333.hateblo.jp/entry/2017/04/10/224413
https://qiita.com/keymoon/items/2a52f1b0fb7ef67fb89e
問題
https://atcoder.jp/contests/dp/tasks/dp_v
"""


import sys


class Rerooting:
    def __init__(self, n: int, decrement: int = 1):
        """nは頂点数 decrementはノードの番号が0-indexedの場合0でよい"""
        self.n = n
        self.adj = [[] for _ in range(n)]
        self.root = None  # 一番最初に根とする頂点
        self.decrement = decrement

    def add_edge(self, u: int, v: int):
        """辺を追加する u,vは元々の値で良い"""
        u -= self.decrement
        v -= self.decrement
        self.adj[u].append(v)
        self.adj[v].append(u)

    def rerooting(self, op, merge, e, root: int = 1) -> list:
        """
        <概要>
        1.rootを根としてまず一度木構造をbfsで求める 多くの場合rootは任意
        2.自身の部分木のdpの値をdp1に、自身を含まない兄弟のdpの値のmergeをdp2に入れる
          木構造が定まっていることからこれが効率的に求められる。 葉側からボトムアップに実行する
        3.任意の頂点を新たに根にしたとき、部分木は
          ①元の部分木 ②兄弟を親とした部分木 ③元の親を親とした(元の根の方向に伸びる)部分木の三つに分かれる。
          ①はstep2のdp1であり、かつdp2はstep3において、②から②と③をmergeした値へと更新されているので
          ②も③も分かっている。 根側からトップダウンに実行する(このことが上記の更新において重要)
        計算量 O(|V|) (|Vは頂点数)
        参照 https://qiita.com/keymoon/items/2a52f1b0fb7ef67fb89e
        """
        # step1
        root -= self.decrement
        assert 0 <= root < self.n
        self.root = root
        self.parent = [-1] * self.n  # 親の番号を記録
        self.order = [root]  # bfsの訪問順を記録 深さが広義単調増加している
        stack = [root]
        while stack:
            from_node = stack.pop()
            for to_node in self.adj[from_node]:
                if to_node == self.parent[from_node]:
                    continue
                self.parent[to_node] = from_node
                self.order.append(to_node)
                stack.append(to_node)
        # step2
        dp1 = [e] * self.n
        dp2 = [e] * self.n
        for from_node in self.order[::-1]:
            t = e
            for to_node in self.adj[from_node]:
                if self.parent[from_node] == to_node:
                    continue
                dp2[to_node] = t
                t = merge(t, op(dp1[to_node], from_node, to_node))
            t = e
            for to_node in self.adj[from_node][::-1]:
                if self.parent[from_node] == to_node:
                    continue
                dp2[to_node] = merge(t, dp2[to_node])
                t = merge(t, op(dp1[to_node], from_node, to_node))
            dp1[from_node] = t
        # step3
        for new_root in self.order[1:]:  # 元の根に関するdp1は既に求まっている
            par = self.parent[new_root]
            dp2[new_root] = op(merge(dp2[new_root], dp2[par]), new_root, par)
            dp1[new_root] = merge(dp1[new_root], dp2[new_root])
        return dp1


def op(a, u, v):
    # dpをmergeする前段階で実行する演算
    # 例:最も遠い点までの距離を求める場合 return a+1
    return a + 1


def merge(a, b):
    # モノイドの性質を満たす演算を定義する それが全方位木DPをする条件
    # 例:最も遠い点までの距離を求める場合 return max(a,b)
    return a * b % M


# mergeの単位元
# 例:最も遠い点までの距離を求める場合 e=0
e = 1

input = sys.stdin.readline
N, M = map(int, input().split())
T = Rerooting(N)
for _ in range(N - 1):
    x, y = map(int, input().split())
    T.add_edge(x, y)

dp = T.rerooting(op=op, merge=merge, e=e, root=1)
print(*dp, sep="\n")
