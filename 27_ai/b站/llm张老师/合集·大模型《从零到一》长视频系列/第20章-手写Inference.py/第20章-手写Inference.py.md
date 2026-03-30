推理就是：加载模型 → 输入 prompt → 自回归生成 → 解码输出。代码只有 30 行，但这是模型"开口说话"的时刻。

```py
1. 加载模型
   checkpoint = torch.load('model.ckpt')
   model.load_state_dict(checkpoint['model_state_dict'])
   model.eval()

2. 编码 prompt
   start_ids = encoding.encode(prompt)
   x = torch.tensor(start_ids)[None, ...]

3. 生成
   with torch.no_grad():
       y = model.generate(x, max_new_tokens=200)
   output = encoding.decode(y[0].tolist())
```

- 模型学到的是数据中的模式，而不是"理解"内容。
- 常见问题
  - 生成重复内容，模型不断重复相同的词或短语
    原因：
    `Temperature 太低`
    训练数据本身有重复
    模型过拟合

    解决：
    提高 Temperature
    使用 Top-K 或 Top-P 采样
    添加 repetition penalty

  - 生成乱码，输出是乱码或不连贯的文本
    原因：
    模型训练不足
    prompt 不在训练分布内
    Temperature 太高

    解决：
    训练更多步
    使用更合适的 prompt
    降低 Temperature

  - 速度太慢，生成每个 token 都很慢
    原因：
    没有使用 GPU
    没有 KV Cache
    模型太大

    解决：
    使用 GPU（如果有）
    实现 KV Cache
    使用更小的模型

- max_new_tokens 最大生成长度 50-500
  temperature 随机性控制 0.5-0.8
  top_k 限制候选词数量 50-100

---

model.py - 完整模型定义
train.py - 训练脚本
inference.py - 推理脚本
