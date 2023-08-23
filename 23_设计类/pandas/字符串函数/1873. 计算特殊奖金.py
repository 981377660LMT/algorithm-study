import pandas as pd

data = [
    [2, "Meir", 3000],
    [3, "Michael", 3800],
    [7, "Addilyn", 7400],
    [8, "Juan", 6100],
    [9, "Kannon", 7700],
]
Employees = pd.DataFrame(data, columns=["employee_id", "name", "salary"]).astype(
    {"employee_id": "int64", "name": "object", "salary": "int64"}
)

# 计算每个雇员的奖金。如果一个雇员的 id 是 奇数 并且他的名字不是以 'M' 开头，那么他的奖金是他工资的 100% ，否则奖金为 0 。
# 返回的结果按照 employee_id 排序。
# +-------------+-------+
# | employee_id | bonus |
# +-------------+-------+
# | 2           | 0     |
# | 3           | 0     |
# | 7           | 7400  |
# | 8           | 0     |
# | 9           | 7700  |
# +-------------+-------+
# 1. 新建列
# 2. apply() 函数


def calculate_special_bonus(employees: pd.DataFrame) -> pd.DataFrame:
    employees["bonus"] = employees.apply(
        lambda row: row["salary"] if row["employee_id"] & 1 and row["name"][0] != "M" else 0,
        # !axis : {0 or 'index', 1 or 'columns'}, default 0
        #     Axis along which the function is applied:
        #     0 or 'index': apply function to each column => 沿着行的方向apply, 即每一列
        #     1 or 'columns': apply function to each row => 沿着列的方向apply, 即每一行
        axis=1,
    )
    df = employees[["employee_id", "bonus"]].sort_values(by=["employee_id"])
    return df
