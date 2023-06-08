# https://leetcode.cn/problems/subrectangle-queries/


# !1.二维线段树(因为是区间修改,所以树套树)
# !2.保存修改倒序查询的做法
#   优化 => 保存修改，当修改超过一个阈值，比如1000，批量处理这些修改，扫描线+一维线段树更新
#   更新控制在O(sqrt)次
