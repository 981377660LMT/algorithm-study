# 我们希望「恶魔」无法到达城堡，即「城堡」和「恶魔」之间不连通，对应到图论模型上就是「割」问题。


# https://leetcode.cn/problems/7rLGCR/solution/lcp-38-shou-wei-cheng-bao-by-zerotrac2-kgv2/
# https://leetcode.cn/problems/7rLGCR/solution/czui-xiao-ge-jian-mo-dai-zhu-shi-by-litt-g77h/


# !建立源点向恶魔出生点连边，再将城堡向汇点连边；再将图中相邻的点之间连边即可
# 因为要移除点而不是移除边 所以要把所有点拆成in 和 out
# "."的自身流量为 1，其余边的流量为 inf,
# !这样最小割一定发生在空地上，且最小割的大小就是建立障碍的个数，也就是答案
# 传送门互相连通，建立一个特殊点把它们都免费连起来即可。
