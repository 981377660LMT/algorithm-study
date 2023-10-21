# n<=2000
# 给定n天以及初始资产m，然后输入一个数组表示每天股票价格。
# 一条只能买或卖1只股，但是可以无限持股，不能贷款即资产必须≥0。
# 问你到最后一天的时候最大资产是多少，资产总数为[现金 + 持股数 * 股票价格]
# 每天可以不买/卖/买 钱不能变为负数

# dp的本质是在DAG上求最短(长)路
# 复杂度是O(index*count)
# dp[index][count]=max(dp[index-1][count],dp[index-1][count+1]+cost,dp[index-1][count-1]-cost)
# `注意状态转移是否合法`

# 也可以用dijk dist[index][count] 表示第index天持有count份股票时的最大收益
# 用dijk些bfs逻辑更加清晰 但是复杂度会多一个logn(如果用dijk的话，当然也可以写普通的bfs+deque)
