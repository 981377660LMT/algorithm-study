import pandas as pd

data = [
    [1, "Joe", 70000, 1],
    [2, "Jim", 90000, 1],
    [3, "Henry", 80000, 2],
    [4, "Sam", 60000, 2],
    [5, "Max", 90000, 1],
]
Employee = pd.DataFrame(data, columns=["id", "name", "salary", "departmentId"]).astype(
    {"id": "Int64", "name": "object", "salary": "Int64", "departmentId": "Int64"}
)
data = [[1, "IT"], [2, "Sales"]]
Department = pd.DataFrame(data, columns=["id", "name"]).astype({"id": "Int64", "name": "object"})


# Employee 表:
# +----+-------+--------+--------------+
# | id | name  | salary | departmentId |
# +----+-------+--------+--------------+
# | 1  | Joe   | 70000  | 1            |
# | 2  | Jim   | 90000  | 1            |
# | 3  | Henry | 80000  | 2            |
# | 4  | Sam   | 60000  | 2            |
# | 5  | Max   | 90000  | 1            |
# +----+-------+--------+--------------+
# Department 表:
# +----+-------+
# | id | name  |
# +----+-------+
# | 1  | IT    |
# | 2  | Sales |
# +----+-------+
# 输出：
# +------------+----------+--------+
# | Department | Employee | Salary |
# +------------+----------+--------+
# | IT         | Jim      | 90000  |
# | Sales      | Henry    | 80000  |
# | IT         | Max      | 90000  |
# +------------+----------+--------+


# https://leetcode.cn/problems/department-highest-salary/solutions/2366207/bu-men-gong-zi-zui-gao-de-yuan-gong-by-l-jo1i/?envType=study-plan-v2&envId=30-days-of-pandas&lang=pythondata
def department_highest_salary(employee: pd.DataFrame, department: pd.DataFrame) -> pd.DataFrame:
    df = employee.merge(department, left_on="departmentId", right_on="id", how="left")
    # id_x	name_x	   salary	 departmentId	      id_y	name_y
    # 1	     Joe	     70000	     1	            1	   IT
    # 2	     Jim	     90000	     1	            1	   IT
    # 3	     Henry	   80000	     2	            2	   Sales
    # 4	     Sam	     60000	     2	            2	   Sales
    # 5	     Max	     90000	     1	            1	   IT
    # 内连接（'inner'）：只保留两个DataFrame中连接键都有的行。这是merge()函数的默认连接类型。
    # 左连接（'left'）：保留左侧DataFrame中的所有行，如果右侧DataFrame中没有匹配的连接键，则结果中的对应行将为NaN。
    # 右连接（'right'）：保留右侧DataFrame中的所有行，如果左侧DataFrame中没有匹配的连接键，则结果中的对应行将为NaN。
    # 全连接（'outer'）：保留两个DataFrame中的所有行，如果某个DataFrame中没有匹配的连接键，则结果中的对应行将为NaN。
    # !在合并后，原始表中具有相同名称的列（例如 name）会被重命名（作为 name_x 和 name_y），因此我们需要进行列重命名。
    df.rename(
        columns={"name_x": "Employee", "name_y": "Department", "salary": "Salary"}, inplace=True
    )

    # 我们根据 Department 列对 df 进行分组，并对 Salary 列应用 transform('max') 函数，这将为每个部门计算最高薪水
    maxSalary = df.groupby("Department")["Salary"].transform("max")
    df = df[df["Salary"] == maxSalary]
    return df[["Department", "Employee", "Salary"]]
