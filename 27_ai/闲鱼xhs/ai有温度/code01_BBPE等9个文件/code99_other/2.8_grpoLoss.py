def compute_loss(self, model, inputs, return_outputs=False, num_items_in_batch=None):
    """
    参数:
        model: 当前训练的策略模型
        inputs: 包含输入数据的字典，包括prompt、completion等信息
        num_items_in_batch: 批次中的项目数量

    返回:
        计算得到的损失值
    """
    if return_outputs:
        raise ValueError("GRPOTrainer不支持返回输出结果")
    
    # 提取输入数据中的prompt和completion部分
    prompt_ids, prompt_mask = inputs["prompt_ids"], inputs["prompt_mask"]
    completion_ids, completion_mask = inputs["completion_ids"], inputs["completion_mask"]
    
    # 拼接prompt和completion形成完整输入
    input_ids = torch.cat([prompt_ids, completion_ids], dim=1)
    attention_mask = torch.cat([prompt_mask, completion_mask], dim=1)
    # 只需要计算completion部分的logits
    logits_to_keep = completion_ids.size(1)
    
    # 获取每个token的对数概率
    per_token_logps = self._get_per_token_logps(model, input_ids, attention_mask, logits_to_keep)
    
    # 计算模型与参考模型之间的KL散度（用于约束策略更新）
    ref_per_token_logps = inputs["ref_per_token_logps"].detach()
    # 使用KL散度的近似计算：f(x) = e^x - x - 1，其中x = ref_logp - current_logp
    per_token_kl = torch.exp(ref_per_token_logps - per_token_logps) - (ref_per_token_logps - per_token_logps) - 1
    
    # 计算优势函数（指导策略更新的信号）
    advantages = inputs["advantages"]  # 形状: (批次大小*生成数量,)
    
    # 1. 计算策略比（当前策略概率 / 参考策略概率）= exp(current_logp - ref_logp)
    ratio = torch.exp(per_token_logps - ref_per_token_logps.detach())  # 形状: (批次大小*生成数量, completion长度)
    
    # 2. 计算clip后的策略比（限制在 [1-ε, 1+ε] 区间）
    # 注意：需要将advantages的维度扩展以匹配ratio（增加completion长度维度）
    advantages_unsqueezed = advantages.unsqueeze(1)  # 形状: (批次大小*生成数量, 1)
    epsilon = self.epsilon  # 从类参数中获取clip系数（需在初始化时定义，如self.epsilon=0.2）
    clipped_ratio = torch.clamp(ratio, 1 - epsilon, 1 + epsilon)  # 形状同上
    
    # 3. 计算原始策略比×优势 和 clip后策略比×优势，取两者最小值（对应公式中的min(...)）
    surr1 = ratio * advantages_unsqueezed  # 原始项
    surr2 = clipped_ratio * advantages_unsqueezed  # clip后的项
    per_token_surrogate = torch.min(surr1, surr2)  # 取最小值，形状同上
    
    # 结合KL惩罚项，取负值是因为要最小化损失（对应最大化目标）
    per_token_loss = -(per_token_surrogate - self.beta * per_token_kl)
    
    # 应用掩码计算批次平均损失：只计算completion部分的有效token
    loss = ((per_token_loss * completion_mask).sum(dim=1) / completion_mask.sum(dim=1)).mean()
    
    # 记录指标：完成序列的平均长度
    completion_length = self.accelerator.gather_for_metrics(completion_mask.sum(1)).float().mean().item()
    self._metrics["completion_length"].append(completion_length)
    
    # 记录指标：平均KL散度
    mean_kl = ((per_token_kl * completion_mask).sum(dim=1) / completion_mask.sum(dim=1)).mean()
    self._metrics["kl"].append(self.accelerator.gather_for_metrics(mean_kl).mean().item())
    
    # 新增：记录clip相关指标（可选，用于监控clip生效情况）
    clip_fraction = (ratio < (1 - epsilon)).float().mean() + (ratio > (1 + epsilon)).float().mean()
    self._metrics["clip_fraction"].append(self.accelerator.gather_for_metrics(clip_fraction).mean().item())
    
    return loss

