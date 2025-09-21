我们将以一个简单的绘图应用命令序列验证为例。假设我们有以下命令：

- `set_color`: 设置颜色。
- `begin_path`: 开始一个新路径。
- `move_to`: 移动画笔到一个点。
- `line_to`: 从当前点画一条线到新点。
- `close_path`: 闭合当前路径。

我们希望强制执行以下规则：一个有效的路径必须以 `begin_path` 开始，后面跟着至少一个 `move_to` 或 `line_to`，最后可以选择性地以 `close_path` 结束。

### 第 1 步：定义你的对象和类型

首先，定义你项目中要匹配的对象的类，并为每个类创建一个实现了 `MatchableType` 接口的“类型定义”对象。

```typescript
// 导入我们之前创建的引擎代码
import { ContentMatch, MatchableType } from './content_engine'

// --- 1a. 定义你的对象类 ---
class SetColorCmd {
  readonly type = SetColorType
  constructor(public color: string) {}
}
class BeginPathCmd {
  readonly type = BeginPathType
}
class MoveToCmd {
  readonly type = MoveToType
  constructor(public x: number, public y: number) {}
}
class LineToCmd {
  readonly type = LineToType
  constructor(public x: number, public y: number) {}
}
class ClosePathCmd {
  readonly type = ClosePathType
}

// --- 1b. 为每个类实现 MatchableType 接口 ---

// 注意：contentMatcher 对于叶子节点总是 ContentMatch.empty
// isGeneratable 为 true 的类型可以被 fillBefore/findWrapping 自动创建

const SetColorType: MatchableType<any> = {
  name: 'set_color',
  isGeneratable: false, // 不能自动生成，因为它需要颜色参数
  contentMatcher: ContentMatch.empty,
  isInGroup: g => g === 'command',
  createDefault: () => {
    throw new Error('Cannot create default set_color')
  },
  isLeaf: true
}

const BeginPathType: MatchableType<any> = {
  name: 'begin_path',
  isGeneratable: true, // 可以自动生成来修复序列
  contentMatcher: ContentMatch.empty,
  isInGroup: g => g === 'command',
  createDefault: () => new BeginPathCmd(),
  isLeaf: true
}

const MoveToType: MatchableType<any> = {
  name: 'move_to',
  isGeneratable: false, // 需要坐标
  contentMatcher: ContentMatch.empty,
  isInGroup: g => g === 'command' || g === 'path_segment',
  createDefault: () => {
    throw new Error('Cannot create default move_to')
  },
  isLeaf: true
}

const LineToType: MatchableType<any> = {
  name: 'line_to',
  isGeneratable: false, // 需要坐标
  contentMatcher: ContentMatch.empty,
  isInGroup: g => g === 'command' || g === 'path_segment',
  createDefault: () => {
    throw new Error('Cannot create default line_to')
  },
  isLeaf: true
}

const ClosePathType: MatchableType<any> = {
  name: 'close_path',
  isGeneratable: true,
  contentMatcher: ContentMatch.empty,
  isInGroup: g => g === 'command',
  createDefault: () => new ClosePathCmd(),
  isLeaf: true
}

// --- 1c. 创建一个类型映射表 ---
const commandTypes = {
  set_color: SetColorType,
  begin_path: BeginPathType,
  move_to: MoveToType,
  line_to: LineToType,
  close_path: ClosePathType
}
```

### 第 2 步：定义模式并编译

使用类似正则表达式的语法定义你的序列模式，然后调用 `ContentMatch.parse` 将其编译成一个状态机。

```typescript
// 模式：一个路径必须以 begin_path 开始，
// 后面跟着一个或多个 "path_segment" 组的成员（即 move_to 或 line_to），
// 最后可以选择性地以 close_path 结束。
const pathPattern = 'begin_path path_segment+ close_path?'

// 编译模式，得到状态机的入口点
const pathMatcher = ContentMatch.parse(pathPattern, commandTypes)
```

### 第 3 步：使用状态机

现在你可以使用编译好的 `pathMatcher` 来执行各种操作了。

#### 用法 A：验证一个完整的序列

你可以遍历一个命令序列，检查它是否完全匹配模式。

```typescript
function validateSequence(sequence: { type: MatchableType<any> }[]): boolean {
  let currentMatch: ContentMatch<any> | null = pathMatcher
  for (const command of sequence) {
    currentMatch = currentMatch?.matchType(command.type) ?? null
  }
  // 序列有效，当且仅当所有命令都被匹配，且最终状态是一个合法的终点。
  return !!currentMatch && currentMatch.validEnd
}

// 示例
const validSequence = [new BeginPathCmd(), new MoveToCmd(10, 10), new LineToCmd(20, 20)]
const invalidSequence = [new MoveToCmd(10, 10), new LineToCmd(20, 20)] // 缺少 begin_path

console.log('序列1是否有效:', validateSequence(validSequence)) // 输出: true
console.log('序列2是否有效:', validateSequence(invalidSequence)) // 输出: false
```

#### 用法 B：自动填充缺失的前缀 (`fillBefore`)

如果你有一个不完整的序列（例如，只有 `line_to`），你可以用 `fillBefore` 来计算需要添加什么前缀才能使其合法。

```typescript
// 这是一个不完整的片段，我们想在它前面自动填充内容
const partialSequence = {
  children: [new LineToCmd(50, 50)],
  childCount: 1,
  child: (i: number) => partialSequence.children[i]
}

// 尝试在 pathMatcher 的初始状态填充，以使 partialSequence 合法
const result = pathMatcher.fillBefore(partialSequence, true) // toEnd=true 确保填充后整个序列是完整的

if (result) {
  console.log('自动填充结果:')
  // result.fragment 是一个数组，包含需要插入的对象实例
  console.log(result.fragment.map(cmd => cmd.constructor.name).join(', ')) // 输出: BeginPathCmd
  // result.match 是填充并匹配 partialSequence 后的新状态
}
```

#### 用法 C：寻找合适的包装 (`findWrapping`)

假设我们有一个 `line_to` 命令，但我们所处的上下文（例如，一个只接受完整路径的 `layer` 节点）不允许直接插入它。我们可以用 `findWrapping` 来找到能包裹它的节点。

_这个例子稍微复杂，因为它需要一个“容器”类型。假设我们有一个 `path` 类型，它的 `content` 就是我们上面定义的 `pathPattern`。_

```typescript
const PathType: MatchableType<any> = {
  name: 'path',
  isGeneratable: true,
  // 它的内容由 pathMatcher 定义
  contentMatcher: pathMatcher,
  isInGroup: g => g === 'drawable',
  createDefault: () => ({ type: PathType, content: [] }),
  isLeaf: false // 它是容器，不是叶子
}

// 假设我们有一个 layer，它的模式是 "path+" (一个或多个路径)
const layerMatcher = ContentMatch.parse('path+', { path: PathType })

// 在 layer 的开头，我们想插入一个 line_to，这显然不合法
// 我们需要找到什么东西来包裹它
const wrapping = layerMatcher.findWrapping(LineToType)

if (wrapping) {
  // wrapping 是一个类型数组
  console.log('找到的包装类型:', wrapping.map(t => t.name).join(' -> ')) // 输出: path
}
// 这意味着：为了在 layer 中放入一个 line_to，你需要先创建一个 path 节点来包裹它。
```

通过这几个例子，你可以看到这个引擎的强大之处：它不仅能做简单的验证，还能智能地推断如何修复和构建合法的对象序列。
