# ROC 与 AUC

AUC 是评价二分类的指标。

- **ROC 曲线**：不管你把“判定阈值”设得多高或多低，模型在**“抓坏人”（召回率）**和**“误伤好人”（误报率）**之间的权衡能力的**全景图**。
- **AUC 值**：给 ROC 曲线打的一个分（0 到 1 之间）。**分数越高，模型区分“好人”和“坏人”的能力越强**。AUC = 0.8 意味着：随机挑一个正样本和一个负样本，模型给正样本打分更高的概率是 80%。

---

### 1. 核心痛点：为什么不用准确率（Accuracy）？

如果有 99 个好人，1 个坏人。模型只要无脑喊“全是好人”，准确率就有 99%。但这个模型是废的，因为它没抓住那 1 个坏人。
**ROC/AUC 就是为了解决样本不平衡和阈值选择困难而生的。**

### 2. ROC 曲线怎么画？

它是二维坐标系里的一条线：

- **横轴 (FPR)**：假正率（误报率）。意思是：**没病的被你说成有病的**比例。（越小越好）
- **纵轴 (TPR)**：真正率（召回率）。意思是：**真有病的被你查出来**的比例。（越大越好）

**怎么动？**
你调整分类器的**阈值**（比如从 0.1 到 0.9），每调一次，就会算出一对 (FPR, TPR)，在图上点一个点。把这些点连起来，就是 ROC 曲线。

### 3. AUC (Area Under Curve) 怎么看？

AUC 就是 ROC 曲线下面的**面积**。

- **AUC = 0.5**：瞎猜（抛硬币）。曲线是一条对角线。
- **0.5 < AUC < 1**：比瞎猜好。曲线越靠近**左上角**（误报少，召回高），面积越大，模型越牛。
- **AUC = 1**：完美神级模型（现实中几乎不存在）。

### 4. 代码实战 (Python/Sklearn)

```python
import numpy as np
from sklearn.metrics import roc_curve, auc
from sklearn.model_selection import train_test_split
from sklearn.linear_model import LogisticRegression
from sklearn.datasets import make_classification
import matplotlib.pyplot as plt

# 1. 搞点假数据
X, y = make_classification(n_samples=1000, n_classes=2, random_state=42)
X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.3, random_state=42)

# 2. 训练模型
model = LogisticRegression()
model.fit(X_train, y_train)

# 3. 预测概率 (关键！要的是 predict_proba 拿出来的概率值，而不是 predict 出来的 0/1)
y_score = model.predict_proba(X_test)[:, 1]

# 4. 计算 FPR, TPR, 阈值
fpr, tpr, thresholds = roc_curve(y_test, y_score)

# 5. 计算 AUC
roc_auc = auc(fpr, tpr)
print(f"一针见血的结果 -> AUC: {roc_auc:.4f}")

# (可选) 简单的绘图逻辑
plt.figure()
plt.plot(fpr, tpr, color='darkorange', lw=2, label=f'ROC curve (area = {roc_auc:.2f})')
plt.plot([0, 1], [0, 1], color='navy', lw=2, linestyle='--') # 瞎猜基准线
plt.xlim([0.0, 1.0])
plt.ylim([0.0, 1.05])
plt.xlabel('False Positive Rate (误报率)')
plt.ylabel('True Positive Rate (召回率)')
plt.title('Receiver Operating Characteristic')
plt.legend(loc="lower right")
# plt.show() # 如果在 Jupyter 或本地运行可打开注释
```
