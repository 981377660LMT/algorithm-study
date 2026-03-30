# CH06: Retrieval-Augmented Generation / RAG檢索增強生成

## 1. 簡介与背景 (Introduction)

**Retrieval-Augmented Generation (RAG)** 是一種結合了**信息檢索 (Information Retrieval)** 和 **大型語言模型生成 (LLM Generation)** 的技術架構。

### 為什麼需要 RAG？

LLM（如 GPT-4）存在以下核心局限性：

- **知識截止 (Knowledge Cutoff)**：模型訓練數據是靜態的，無法回答最新的時事或動態數據。
- **幻覺 (Hallucination)**：當模型不知道答案時，傾向於一本正經地胡說八道。
- **私有數據缺乏 (Lack of Private Data)**：通用模型沒有企業內部的文檔或個人數據。

RAG 通過在生成答案前先去外部知識庫檢索相關信息，並將這些信息作為上下文提供給 LLM，從而解決上述問題。

---

## 2. 產品視角 (Product Perspective)

### 核心價值主張

- **準確性與可信度**：回答基於真實引用的文檔，而非模型臆測。
- **實時性**：無需重新訓練模型即可更新知識庫。
- **成本效益**：相比微調 (Fine-tuning)，構建索引的成本更低且更靈活。
- **數據隱私與權限**：可以實現基於用戶權限的檢索，確保數據安全。

### 典型應用場景 (Use Cases)

1.  **企業知識庫問答 (Enterprise Knowledge Q&A)**：員工查詢 HR 政策、技術文檔、銷售話術。
2.  **客戶支持機器人 (Customer Support Bot)**：基於最新的產品手冊和幫助中心文檔回答用戶提問。
3.  **法律/醫療助手**：引用具體的法條或醫學文獻，要求高度嚴謹的場景。
4.  **代碼助手**：基於特定代碼庫的上下文進行問答。

---

## 3. 技術架構與設計 (Technical Design)

RAG 系統通常包含三個主要階段：**Indexing (索引)**、**Retrieval (檢索)** 和 **Generation (生成)**。

### 3.1 數據準備與索引階段 (Indexing Pipeline)

這是離線或定期執行的過程，將原始數據轉化為可檢索的向量。

1.  **數據加載 (Load)**：從 PDF, Word, HTML, Database 等來源讀取數據。
2.  **文本切割 (Chunking)**：
    - _原因_：LLM 有上下文窗口限制 (Context Window)，且過長的文本會稀釋關鍵信息。
    - _策略_：固定字符數 (Fixed-size)、遞歸字符分割 (Recursive)、語義分割 (Semantic Chunking)。
    - _Key Point_：Chunk size 和 Overlap (重疊) 的設定對檢索質量影響很大。
3.  **向量化 (Embedding)**：使用 Embedding 模型 (如 OpenAI text-embedding-3, BGE) 將文本轉化為高維向量。
4.  **存儲 (Storage)**：將向量和元數據 (Metadata) 存入 **向量數據庫 (Vector Database)** (如 Pinecone, Milvus, Weaviate, ElasticSearch)。

### 3.2 檢索與生成階段 (Retrieval & Generation Pipeline)

這是主要與用戶交互的在線流程。

1.  **用戶查詢處理 (Query Processing)**：
    - 用戶輸入 Query。
    - (可選) **Query Rewriting**：將用戶模糊的問題重寫為更適合檢索的形式。
2.  **向量搜索 (Vector Search)**：
    - 將 Query 轉化為 Embedding 向量。
    - 在向量數據庫中進行 **ANN (Approximate Nearest Neighbor)** 搜索，找出最相似的 Top-K 個 Chunk。
3.  **重排序 (Reranking) (Advanced)**：
    - 向量搜索主要基於語義相似度，但不一定精準匹配問題意圖。
    - 使用 Cross-Encoder 模型對 Top-K 結果進行精細打分和重新排序，選出質量最高的 Top-N。
4.  **構建 Prompt (Prompt Engineering)**：
    - 將 System Prompt + 檢索到的上下文 (Context) + 用戶問題 (Question) 組合成最終的 Prompt。
    - _Prompt 模板示例_：

      ```text
      You are a helpful assistant. Use the following pieces of context to answer the question at the end.
      If you don't know the answer, just say that you don't know, don't try to make up an answer.

      Context:
      {retrieved_chunks}

      Question:
      {user_query}
      ```

5.  **LLM 生成 (Generation)**：將 Prompt 發送給 LLM，生成最終回答。

### 3.3 高級優化技術 (Advanced RAG Techniques)

為了提升召回率 (Recall) 和準確率 (Precision)，通常會採用以下技術：

- **Hybrid Search (混合檢索)**：結合 **關鍵詞搜索 (Keyword Search/BM25)** 和 **向量搜索 (Vector Search)**。關鍵詞對專有名詞匹配更好，向量對語義理解更好。
- **Metadata Filtering (元數據過濾)**：先根據時間、部門、文件類型過濾，再進行向量搜索。
- **Parent Document Retriever**：索引時切分小塊 (Small Chunk) 以保證語義精準，檢索時返回該小塊所屬的父文檔塊 (Large Chunk) 以提供更完整的上下文。
- **Multi-Query / Query Expansion**：將一個問題拆解或擴展成多個子問題去檢索，匯總結果。

## 4. API 與組件設計示例 (API Design)

### POST /v1/chat/completions

除了標準的 chat 接口，RAG 接口通常需要支持引用源的返回。

```json
// Response Payload Example
{
  "id": "chatcmpl-123",
  "object": "chat.completion",
  "created": 1677652288,
  "choices": [
    {
      "index": 0,
      "message": {
        "role": "assistant",
        "content": "根據公司政策，年假為 15 天..."
      },
      "finish_reason": "stop"
    }
  ],
  // 擴展字段：引用來源
  "citations": [
    {
      "document_id": "doc_882",
      "source": "employee_handbook_v2.pdf",
      "page": 12,
      "score": 0.89
    }
  ]
}
```

## 5. 面試重點總結 (Interview Takeaways)

1.  **Context Window Limit**：即使現在 LLM 支持 128k+ 上下文，RAG 依然重要，因為將所有知識塞入 Prompt 會導致 "Lost in the Middle" 現象，且成本過高、延遲過大。
2.  **Embedding Model 的選擇**：對中文場景，選擇支持中文優良的模型 (如 BGE-M3) 至關重要。
3.  **評估 (Evaluation)**：如何評估 RAG 系統？常用框架如 **RAGAS**，指標包括 Context Recall (召回率), Faithfulness (忠實度), Answer Relevance (相關性)。
