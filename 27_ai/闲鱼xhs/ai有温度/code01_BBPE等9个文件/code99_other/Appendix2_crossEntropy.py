import torch
import numpy as np

# 定义二分类交叉熵损失函数
def binary_cross_entropy_loss(y_pred, y_true):
    # 1. Sigmoid 转换：将模型输出映射到 [0,1] 概率区间
    y_pred = 1 / (1 + torch.exp(-y_pred))  
    
    # 2. 数值稳定：限制预测值范围，避免 log(0) 导致梯度爆炸/NaN
    epsilon = 1e-7
    y_pred = torch.clamp(y_pred, min=epsilon, max=1.0 - epsilon)  
    
    # 3. 计算交叉熵损失：逐样本计算后求平均
    loss = -(y_true * torch.log(y_pred) + (1 - y_true) * torch.log(1 - y_pred))
    return loss.mean()  

class BinaryClassifier(torch.nn.Module):
    def __init__(self):
        super(BinaryClassifier, self).__init__()
        self.layer = torch.nn.Linear(2, 1)  # 输入维度 2，输出维度 1（二分类）
    
    def forward(self, x):
        return self.layer(x)

model = BinaryClassifier()
features = np.array([[0.2, 0.4], [0.6, 0.8], [0.1, 0.3], [0.5, 0.7]])
labels = np.array([1, 0, 1, 0])

features = torch.tensor(features, dtype=torch.float32)
labels = torch.tensor(labels, dtype=torch.float32).unsqueeze(1)  # 适配模型输出维度

# 前向计算与损失求解
outputs = model(features)
loss = binary_cross_entropy_loss(outputs, labels)
print("二分类交叉熵损失：", loss.item())