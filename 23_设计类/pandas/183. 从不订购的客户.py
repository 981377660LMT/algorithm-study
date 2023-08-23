import pandas as pd

data = [[1, "Joe"], [2, "Henry"], [3, "Sam"], [4, "Max"]]
Customers = pd.DataFrame(data, columns=["id", "name"]).astype({"id": "Int64", "name": "object"})
data = [[1, 3], [2, 1]]
Orders = pd.DataFrame(data, columns=["id", "customerId"]).astype(
    {"id": "Int64", "customerId": "Int64"}
)

# 找出所有从不点任何东西的顾客。
# 以 任意顺序 返回结果表。
# +-----------+
# | Customers |
# +-----------+
# | Henry     |
# | Max       |
# +-----------+


def find_customers(customers: pd.DataFrame, orders: pd.DataFrame) -> pd.DataFrame:
    df = customers[~customers["id"].isin(orders["customerId"])]
    return df[["name"]].rename(columns={"name": "Customers"})


print(find_customers(Customers, Orders))
