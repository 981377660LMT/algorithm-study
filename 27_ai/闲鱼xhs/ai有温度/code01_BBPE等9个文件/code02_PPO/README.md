> 小红书 @AI有温度

此PPO项目约消耗12G显存

### 流程概述

- 为了节省显存，Actor 与 Critic 模型共享底座（qwen0.5b + LoRA 适配器）共同更新，仅通过不同的输出头实现功能区分：
    - Actor 通过底座模型的logits计算生成 token 的概率分布，用于策略更新
    - Critic 通过在底座后接一个线性层价值头 v_head用于估计状态价值，用于评估策略优劣
- Reward 模型使用Erlangshen-Roberta-330M-Sentiment
- Reference 使用qwen0.5b

### 模型下载

此项目为了让大家都可以把模型跑起来，所以actor模型使用qwen2.5-0.5B（如果你的卡多可以用大模型），reward模型使用Erlangshen-Roberta-330M-Sentiment。首先需要下载这两个模型：

- [https://www.modelscope.cn/models/Qwen/Qwen2.5-0.5B-Instruct](https://www.modelscope.cn/models/Qwen/Qwen2.5-0.5B-Instruct)
- [https://huggingface.co/IDEA-CCNL/Erlangshen-Roberta-330M-Sentiment](https://huggingface.co/IDEA-CCNL/Erlangshen-Roberta-330M-Sentiment)

### 训练数据

- data/train_data.json 此数据是由大模型直接生成，仅用于学习使用。
- 大家也可以自行下载数据集进行训练调试

### config.py配置文件

使用前需要修改模型地址与模型储存地址

### 运行顺序

1. main.py进行训练
2. inference.py进行推理
3. ppo.py实现了训练的过程与步骤