def _get_per_token_logps(self, model, input_ids, attention_mask, logits_to_keep):
    """
    计算每个token的对数概率
    
    参数:
        model: 用于计算logits的模型
        input_ids: 输入的token ids
        attention_mask: 注意力掩码
        logits_to_keep: 需要保留的logits数量（completion部分长度）
        
    返回:
        每个token的对数概率张量
    """
    # 加1是因为最后一个logit会被排除（不需要预测下一个token）
    logits = model(
        input_ids=input_ids, attention_mask=attention_mask, logits_to_keep=logits_to_keep + 1
    ).logits  # 形状: (批次大小, 序列长度, 词汇表大小)
    
    # 排除最后一个logit（因为它对应下一个token的预测，我们不需要）
    logits = logits[:, :-1, :]  # 形状: (批次大小, 序列长度-1, 词汇表大小)
    
    # 计算输入token的对数概率，使用循环减少内存峰值
    per_token_logps = []
    # 遍历每个样本的logits和对应的completion token ids
    for logits_row, input_ids_row in zip(logits, input_ids[:, -logits_to_keep:]):
        # 对logits进行softmax得到概率分布，再取对数
        log_probs = logits_row.log_softmax(dim=-1)
        # 提取每个token对应的对数概率
        token_log_prob = torch.gather(log_probs, dim=1, index=input_ids_row.unsqueeze(1)).squeeze(1)
        per_token_logps.append(token_log_prob)
    
    # 堆叠所有样本的结果
    return torch.stack(per_token_logps)

