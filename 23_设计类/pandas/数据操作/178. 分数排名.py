import pandas as pd

data = [[1, 3.5], [2, 3.65], [3, 4.0], [4, 3.85], [5, 4.0], [6, 3.65]]
Scores = pd.DataFrame(data, columns=["id", "score"]).astype({"id": "Int64", "score": "Float64"})


# 输入:
# Scores 表:
# +----+-------+
# | id | score |
# +----+-------+
# | 1  | 3.50  |
# | 2  | 3.65  |
# | 3  | 4.00  |
# | 4  | 3.85  |
# | 5  | 4.00  |
# | 6  | 3.65  |
# +----+-------+
# 输出:
# +-------+------+
# | score | rank |
# +-------+------+
# | 4.00  | 1    |
# | 4.00  | 1    |
# | 3.85  | 2    |
# | 3.65  | 3    |
# | 3.65  | 3    |
# | 3.50  | 4    |
# +-------+------+
# 查询并对分数进行排序。排名按以下规则计算:

# 分数应按从高到低排列。
# 如果两个分数相等，那么两个分数的排名应该相同。
# 在排名相同的分数后，排名数应该是下一个连续的整数。换句话说，排名之间不应该有空缺的数字。
# 按 score 降序返回结果表。


# !dense_rank


def order_scores(scores: pd.DataFrame) -> pd.DataFrame:
    scores["rank"] = scores["score"].rank(method="dense", ascending=False)  # 每个分数获得唯一的排名
    return scores[["score", "rank"]].sort_values(by=["score"], ascending=False)  # 按排名升序排列
