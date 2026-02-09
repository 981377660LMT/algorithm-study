https://my.feishu.cn/wiki/M1Cew0iYaiH9jfkeD5XcXEuZn3d

1. BPE、BBPE、WordPiece、Unigram
   LLM 使用的 Subword 分词
   gpt2：https://github.com/huggingface/transformers/blob/05260a1fc1c8571a2b421ce72b680d5f1bc3e5a4/src/transformers/models/gpt2/tokenization_gpt2.py#L75
2. transformer
   https://my.feishu.cn/wiki/YL6XwwuDeitwPokJrzvcz5rknsf
   BERT、GPT
3. 解码技术 Decoding
   https://my.feishu.cn/wiki/I4OVwA7dYiC3lTk0ffpczFlSnxe
   从概率分布选词

   几种常见的解码方式
   1. Greedy Decoding（贪心解码）：每一步都选取概率最高的词。
   2. Beam Search：保留若干条最有希望的候选序列同步扩展。
   3. Sampling（随机采样）：按概率分布直接采样。
   4. Top-k 采样：在概率最高的k个候选中做随机采样。
   5. Top-p（核采样）：在累积概率达到p的最小词集合中做随机采样。
   6. Temperature 调整：对logits（未经过 softmax 的概率）做缩放，控制模型“确定性”或
      “随机性”。(`在softmax之前除以temperature`)

4. 大语言模型架构
   根据训练任务分类
   - `Encoder-only（编码器）：如BERT，适合理解任务`。
     https://my.feishu.cn/wiki/QIRDwHkqKiEStbkLczjcHn2hnQh
   - Encoder-Decoder（编码器-解码器）：如T5，适合生成任务。
   - Non-causal Decoder-only（非自回归解码器）
   - **主流：**`Causal Decoder-only（自回归解码器）：如GPT，适合生成任务。`

5. LLM 预训练