def _prepare_inputs(self, inputs: dict[str, Union[torch.Tensor, Any]]) -> dict[str, Union[torch.Tensor, Any]]:
    """
    准备模型输入数据，处理prompt和生成的completion，计算奖励和优势函数
    
    参数:
        inputs: 原始输入数据字典
        
    返回:
        处理后的输入字典，包含训练所需的所有信息
    """
    device = self.accelerator.device
    
    # 提取prompt文本
    prompts = [x["prompt"] for x in inputs]
    # 应用聊天模板处理prompt（如果是对话格式）
    prompts_text = [maybe_apply_chat_template(example, self.processing_class)["prompt"] for example in inputs]
    
    # 对prompt进行tokenize处理
    prompt_inputs = self.processing_class(
        prompts_text, return_tensors="pt", padding=True, padding_side="left", add_special_tokens=False
    )
    prompt_inputs = super()._prepare_inputs(prompt_inputs)
    prompt_ids, prompt_mask = prompt_inputs["input_ids"], prompt_inputs["attention_mask"]
    
    # 如果设置了最大prompt长度，截断过长的prompt
    if self.max_prompt_length is not None:
        prompt_ids = prompt_ids[:, -self.max_prompt_length :]
        prompt_mask = prompt_mask[:, -self.max_prompt_length :]
    
    # 生成completion（响应内容）
    with unwrap_model_for_generation(self.model, self.accelerator) as unwrapped_model:
        prompt_completion_ids = unwrapped_model.generate(
            prompt_ids, attention_mask=prompt_mask, generation_config=self.generation_config
        )
    
    # 计算prompt长度并提取completion部分
    prompt_length = prompt_ids.size(1)
    prompt_ids = prompt_completion_ids[:, :prompt_length]  # 提取prompt部分
    completion_ids = prompt_completion_ids[:, prompt_length:]  # 提取completion部分
    # 为每个生成结果重复prompt_mask（因为一个prompt可能生成多个completion）
    prompt_mask = prompt_mask.repeat_interleave(self.num_generations, dim=0)
    
    # 生成completion的掩码：在第一个EOS token后截断
    is_eos = completion_ids == self.processing_class.eos_token_id  # 标记EOS token位置
    # 初始化EOS位置为completion长度（默认没有EOS）
    eos_idx = torch.full((is_eos.size(0),), is_eos.size(1), dtype=torch.long, device=device)
    # 找到有EOS的样本，并记录第一个EOS的位置
    eos_idx[is_eos.any(dim=1)] = is_eos.int().argmax(dim=1)[is_eos.any(dim=1)]
    # 生成序列索引，用于创建掩码
    sequence_indices = torch.arange(is_eos.size(1), device=device).expand(is_eos.size(0), -1)
    # 创建completion掩码：EOS之前的token有效，之后无效
    completion_mask = (sequence_indices <= eos_idx.unsqueeze(1)).int()
    
    # 拼接prompt_mask和completion_mask用于logit计算
    attention_mask = torch.cat([prompt_mask, completion_mask], dim=1)  # 形状: (批次大小*生成数量, prompt长度+completion长度)
    
    # 需要保留的logits数量（completion部分长度）
    logits_to_keep = completion_ids.size(1)
    
    # 使用参考模型计算log概率（不计算梯度）
    with torch.inference_mode():
        ref_per_token_logps = self._get_per_token_logps(
            self.ref_model, prompt_completion_ids, attention_mask, logits_to_keep
        )
    
    # 解码生成的completion文本
    completions = self.processing_class.batch_decode(completion_ids, skip_special_tokens=True)
    # 如果是对话格式，包装成对话结构
    if is_conversational(inputs[0]):
        completions = [[{"role": "assistant", "content": completion}] for completion in completions]
    
    # 为每个prompt重复对应次数（因为一个prompt生成了多个completion）
    prompts = [prompt for prompt in prompts for _ in range(self.num_generations)]
    
    # 初始化奖励张量：(批次大小*生成数量, 奖励函数数量)
    rewards_per_func = torch.zeros(len(prompts), len(self.reward_funcs), device=device)
    # 计算每个奖励函数的奖励值
    for i, (reward_func, reward_processing_class) in enumerate(
        zip(self.reward_funcs, self.reward_processing_classes)
    ):
        # 如果是神经网络奖励模型
        if isinstance(reward_func, nn.Module):
            if is_conversational(inputs[0]):
                # 处理对话格式：拼接prompt和completion作为消息
                messages = [{"messages": p + c} for p, c in zip(prompts, completions)]
                texts = [apply_chat_template(x, reward_processing_class)["text"] for x in messages]
            else:
                # 非对话格式：直接拼接prompt和completion
                texts = [p + c for p, c in zip(prompts, completions)]
            
            # 对奖励模型的输入进行tokenize
            reward_inputs = reward_processing_class(
                texts, return_tensors="pt", padding=True, padding_side="right", add_special_tokens=False
            )
            reward_inputs = super()._prepare_inputs(reward_inputs)
            
            # 计算奖励值（不计算梯度）
            with torch.inference_mode():
                rewards_per_func[:, i] = reward_func(** reward_inputs).logits[:, 0]  # 形状: (批次大小*生成数量,)
        else:
            # 处理自定义奖励函数（非神经网络）
            # 准备奖励函数需要的其他参数
            reward_kwargs = {key: [] for key in inputs[0].keys() if key not in ["prompt", "completion"]}
            for key in reward_kwargs:
                for example in inputs:
                    # 每个参数重复num_generations次（与生成数量匹配）
                    reward_kwargs[key].extend([example[key]] * self.num_generations)
            
            # 调用自定义奖励函数
            output_reward_func = reward_func(prompts=prompts, completions=completions, **reward_kwargs)
            rewards_per_func[:, i] = torch.tensor(output_reward_func, dtype=torch.float32, device=device)
    
    # 汇总所有奖励函数的奖励值
    rewards = rewards_per_func.sum(dim=1)
    
    # 计算每组（同一prompt的多个生成结果）的奖励均值和标准差
    mean_grouped_rewards = rewards.view(-1, self.num_generations).mean(dim=1)
    std_grouped_rewards = rewards.view(-1, self.num_generations).std(dim=1)
    
    # 标准化奖励以计算优势函数（使每组内的奖励均值为0，标准差为1）
    mean_grouped_rewards = mean_grouped_rewards.repeat_interleave(self.num_generations, dim=0)
    std_grouped_rewards = std_grouped_rewards.repeat_interleave(self.num_generations, dim=0)
    advantages = (rewards - mean_grouped_rewards) / (std_grouped_rewards + 1e-4)  # 加小epsilon避免除零
    
    # 返回处理后的所有输入数据
    return {
        "prompt_ids": prompt_ids,
        "prompt_mask": prompt_mask,
        "completion_ids": completion_ids,
        "completion_mask": completion_mask,
        "ref_per_token_logps": ref_per_token_logps,
        "advantages": advantages,
    }