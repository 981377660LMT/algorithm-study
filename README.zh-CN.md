<div align='center'>
  <img src='./logo.svg' width='300' alt="Algorithm"/>
  <p>
    算法笔记与模板
  </p>
</div>

<div align='center'>
  <a href='./README.md'>English</a> | 中文
</div>

[![Build Status](https://app.travis-ci.com/981377660LMT/algorithm-study.svg?branch=master)](https://app.travis-ci.com/981377660LMT/algorithm-study)

---

## 📖 模板

0. **字符串**

   - [子序列匹配](17_%E6%A8%A1%E5%BC%8F%E5%8C%B9%E9%85%8D/isSubsequence.py)
   - [Rabin-Karp 字符串哈希](18_%E5%93%88%E5%B8%8C/%E5%AD%97%E7%AC%A6%E4%B8%B2%E5%93%88%E5%B8%8C/StringHasher.py)
   - [最小表示法](17_%E6%A8%A1%E5%BC%8F%E5%8C%B9%E9%85%8D/%E6%9C%80%E5%B0%8F%E8%A1%A8%E7%A4%BA%E6%B3%95/%E6%9C%80%E5%B0%8F%E8%A1%A8%E7%A4%BA%E6%B3%95.py)
   - [KMP 算法](17_%E6%A8%A1%E5%BC%8F%E5%8C%B9%E9%85%8D/kmp/kmp.py)
   - [Z 算法 (扩展 KMP) ](17_%E6%A8%A1%E5%BC%8F%E5%8C%B9%E9%85%8D/kmp/z%E5%87%BD%E6%95%B0-%E6%89%A9%E5%B1%95kmp.py)
   - [马拉车 (Manacher) 算法](17_%E6%A8%A1%E5%BC%8F%E5%8C%B9%E9%85%8D/%E9%A9%AC%E6%8B%89%E8%BD%A6%E6%8B%89%E9%A9%AC/Manacher.py)
   - [后缀数组](17_%E6%A8%A1%E5%BC%8F%E5%8C%B9%E9%85%8D/%E5%90%8E%E7%BC%80%E6%95%B0%E7%BB%84/SA.py)

1. **栈**

   - [单调栈](1_stack/%E5%8D%95%E8%B0%83%E6%A0%88/%E6%AF%8F%E4%B8%AA%E5%85%83%E7%B4%A0%E4%BD%9C%E4%B8%BA%E6%9C%80%E5%80%BC%E7%9A%84%E5%BD%B1%E5%93%8D%E8%8C%83%E5%9B%B4.py)

2. **队列**

   - [队列](2_queue/Deque/Queue.ts)
   - [优先队列 (堆)](2_queue/PriorityQueue.ts)
   - 双端队列
     - [ArrayDeque](2_queue/Deque/ArrayDeque.ts)
     - [LinkedList](3_linkedList/LinkedList.ts)
   - [单调队列](2_queue/%E5%8D%95%E8%B0%83%E9%98%9F%E5%88%97Monoqueue/MonoQueue.py)

3. **链表**

   - [LinkedListNode](3_linkedList/LinkedListNode.py)
   - [LinkedList](3_linkedList/LinkedList.ts)

4. **树**

   - [DFS 序](6_tree/%E6%A0%91%E7%9A%84%E6%80%A7%E8%B4%A8/dfs%E5%BA%8F/DFSOrder.py)
   - [Trie (字典树)](6_tree/%E5%89%8D%E7%BC%80%E6%A0%91trie/Trie.py)
   - [01 Trie](6_tree/%E5%89%8D%E7%BC%80%E6%A0%91trie/%E6%9C%80%E5%A4%A7%E5%BC%82%E6%88%96%E5%89%8D%E7%BC%80%E6%A0%91/XORTrie.py)
   - [树状数组](6_tree/%E6%A0%91%E7%8A%B6%E6%95%B0%E7%BB%84/%E7%BB%8F%E5%85%B8%E9%A2%98/BIT.py)
   - [倍增法求 LCA](6_tree/LCA%E9%97%AE%E9%A2%98/TreeManager.py)
   - [线段树](6_tree/%E7%BA%BF%E6%AE%B5%E6%A0%91/template)
   - [持久化线段树](6_tree/%E5%8F%AF%E6%8C%81%E4%B9%85%E5%8C%96%E7%BA%BF%E6%AE%B5%E6%A0%91/255.%20%E7%AC%ACK%E5%B0%8F%E6%95%B0-%E6%9F%A5%E8%AF%A2%E5%8C%BA%E9%97%B4%E7%AC%ACk%E5%B0%8F%E6%95%B0.ts)
   - [树链剖分](24_%E9%AB%98%E7%BA%A7%E6%95%B0%E6%8D%AE%E7%BB%93%E6%9E%84/%E6%A0%91%E9%93%BE%E5%89%96%E5%88%86/2568.%20%E6%A0%91%E9%93%BE%E5%89%96%E5%88%86.py)
   - [Treap](4_set/%E6%9C%89%E5%BA%8F%E9%9B%86%E5%90%88/js/Treap.ts)
   - [TreeSet](4_set/%E6%9C%89%E5%BA%8F%E9%9B%86%E5%90%88/js/TreeSet.ts)
   - [SortedList (py)](4_set/%E6%9C%89%E5%BA%8F%E9%9B%86%E5%90%88/ATC-SortedList.py)
     [SortedList (ts)](4_set/%E6%9C%89%E5%BA%8F%E9%9B%86%E5%90%88/js/Treap.ts)

5. **图**

   - [拓扑排序](7_graph/%E6%8B%93%E6%89%91%E6%8E%92%E5%BA%8F/toposort.py)
   - [Dijkstra](7_graph/%E5%B8%A6%E6%9D%83%E5%9B%BE%E6%9C%80%E7%9F%AD%E8%B7%AF%E5%92%8C%E6%9C%80%E5%B0%8F%E7%94%9F%E6%88%90%E6%A0%91/dijkstra%E5%8D%95%E6%BA%90/dijkstra%E6%A8%A1%E6%9D%BF.py)
   - [BellmanFord](7_graph/%E5%B8%A6%E6%9D%83%E5%9B%BE%E6%9C%80%E7%9F%AD%E8%B7%AF%E5%92%8C%E6%9C%80%E5%B0%8F%E7%94%9F%E6%88%90%E6%A0%91/Bellman_ford/bellmanford%E5%88%A4%E6%96%AD%E8%B4%9F%E7%8E%AF.py)
   - [SPFA](7_graph/%E5%B8%A6%E6%9D%83%E5%9B%BE%E6%9C%80%E7%9F%AD%E8%B7%AF%E5%92%8C%E6%9C%80%E5%B0%8F%E7%94%9F%E6%88%90%E6%A0%91/Bellman_ford/spfa/spfa%E6%B1%82%E6%9C%80%E7%9F%AD%E8%B7%AF.py)
   - [Floyd](7_graph/%E5%B8%A6%E6%9D%83%E5%9B%BE%E6%9C%80%E7%9F%AD%E8%B7%AF%E5%92%8C%E6%9C%80%E5%B0%8F%E7%94%9F%E6%88%90%E6%A0%91/floyd%E5%A4%9A%E6%BA%90/Floyd.py)
   - [二分图检测](7_graph/%E4%BA%8C%E5%88%86%E5%9B%BE/%E4%BA%8C%E5%88%86%E5%9B%BE%E6%A3%80%E6%B5%8B.ts)
   - [匈牙利算法 (无权二分图最大匹配)](7_graph/%E4%BA%8C%E5%88%86%E5%9B%BE/%E6%97%A0%E6%9D%83%E4%BA%8C%E9%83%A8%E5%9B%BE%E6%9C%80%E5%A4%A7%E5%8C%B9%E9%85%8D%E9%97%AE%E9%A2%98/%E5%8C%88%E7%89%99%E5%88%A9%E7%AE%97%E6%B3%95%E6%A8%A1%E6%9D%BF.ts)
   - [Kuhn-Munkres 算法 (带权二分图最大权匹配)](7_graph/%E4%BA%8C%E5%88%86%E5%9B%BE/%E5%B8%A6%E6%9D%83%E4%BA%8C%E5%88%86%E5%9B%BE%E7%9A%84%E6%9C%80%E5%A4%A7%E6%9D%83%E5%8C%B9%E9%85%8D%E9%97%AE%E9%A2%98/KM%E7%AE%97%E6%B3%95%E6%A8%A1%E6%9D%BF.py)
   - [欧拉回路](7_graph/%E6%AC%A7%E6%8B%89/getEulerLoop.py)
   - [欧拉路径](7_graph/%E6%AC%A7%E6%8B%89/getEulerPath.py)
   - [最大流 (Dinic)](7_graph/acwing/%E7%BD%91%E7%BB%9C%E6%B5%81/0-%E6%9C%80%E5%A4%A7%E6%B5%81%E6%A8%A1%E6%9D%BF/Dinic%E6%B1%82%E6%9C%80%E5%A4%A7%E6%B5%81.py)
   - [最小费用最大流](7_graph/acwing/%E7%BD%91%E7%BB%9C%E6%B5%81/4-%E8%B4%B9%E7%94%A8%E6%B5%81/MinCostMaxFlow.py)
   - [Tarjan 算法](7_graph/acwing/%E6%97%A0%E5%90%91%E5%9B%BE%E7%9A%84%E5%8F%8C%E8%BF%9E%E9%80%9A%E5%88%86%E9%87%8F/Tarjan.py)

6. **并查集**

   - [并查集](14_%E5%B9%B6%E6%9F%A5%E9%9B%86/UnionFind.py)
   - [维护到根结点距离的并查集](14_%E5%B9%B6%E6%9F%A5%E9%9B%86/%E7%BB%8F%E5%85%B8%E9%A2%98/%E7%BB%B4%E6%8A%A4%E5%88%B0%E6%A0%B9%E8%8A%82%E7%82%B9%E8%B7%9D%E7%A6%BB/UnionFindMapWithDist.py)

7. **位运算**

   - [布隆过滤器](18_%E5%93%88%E5%B8%8C/%E5%B8%83%E9%9A%86%E8%BF%87%E6%BB%A4%E5%99%A8.ts)
   - [BitSet (位集)](18_%E5%93%88%E5%B8%8C/BitSet/BitSet.py)
   - [子集](21_%E4%BD%8D%E8%BF%90%E7%AE%97/%E4%BA%8C%E8%BF%9B%E5%88%B6%E6%9E%9A%E4%B8%BE%E4%B8%8E%E4%B8%89%E8%BF%9B%E5%88%B6%E6%9E%9A%E4%B8%BE/%E6%9E%9A%E4%B8%BE%E5%AD%90%E9%9B%86/powerset.py)

8. **数学**

   - [凸包](19_%E6%95%B0%E5%AD%A6/%E8%AE%A1%E7%AE%97%E5%87%A0%E4%BD%95/%E5%87%B8%E5%8C%85/587.%20%E5%AE%89%E8%A3%85%E6%A0%85%E6%A0%8F.py)
   - [多边形面积](19_%E6%95%B0%E5%AD%A6/%E8%AE%A1%E7%AE%97%E5%87%A0%E4%BD%95/%E5%87%B8%E5%8C%85/%E5%A4%9A%E8%BE%B9%E5%BD%A2%E9%9D%A2%E7%A7%AF%E5%85%AC%E5%BC%8F.py)
   - [直线方程](19_%E6%95%B0%E5%AD%A6/%E8%AE%A1%E7%AE%97%E5%87%A0%E4%BD%95/%E5%A4%9A%E7%82%B9%E5%85%B1%E7%BA%BF%E9%97%AE%E9%A2%98/%E4%B8%A4%E7%82%B9%E6%B1%82%E7%9B%B4%E7%BA%BF%E6%96%B9%E7%A8%8B.py)
   - [斯特林数](19_%E6%95%B0%E5%AD%A6/%E7%BB%84%E5%90%88/%E6%96%AF%E7%89%B9%E6%9E%97%E6%95%B0)
   - [康托展开](19_%E6%95%B0%E5%AD%A6/%E6%95%B0%E8%AE%BA/%E5%BA%B7%E6%89%98%E5%B1%95%E5%BC%80/%E5%BA%B7%E6%89%98%E5%B1%95%E5%BC%80.ts)
   - [质数](19_%E6%95%B0%E5%AD%A6/%E5%9B%A0%E6%95%B0%E7%AD%9B/prime.py)
   - [组合](19_%E6%95%B0%E5%AD%A6/acwing%E4%B8%93%E9%A1%B9%E8%AE%AD%E7%BB%83/%E7%BB%84%E5%90%88%E8%AE%A1%E6%95%B0/%E6%B1%82%E7%BB%84%E5%90%88%E6%8E%92%E5%88%97%E9%98%B6%E4%B9%98)
   - [线性基](21_%E4%BD%8D%E8%BF%90%E7%AE%97/%E6%8C%89%E4%BD%8D%E5%BC%82%E6%88%96/%E7%BA%BF%E6%80%A7%E5%9F%BA/%E7%BA%BF%E6%80%A7%E5%9F%BA.py)
   - [卷积](19_%E6%95%B0%E5%AD%A6/%E5%8D%B7%E7%A7%AF/Convolution.py)
   - [快速幂](19_%E6%95%B0%E5%AD%A6/%E6%95%B0%E8%AE%BA/%E5%BF%AB%E9%80%9F%E5%B9%82/qpow.ts)
   - [矩阵快速幂](19_%E6%95%B0%E5%AD%A6/%E7%9F%A9%E9%98%B5%E8%BF%90%E7%AE%97/%E7%9F%A9%E9%98%B5%E5%BF%AB%E9%80%9F%E5%B9%82/matqpow.py)

9. **杂项**

   - [二维前缀和](22_%E4%B8%93%E9%A2%98/%E5%89%8D%E7%BC%80%E4%B8%8E%E5%B7%AE%E5%88%86/%E5%B7%AE%E5%88%86%E6%95%B0%E7%BB%84/%E4%BA%8C%E7%BB%B4%E5%B7%AE%E5%88%86/%E4%BA%8C%E7%BB%B4%E5%B7%AE%E5%88%86%E6%A8%A1%E6%9D%BF.py)
   - [二维差分](22_%E4%B8%93%E9%A2%98/%E5%89%8D%E7%BC%80%E4%B8%8E%E5%B7%AE%E5%88%86/%E5%B7%AE%E5%88%86%E6%95%B0%E7%BB%84/%E4%BA%8C%E7%BB%B4%E5%B7%AE%E5%88%86/%E4%BA%8C%E7%BB%B4%E5%B7%AE%E5%88%86%E6%A8%A1%E6%9D%BF.py)
   - [ST 表](22_%E4%B8%93%E9%A2%98/RMQ%E9%97%AE%E9%A2%98/SparseTable.py)
   - [Bisect](9_%E6%8E%92%E5%BA%8F%E5%92%8C%E6%90%9C%E7%B4%A2/%E4%BA%8C%E5%88%86/bisect.ts)
   - [Trisect](19_%E6%95%B0%E5%AD%A6/%E6%A8%A1%E6%8B%9F%E9%80%80%E7%81%AB%E4%B8%8E%E7%88%AC%E5%B1%B1%E6%B3%95/%E4%B8%89%E5%88%86%E6%B3%95%E6%B1%82%E5%87%B8%E5%87%BD%E6%95%B0%E6%9E%81%E5%80%BC.py)
   - [回文生成器](22_%E4%B8%93%E9%A2%98/%E6%9E%9A%E4%B8%BE/%E6%9E%9A%E4%B8%BE%E5%9B%9E%E6%96%87/%E6%9E%9A%E4%B8%BE%E5%9B%9E%E6%96%87.py)
   - [NextPermutation](12_%E8%B4%AA%E5%BF%83%E7%AE%97%E6%B3%95/%E7%BB%8F%E5%85%B8%E9%A2%98/%E6%8E%92%E5%88%97/api/nextPermutation.py)
   - [莫队算法 (静态查询)](22_%E4%B8%93%E9%A2%98/%E7%A6%BB%E7%BA%BF%E6%9F%A5%E8%AF%A2/%E8%8E%AB%E9%98%9F/MoAlgo.py)
   - Itertools
     - [product](13_%E5%9B%9E%E6%BA%AF%E7%AE%97%E6%B3%95/itertools/product.ts)
     - [permutations](13_%E5%9B%9E%E6%BA%AF%E7%AE%97%E6%B3%95/itertools/permutations.ts)
     - [combinations](13_%E5%9B%9E%E6%BA%AF%E7%AE%97%E6%B3%95/itertools/combinations.ts)
     - [combinations_with_replacement](13_%E5%9B%9E%E6%BA%AF%E7%AE%97%E6%B3%95/itertools/combinationsWithReplacement.ts)

## ❤️ 感谢

[contest.js](https://github.com/harttle/contest.js)
