import torch
from transformers import AutoModelForCausalLM, AutoTokenizer
from config import Config
from peft import LoraConfig, PeftModel
from config import LoraArguments


class LoraModel(PeftModel):
    def __init__(self, config: Config, model):
        lora_args = LoraArguments()
        lora_config = LoraConfig(
            r=lora_args.lora_r,
            lora_alpha=lora_args.lora_alpha,
            target_modules=lora_args.lora_target_modules,
            lora_dropout=lora_args.lora_dropout,
            task_type="CAUSAL_LM",
        )
        super().__init__(model, lora_config)
        if lora_args.is_reload_trained_params:
            super().from_pretrained(model, config.save_lora_path)
        for name, module in self.named_modules():
            if 'lora_' in name:
                for param in module.parameters():
                    param.requires_grad = True

    def forward(self, input_ids, attention_mask):
        res = super().forward(input_ids, attention_mask, output_hidden_states=True)
        return res.logits


class Model(torch.nn.Module):
    def __init__(self, config: Config):
        super().__init__()
        model = AutoModelForCausalLM.from_pretrained(config.gpt_model).to(config.device).eval()
        self.model = LoraModel(config, model)
        self.tokenizer = AutoTokenizer.from_pretrained(config.gpt_model)

    def forward(self, input_ids, attention_mask):
        logits = self.model(input_ids, attention_mask)
        return logits

