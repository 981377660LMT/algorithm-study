# CH10: Personalized Headshot Generation / 個人化頭像生成

## 1. 簡介与背景 (Introduction)

**Personalized Headshot Generation** 指的是通過微調（Fine-tuning）生成模型，使其學習特定用戶的面部特徵，從而生成該用戶在不同風格、場景或服裝下的高質量照片。

這與通用的 Text-to-Image（CH09）的最大區別在於：**一致性（Identity Consistency）**。通用模型可以生成"一個漂亮的女人"，但無法生成"你"。

### 技術演進與核心

該領域的爆發主要歸功於以下技術的普及：

- **Dreambooth (Google)**：讓模型將一個特定的標識符（如 `sks person`）綁定到幾張輸入圖片的主體上。
- **LoRA (Low-Rank Adaptation)**：一種高效的參數微調技術，使得訓練一個個人模型只需幾分鐘，且模型文件極小（幾十 MB 而非幾 GB）。

---

## 2. 產品視角 (Product Perspective)

### 核心價值主張 (Value Proposition)

- **替代專業攝影**：只需幾美元，即可獲得媲美影樓几百美元的 LinkedIn 頭像、商務照或寫真集。
- **風格多樣化**：用戶可以輕易獲得自己穿著宇航服、或是 90 年代復古風格的照片，這在現實中難以拍攝。
- **隱私與便捷**：無需出門，無需面對攝影師的尷尬，在家上傳幾張自拍即可。

### 典型應用場景 (Use Cases)

1.  **職業頭像 (Professional Headshots)**：用於 LinkedIn、簡歷、公司官網的商務正裝照。
2.  **社交媒體頭像 (Social Profiles)**：風格化的卡通、油畫或賽博朋克風格頭像。
3.  **約會應用照片 (Dating Apps)**：生成高質量的生活照，提升匹配率。
4.  **虛擬試穿/髮型預覽**：看自己留長髮或染髮的樣子。

### 二、用戶交互流程 (User Flow)

1.  **上傳**：用戶上傳 10-20 張不同角度、光線清晰的自拍照。
2.  **訓練 (Training)**：系統在後台訓練專屬的 LoRA 模型（耗時 10-30 分鐘）。
3.  **生成 (Inference)**：用戶選擇模板（Prompt），系統生成圖片。
4.  **挑選與下載**：由於 AI 生成存在隨機性，通常會生成大量圖片（如 100 張），讓用戶挑選其中最像的幾張。

---

## 3. 技術架構與設計 (Technical Design)

系統設計的核心在於如何構建一個**自動化、低成本、高併發**的訓練與推理流水線。

### 3.1 訓練流水線 (Fine-tuning Pipeline)

這是區別於普通生成應用的獨有環節。

1.  **預處理 (Preprocessing)**：
    - **人臉檢測與裁剪 (Face Detection & Cropping)**：自動識別照片中的人臉，並將其裁剪並對齊到 $512 \times 512$ 或 $768 \times 768$。
    - **圖像標註 (Captioning)**：使用 BLIP 等模型為訓練圖生成描述文字（如 "a photo of a sks person wearing a t-shirt"）。
2.  **Dreambooth / LoRA 訓練**：
    - **Rare Token 綁定**：在 Prompt 中插入一個生僻詞（如 `sks` 或 `ohwx`）代表用戶。訓練目標是讓模型學會：當 Prompt 出現 `sks person` 時，就要畫出這個特定的用戶。
    - **正則化數據 (Prior Preservation Loss)**：為了防止模型"學傻了"（例如學會了你的臉，但忘記了什麼是"人"），訓練時會混入通用的生成的"人"的圖片作為正則化數據，保持模型的通用能力。
    - **LoRA 注入**：不更新整個 SD 模型的幾十億參數，只更新注意力層中插入的低秩矩陣（Rank $r=4$ to $128$）。

### 3.2 推理流水線 (Inference Pipeline)

1.  **加載適配器**：基礎模型（如 Stable Diffusion XL）常駐顯存，根據用戶請求動態加載對應的 LoRA 權重文件。
2.  **Prompt Engineering**：
    - 將用戶選擇的風格（如 "Business suit, office background"）與用戶標識符（"photo of sks person"）結合。
    - 大量使用 **Negative Prompts** 來抑制變形、模糊等問題。
3.  **ControlNet (可選)**：為了保證生成的姿勢自然或符合模板，可能會使用 OpenPose ControlNet 來固定姿勢。
4.  **人臉修復 (Face Restoration)**：
    - **問題**：Stable Diffusion 生成全身照時，臉部像素佔比小，容易崩壞（眼睛不對稱、嘴巴奇怪）。
    - **解決**：使用 **CodeFormer** 或 **GFPGAN** 對生成圖的人臉區域進行裁剪、超分重建、再貼回原圖。這是生成高質量頭像的**關鍵一步**。

---

## 4. API 設計示例 (API Design)

由於訓練是耗時任務，API 必須是異步的。

### POST /v1/models/train (創建個人模型)

```json
// Request
{
  "user_id": "user_123",
  "instance_name": "sks_man",
  "images": [
    "base64_img1...",
    "base64_img2..."
  ],
  "class_type": "man" // 用於生成正則化數據
}

// Response
{
  "training_job_id": "job_train_999",
  "status": "queued",
  "estimated_time": 1200 // 秒
}
```

### POST /v1/images/generate_avatar (生成頭像)

```json
// Request
{
  "model_id": "lora_weights_user_123_v1",
  "style_template": "professional_linkedin", // 預設的 Prompt 模板 ID
  "num_images": 8
}

// Response
{
  "generation_job_id": "job_gen_888",
  "status": "processing"
}
```

---

## 5. 面試重點總結 (Interview Takeaways)

1.  **Catastrophic Forgetting (災難性遺忘)**：
    - 微調模型時，模型容易學會新知識（你的臉）但忘記舊知識（什麼是人）。
    - **解決方案**：Dreambooth 的 Prior Preservation Loss，即在訓練數據中加入機器生成的通用人臉圖來"複習"。

2.  **為什麼選擇 LoRA 而不是全量微調 (Full Fine-tuning)？**
    - **存儲成本**：全量微調一個 SDXL 模型需要 6GB+ 存儲空間，而 LoRA 只需要 10MB-200MB。對於百萬用戶級別的產品，存儲成本差異巨大。
    - **切換速度**：加載 LoRA 權重非常快，可以實現多租戶共享同一個 Base Model，動態插拔 LoRA，提高 GPU 利用率。

3.  **恐怖谷效應與人臉相似度**：
    - 如果訓練過度（Overfitting），生成的圖片會像簡單的 PS 貼圖，光影不自然。
    - 如果訓練不足（Underfitting），生成的圖片不像用戶。
    - 這是產品調節的難點：需要在 `Training Steps` 和 `LoRA Rank` 之間找到平衡。
