ac-predictor

ac-predictor 是 AtCoder 上的一个用户脚本，用于在比赛过程中预测用户的表现分（Performance）和评级(Rating)变化，提升比赛体验。代码模块化组织，主要分为以下几个部分：

# ac-predictor.js 代码详细解析

该文件是 AtCoder 上的一个用户脚本，用于在比赛过程中预测用户的 Perf（Performance）和评级变化，提升比赛体验。代码模块化组织，主要分为以下几个部分：

## 1. 多语言支持

```javascript
// 多语言支持代码
```

- **功能**：根据用户当前的语言环境（日本语或英语），自动加载对应的文本内容，提供多语言支持。
- **实现**：
  - 定义了日语 (jaJson) 和英语 (enJson) 两个对象，包含界面所需的文本。
  - 通过 getCurrentLanguage() 函数检测当前语言环境：
    - 检查导航栏中的语言选项，判断是日语还是英语。
  - getTranslation(label) 函数根据当前语言返回对应的翻译。
  - substitute(input) 函数用于替换模板中的占位符为对应的翻译文本。

## 2. 用户配置

```javascript
// 用户配置代码
```

- **功能**：提供一个设置界面，允许用户自定义脚本的行为。
- **实现**：
  - 使用 localStorage 存储用户配置，键为 ac-predictor-config。
  - defaultConfig 定义了默认配置，包括：
    - useResults：是否使用比赛结果。
    - hideDuringContest：是否在比赛期间隐藏预测。
    - hideUntilFixed：是否在 Perf 确定前隐藏预测。
    - useFinalResultOnVirtual：虚拟参赛时是否使用最终结果。
  - 提供了 getConfig、setConfig、getConfigObj、storeConfigObj 等函数来获取和设置配置。
  - ConfigView 类用于生成设置界面，添加配置选项。
  - ConfigController 类负责注册配置视图，将其添加到页面中。

## 3. 数据获取与缓存

```javascript
// 数据获取与缓存代码
```

- **功能**：从 AtCoder 网站或指定的 API 获取比赛相关的数据，并对频繁请求的数据进行缓存，减少网络请求次数。
- **实现**：
  - 定义了 Cache 类，用于缓存数据，具有过期时间：
    - has(key)：检查缓存中是否存在未过期的数据。
    - set(key, content)：缓存数据并设置过期时间。
    - get(key)：获取缓存的数据。
  - 提供了多个数据获取函数，例如：
    - getAPerfs(contestScreenName)：获取指定比赛的用户 Perf 数据。
    - getContestDetails()：获取比赛的详细信息。
    - `getHistory(userScreenName, contestType)`：获取用户的历史比赛数据。
    - getStandings(contestScreenName)：获取比赛的排名数据。
    - getResults(contestScreenName)：获取比赛的结果数据。
  - 使用了 addHandler 函数，通过监听 AJAX 请求，缓存响应的数据，防止重复请求。

## 4. Perf 预测

### 4.1 EloPerformanceProvider 类

```javascript
class EloPerformanceProvider {
  // 类定义
}
```

- **功能**：基于 Elo 算法，根据用户排名和其他用户的已知 Perf，预测用户的 Perf。
- **实现**：
  - **属性**：
    - ranks：用户的排名映射。
    - ratings：其他用户的已知 Perf 列表。
    - cap：Perf 上限，预测的 Perf 不会超过该值。
    - rankMemo：缓存 Perf 对应的排名，提升计算效率。
  - **方法**：
    - availableFor(userScreenName)：检查是否有指定用户的排名数据。
    - getPerformance(userScreenName)：获取指定用户的预测 Perf。
    - getPerformanceForRank(rank)：根据排名计算预测 Perf，使用二分查找反推 Perf 值。
    - getRankForPerformance(performance)：根据 Perf 计算对应的预测排名，使用 Elo 公式。

### 4.2 InterpolatePerformanceProvider 类

```javascript
class InterpolatePerformanceProvider {
  // 类定义
}
```

- **功能**：当部分用户的 Perf 不可用时，基于已知的 Perf 数据，对缺失的数据进行插值。
- **实现**：
  - **属性**：
    - ranks：用户的排名映射。
    - baseProvider：基础的 Perf 提供者，例如 EloPerformanceProvider。
  - **方法**：
    - getPerformance(userScreenName)：获取 �� 户的预测 Perf，对缺失的数据进行插值。
    - getPerformances()：获取所有用户的预测 Perf。

### 4.3 FixedPerformanceProvider 类

```javascript
class FixedPerformanceProvider {
  // 类定义
}
```

- **功能**：在比赛结果已知的情况下，直接提供每个用户的实际 Perf，不进行预测。

## 5. Rating 计算

### 5.1 算法竞赛评级计算

```javascript
// 相关函数
```

- **功能**：根据用户的历史 Perf 和新 Perf，计算未正数化的评级，并应用正数化函数得到最终评级。
- **实现**：
  - calcAlgRatingFromHistory(history)：根据历史 Perf 计算未正数化的评级。
  - calcAlgRatingFromLast(last, perf, ratedMatches)：基于上次的未正数化评级和新 Perf，增量计算新的未正数化评级。
  - positivizeRating(rating)：将未正数化的评级转换为正数化的评级。

