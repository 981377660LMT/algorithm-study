多重背包问题转化为 01 背包问题
多重背包问题是这样描述的:
有 n 种物品，物品 j 的体积为 vj，价值为 w;，有一个体积限制 V。`每种物品还有一个 c，表示每种物品的个数`，此问题称为多重背包问题。

思路 1 ：将物品展开，全拆成 01

思路 2 ：2 进制分解
有这样一个事实:任意一个数 n，它一定可以用 1,2,4.8.”，以及 P 到 2+1 之间的某个数表示。例如 13 以内的所有数都可以用 1,2,4,6 表示。
所以对于物品 i, 数量限制是 c,可以将其分成若干物品，它们的价值和体积为:(v， v)，(2*u;，2* ;)，.….

https://tjkendev.github.io/procon-library/python/dp/knapsack2.html

dp/ndp/set 比较容易写

- 如果物品价值之和很小,背包体积很大,那就将价值之和作为 dp 的维度
- 如果物品价值之和很大,背包体积很小,那就将背包体积作为 dp 的维度
