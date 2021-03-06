# 小美和小团所在公司的食堂有N张餐桌，从左到右摆成一排，
# 每张餐桌有2张餐椅供至多2人用餐，公司职员排队进入食堂用餐。
# 小美发现职员用餐的一个规律并告诉小团：当男职员进入食堂时，
# 他会优先选择已经坐有1人的餐桌用餐，只有当每张餐桌要么空着要么坐满2人时，
# 他才会考虑空着的餐桌；
# 当女职员进入食堂时，她会优先选择未坐人的餐桌用餐，
# 只有当每张餐桌都坐有至少1人时，她才会考虑已经坐有1人的餐桌；
# 无论男女，当有多张餐桌供职员选择时，`他会选择最靠左的餐桌用餐`。
# 现在食堂内已有若干人在用餐，另外M个人正排队进入食堂，
# 小团会根据小美告诉他的规律预测排队的每个人分别会坐哪张餐桌。


# 使用三个小根堆，分别存储当前人数为0,1,2的三种桌子的桌号，记为pq0,pq1,pq2
# 以男职员为例：
# 先尝试坐人数为1的桌子，该桌子人数就变成了2，等价于：将pq1的堆顶弹出，同时推入pq2
# 如果没有人数为1的桌子了，等价于pq1为空，就去坐人数为0的桌子，等价于：将pq0的堆顶弹出，同时推入pq1
# 因为桌号存储在优先队列，所以堆顶的桌号总是最小的，保证每个人有多个选择时优先坐最左边的桌子。
