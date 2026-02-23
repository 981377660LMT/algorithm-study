# V8 v8.1 发布：Intl.DisplayNames API 简介

## V8 release v8.1

- **Original Link**: [https://v8.dev/blog/v8-release-81](https://v8.dev/blog/v8-release-81)
- **Publication Date**: 2020-02-25
- **Summary**: V8 v8.1 引入了 `Intl.DisplayNames` API，支持获取本地化的语言、区域和货币名称，标志着国际化能力的进一步增强。

---

### 1. 新特性：Intl.DisplayNames

这是国际化（i18n）领域的一个重大补充。

- **痛点**：以前如果你想在网页上显示“法语”在不同语言下的名称（如 English 下显示 "French"，中文下显示“法文”），你必须自己维护庞大的翻译 JSON 文件。
- **解决方案**：`Intl.DisplayNames` 直接利用操作系统和浏览器内置的 ICU 数据。

```javascript
const zhLanguageNames = new Intl.DisplayNames(['zh-Hans'], { type: 'language' })
zhLanguageNames.of('fr') // "法文"

const enRegionNames = new Intl.DisplayNames(['en'], { type: 'region' })
enRegionNames.of('US') // "United States"
```

### 2. 性能与引擎内部

- **更轻量的数据加载**：V8 优化了国际化数据在内存中的表示，减少了因为加载多种区域设置（Locales）而导致的内存剧增。
- **JavaScript 语法增强**：v8.1 继续打磨 8.0 引入的 `?.` 和 `??` 语法，确保这些操作符在解释器和编译器中都能达到最优路径。

### 3. V8 API 变动

- 发布流程：此版本开始，V8 的发布节奏与 Chrome 浏览器更加紧密同步，确保了每次 Chrome 更新都能搭载最稳定的 V8 引擎。
- 弃用提醒：开始清理一些旧的、非规范的微任务接口，向更符合 TC39 规范的 `queueMicrotask` 靠拢。

### 4. 总结

v8.1 虽然是一个较小的版本，但 `Intl.DisplayNames` 的加入极大地降低了全球化应用的资源体积，是 Web 国际化标准化进程中的重要一步。
