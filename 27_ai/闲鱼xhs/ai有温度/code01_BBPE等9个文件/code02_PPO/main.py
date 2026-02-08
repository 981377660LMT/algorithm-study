import torch
from config import Config
from torch.utils.data import DataLoader
from torch.optim import Adam
from actor_critic_model import ActorCriticLoraModel
from reward_model import RewardModel
from data_load import CustomDataset
from reference_model import ReferenceModel
from ppo import PPO
from tools import Tools


class TrainPpo:
    def __init__(self):
        self.config = Config()
        # 演员和评论家模型
        self.actor_critic_model = ActorCriticLoraModel(self.config).to(self.config.device)
        self.tokenizer = self.actor_critic_model.tokenizer
        # 获得演员和评论家模型优化器, 这里使用的是lora, 不优化全量数据
        self.actor_critic_opt = Adam(self.actor_critic_model.parameters(), lr=self.config.lr)
        # 参考模型
        self.reference_model = ReferenceModel(self.config).to(self.config.device)
        # 奖励模型
        self.reward_model = RewardModel(self.config).to(self.config.device)
        # 训练数据
        dataset = CustomDataset(self.config.data_path, self.tokenizer)
        
        self.data_loader = DataLoader(dataset, batch_size=self.config.batch_size, shuffle=True,
                                      collate_fn=dataset.collate_fn)
        self.ppo = PPO(self.actor_critic_model, self.config, self.actor_critic_opt)

    def train_ppo(self):
        self.save_model()
        for epoch in range(self.config.epochs):
            for batch_data in self.data_loader:
                # 修复：如果batch_data是元组，直接使用索引访问，而不是当作字典处理
                # 假设batch_data的第一个元素是输入数据
                batch_tensor = batch_data[0].to(self.config.device)
                
                # 获得演员模型生成的结果
                prompt_generate, prompt_generate_ids, generate_ids = self.actor_critic_model.actor_generate(
                    batch_tensor)  # 使用处理好的张量
                
                attention_mask = (prompt_generate_ids != self.tokenizer.pad_token_id)
                generate_ids_mask = (generate_ids[:, :-1] != self.tokenizer.pad_token_id)
                response_shape = generate_ids.shape[1] - 1
                tools = Tools(response_shape, generate_ids_mask)
                
                pure_generate = [one.split("assistant\n")[1] for one in prompt_generate]
                reward = self.reward_model(pure_generate)
                prob_refs = self.reference_model(prompt_generate_ids, attention_mask, tools)
                self.ppo.train(prompt_generate_ids, attention_mask, prob_refs, reward, tools)
        
        self.save_model()
    

    def save_model(self):
        # 保存lora参数
        self.actor_critic_model.model.save_pretrained(self.config.save_lora_path, safe_serialization=False)
        # 保存价值模型参数
        torch.save(self.actor_critic_model.model.v_head.state_dict(), self.config.save_v_head_path)


if __name__ == '__main__':
    train_ppo = TrainPpo()
    train_ppo.train_ppo()