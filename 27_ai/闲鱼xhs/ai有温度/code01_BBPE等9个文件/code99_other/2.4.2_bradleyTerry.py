import torch
import torch.nn as nn
import torch.optim as optim


class BTModel(nn.Module):
    def __init__(self, N):
        super(BTModel, self).__init__()
        self.reward = nn.Parameter(torch.ones(N))

    def forward_exp(self, chosen_id, rejected_id):
        reward_chosen = torch.exp(self.reward[chosen_id])
        reward_rejected = torch.exp(self.reward[rejected_id])
        return reward_chosen / (reward_chosen + reward_rejected)

    def forward_sigmoid(self, chosen_id, rejected_id):
        reward_chosen = self.reward[chosen_id]
        reward_rejected = self.reward[rejected_id]
        return torch.sigmoid(reward_chosen - reward_rejected)

    def loss(self, pred, label):
        return -torch.log(pred) if label == 1 else -torch.log(1 - pred)

# 给出4个选手，
N = 4
model = BTModel(4)
print('reward:', model.reward)
# 0 > 1
# 2 > 3
# 1 > 3
datas = [(0, 1, 1), (2, 3, 1), (1, 3, 1)] # 比赛数据，也可以认为是偏好数据
optimizer = optim.SGD(model.parameters(), lr=0.01)

# 训练模型
loss_fn = nn.BCELoss()
for i in range(100):
    total_loss = 0
    for data in datas:
        id_i, id_j, label = data
        optimizer.zero_grad()
        pred = model.forward_sigmoid(id_i, id_j)
        # pred = model.forward_exp(id_i, id_j)
        loss = model.loss(pred, torch.tensor(label, dtype=torch.float32))
        loss.backward()
        optimizer.step()

        total_loss += loss.item()
    if i%10==0 : print(f"Epoch {i}, Loss: {total_loss}")

# 输出每个选手的强度参数
print(model.reward)