### 5.2 启发式竞赛评级计算

```javascript
// 相关函数
```

- **功能**：使用特殊的算法，根据用户的历史 Perf 数据，计算启发式竞赛的评级。
- **实现**：
  - calcHeuristicRatingFromHistory(history)：根据历史 Perf 计算未正数化的评级。

## 5.3 IncrementalAlgRatingProvider 类

```javascript:atcoder/ac-predictor/acPredicator.js
class IncrementalAlgRatingProvider {
  unpositivizedRatingMap
  competitionsMap
  constructor(unpositivizedRatingMap, competitionsMap) {
    this.unpositivizedRatingMap = unpositivizedRatingMap
    this.competitionsMap = competitionsMap
  }
  availableFor(userScreenName) {
    return this.unpositivizedRatingMap.has(userScreenName)
  }
  async getRating(userScreenName, newPerformance) {
    if (!this.availableFor(userScreenName)) {
      throw new Error(`rating not available for ${userScreenName}`)
    }
    const rating = this.unpositivizedRatingMap.get(userScreenName)
    const competitions = this.competitionsMap.get(userScreenName)
    return Math.round(positivizeRating(calcAlgRatingFromLast(rating, newPerformance, competitions)))
  }
}
```

- **功能**：根据用户的未正数化评级和参加的比赛次数，增量计算新的评级。

## 5.4 ConstRatingProvider 类

```javascript:atcoder/ac-predictor/acPredicator.js
class ConstRatingProvider {
  ratings
  constructor(ratings) {
    this.ratings = ratings
  }
  availableFor(userScreenName) {
    return this.ratings.has(userScreenName)
  }
  async getRating(userScreenName, newPerformance) {
    if (!this.availableFor(userScreenName)) {
      throw new Error(`rating not available for ${userScreenName}`)
    }
    return this.ratings.get(userScreenName)
  }
}
```

- **功能**：直接提供每个用户的固定评级，不进行计算。

## 5.5 FromHistoryHeuristicRatingProvider 类

```javascript:atcoder/ac-predictor/acPredicator.js
class FromHistoryHeuristicRatingProvider {
  performancesProvider
  constructor(performancesProvider) {
    this.performancesProvider = performancesProvider
  }
  availableFor(userScreenName) {
    return true
  }
  async getRating(userScreenName, newPerformance) {
    const performances = await this.performancesProvider(userScreenName)
    performances.push(newPerformance)
    return Math.round(positivizeRating(calcHeuristicRatingFromHistory(performances)))
  }
}
```

- **功能**：启发式比赛的评级计算，基于用户的历史 Perf 数据。

## 6. 界面增强

### 6.1 StandingsTableView 类

```javascript
class StandingsTableView {
  // 类定义
}
```

- **功能**：修改比赛排名页面，添加 Perf 和评级变化列，显示预测结果。
- **实现**：
  - 监听排名表格的变化，实时更新预测结果。
  - 在表头添加新的列，显示 “perf” 和 “rating delta”。
  - 为每一行的用户添加预测 Perf 和评级变化信息。

### 6.2 EstimatorElement 类

```javascript
class EstimatorElement {
  // 类定义
}
```

- **功能**：提供一个工具，用户可以输入目标评级或 Perf，计算所需的 Perf 或可达到的评级。
- **实现**：
  - 提供输入框，用户输入目标值。
  - 实时计算并显示结果。
  - 支持在评级和 Perf 之间切换计算模式。

## 7. 页面控制器

### 7.1 StandingsPageController 类

```javascript
class StandingsPageController {
  // 类定义
}
```

- **功能**：控制普通排名页面的行为，加载数据，更新视图。
- **实现**：
  - 初始化时，获取比赛详情和排名数据。
  - 根据用户配置，决定是否在比赛期间或结果未确定时隐藏预测。
  - 使用 Perf 提供者和评级计算器，获取预测结果并更新视图。

### 7.2 VirtualStandingsPageController 类

```javascript
class VirtualStandingsPageController {
  // 类定义
}
```

- **功能**：控制虚拟参赛排 �� 页面的行为。
- **实现**：
  - 判断用户是否在进行虚拟参赛。
  - 根据配置，决定是否使用最终结果作为参考。
  - 合并原始参赛者和虚拟参赛者的排名，进行 Perf 预测。

### 7.3 ExtendedStandingsPageController 类

```javascript
class ExtendedStandingsPageController {
  // 类定义
}
```

- **功能**：控制扩展排名页面的行为，预测未参加比赛的用户 Perf（如解题讲解者）。
- **实现**：
  - 获取扩展排名数据和用户 Perf 数据。
  - 对于没有成绩的用户，进行 Perf 插值预测。
  - 更新视图，显示预测的 Perf。

## 8. 脚本入口

```javascript
// 脚本入口代码
```

- **功能**：根据当前页面，初始化对应的页面控制器。
- **实现**：
  - 检查当前页面的 URL，判断是普通排名页面、虚拟排名页面还是扩展排名页面。
  - 根据页面类型，创建并注册对应的控制器。
  - 注册用户配置界面。

---

以上就是 `ac-predictor.js` 文件的详细代码解析。该脚本通过模块化的设计，实现了在 AtCoder 上预测用户 Perf 和评级变化的功能，增强了用户的比赛体验。
