卡特兰数 — 计数的映射方法的伟大胜利

1. 进出栈
   一个足够大的栈的进栈序列 1,2,3...n 为时有多少个不同的出栈序列？
2. 求(n+1)个叶子的满二叉树的个数
   向左记为+1，向右记为-1，按照向左优先的原则
   四个叶子结点，n=3 时 +1+1+1-1-1-1
3. 小兔的路径
4. 不相交的握手
5. 凸多边形的剖分
6. n 对括号形成的合法括号表达式

二叉树的计数问题:已知二叉树有 n 个结点，求能构成多少种不同的二叉树。
括号化问题:一个合法的表达式由()包围，()可以嵌套和连接，如:(O))()也是合法表达式，现给出 n 对括号，求可以组成的合法表达式的个数。
划分问题:将一个凸 n ＋ 2 多边形区域分成三角形区域的方法数。
出栈问题:一个栈的进栈序列为 1,2,3,..n，求不同的出栈序列有多少种。
路径问题:在 n·n 的方格地图中，从一个角到另外一个角，求不跨越对角线的路径数有多少种。
握手问题:2n 个人均匀坐在一个圆桌边上，某个时刻所有人同时与另一个人握手，要求手之间不能交叉，求共有多少种握手方法。

所有问题可以转化为`不跨越对角线的路径问题`

```Python
@lru_cache(None)
def dp(i, j):
    if i == 0 and j == 0:
        return 1
    if i < 0 or j < 0:
        return 0
    if i < j:
        return 0
    return dp(i - 1, j) + dp(i, j - 1)

```
