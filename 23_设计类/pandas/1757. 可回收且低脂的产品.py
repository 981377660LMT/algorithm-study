import pandas as pd

data = [["0", "Y", "N"], ["1", "Y", "Y"], ["2", "N", "Y"], ["3", "Y", "Y"], ["4", "N", "N"]]
Products = pd.DataFrame(data, columns=["product_id", "low_fats", "recyclable"]).astype(
    {"product_id": "int64", "low_fats": "category", "recyclable": "category"}
)


def find_products(products: pd.DataFrame) -> pd.DataFrame:
    df = products[(products["low_fats"] == "Y") & (products["recyclable"] == "Y")]
    return df[["product_id"]]
