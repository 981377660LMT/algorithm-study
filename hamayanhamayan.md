## hamayanhamayan 的题集整理

https://blog.hamayanhamayan.com/entry/2100/01/01/000000
每章下面有很多例题，最好只做标有`解说`的例题

- 计数问题(atcoder 最难的类型)

https://compro.tsutaj.com//archive/181015_incexc.pdf (容斥原理)
https://drive.google.com/file/d/1WC7Y2Ni-8elttUgorfbix9tO1fvYN3g3/view

1. 设计状态 (根据复杂度猜，有时非常难想)
2. 改变搜索的顺序
   - 求序列数=>插入 dp, 只与元素大小有关, 先排序，再从大到小填
   - 排序+贪心
     https://drken1215.hatenablog.com/entry/2020/10/16/150700
   - 区间问题按照结束点排序
3. 改变条件的说法
   ad-hoc (特殊案例特殊分析)
   - k 进制变换
4. 贪心
5. 分类讨论
6. 所有 xxx 的所有 xxx 的值之和
   - 计算贡献
   - 分解成一次和

- TODO:/ntt 卷积/母函数(fps 知识)/数论积性函数(multiplicative-function)/矩阵知识
  https://hitonanode.github.io/cplib-cpp/formal_power_series/coeff_of_rational_function.hpp
  https://hitonanode.github.io/cplib-cpp/formal_power_series/sum_of_exponential_times_polynomial_limit.hpp
  https://hitonanode.github.io/cplib-cpp/formal_power_series/sum_of_exponential_times_polynomial.hpp

多项式全家桶(准备放弃)
https://maspypy.github.io/library/ (跟着 verify 的题做)
https://www.luogu.com.cn/training/3015
https://hackmd.io/@tatyam-prime/ryU4Ujup9#xKfracPxQx-%E3%82%92%E6%B1%82%E3%82%81%E3%82%8B
math/convolution/polynomial/matrix (先从卷积开始写库)
https://maspypy.com/category/algorithm_math
https://judge.yosupo.jp/

- ？？？

二元有理数(Dyadic Number) => 非公平博弈里使用
dual Problem
https://emthrm.github.io/cp-library/dual_problem.html

- **数学**

  - [ ] 总览
        https://blog.hamayanhamayan.com/entry/2017/10/14/125941
  - [ ] 容斥原理
        https://blog.hamayanhamayan.com/entry/2017/04/17/161345
  - [ ] 联立方程
        https://blog.hamayanhamayan.com/entry/2017/03/15/221719
  - [ ] 概率与期望
        https://blog.hamayanhamayan.com/entry/2016/11/14/223727
  - [ ] 卷积(畳み込み)
        https://blog.hamayanhamayan.com/entry/2017/05/20/125607
    1. 什么东西可以用卷积加速计算?(什么样的情形)
       - 01 序列卷积,计算配对个数时可以用卷积加速
         但是要化成`卷积形式`,可能要翻转序列
  - [ ] xor 问题
        https://blog.hamayanhamayan.com/entry/2017/05/20/145021
    - 01 の場合、a xor b = a(1-b) + b(1-a)となる => 卷积
  - [ ] 线性规划
        https://blog.hamayanhamayan.com/entry/2017/05/31/131424
  - [ ] 斐波那契数列
        https://blog.hamayanhamayan.com/entry/2018/07/28/091611

- 木構造

  - [x] LCA
  - [x] HL 分解
  - [x] 欧拉路径
  - [x] 树的同型判定
  - [x] 重心分解
  - [x] LCT

- 图论

  - [ ] 网络流
    1. 最大流
    2. 最小费用最大流
       コストを損失と考えて最大化問題を解く
    3. 最大匹配
  - [ ] 无向图上的计数问题
        https://blog.hamayanhamayan.com/entry/2100/01/01/000000
  - [ ] 无向图相关问题
    - 只有偶数长度的环=>二分图
    - `只有奇数长度的环`=>仙人掌(cactus), 点双缩点后各个连通分量为桥或者环=>dp
    - 每条边选一个端点,使得没有重合的端点
      => 所有的连通分量中,`边数<=顶点数(不存在弦)`
  - [x] 最短路
  - [x] 欧拉
  - [ ] 有向图相关问题
  - [ ] 最大独立集/最大团
  - [ ] bfs
  - [ ] 最小生成树

