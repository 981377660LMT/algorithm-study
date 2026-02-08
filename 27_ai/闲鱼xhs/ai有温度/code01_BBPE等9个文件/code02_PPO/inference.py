from transformers import AutoModelForCausalLM, AutoTokenizer
from config import Config
from peft import LoraConfig, PeftModel
from config import LoraArguments


class LoraPPOModel(PeftModel):
    def __init__(self, config: Config):
        self.config = config
        model = AutoModelForCausalLM.from_pretrained(config.gpt_model).to(config.device)
        self.tokenizer = AutoTokenizer.from_pretrained(config.gpt_model)
        lora_args = LoraArguments()
        lora_config = LoraConfig(
            r=lora_args.lora_r,
            lora_alpha=lora_args.lora_alpha,
            target_modules=lora_args.lora_target_modules,
            lora_dropout=lora_args.lora_dropout,
            task_type="CAUSAL_LM",
        )
        super().__init__(model, lora_config)
        model = super().from_pretrained(model, config.save_lora_path)
        self.lora_ppo_model = model.merge_and_unload()
        self.raw_model = AutoModelForCausalLM.from_pretrained(config.gpt_model).to(config.device)
        print()

    def forward(self, query, system_content):
        messages = [
            {"role": "system", "content": system_content},
            {"role": "user", "content": query}
        ]
        text = self.tokenizer.apply_chat_template(
            messages,
            tokenize=False,
            add_generation_prompt=True
        )
        model_inputs = self.tokenizer([text], return_tensors="pt").to(self.config.device)
        lora_ppo_response = self.predict(model_inputs, self.lora_ppo_model, self.tokenizer)
        raw_response = self.predict(model_inputs, self.raw_model, self.tokenizer)
        return lora_ppo_response, raw_response

    @staticmethod
    def predict(model_inputs, model, tokenizer):
        generated_ids = model.generate(
            model_inputs.input_ids,
            max_new_tokens=512,
            num_beams=1,
            do_sample=False
        )
        generated_ids = [
            output_ids[len(input_ids):] for input_ids, output_ids in zip(model_inputs.input_ids, generated_ids)
        ]

        response = tokenizer.batch_decode(generated_ids, skip_special_tokens=True)[0]
        return response


if __name__ == '__main__':
    lora_ppo_model = LoraPPOModel(Config())
    lora_ppo_response, raw_response = lora_ppo_model("饭店服务员的态度太差，使用委婉积极的态度投诉", "你是一个有文化的文明人")
    print(f"lora_ppo_response:{lora_ppo_response}")
    print(f"raw_response:{raw_response}")
