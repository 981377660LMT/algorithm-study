import pandas as pd

data = [[1, "aLice"], [2, "bOB"]]
Users = pd.DataFrame(data, columns=["user_id", "name"]).astype(
    {"user_id": "Int64", "name": "object"}
)


# 编写解决方案，修复名字，使得只有第一个字符是大写的，其余都是小写的。
# 返回按 user_id 排序的结果表。
# +---------+-------+
# | user_id | name  |
# +---------+-------+
# | 1       | Alice |
# | 2       | Bob   |
# +---------+-------+


def fix_names(users: pd.DataFrame) -> pd.DataFrame:
    users["name"] = users["name"].str.capitalize()
    return users.sort_values(by=["user_id"])
