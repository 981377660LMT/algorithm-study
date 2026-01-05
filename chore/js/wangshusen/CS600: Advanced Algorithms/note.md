https://github.com/wangshusen/AdvancedAlgorithms

# 高级算法 (Advanced Algorithms)

## 1. 基础数据结构 (Basic Data Structures)

- [数组、向量与链表 (Array, vector, and linked list)](#11-数组向量与链表-array-vector-and-linked-list)
- [二分查找 (Binary search)](#12-二分查找-binary-search)
- [跳表 (Skip list)](#13-跳表-skip-list)

## 2. 排序 (Sorting)

- [插入排序 (Insertion sort)](#21-插入排序-insertion-sort)
- [冒泡排序 (Bubble sort)](#22-冒泡排序-bubble-sort)
- [归并排序 (Merge sort)](#23-归并排序-merge-sort)
- [快速排序 (Quick sort)](#24-快速排序-quick-sort)
- [快速选择 (Quick select)](#25-快速选择-quick-select)

## 3. 分治法 (Divide and Conquer)

- [主定理 (Master theorem)](#31-主定理-master-theorem)

## 4. 矩阵数据结构与算法 (Matrix Data Structure and Algorithms)

- [矩阵加法与乘法 (Addition and multiplication)](#41-矩阵加法与乘法-addition-and-multiplication)
- [稠密矩阵数据结构 (Dense matrix data structures)](#42-稠密矩阵数据结构-dense-matrix-data-structures)
- [稀疏矩阵数据结构 (Sparse matrix data structures)](#43-稀疏矩阵数据结构-sparse-matrix-data-structures)
- [快速矩阵乘法与 Strassen 算法 (Fast matrix multiplication and Strassen algorithm)](#44-快速矩阵乘法与-strassen-算法-fast-matrix-multiplication-and-strassen-algorithm)

## 5. 二叉树 (Binary Trees)

- [树的基础知识 (Tree basics)](#51-树的基础知识-tree-basics)
- [二叉搜索树：查找与插入 (BST: search and insertion)](#52-二叉搜索树查找与插入-bst-search-and-insertion)
- [二叉搜索树：遍历 (BST: traversal)](#53-二叉搜索树遍历-bst-traversal)
- [二叉搜索树：删除 (BST: deletion)](#54-二叉搜索树删除-bst-deletion)

## 6. 平衡树 (Balanced Trees)

## 7. 优先队列 (Priority Queues)

- [优先队列概览 (Priority queues)](#71-优先队列概览-priority-queues)
- [二叉堆 (Binary heaps)](#72-二叉堆-binary-heaps)

## 8. 并查集 (Disjoint Sets)

## 9. 图论基础 (Graph Basics)

- [图的数据结构 (Graph data structures)](#91-图的数据结构-graph-data-structures)
- [拓扑排序 (Topological sort)](#92-拓扑排序-topological-sort)

## 10. 最短路径 (Shortest Paths)

- [最短路径问题 (Shortest path problem)](#101-最短路径问题-shortest-path-problem)
- [无权图的单源最短路径 (Single-source shortest path in unweighted graphs)](#102-无权图的单源最短路径-single-source-shortest-path-in-unweighted-graphs)
- [有权图的单源最短路径 (Single-source shortest path in weighted graphs)](#103-有权图的单源最短路径-single-source-shortest-path-in-weighted-graphs)

## 11. 最小生成树 (Minimum Spanning Trees)

- [生成树基础 (Spanning trees)](#111-生成树基础-spanning-trees)
- [Prim 算法 (Prim's algorithm)](#112-prim-算法-prims-algorithm)
- [Kruskal 算法 (Kruskal's algorithm)](#113-kruskal-算法-kruskals-algorithm)

## 12. 网络流问题 (Network Flow Problems)

- [最大流问题 (Maximum flow problem)](#121-最大流问题-maximum-flow-problem)
- [Ford-Fulkerson 算法](#122-ford-fulkerson-算法)
- [Edmonds–Karp 算法](#123-edmonds-karp-算法)
- [Dinic 算法](#124-dinic-算法)
- [最大流最小割 (Max-flow and min-cut)](#125-最大流最小割-max-flow-and-min-cut)

## 13. 二分图 (Bipartite Graphs)

- [二分性检测 (Testing bipartiteness)](#131-二分性检测-testing-bipartiteness)
- [二分图最大匹配 (Maximum cardinality bipartite matching)](#132-二分图最大匹配-maximum-cardinality-bipartite-matching)
- [二分图最大权匹配 (Maximum weight bipartite matching)](#133-二分图最大权匹配-maximum-weight-bipartite-matching)
- [匈牙利算法 (Hungarian algorithm)](#134-匈牙利算法-hungarian-algorithm)
- [稳定婚姻问题 (Stable marriage problem)](#135-稳定婚姻问题-stable-marriage-problem)
- [Gale-Shapley 算法](#136-gale-shapley-算法)

## 14. 动态规划 (Dynamic Programming)

- [斐波那契数列 (Fibonacci numbers)](#141-斐波那契数列-fibonacci-numbers)
- [爬楼梯 (Climbing stairs)](#142-爬楼梯-climbing-stairs)
- [最长公共子串 (Longest common substring)](#143-最长公共子串-longest-common-substring)
- [编辑距离 (Edit distance)](#144-编辑距离-edit-distance)
- [硬币组合 (Combinations of coins)](#145-硬币组合-combinations-of-coins)
- [背包问题 (Knapsack)](#146-背包问题-knapsack)

## 15. 字符串 (Strings)

- [前缀树 (Tries)](#151-前缀树-tries)
- [KMP 算法](#152-kmp-算法)

## 16. 随机算法 (Randomized Algorithms)

- [蒙特卡罗算法 (Monte Carlo algorithms)](#161-蒙特卡罗算法-monte-carlo-algorithms)
- [集中不等式 (Concentration inequalities)](#162-集中不等式-concentration-inequalities)
- [伪随机数生成器 (Pseudo random number generators)](#163-伪随机数生成器-pseudo-random-number-generators)
- [随机打乱 (Random shuffling)](#164-随机打乱-random-shuffling)

## 17. 指纹与哈希 (Fingerprinting & Hashing)

- [哈希表 (Hash table)](#171-哈希表-hash-table)
- [抗碰撞哈希 (Collision-resistant hash)](#172-抗碰撞哈希-collision-resistant-hash)
- [局部敏感哈希 (Locality sensitive hashing)](#173-局部敏感哈希-locality-sensitive-hashing)

## 18. 加密算法 (Cryptographic Algorithms)

- [RSA 算法](#181-rsa-算法)
- [数字签名 (Digital signature)](#182-数字签名-digital-signature)
- [同态加密 (Homomorphic encryption)](#183-同态加密-homomorphic-encryption)

## 19. 加密货币 (Crypto Currency)

- [去中心化支付系统 (Decentralized payment system)](#191-去中心化支付系统-decentralized-payment-system)
- [区块链 (Blockchain)](#192-区块链-blockchain)
- [默克尔树 (Merkel tree)](#193-默克尔树-merkel-tree)

## 20. 回溯法 (Backtracking)

- [收费站重建问题 (Turnpike reconstruction)](#201-收费站重建问题-turnpike-reconstruction)
- [八皇后问题 (Eight queens problem)](#202-八皇后问题-eight-queens-problem)
- [哈密顿回路 (Hamiltonian cycle)](#203-哈密顿回路-hamiltonian-cycle)
