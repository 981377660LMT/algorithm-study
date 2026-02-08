import torch
from config import Config
from tools import Tools

class PPO:
    def __init__(self, actor_critic_model, config: Config, actor_critic_opt):
        self.actor_critic_model = actor_critic_model
        self.config = config
        self.actor_critic_opt = actor_critic_opt

    def train(self, prompt_generate_ids, attention_mask, prob_refs, reward, tools: Tools):
        with torch.no_grad():
            _, old_values = self.actor_critic_model(prompt_generate_ids, attention_mask, tools)  # 计算每个token的价值
        for _ in range(self.config.ppo_epochs):
            # 获得actor_critic模型新的probs和token对应的价值
            new_probs, new_values = self.actor_critic_model(prompt_generate_ids, attention_mask, tools)
            # 计算奖励值
            rewards, non_score_rewards = self.compute_rewards(reward, new_probs, prob_refs)  # 计算reward
            loss = self.loss(new_probs=new_probs, old_values=old_values, new_values=new_values,
                             rewards=rewards, old_probs=prob_refs)

            self.actor_critic_opt.zero_grad()
            loss.backward()
            self.actor_critic_opt.step()
            print(loss)

    def compute_rewards(self, scores, probs, ref_probs):
        """计算reward值,确保奖励张量在正确设备上"""
        rewards, non_score_rewards = [], []
        # 获取设备（从probs中获取，因为probs来自模型输出，在正确设备上）
        device = probs[0].device if probs else torch.device('cpu')
        
        for score, prob, ref_prob in zip(scores, probs, ref_probs):
            kl = prob - ref_prob  # (seq_len, )
            non_score_reward = -self.config.kl_ctl_value * kl  # (seq_len, )
            # 确保score转移到正确设备并扩展维度
            score_tensor = torch.tensor(score, device=device).expand_as(non_score_reward)
            reward = non_score_reward.clone()
            # 最后一位加上奖励模型给出的分数
            reward[-1] += score_tensor[-1]
            non_score_rewards.append(non_score_reward)
            rewards.append(reward)
        return rewards, non_score_rewards
    
    def loss(self, new_probs, old_values, new_values, rewards, old_probs):
        """计算损失，确保所有张量设备一致"""
        loss = torch.tensor(0.0, device=new_probs[0].device)  # 确保loss在正确设备上
        for new_prob, old_value, new_value, reward, old_prob in zip(new_probs, old_values, new_values, rewards, old_probs):
            # 确保所有输入张量设备一致
            reward = reward.to(new_prob.device)
            
            new_prob = new_prob.unsqueeze(0)
            old_value = old_value.unsqueeze(0)
            new_value = new_value.unsqueeze(0)
            reward = reward.unsqueeze(0)
            old_prob = old_prob.unsqueeze(0)
            
            # 后续计算保持不变...
            last_gae_lam = 0
            advantages_reversed = []
            gen_len = new_prob.shape[1]
            
            for t in reversed(range(gen_len)):
                next_values = old_value[:, t + 1] if t < gen_len - 1 else 0.0
                delta = reward[:, t] + self.config.gamma * next_values - old_value[:, t]
                last_gae_lam = delta + self.config.gamma * self.config.lam * last_gae_lam
                advantages_reversed.append(last_gae_lam)
            
            advantages = torch.stack(advantages_reversed[::-1]).transpose(0, 1)
            returns = advantages + old_value
            advantages = self.whiten(advantages)
            advantages = advantages.detach()
            
            value_clipped = torch.clamp(new_value,
                                        old_value - self.config.cliprange_value,
                                        old_value + self.config.cliprange_value)
            vf_loss1 = (new_value - returns) **2
            vf_loss2 = (value_clipped - returns)** 2
            vf_loss = torch.mean(torch.max(vf_loss2, vf_loss1))
            
            ratio = torch.exp(new_prob - old_prob)
            pg_losses = -advantages * ratio
            pg_losses2 = -advantages * torch.clamp(ratio,
                                                   1.0 - self.config.cliprange,
                                                   1.0 + self.config.cliprange)
            pg_loss = torch.mean(torch.max(pg_losses, pg_losses2))
            loss += pg_loss + self.config.vf_coef * vf_loss
        return loss

    @staticmethod
    def whiten(values, shift_mean=True):
        """
        归一化
        :param values: 要归一化的值
        :param shift_mean: 负一化方式
        :return: 返回归一化之后的结果
        """
        mean, var = torch.mean(values), torch.var(values)
        whitened = (values - mean) * torch.rsqrt(var + 1e-8)
        if not shift_mean:
            whitened += mean
        return whitened
