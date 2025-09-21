import { ContentMatch, MatchableType } from './parse'

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

// 模式：一个路径必须以 begin_path 开始，
// 后面跟着一个或多个 "path_segment" 组的成员（即 move_to 或 line_to），
// 最后可以选择性地以 close_path 结束。
const pathPattern = 'begin_path path_segment+ close_path?'

// 编译模式，得到状态机的入口点
const pathMatcher = ContentMatch.parse(pathPattern, commandTypes)

function validateSequence(sequence: { type: MatchableType<any> }[]): boolean {
  let currentMatch: ContentMatch<any> | null = pathMatcher
  for (const command of sequence) {
    currentMatch = currentMatch?.matchType(command.type) ?? null
  }
  // 序列有效，当且仅当所有命令都被匹配，且最终状态是一个合法的终点。
  return !!currentMatch && currentMatch.validEnd
}

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
