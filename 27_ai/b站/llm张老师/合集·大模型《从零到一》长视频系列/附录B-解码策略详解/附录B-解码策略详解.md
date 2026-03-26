解码策略决定了模型如何从概率分布中选择下一个 token。Greedy 选最高概率，Sampling 随机采样，Top-K/Top-P 限制候选范围，Temperature 调节分布锐度。不同策略适合不同场景。

- 现代 LLM（ChatGPT 等）通常不用 Beam Search，而用 Sampling + Top-P/K

---

Beam Search 翻译、摘要（需要准确性）
Sampling + Top-P 对话、创作（需要多样性）
Greedy (T=0) 代码、数学（需要确定性）
现代 LLM（ChatGPT、Claude）主要用 Sampling + Top-P。
