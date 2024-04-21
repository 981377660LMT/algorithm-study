# from titan_pylib.data_structures.dynamic_connectivity.euler_tour_tree import EulerTourTree
from typing import (
    Generator,
    Generic,
    TypeVar,
    Callable,
    Iterable,
    Optional,
    Union,
    Tuple,
    List,
    Dict,
)
from types import GeneratorType

T = TypeVar("T")
F = TypeVar("F")


class EulerTourTree(Generic[T, F]):
    """``Euler Tour Tree`` です。部分木クエリの強さに定評があります。

    森です。各連結成分(木)は独立なのでそれごとに見ていきます。

    基本戦術はオイラーツアー(Euler tour technique)です。これはググってください。
    実は、ある木をオイラーツアーした列と、その木に対して 根の変更 / 辺の追加 / 辺の削除 をした木のオイラーツアーした列との差分は多くはありません。実際に紙に書くとよいでしょう。これをうまいこと管理します。

    列を平衡二分木で管理します。ここで、オイラーツアーの列はでは頂点ではなく辺の列とします。例えば、
    `[0, 1, 0, 2, 0]`
    の頂点列は辺列
    `[{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 2}, {2, 2}, {2, 0}]`
    などとなります。
    また、補助データ構造として辺からノードのポインタをたどれる辞書を保持します。
    各処理の流れは以下のようになります。

    - 根の変更: `reroot(v)`
      - 辺 `{v, v}` の頂点の直前で `split` し、それを順に `A, B` とする
      - `B` と `A` をこの順にマージする
    - 辺の追加: `link(u, v)`
      - `reroot(u); reroot(v)`
      - `u, v` が属する木(のオイラーツアーした辺の列)をそれぞれ `E1, E2` とする
      - `E1, [{u, v}], E2, [{v, u}]` をこの順にマージする
    - 辺の削除: `cut(u, v)`
      - `reroot(v); reroot(u)`
      - `{u, v}, {v, u}` で `split` してできたものを順に `A, B, C` とする。ただし、 `{u, v}` は `A` に含まれ、 `{v, u}` は `C` に含まれる。
      - `A` の末尾と `C` の先頭を削除し、 `A` と `C` をこの順にマージする
    - 連結性判定: `same(u, v)`
      - `u, v` を `splay` して、 `u` の親が `None` じゃなければ連結。 `u == v` のときは別途処理をする。

    計算量は、オイラーツアーの管理と辺→ノードの管理に赤黒木を使えば最悪 `O(logN)` です。この実装ではsplay木とハッシュテーブルで管理するので償却 `O(logN)` +期待 `O(1)` です(もとのグラフの頂点数に対してオイラーツアーすると `(元の頂点数)+(辺数)*2?` の頂点ができるので、たしかに `O(元のグラフの頂点数)` ではありますが、定数倍がバカです)。
    オイラーツアーがあれば、部分木クエリは簡単に処理できます。
    """

    class _Node:
        def __init__(self, key: T, lazy: F):
            self.key: T = key
            self.data: T = key
            self.lazy: F = lazy
            self.par: Optional[EulerTourTree._Node] = None
            self.left: Optional[EulerTourTree._Node] = None
            self.right: Optional[EulerTourTree._Node] = None

        def __str__(self):
            if self.left is None and self.right is None:
                return f"(key,par):{self.key,self.data,self.lazy,(self.par.key if self.par else None)}\n"
            return f"(key,par):{self.key,self.data,self.lazy,(self.par.key if self.par else None)},\n left:{self.left},\n right:{self.right}\n"

        __repr__ = __str__

    def __init__(
        self,
        n_or_a: Union[int, Iterable[T]],
        op: Callable[[T, T], T],
        mapping: Callable[[F, T], T],
        composition: Callable[[F, F], F],
        e: T,
        id: F,
    ) -> None:
        self.op = op
        self.mapping = mapping
        self.composition = composition
        self.e = e
        self.id = id
        a = [e for _ in range(n_or_a)] if isinstance(n_or_a, int) else list(n_or_a)
        self.n: int = len(a)
        self.ptr_vertex: List[EulerTourTree._Node] = [
            EulerTourTree._Node(elem, id) for i, elem in enumerate(a)
        ]
        self.ptr_edge: Dict[Tuple[int, int], EulerTourTree._Node] = {}
        self._group_numbers: int = self.n

    @staticmethod
    def antirec(func, stack=[]):
        # 参考: https://github.com/cheran-senthil/PyRival/blob/master/pyrival/misc/bootstrap.py
        def wrappedfunc(*args, **kwargs):
            if stack:
                return func(*args, **kwargs)
            to = func(*args, **kwargs)
            while True:
                if isinstance(to, GeneratorType):
                    stack.append(to)
                    to = next(to)
                else:
                    stack.pop()
                    if not stack:
                        break
                    to = stack[-1].send(to)
            return to

        return wrappedfunc

    def build(self, G: List[List[int]]) -> None:
        """隣接リスト ``G`` をもとにして、辺を張ります。

        :math:`O(n)` です。

        Args:
          G (List[List[int]]): 隣接リストです。

        Note:
          ``build`` メソッドを使用する場合は他のメソッドより前に使用しなければなりません。
        """
        n, ptr_vertex, ptr_edge, e, id = self.n, self.ptr_vertex, self.ptr_edge, self.e, self.id
        seen = [0] * n
        _Node = EulerTourTree._Node

        @EulerTourTree.antirec
        def dfs(v: int, p: int = -1) -> Generator:
            a.append(v * n + v)
            for x in G[v]:
                if x == p:
                    continue
                a.append(v * n + x)
                yield dfs(x, v)
                a.append(x * n + v)
            yield

        @EulerTourTree.antirec
        def rec(l: int, r: int) -> Generator:
            mid = (l + r) >> 1
            u, v = divmod(a[mid], n)
            node = ptr_vertex[u] if u == v else _Node(e, id)
            if u == v:
                seen[u] = 1
            else:
                ptr_edge[u * n + v] = node
            if l != mid:
                node.left = yield rec(l, mid)
                node.left.par = node
            if mid + 1 != r:
                node.right = yield rec(mid + 1, r)
                node.right.par = node
            self._update(node)
            yield node

        for root in range(self.n):
            if seen[root]:
                continue
            a: List[int] = []
            dfs(root)
            rec(0, len(a))

    def _popleft(self, v: _Node) -> Optional[_Node]:
        v = self._left_splay(v)
        if v.right:
            v.right.par = None
        return v.right

    def _pop(self, v: _Node) -> Optional[_Node]:
        v = self._right_splay(v)
        if v.left:
            v.left.par = None
        return v.left

    def _split_left(self, v: _Node) -> Tuple[_Node, Optional[_Node]]:
        # x, yに分割する。ただし、xはvを含む
        self._splay(v)
        x, y = v, v.right
        if y:
            y.par = None
        x.right = None
        self._update(x)
        return x, y

    def _split_right(self, v: _Node) -> Tuple[Optional[_Node], _Node]:
        # x, yに分割する。ただし、yはvを含む
        self._splay(v)
        x, y = v.left, v
        if x:
            x.par = None
        y.left = None
        self._update(y)
        return x, y

    def _merge(self, u: Optional[_Node], v: Optional[_Node]) -> None:
        if u is None or v is None:
            return
        u = self._right_splay(u)
        self._splay(v)
        u.right = v
        v.par = u
        self._update(u)

    def _splay(self, node: _Node) -> None:
        self._propagate(node)
        while node.par is not None and node.par.par is not None:
            pnode = node.par
            gnode = pnode.par
            self._propagate(gnode)
            self._propagate(pnode)
            self._propagate(node)
            node.par = gnode.par
            if (gnode.left is pnode) == (pnode.left is node):
                if pnode.left is node:
                    tmp1 = node.right
                    pnode.left = tmp1
                    node.right = pnode
                    pnode.par = node
                    tmp2 = pnode.right
                    gnode.left = tmp2
                    pnode.right = gnode
                    gnode.par = pnode
                else:
                    tmp1 = node.left
                    pnode.right = tmp1
                    node.left = pnode
                    pnode.par = node
                    tmp2 = pnode.left
                    gnode.right = tmp2
                    pnode.left = gnode
                    gnode.par = pnode
                if tmp1:
                    tmp1.par = pnode
                if tmp2:
                    tmp2.par = gnode
            else:
                if pnode.left is node:
                    tmp1 = node.right
                    pnode.left = tmp1
                    node.right = pnode
                    tmp2 = node.left
                    gnode.right = tmp2
                    node.left = gnode
                    pnode.par = node
                    gnode.par = node
                else:
                    tmp1 = node.left
                    pnode.right = tmp1
                    node.left = pnode
                    tmp2 = node.right
                    gnode.left = tmp2
                    node.right = gnode
                    pnode.par = node
                    gnode.par = node
                if tmp1:
                    tmp1.par = pnode
                if tmp2:
                    tmp2.par = gnode
            self._update(gnode)
            self._update(pnode)
            self._update(node)
            if node.par is None:
                return
            if node.par.left is gnode:
                node.par.left = node
            else:
                node.par.right = node
        if node.par is None:
            return
        pnode = node.par
        self._propagate(pnode)
        self._propagate(node)
        if pnode.left is node:
            pnode.left = node.right
            if pnode.left:
                pnode.left.par = pnode
            node.right = pnode
        else:
            pnode.right = node.left
            if pnode.right:
                pnode.right.par = pnode
            node.left = pnode
        node.par = None
        pnode.par = node
        self._update(pnode)
        self._update(node)

    def _left_splay(self, node: _Node) -> _Node:
        self._splay(node)
        while node.left is not None:
            node = node.left
        self._splay(node)
        return node

    def _right_splay(self, node: _Node) -> _Node:
        self._splay(node)
        while node.right is not None:
            node = node.right
        self._splay(node)
        return node

    def _propagate(self, node: Optional[_Node]) -> None:
        if node is None or node.lazy == self.id:
            return
        if node.left:
            node.left.key = self.mapping(node.lazy, node.left.key)
            node.left.data = self.mapping(node.lazy, node.left.data)
            node.left.lazy = self.composition(node.lazy, node.left.lazy)
        if node.right:
            node.right.key = self.mapping(node.lazy, node.right.key)
            node.right.data = self.mapping(node.lazy, node.right.data)
            node.right.lazy = self.composition(node.lazy, node.right.lazy)
        node.lazy = self.id

    def _update(self, node: _Node) -> None:
        self._propagate(node.left)
        self._propagate(node.right)
        node.data = node.key
        if node.left:
            node.data = self.op(node.left.data, node.data)
        if node.right:
            node.data = self.op(node.data, node.right.data)

    def link(self, u: int, v: int) -> None:
        """辺 ``{u, v}`` を追加します。
        :math:`O(\\log{n})` です。

        Note:
          ``u`` と ``v`` が同じ連結成分であってはいけません。
        """
        # add edge{u, v}
        self.reroot(u)
        self.reroot(v)
        assert u * self.n + v not in self.ptr_edge, f"EulerTourTree.link(), {(u, v)} in ptr_edge"
        assert v * self.n + u not in self.ptr_edge, f"EulerTourTree.link(), {(v, u)} in ptr_edge"
        uv_node = EulerTourTree._Node(self.e, self.id)
        vu_node = EulerTourTree._Node(self.e, self.id)
        self.ptr_edge[u * self.n + v] = uv_node
        self.ptr_edge[v * self.n + u] = vu_node
        u_node = self.ptr_vertex[u]
        v_node = self.ptr_vertex[v]
        self._merge(u_node, uv_node)
        self._merge(uv_node, v_node)
        self._merge(v_node, vu_node)
        self._group_numbers -= 1

    def cut(self, u: int, v: int) -> None:
        """辺 ``{u, v}`` を削除します。
        :math:`O(\\log{n})` です。

        Note:
          辺 ``{u, v}`` が存在してなければいけません。
        """
        # erace edge{u, v}
        self.reroot(v)
        self.reroot(u)
        assert u * self.n + v in self.ptr_edge, f"EulerTourTree.cut(), {(u, v)} not in ptr_edge"
        assert v * self.n + u in self.ptr_edge, f"EulerTourTree.cut(), {(v, u)} not in ptr_edge"
        uv_node = self.ptr_edge.pop(u * self.n + v)
        vu_node = self.ptr_edge.pop(v * self.n + u)
        a, _ = self._split_left(uv_node)
        _, c = self._split_right(vu_node)
        a = self._pop(a)
        c = self._popleft(c)
        self._merge(a, c)
        self._group_numbers += 1

    def merge(self, u: int, v: int) -> bool:
        """
        頂点 ``u`` と ``v`` が同じ連結成分にいる場合はなにもせず ``False`` を返します。
        そうでない場合は辺 ``{u, v}`` を追加し ``True`` を返します。
        :math:`O(\\log{n})` です。
        """
        if self.same(u, v):
            return False
        self.link(u, v)
        return True

    def split(self, u: int, v: int) -> bool:
        """
        辺 ``{u, v}`` が存在しない場合はなにもせず ``False`` を返します。
        そうでない場合は辺 ``{u, v}`` を削除し ``True`` を返します。
        :math:`O(\\log{n})` です。
        """
        if u * self.n + v not in self.ptr_edge or v * self.n + v not in self.ptr_edge:
            return False
        self.cut(u, v)
        return True

    def leader(self, v: int) -> _Node:
        """頂点 ``v`` を含む木の代表元を返します。
        :math:`O(\\log{n})` です。

        Note:
          ``reroot`` すると変わるので注意です。
        """
        # vを含む木の代表元
        # rerootすると変わるので注意
        return self._left_splay(self.ptr_vertex[v])

    def reroot(self, v: int) -> None:
        """頂点 ``v`` を含む木の根を ``v`` にします。

        :math:`O(\\log{n})` です。
        """
        node = self.ptr_vertex[v]
        x, y = self._split_right(node)
        self._merge(y, x)
        self._splay(node)

    def same(self, u: int, v: int) -> bool:
        """
        頂点 ``u`` と ``v`` が同じ連結成分にいれば ``True`` を、
        そうでなければ ``False`` を返します。

        :math:`O(\\log{n})` です。
        """
        u_node = self.ptr_vertex[u]
        v_node = self.ptr_vertex[v]
        self._splay(u_node)
        self._splay(v_node)
        return u_node.par is not None or u_node is v_node

    def _show(self) -> None:
        # for debug
        print("+++++++++++++++++++++++++++")
        for i, v in enumerate(self.ptr_vertex):
            print((i, i), v, end="\n\n")
        for k, v in self.ptr_edge.items():
            print(k, v, end="\n\n")
        print("+++++++++++++++++++++++++++")

    def subtree_apply(self, v: int, p: int, f: F) -> None:
        """頂点 ``v`` を根としたときの部分木に ``f`` を作用します。

        ``v`` の親は ``p`` です。
        ``v`` の親が存在しないときは ``p=-1`` として下さい。

        :math:`O(\\log{n})` です。

        Args:
          v (int): 根です。
          p (int): ``v`` の親です。
          f (F): 作用素です。
        """
        if p == -1:
            v_node = self.ptr_vertex[v]
            self._splay(v_node)
            v_node.key = self.mapping(f, v_node.key)
            v_node.data = self.mapping(f, v_node.data)
            v_node.lazy = self.composition(f, v_node.lazy)
            return
        self.reroot(v)
        self.reroot(p)
        assert (
            p * self.n + v in self.ptr_edge
        ), f"EulerTourTree.subtree_apply(), {(p, v)} not in ptr_edge"
        assert (
            v * self.n + p in self.ptr_edge
        ), f"EulerTourTree.subtree_apply(), {(v, p)} not in ptr_edge"
        v_node = self.ptr_vertex[v]
        a, b = self._split_right(self.ptr_edge[p * self.n + v])
        b, d = self._split_left(self.ptr_edge[v * self.n + p])
        self._splay(v_node)
        v_node.key = self.mapping(f, v_node.key)
        v_node.data = self.mapping(f, v_node.data)
        v_node.lazy = self.composition(f, v_node.lazy)
        self._propagate(v_node)
        self._merge(a, b)
        self._merge(b, d)

    def subtree_sum(self, v: int, p: int) -> T:
        """頂点 ``v`` を根としたときの部分木の総和を返します。

        ``v`` の親は ``p`` です。
        ``v`` の親が存在しないときは ``p=-1`` として下さい。

        :math:`O(\\log{n})` です。

        Args:
          v (int): 根です。
          p (int): ``v`` の親です。
        """
        if p == -1:
            v_node = self.ptr_vertex[v]
            self._splay(v_node)
            return v_node.data
        self.reroot(v)
        self.reroot(p)
        assert (
            p * self.n + v in self.ptr_edge
        ), f"EulerTourTree.subtree_sum(), {(p, v)} not in ptr_edge"
        assert (
            v * self.n + p in self.ptr_edge
        ), f"EulerTourTree.subtree_sum(), {(v, p)} not in ptr_edge"
        v_node = self.ptr_vertex[v]
        a, b = self._split_right(self.ptr_edge[p * self.n + v])
        b, d = self._split_left(self.ptr_edge[v * self.n + p])
        self._splay(v_node)
        res = v_node.data
        self._merge(a, b)
        self._merge(b, d)
        return res

    def group_count(self) -> int:
        """連結成分の個数を返します。
        :math:`O(1)` です。
        """
        return self._group_numbers

    def get_vertex(self, v: int) -> T:
        """頂点 ``v`` の ``key`` を返します。
        :math:`O(\\log{n})` です。
        """
        node = self.ptr_vertex[v]
        self._splay(node)
        return node.key

    def set_vertex(self, v: int, val: T) -> None:
        """頂点 ``v`` の ``key`` を ``val`` に更新します。
        :math:`O(\\log{n})` です。
        """
        node = self.ptr_vertex[v]
        self._splay(node)
        node.key = val
        self._update(node)

    def __getitem__(self, v: int) -> T:
        return self.get_vertex(v)

    def __setitem__(self, v: int, val: T) -> None:
        return self.set_vertex(v, val)
