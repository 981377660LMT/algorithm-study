from transformers import AutoModelForCausalLM, AutoTokenizer
from config import Config
from peft import LoraConfig, PeftModel
from config import LoraArguments


class LoraDPOModel(PeftModel):
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
        self.lora_dpo_model = model.merge_and_unload()
        self.raw_model = AutoModelForCausalLM.from_pretrained(config.gpt_model).to(config.device)
        print()

    def forward(self, query, instruction):
        messages = [
            {"role": "system", "content": "你是一个非常有帮助和智能的助手。"},
            {"role": "instrution", "content": instruction},
            {"role": "user", "content": query}
        ]
        text = self.tokenizer.apply_chat_template(
            messages,
            tokenize=False,
            add_generation_prompt=True
        )
        model_inputs = self.tokenizer([text], return_tensors="pt").to(self.config.device)
        lora_dpo_response = self.predict(model_inputs, self.lora_dpo_model, self.tokenizer)
        raw_response = self.predict(model_inputs, self.raw_model, self.tokenizer)
        return lora_dpo_response, raw_response

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
    lora_dpo_model = LoraDPOModel(Config())
    lora_dpo_response, raw_response = lora_dpo_model("操作说明：为一个关于回收的网站设计一个标志。", "以下是描述任务的说明。编写适当地完成请求的响应。")
    print(f"lord_dpo_response:{lora_dpo_response}")
    print(f"raw_response:{raw_response}")
