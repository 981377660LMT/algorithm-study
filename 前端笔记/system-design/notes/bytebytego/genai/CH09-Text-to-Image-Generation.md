# CH09: Text-to-Image Generation / 文字轉圖片生成

## 1. 簡介与背景 (Introduction)

**Text-to-Image Generation** 是一項將自然語言描述（Prompt）轉換為視覺圖像的技術。這是多模態（Multimodal）AI 的核心應用之一。

### 技術演進路線

- **早期階段 (GANs)**：如 StackGAN，雖然能生成圖片，但僅限於簡單場景（如"一隻紅色的鳥"），解析度低且語義理解有限。
- **自回歸模型 (Auto-regressive)**：如 DALL-E 1、Parti。將圖像視為一串 token 序列，像預測下一個單詞一樣預測下一個像素塊。
- **擴散模型 (Diffusion Models)**：如 DALL-E 2/3, Stable Diffusion, Midjourney。這是當前的**統治級架構**，它通過學習如何從純噪聲中去除噪聲來還原圖像，具備極強的語義理解和細節生成能力。

---

## 2. 產品視角 (Product Perspective)

### 核心價值主張

- **創意平權**：降低視覺創作門檻，讓沒有繪畫技能的人也能通過文字創作高質量藝術。
- **效率革命**：將概念設計、廣告素材、插畫的製作時間從幾天縮短到幾秒。
- **按需生成**：替代部份傳統圖庫（Stock Photos），生成不存在的、定製化的場景。

### 典型應用場景 (Use Cases)

1.  **概念藝術與設計 (Concept Art)**：遊戲、電影的前期視覺開發，快速產出多種風格草圖。
2.  **營銷與廣告 (Marketing)**：生成社交媒體配圖、產品海報背景、個性化營銷郵件圖片。
3.  **UI/UX 設計**：快速生成佔位圖（Placeholder）或圖標資產。
4.  **電商應用**：虛擬模特試衣（結合 Inpainting 技術）。

### 產品設計挑戰

- **Prompt Engineering**：用戶通常不知道如何寫出好 Prompt。產品需要提供 Prompt 輔助、優化或模板功能。
- **可控性 (Controllability)**：用戶往往不僅想"生成一張圖"，還想控制構圖、姿勢、色調（需要 ControlNet 等技術輔助）。
- **安全性與版權**：防止生成 NSFW（不雅內容）、暴力內容或侵犯版權的圖像。

---

## 3. 技術架構與設計 (Technical Design)

目前業界標準架構是 **Latent Diffusion Models (LDM)**（例如 Stable Diffusion）。其核心思想是在低維的潛在空間（Latent Space）進行擴散過程，以平衡計算效率與圖像質量。

### 3.1 核心組件 (Key Components)

文字轉圖片模型通常由三個主要部分組成：

1.  **Text Encoder (文本編碼器)**：
    - _作用_：理解用戶的 Prompt。
    - _模型_：通常使用預訓練的 **CLIP (Contrastive Language-Image Pre-training)** 或 **T5**。
    - _原理_：CLIP 將文字映射到一個與圖片共享的向量空間。如果"一隻貓"的文字向量與一張真貓圖片的向量很接近，模型就能理解文字的視覺含義。

2.  **Image Information Creator (圖像生成核心 / U-Net)**：
    - _作用_：這是模型的"引擎"。它在潛在空間中執行去噪（Denoising）過程。
    - _架構_：通常是 **U-Net** 結構加上 **Cross-Attention (交叉注意力)** 機制。
    - _工作流_：
      - 輸入：帶噪聲的圖像潛在向量 (Noisy Latent)。
      - 條件：Text Encoder 輸出的文本向量 (Conditioning)。
      - Cross-Attention：U-Net 的每一層都會通過 Attention 機制"查看"文本向量，將文字信息注入到圖像特徵中（例如，當處理到圖像的上方區域時，Attention 機制注意到文本中的"藍天"，從而引導該區域去噪為藍色）。

