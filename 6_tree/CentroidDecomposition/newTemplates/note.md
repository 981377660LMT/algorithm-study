1. 基于顶点的重心分解(点分治)，每个顶点恰好作为子树的重心一次，并可以获得每个子树对应的节点.
   所有子树的节点数之和为 O(nlogn)
   [CentroidDecomposition1](CentroidDecomposition0.go)
2. 1/3 重心分解(边分治)，每次将树分成颜色为蓝、红的两部分。
   类似于点分治，不过钦定的并非经过子树重心而是一条使得子树尽量平衡的边
   [CentroidDecomposition1](CentroidDecomposition1.go)
3. 三度化边分治
   一方で 1/3 重心分解では，木の分割の際に頂点が分裂してしまうため，ひとつの頂点 v
   に注目したときに v が現れるノードの個数が抑えられないという欠点があります（例えばスターグラフを考えると，`ひとつの頂点がすべてのノードに現れます`）．
   このことは，`「始点を固定したパス全体」`を扱う際に障害となります

   [CentroidDecomposition2](CentroidDecomposition2.go)
   **始点ごとのパスを扱う問題でも重心分解と同様の計算量が達成できるようになります**

---

https://maspypy.com/%E9%87%8D%E5%BF%83%E5%88%86%E8%A7%A3%E3%83%BB1-3%E9%87%8D%E5%BF%83%E5%88%86%E8%A7%A3%E3%81%AE%E3%81%8A%E7%B5%B5%E6%8F%8F%E3%81%8D

各ノードの根同士を結んで得られる木のことを指して重心分解と呼ぶ流儀もあるようです．筆者の経験上，競技プログラミングにおいては，各ノードとして現れる部分木のリストのみが重要であることが多く，重心を結んで得られる木の構造（例えば各頂点の親がどの頂点であるかといった情報）そのものは不要なことが多いです．ただし，この木の構築を題材とした出題も稀にあります（https://atcoder.jp/contests/abc291/tasks/abc291_h）

`重心分解通常关注的是通过分解得到的各个子树的列表，而不是这些重心之间构成的树的具体结构`（例如，每个顶点的父节点是哪个顶点）。这是因为在很多问题中，我们更关心如何通过分解树来简化问题，而不是分解后各个重心之间的具体连接方式。

重心分解の競技プログラミングにおける主要な用途は，「木上のすべてのパスに対する何かを計算せよ」という種類のものです．
重心分解在競技编程中的主要用途是解决“计算树上所有路径上的某些属性”的问题。

あるノードの表す木おいて，パスを「根（重心）を通るか否か」で場合分けして考えましょう．「根を通らないパス」は，「子ノードの表す木のパス」であると言い換えられます．したがって各ノードにおいては，根を通るパスの計算に専念すればよくなります．$n$ 頂点のノードでこの計算が $O(n)$ 時間で行えるなら全体で $O(N\log N)$ 時間，$O(n\log n)$ 時間で行えるなら全体で $O(N\log^2N)$ 時間などとなります．

在树的重心分解中，我们可以将路径分为两类：“经过根（重心）的路径”和“不经过根的路径”。这种分类方法可以简化问题的求解过程。具体来说：
经过根的路径：这些路径必须经过重心。我们只需要专注于计算这些路径。
不经过根的路径：这些路径完全位于某个子树中。因此，我们可以递归地在每个子树中计算这些路径。
通过这种分类，我们可以将复杂的问题分解为更小的子问题，从而简化计算过程。假设在一个包含$n$个顶点的子树中，我们可以在$O(n)$时间内计算经过根的路径，那么整个树的计算复杂度可以达到$O(N \log N)$。如果计算复杂度是$O(n \log n)$，那么整体复杂度将是$O(N \log^2 N)$。

結局重心分解を行うと，$n$ 頂点の根付き木に対して，「根を通るパス」全体に対する計算を $O(n)$ 時間などで行う問題に帰着できることになります．
最后，如果我们分解重心，我们可以归结为计算顶点为 $n$ 的根树的整个“通过根的路径”的问题，例如 $O（n）$ 时间。

---

各部分木の頂点に何らかの色をつける（最初の重心分解の図を参照）ことにすると，
如果我们决定给每个部分树的顶点赋予一些颜色（参见质心分解的第一个图），我们可以得到

2 頂点の組 $u,v$ であって色が異なるものすべてに対する $f(a_u,b_v)$ に対して何かの処理をせよ．
对所有两个顶点对 $u，v$ 的不同颜色的 $f（a_u，b_v）$ 做一些事情。

---

1/3 重心分解

これによって「すべてのパスに対する計算」を「赤色の頂点 $u$, 青色の頂点 $v$ に対する $f(a_u,b_v)$ の計算」に変換可能になります．
这样就可以**将“所有路径的计算”转换为“红色顶点 $u$ 和蓝色顶点 $v$ 的 $f（a_u，b_v）$ 的计算”。**
