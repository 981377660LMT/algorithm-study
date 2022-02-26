# n<=5000  q<=1000 完全可以用dfs序加counter暴力统计 最多5000*1000
# 小团只能从某一个节点的子树中选取某一种颜色的拿。
# 小团想要拿到数量最多的那种颜色的所有苹果，请帮帮她

# 解法1：
# 一个子树可以转化为一个dfs序连续的区间，然后用前缀和统计区间内每种颜色有多少(这样做只限于颜色很少的情况)，颜色多直接counter
# (所以本题就是树上莫队模板题。)

# 解法2：
# 后序dfs返回Counter
