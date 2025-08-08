# abc414-G - Get Many Cola-可乐(超大容量完全背包)
# https://atcoder.jp/contests/abc415/tasks/abc415_g
#
# 有一个神奇的可乐店，只能用空瓶换新可乐。
# 初始有 N 瓶可乐，0 空瓶。每次可以：
#
# 喝掉一瓶可乐，可乐-1，空瓶+1。
# 选择某个 i（1≤i≤M），用 Ai 个空瓶换 Bi 瓶可乐（需空瓶≥Ai，且 Bi<Ai）。 问最多能喝多少瓶可乐。
#
# 1<=N<=1e15
# 1<=Bi<Ai<=300
#
# !总喝数 = N + max ∑Bi s.t. ∑(Ai−Bi) ≤ N，是一个容量为 N 的完全背包，物品数 M，物品 i 的“重量”Ci=Ai−Bi，价值Vi=Bi。
# !O(maxW**3)
