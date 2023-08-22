https://oi-wiki.org/dp/opt/quadrangle/
如果区间权值的函数 **w(left,right)** 满足 `区间包含单调性`以及`四边形不等式`
那么可以利用决策的单调性进行优化

很多区间 DP 都能用四边形不等式优化

- https://beet-aizu.github.io/library/algorithm/knuthyao.cpp
- https://www.luogu.com.cn/blog/command-block/dp-di-jue-ce-dan-diao-xing-you-hua-zong-jie
- https://blog.csdn.net/weixin_43914593/article/details/105150937

---

Monge 性(四边形不等式)
https://hackmd.io/@tatyam-prime/monge1

形象地理解:给定一个矩阵和 row1<=row2,col1<=col2
如果`左上角+右下角<=右上角+左下角`,那么这个矩阵就是 Monge 性的
`左上角+右下角<=右上角+左下角` => 变形
`右下角-右上角<=左下角-左上角`
固定(row2,col1)的左端,向右移动时增量具有单调性
Monge 性等价条件:`对任意一个2*2的子矩阵,左上角+右下角<=右上角+左下角`

- 如果 `n*m` 的矩阵是 Monge 性的,那么每行的最小值可以在 O(n+m)内求出来
- Monge 图(邻接矩阵是 Monge 性的)的单源最短路可以 O(n)求
- 优化区间 dp
- Monge 是 totally monotone 的

  > monotone: `每一行取得最小值的列数是单调的`
  > totally monotone: 任意的子矩阵具有 monotone 性 （完全单调矩阵）

- Monge 矩阵的例子

1. A[i][j] = f(j-i)
2. A[i][j] = `(j-i)^2 +∑ak (k=1..j-i)`
3. A[i][j]=`(ai-bj)^2`(a,b 都是单增的数列)
