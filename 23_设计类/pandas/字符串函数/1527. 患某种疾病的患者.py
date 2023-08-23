import pandas as pd

data = [
    [1, "Daniel", "YFEV COUGH"],
    [2, "Alice", ""],
    [3, "Bob", "DIAB100 MYOP"],
    [4, "George", "ACNE DIAB100"],
    [5, "Alain", "DIAB201"],
]
Patients = pd.DataFrame(data, columns=["patient_id", "patient_name", "conditions"]).astype(
    {"patient_id": "int64", "patient_name": "object", "conditions": "object"}
)


# 查询患有 I 类糖尿病的患者 ID （patient_id）、患者姓名（patient_name）以及其患有的所有疾病代码（conditions）。
# I 类糖尿病的代码总是包含前缀 DIAB1 。
# +------------+--------------+--------------+pandas
# | patient_id | patient_name | conditions   |
# +------------+--------------+--------------+
# | 3          | Bob          | DIAB100 MYOP |
# | 4          | George       | ACNE DIAB100 |
# +------------+--------------+--------------+
# 解释：Bob 和 George 都患有代码以 DIAB1 开头的疾病。


def find_patients(patients: pd.DataFrame) -> pd.DataFrame:
    ...