- 数据结构

  - [x] 并查集
    1. 森の連結成分は「頂点数-辺数」
  - [x] 线段树/树状数组
    1. [二维线段树和线段树套线段树减少空间的技巧](https://blog.hamayanhamayan.com/entry/2017/12/09/015937)
  - [x] 持久化数据结构
    - 部分永続。最新版のみ変更可能、Undo ができる機構がある。履歴は一直線。
    - 完全永続。昔のどのバージョンからでも変更でき、履歴が全て残っている。履歴は木。
    - 永続配列があれば永続 Union-Find が作れるらしい
    - 永続配列は永続木があれば作れるらしい
    - 部分永続はスタックで履歴を残しておいて逆操作をしていって戻すやり方がある
      https://kmjp.hatenablog.jp/entry/2017/06/13/1000
  - [x] 平衡树
  - [x] ST 表
    - 结合律和幂等律 => SparseTable
    - 结合律 => DisjointSparseTable
  - [x] WaveletMatrix

- dp

  - [ ] dp 总览
    - 排序后再 dp
    - `遷移が多い場合は貪欲法`を使うことで最適であろう遷移を絞ることができる
    - 全体で遷移数がそんなに多くないことを見抜く
  - [ ] 插入 dp
    - 能不能模板化? https://blog.hamayanhamayan.com/entry/2017/10/04/112825
  - [x] 区间 dp
  - [x] 数位 dp
  - [x] dp 优化
  - [x] 状压 dp(bit dp)
  - [x] 连通性 dp(連結 dp)
    - 非常难,跳过
  - [x] 树形 dp

    1. dp[i] := 頂点 i の部分木についての何か
    2. 全方位木 DP
    3. 二乗の木 DP という、頂点集合の DP をマージする時に部分木の要素数の個数分だけ使ってマージするようにすると O(N^3)が O(N^2)に落ちるテクがある
    4. マージテク(木 DP っぽいことをしながらマージテクを使うものもある（配列で持たず map で持ってマージしていく）)

- 杂项

  - [ ] 筛法
    - 埃氏筛
    - 区间筛(R<=10^12, R-L<=10^6 という制約がヒント)
    - つまるところは「`rep(i,1,N) for(j=i;j<=N;j+=i)` というループ構造は O(NlogN)で行えるという所を応用して問題を解く」
    - SOSDp 是 zeta 变换的下位替换(不如 zeta 变换)
  - [x] 剪枝
  - [x] 曼哈顿距离
    - 45 度逆时针旋转
    - 两点的曼哈顿距离在仿射变换(平移、旋转)后不发生变化
  - [x] 滑窗
  - [x] 二分/三分探索
  - [x] 前缀和/差分
    - 删除子数组后的性质=>前缀+后缀
  - [x] LIS
  - [x] 莫队
  - [x] 根号分解
  - [x] 字符串
    - 模式匹配可以 bitset 加速
    - 最大字典序需要从后往前确定(适合记忆化 dfs )
  - [x] game
    - Grundy 数
    - A と B の差の偶奇によって勝敗が決まる。
  - [x] 2-sat
    - ダメな組合せが見つかったときに制約を作る
    - `2^n となるような場合は 2-SAT かもと考えれば思いつくのかな?`
  - [x] merge 技巧
  - [x] 分治
    - 「列の分割統治法」、「木の分割統治法」、「平面の分割統治法」
  - [x] 凸包
  - [x] 随机
    - 状態集合の中で`正解となるパターン数が非常に多い`場合は乱択アルゴリズムが使える
  - [x] meet in the middle
    - 2 グループに分けて全列挙をして、1 つのグループは全探索し、もう一方のグループに関しては二分探索などで高速に処理する
  - [x] 最近点対
  - [x] 几何
    - 全ての座標が異なる x,y 座標が 10^5 くらいの場合は、x^2+y^2 が同じ座標が 144 個くらいしかない
  - [x] 构文解析
  - [x] 扫描线(平面走查)
  - [x] 交互问题
    - 最优战略
    - 二分/三分
    - 按位计算
    - 随机

- 競技プログラミングにおける細かな話題まとめ

  - [ ] 最大長方形
  - [ ] Meldable Heap
  - [ ] 二分图判定
  - [ ] 树的重心
  - [ ] 立方体の和集合の体積の総和を求める
  - [ ] 全通りの組合せを考えて答える
  - [ ] bitset 64 倍高速化
    1. 是否到达
    2. dp 值为布尔值的情况
    3. 64 倍 SCC
  - [ ] 最大化最小字典序
  - [ ] 回文
  - [ ] 基环树 (N 頂点 N 辺のグラフ、Functional Graph)
  - [ ] 最短路径树
  - [ ] 一端を固定して、もう片方を伸ばしていくと状態変化がそんなに無いので分割統治っぽくでき、高速に処理できる
    1. ある配列の左端を固定して、右端を動かしながら AND や OR を取るとする。すると、`求まる AND,OR はそれぞれ高々 32 通りしかない`
    2. ある配列の左端を固定して、右端を動かしながら gcd を取っていき、`gcd が同じ要素をまとめていくと O(logN)グループしか無い`
  - [ ] 自動ベクトル化
  - [ ] 输入随机的情况下的技巧
  - [ ] 約数を使って再帰していく問題
  - [ ] 无向图最小路径覆盖
  - [ ] 特殊なソートで最適解を得るテク
  - [ ] 実はそんなに数が無い
  - [ ] 不変量を使った問題
        一次会改变多个数的题目，往往入手点在「不变量」上，也就是操作不会改变什么
        一般要用前缀和数组/差分数组来构造出不变量
  - [ ] 贪心+背包 dp
  - [ ] 有理数、ファレイ数列、Stern-Brocot Tree
  - [ ] Sliding Window Aggregation
  - [ ] 小ネタ、小テク
    1. 特殊的坐标压缩
    2. 取模后数字变为一半
    3. 最大栈/最小栈 维护前缀的最大值/最小值
    4. 特殊な制約により計算量が抑えられる系
    5. クエリの要素数によって場合分けを行うことでならしで間に合わせることができたりする
    6. 条件を少しずつ変えて全探索は差分を計算することで全てを計算し直さずに全探索できる
    7. より制約の緩い問題を解くことに帰着する
    8. 数列の中に順列が含まれていて、順列がどんな形でも関係ない場合に、特定の順列のパターンは全てのパターン ÷ 順列のパターンで求められる
    9. 「ある盤面が列または行のスワップ操作だけで全部白にできますか？」→ 全ての 2x2 の subrectangle について黒マスの個数が偶数
    10. 正难则反
    11. set で区間を保持しながら解く
    12. 右手法という方針がある
    13. クエリ先読み(离线查询)
    14. 网格图(grid) 是二分图

---

https://qiita.com/e869120
https://qiita.com/drken
