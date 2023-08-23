import pandas as pd

data = [
    [1, "Winston", "winston@leetcode.com"],
    [2, "Jonathan", "jonathanisgreat"],
    [3, "Annabelle", "bella-@leetcode.com"],
    [4, "Sally", "sally.come@leetcode.com"],
    [5, "Marwan", "quarz#2020@leetcode.com"],
    [6, "David", "david69@gmail.com"],
    [7, "Shapiro", ".shapo@leetcode.com"],
]
Users = pd.DataFrame(data, columns=["user_id", "name", "mail"]).astype(
    {"user_id": "int64", "name": "object", "mail": "object"}
)


# 编写一个解决方案，以查找具有有效电子邮件的用户。
# 一个有效的电子邮件具有前缀名称和域，其中：
# 1.前缀 名称是一个字符串，可以包含字母（大写或小写），数字，下划线 '_' ，点 '.' 和/或破折号 '-' 。前缀名称 必须 以字母开头。
# 2.域 为 '@leetcode.com' 。
# 以任何顺序返回结果表。
# +---------+-----------+-------------------------+
# | user_id | name      | mail                    |
# +---------+-----------+-------------------------+
# | 1       | Winston   | winston@leetcode.com    |
# | 3       | Annabelle | bella-@leetcode.com     |
# | 4       | Sally     | sally.come@leetcode.com |
# +---------+-----------+-------------------------+


def valid_emails(users: pd.DataFrame) -> pd.DataFrame:
    # \w: 匹配字母,数字,下划线
    return users[users["mail"].str.match(r"^[a-zA-Z][\w\.\-]*@leetcode\.com$")]
