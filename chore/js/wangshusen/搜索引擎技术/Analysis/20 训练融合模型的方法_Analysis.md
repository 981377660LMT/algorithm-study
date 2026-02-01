# 20 训练融合模型的方法：pointwise 保值、pairwise 保序、listwise 保“前排”

上一章把融合模型抽象为 $f(x,z;w)$：

- $x$：非个性化特征
- $z$：个性化特征
- $w$：模型参数

本章的核心是：训练融合模型不是“拟合一个数值”，而是“学一个排序”。

因此本章依次介绍三种训练范式：

- pointwise：把每条样本当独立回归/分类
- pairwise：把样本两两成对，优化正逆序
- listwise：把排序列表整体考虑，直接面向 NDCG 等排序指标

---

## 20.1 pointwise：实现简单、容易作为 baseline，但不直接优化排序

pointwise 把每个 $(u,q,d_i)$ 看作独立样本，目标是让：

- $p_i=f(x_i,z_i;w)$ 拟合 $y_i$

常见损失：

- 均方误差：$\frac{1}{k}\sum_i (p_i-y_i)^2$
- 若 $y\in[0,1]$：交叉熵 $\frac{1}{k}\sum_i \big[y_i\ln p_i+(1-y_i)\ln(1-p_i)\big]$

本章给了一个非常关键的直觉：

- 排序关心的是“相对顺序”，不是“绝对分值”
- 给所有文档分数同时加一个常数 $\delta$，排序不变，但 pointwise 会认为误差变大

因此 pointwise 训练的本质是“保值”，并不天然契合排序。

但它的工程价值很高：

- 容易实现、不容易出错
- 在项目初期非常适合作为线上 baseline

---

## 20.2 pairwise：用 RankNet 最大化正序对，最小化逆序对

pairwise 的基本设定：

- 若 $y_i>y_j$，则文档 $d_i$ 应排在 $d_j$ 前
- 模型打分 $p_i, p_j$ 产生一个“顺序是否一致”的判断

### 20.2.1 RankNet 损失：鼓励 $p_i-p_j$ 在正确方向变大

RankNet 对每个正序对 $(i,j):y_i>y_j$ 定义损失：

$$
 c_{ij}=\ln\big(1+\exp(-(p_i-p_j))\big)
$$

- 若 $p_i\ge p_j$（正序对），损失较小
- 若 $p_i<p_j$（逆序对），损失变大并推动梯度纠正

本章进一步推导了梯度可写成“lambda”形式：

- $\lambda_{ij} = \frac{\partial c_{ij}}{\partial p_i}$

最终对每个样本 $i$ 的梯度可以理解为：

- 来自“它作为赢家/输家”参与的所有 pair 的 lambda 累积

### 20.2.2 RankNet 的问题：不看差距、不看位置

本章指出 RankNet 有两个关键缺陷（这也是为什么后面要 listwise）：

1. **不考虑 $y_i$ 与 $y_j$ 的差距**：
   - $y_i$ 只比 $y_j$ 大一点，和大很多，在损失里权重相同

2. **不考虑位置重要性**：
   - top1/top2 的交换对用户体验影响远大于第 200/201 的交换
   - RankNet 对不同位置的 pair 没有天然区分

---

## 20.3 listwise：LambdaRank 用 NDCG 增量给 pair 加权

LambdaRank 的核心思想很“朴素”：

- 仍然用 pairwise 的形式优化
- 但给每个 pair 的损失乘上“交换这两个文档带来的 NDCG 变化”

### 20.3.1 NDCG 与交换增量

DCG：

$$
DCG@n = \sum_{i=1}^n \frac{2^{y_i}-1}{\log_2(i+1)}
$$

NDCG：

$$
NDCG@n = \frac{DCG@n}{IDCG@n}
$$

交换 $i$ 与 $j$ 的位置，NDCG 的变化量可写为 $\Delta_{ij}$（本章给出了推导）。

### 20.3.2 加权后的损失：把优化重心放到“前排 + 大差距”

LambdaRank 的损失形式可以理解为：

$$
 c = \sum_{(i,j):y_i>y_j} |\Delta_{ij}|\;\ln\big(1+\exp(-(p_i-p_j))\big)
$$

这会带来两个效果：

- 前排位置的错误更“贵”，模型更努力修正
- 真实差距大的交换对 NDCG 影响更大，也会得到更大权重

因此当评价指标是 NDCG 时，listwise 往往优于单纯 pairwise。

---

## 20.4 工程落地建议（按“先稳后强”的路线）

- 先上 pointwise 作为 baseline：保证训练链路正确、特征正确、线上可用
- 再尝试 pairwise：验证是否能提升正逆序比与线上指标
- 最后上 listwise（LambdaRank）：在明确要优化 NDCG、且对头部排序要求高时通常更值得

另外要注意训练成本：

- pair/listwise 都会引入“样本成对”带来的计算与采样策略问题
- 工程上通常要做 pair 采样/截断，避免 $O(k^2)$ 爆炸

---

## 20.5 常见坑与检查清单

- 只看离线 loss：排序模型的离线 loss 下降不保证线上业务指标提升
- pair 采样不合理：只采到容易 pair，模型学不到真正难例
- listwise 权重实现有误：$\Delta_{ij}$ 计算错误会导致训练方向偏离 NDCG

---

## 20.6 练习（建议动手）

1. 用同一份训练数据分别训练 pointwise 与 RankNet，比较离线 AUC 与 NDCG 的变化。
2. 对一个 query 的 top50 列表，手算几次交换的 $|\Delta_{ij}|$，直观看到“前排更重要”。
3. 设计一个 pair 采样策略：优先采样相邻位置、或采样跨档位 pair，并解释你的理由。

---

## 20.7 与其它章节的连接

- 与第 19 章：融合模型训练目标 $y$ 的定义（满意度+行为）决定了你最终优化的是什么
- 与第 17 章：召回训练里也会用 pairwise/batch 内负采样，方法论是相通的
