import pandas as pd

data = [[1, "Vote for Biden"], [2, "Let us make America great again!"]]
Tweets = pd.DataFrame(data, columns=["tweet_id", "content"]).astype(
    {"tweet_id": "Int64", "content": "object"}
)


# 查询所有无效推文的编号（ID）。当推文内容中的字符数严格大于 15 时，该推文是无效的。
# 以任意顺序返回结果表。


def invalid_tweets(tweets: pd.DataFrame) -> pd.DataFrame:
    filter_ = tweets["content"].str.len() > 15
    pd = tweets[filter_]
    return pd[["tweet_id"]]
