import pandas as pd

# 查询 Employee 表中第 n 高的工资。如果没有第 n 个最高工资，查询结果应该为 null 。


# Employee table:
# +----+--------+
# | id | salary |
# +----+--------+
# | 1  | 100    |
# | 2  | 200    |
# | 3  | 300    |
# +----+--------+
# n = 2
#
# 输出:
# +------------------------+
# | getNthHighestSalary(2) |
# +------------------------+
# | 200                    |
# +------------------------+


def nth_highest_salary(employee: pd.DataFrame, N: int) -> pd.DataFrame:
    df = employee[["salary"]].drop_duplicates()
    if len(df) < N:
        return pd.DataFrame({f"getNthHighestSalary({N})": [None]})
    return df.sort_values(by=["salary"], ascending=False).iloc[N - 1 : N]
