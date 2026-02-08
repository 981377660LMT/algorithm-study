import torch
from typing import List
from dataclasses import dataclass, field

class Config:
    # 换成自己的目录
    gpt_model = "./down_model/qwen0.5"
    data_path = "./data/train_data.json"
    save_lora_path = "./save_lora"
    device = "cuda:0" if torch.cuda.is_available() else "cpu"
    batch_size = 2
    epochs = 3
    lr = 0.001
    # DPO 参数 ############################
    dpo_epochs = 3
    beta = 0.1


@dataclass
class LoraArguments:
    lora_r: int = 2
    lora_alpha: int = 8
    lora_dropout: float = 0
    lora_target_modules: List[str] = field(
        default_factory=lambda: ['k_proj', 'v_proj']
    )
    # lora_target_modules = None
    lora_weight_path: str = ""
    q_lora: bool = False
    load_in_4bit: bool = False
    load_in_8bit: bool = False
    is_reload_trained_params = False  # 是否接着上次训练模型继续训练