3.  **Image Decoder (圖像解碼器 / VAE)**：
    - _作用_：將潛在空間的向量還原為像素圖片。
    - _模型_：**Variational Autoencoder (VAE)** 的 Decoder 部分。
    - _流程_：Latent Vector ($64 \times 64 \times 4$) $\rightarrow$ Pixel Image ($512 \times 512 \times 3$)。

### 3.2 生成與推理流程 (Inference Process)

當用戶輸入 "A cyberpunk city" 時：

1.  **Text Encoding**：將 Prompt 輸入 CLIP Text Encoder，得到文本 Embedding (如 $77 \times 768$ 維度的矩陣)。
2.  **Noise Generation**：在潛在空間隨機生成一個純噪聲張量 (Random Noise Latent)。
3.  **Denoising Loop (去噪循環)**：
    - 設定步數 (Steps)，例如 20-50 步。
    - 在每一步，U-Net 預測當前的噪聲分佈，但在預測時會參考文本 Embedding。
    - **Scheduler (調度器)**：根據預測的噪聲，從當前圖像中減去一部分噪聲，得到稍微清晰一點的圖像。常見調度器有 Euler a, DDIM, DPM++。
4.  **Decoding**：循環結束後，得到一個乾淨的 Latent，通過 VAE Decoder 轉換成最終的 PNG/JPG圖片。

### 3.3 關鍵技術：Classifier-Free Guidance (CFG)

這是一個在推理階段控制"模型多聽話"的參數。

- **原理**：模型同時進行兩個預測：
  1.  **有條件預測**：基於 Prompt 去噪。
  2.  **無條件預測**：忽略 Prompt，基於空字符串去噪。
- **公式**：$\text{Final Result} = \text{Unconditioned} + \text{Scale} \times (\text{Conditioned} - \text{Unconditioned})$。
- **CFG Scale**：
  - 值低 (1-3)：圖像更有創意，但可能與 Prompt 關聯不緊密。
  - 值高 (7-12)：圖像嚴格遵循 Prompt，對比度通常更高。
  - 值過高 (>15)：圖像容易出現偽影或顏色過飽和（Deep-fried effect）。

---

## 4. API 設計示例 (API Design)

### POST /v1/images/generations

```json
// Request Payload (Python style naming as per convention concepts)
{
  "prompt": "An astronaut riding a horse on mars, photorealistic, 4k",
  "negative_prompt": "cartoon, low resolution, blurry, deformed hands",
  "n": 1, // Number of images
  "size": "1024x1024",
  "response_format": "url",

  // Advanced Parameters
  "steps": 30,            // 推理步數，越高質量越好但越慢
  "guidance_scale": 7.5,  // CFG Scale, 7.5 是常用默認值
  "seed": 123456,         // 固定種子以復現結果
  "sampler": "DPM++ 2M Karras" // 採樣算法
}

// Response
{
  "created": 1677652288,
  "data": [
    {
      "url": "https://cdn.openai.com/..."
    }
  ]
}
```

---

## 5. 面試重點總結 (Interview Takeaways)

1.  **CLIP 的作用是什麼？**
    - CLIP 是連接文本和圖像的橋樑。它不是生成圖片的，而是用來衡量"這張圖和這段文字有多匹配"。在文本生成圖像模型中，它作為 Encoder 提供語義指導。

2.  **為什麼要用 Latent Diffusion 而不是 Pixel Diffusion？**
    - **計算效率**。在 $512 \times 512$ 的像素空間運算太慢且顯存消耗大。Latent Diffusion 先壓縮再擴散，將計算量降低了幾十倍，使得在消費級顯卡上運行成為可能。

3.  **什麼是 Cross-Attention？**
    - 這是將文本信息注入圖像生成過程的關鍵機制。在 U-Net 的每一層，圖像特徵作為 Query，文本 Embedding 作為 Key 和 Value。這讓模型知道"在哪裡畫什麼"。

4.  **Negative Prompt (負向提示詞) 的原理？**
    - 利用 CFG 機制，不僅僅是"靠近正向提示詞"，同時也要"遠離負向提示詞"。公式變為：$\text{Final} = \text{Negative} + \text{Scale} \times (\text{Positive} - \text{Negative})$。
