# 布尔矩阵乘法

输入和输出矩阵的元素均为布尔值。
按矩阵乘法的公式运算时，可以把“乘”看成 and，把“加”看成 or
对矩阵乘法 C[i][j] |= A[i][k] & B[k][j], 它的一个直观意义是把 A 的行和 B 的列看成集合，
A 的第 i 行包含元素 k 当且仅当 A[i][k]=1。
B 的第 j 列包含元素 k 当且仅当 B[k][j]=1。
那么 C[i][j]代表 A 的第 i 行和 B 的第 j 列是否包含公共元素。

一个应用是传递闭包(Transitive Closure)的加速计算。

## 实现

- bitset 加速，`O(n^3/w)`
  https://nyaannyaan.github.io/library/matrix/bitmatrix.hpp
- 分块+位运算加速，`O(n^3/w)`
  https://github.com/spaghetti-source/algorithm/blob/4fdac8202e26def25c1baf9127aaaed6a2c9f7c7/math/matrix_bool.cc#L5

- 最快的布尔矩阵乘法 `O(n^3/(w*logn))`，Method of Four Russians + 压位
  https://www.doc88.com/p-2136480081151.html

# f2 上的矩阵乘法

输入矩阵中的元素只包含 0 或 1，乘积对 2 取模
在 mod2 的意义下，加法可以看成`异或`
和布尔矩阵乘法一样,`把 mat.mul 的 or 换成 xor 就可以了`
https://hitonanode.github.io/cplib-cpp/linear_algebra_matrix/linalg_bitset.hpp

---

传递闭包的分治加速计算
https://leetcode.cn/problems/course-schedule-iv/solution/chuan-di-bi-bao-gui-yue-dao-ju-zhen-chen-47yu/
https://github.com/hqztrue/shared_materials/blob/master/code/matrix_multiplication/%E4%BC%A0%E9%80%92%E9%97%AD%E5%8C%85%26bool_mat_mul.cpp
https://github.com/hqztrue/shared_materials/tree/master/code/matrix_multiplication
