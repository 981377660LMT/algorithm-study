import pandas as pd

data = [
    [1, 3, 5, "2019-08-01"],
    [1, 3, 6, "2019-08-02"],
    [2, 7, 7, "2019-08-01"],
    [2, 7, 6, "2019-08-02"],
    [4, 7, 1, "2019-07-22"],
    [3, 4, 4, "2019-07-21"],
    [3, 4, 4, "2019-07-21"],
]
Views = pd.DataFrame(data, columns=["article_id", "author_id", "viewer_id", "view_date"]).astype(
    {
        "article_id": "Int64",
        "author_id": "Int64",
        "viewer_id": "Int64",
        "view_date": "datetime64[ns]",
    }
)


# 请查询出所有浏览过自己文章的作者
# 结果按照 id 升序排列。
# +------+
# | id   |
# +------+
# | 4    |
# | 7    |
# +------+
#
# !提取满足条件的行
# 布尔索引(Series[bool])允许我们通过使用布尔数组或条件来过滤 DataFrame。
# 这意味着我们可以使用布尔值的 Series 或创建对于 DataFrame 中每一行都评估为 True 或 False 的条件。
# 通过将这些布尔值或条件应用为 DataFrame 的索引，我们可以选择性地`提取满足条件的行`。


def article_views(views: pd.DataFrame) -> pd.DataFrame:
    df = views[views["author_id"] == views["viewer_id"]]
    df.drop_duplicates(subset=["author_id"], inplace=True)  # 去重
    df.sort_values(by=["author_id"], inplace=True)  # 排序
    df.rename(columns={"author_id": "id"}, inplace=True)  # 重命名列
    return df[["id"]]
