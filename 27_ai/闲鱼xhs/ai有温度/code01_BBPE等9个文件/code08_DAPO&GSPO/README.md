# README

GSPO在GRPO的基础上进行了改进，训练流程代码相同，只更改 GRPO 的配置就可以实现 GSPO：


```python title="GRPO 的基础上进行改进"
from trl import GRPOConfig
training_args = GRPOConfig(
    importance_sampling_level="sequence", # 如果注掉这行，就是 DAPO
    loss_type="grpo",
    beta=0.0,  # GSPO set KL regularization to zero: https://github.com/volcengine/verl/pull/2775#issuecomment-3131807306 
    epsilon=3e-4,  # GSPO paper (v2), section 5.1
    epsilon_high=4e-4,  # GSPO paper (v2), section 5.1
    gradient_accumulation_steps=1,
    steps_per_generation=4,  # partition rollout batch into 4 mini-batches. GSPO paper (v2), section 5.1. Must be 4 times gradient_accumulation_steps
)

```


**1. ****`importance_sampling_level="sequence"`****  **

**重要性采样（Importance Sampling）** 是强化学习中用于解决“策略分布偏移”的技术（即训练时的策略与生成样本时的策略可能不同，需要通过权重校正）。 &#x20;

`sequence` 表示在**整个序列级别**进行重要性采样，即对模型生成的完整序列（而非单个 token）计算校正权重。这更符合文本生成任务的特性（序列整体语义更重要）。 &#x20;

**2. ****`loss_type="grpo"`****  **

指定训练时使用的损失函数类型为 **GRPO 损失**。 &#x20;

GRPO 的损失函数核心是结合“奖励信号”和“策略熵/正则项”，引导模型向高奖励的生成方向更新，同时避免策略突变（保证训练稳定性）。 &#x20;

**3. ****`beta=0.0`****  **

`beta` 是 **KL 散度正则化的系数**，用于限制新策略与旧策略的差异（KL 散度越大，惩罚越强，防止策略更新过于激进）。 &#x20;

注释提到“GSPO 将 KL 正则化设为 0”，这是因为 GSPO 更依赖其他机制（如下文的 `epsilon`）控制策略更新幅度，而非 KL 散度，因此将 `beta` 设为 0 以关闭 KL 正则化。 &#x20;

**4. ****`epsilon=3e-4`**** 和 ****`epsilon_high=4e-4`****  **

这两个参数来自 GSPO 论文（v2 版本），用于**控制策略更新的剪辑范围**（类似 PPO 中的 `clip_param`），防止更新幅度过大导致训练不稳定。 &#x20;

`epsilon` 是基础剪辑阈值，`epsilon_high` 是更高的阈值（可能在某些场景下动态调整，平衡探索与利用）。 &#x20;

数值越小，对策略更新的限制越严格（更稳定但可能收敛慢）；数值越大，允许策略更大幅度更新（可能更快收敛但风险更高）。 &#x20;

这些参数的设置本质上是在 GRPO 框架下融入 GSPO 的设计理念：通过关闭 KL 正则化（`beta=0`），改用 `epsilon` 剪辑控制策略更新，并拆分批次处理（`steps_per_generation`），以在文本生成任务中实现更稳定、高效的强化学习训练。

参数的具体数值参考了 GSPO 论文的推荐配置，适用于需要平衡探索与利用的场景。


#### 训练数据可在 GRPO 中获取，数据集相同，只是训练方法略有差异