训练循环就是：准备数据 → 前向传播 → 计算损失 → 反向传播 → 更新参数，不断重复。代码不到 100 行，但让模型从"一无所知"变成能预测下一个词。

```
1. 加载数据 → 分词 → 转 Tensor → 分割 train/val

2. 训练循环：
   for step in range(max_iters):
       x, y = get_batch('train')     # 获取数据
       logits, loss = model(x, y)    # 前向传播
       optimizer.zero_grad()         # 清零梯度
       loss.backward()               # 反向传播
       optimizer.step()              # 更新参数

3. 保存模型 → torch.save()
```

- 超参数

```py
# 超参数配置
h_params = {
    # 模型架构
    "d_model": 80,           # 嵌入维度（小模型用小值）
    "num_blocks": 6,         # Transformer 块数量
    "num_heads": 4,          # 注意力头数

    # 训练配置
    "batch_size": 2,         # 每次训练多少个样本
    "context_length": 128,   # 上下文长度（序列长度）
    "max_iters": 500,        # 训练多少步
    "learning_rate": 1e-3,   # 学习率

    # 正则化
    "dropout": 0.1,          # Dropout 概率

    # 评估配置
    "eval_interval": 50,     # 每多少步评估一次
    "eval_iters": 10,        # 评估时用多少个 batch

    # 设备
    "device": "cuda" if torch.cuda.is_available() else "cpu",

    # 随机种子（可复现）
    "TORCH_SEED": 1337
}
```

- model.train() vs model.eval()
  模式 Dropout BatchNorm
  model.train() 随机丢弃 使用 batch 统计量
  model.eval() 不丢弃 使用全局统计量
  评估时必须用 model.eval()，否则结果会有随机性。

- AdamW 是目前最常用的优化器，结合了：

  Momentum：考虑历史梯度方向
  自适应学习率：每个参数有自己的学习率
  Weight Decay：L2 正则化，防止过拟合

  **现代大模型训练几乎都用 AdamW。**

- epoch

```py
# 训练循环
for step in range(h_params['max_iters']):
  # 定期评估
  if step % h_params['eval_interval'] == 0 or step == h_params['max_iters'] - 1:
      losses = estimate_loss()
      print(f'Step: {step}, '
            f'Training Loss: {losses["train"]:.3f}, '
            f'Validation Loss: {losses["valid"]:.3f}')

  # 1. 获取一个 batch
  xb, yb = get_batch('train')

  # 2. 前向传播
  logits, loss = model(xb, yb)

  # 3. 反向传播
  optimizer.zero_grad(set_to_none=True)  # 清零梯度
  loss.backward()                         # 计算梯度

  # 4. 更新参数
  optimizer.step()
```

损失从 ~10.8 下降到 ~2.8
验证损失略高于训练损失（正常，因为是没见过的数据）
如果验证损失开始上升，说明过拟合了
