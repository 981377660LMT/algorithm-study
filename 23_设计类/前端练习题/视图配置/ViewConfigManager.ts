/* eslint-disable arrow-body-style */

// ==========================================
// 1. 类型定义 (Type Definitions)
// ==========================================

// 递归获取嵌套路径类型 (例如: "user.name" | "order.id")
export type Path<T> = T extends object
  ? {
      [K in keyof T]: K extends string
        ? T[K] extends object
          ? K | `${K}.${Path<T[K]>}`
          : K
        : never
    }[keyof T]
  : never

// 获取路径对应的值类型
export type PathValue<T, P extends Path<T>> = P extends `${infer K}.${infer Rest}`
  ? K extends keyof T
    ? Rest extends Path<T[K]>
      ? PathValue<T[K], Rest>
      : never
    : never
  : P extends keyof T
  ? T[P]
  : never

export type FilterOperator = 'eq' | 'neq' | 'contains' | 'gt' | 'lt' | 'gte' | 'lte' | 'in'

export interface FilterDescriptor<T> {
  id: string
  field: Path<T>
  operator: FilterOperator
  value: any
  enabled: boolean
}

export interface SortDescriptor<T> {
  id: string
  field: Path<T>
  direction: 'asc' | 'desc'
}

export interface GroupDescriptor<T> {
  id: string
  field: Path<T>
  direction: 'asc' | 'desc' // 分组本身的排序
}

// 完整的视图配置状态
export interface ViewConfig<T> {
  filters: FilterDescriptor<T>[]
  sorts: SortDescriptor<T>[]
  groups: GroupDescriptor<T>[]
}

// ==========================================
// 2. 视图输出结构 (View Output Structure)
// ==========================================

export enum ViewNodeType {
  Group = 'GROUP',
  Row = 'ROW'
}

export interface ViewRowNode<T> {
  type: ViewNodeType.Row
  id: string // 唯一标识
  data: T
  originalIndex: number // 在原始数组中的索引
}

export interface ViewGroupNode<T> {
  type: ViewNodeType.Group
  id: string // 分组唯一标识 (例如: "group:country:USA")
  field: Path<T>
  value: any
  depth: number
  children: ViewNode<T>[] // 可能是子分组，也可能是数据行
  count: number // 子项数量
  isExpanded?: boolean // UI 状态辅助
}

export type ViewNode<T> = ViewRowNode<T> | ViewGroupNode<T>

// ==========================================
// 3. 工具函数 (Utils)
// ==========================================

function getNestedValue<T>(obj: T, path: string): any {
  return path.split('.').reduce((acc: any, part) => acc && acc[part], obj)
}

function compareValues(a: any, b: any, direction: 'asc' | 'desc'): number {
  if (a === b) return 0
  if (a === null || a === undefined) return 1
  if (b === null || b === undefined) return -1

  const res = a > b ? 1 : -1
  return direction === 'asc' ? res : -res
}

// ==========================================
// 4. 视图管理器类 (The Manager)
// ==========================================

export class ViewConfigManager<T> {
  private _data: T[] = []
  private _config: ViewConfig<T>

  constructor(data: T[], config?: Partial<ViewConfig<T>>) {
    this._data = data
    this._config = {
      filters: [],
      sorts: [],
      groups: [],
      ...config
    }
  }

  // --- Configuration Setters (Fluent API) ---

  public setFilters(filters: FilterDescriptor<T>[]) {
    this._config.filters = filters
    return this
  }

  public setSorts(sorts: SortDescriptor<T>[]) {
    this._config.sorts = sorts
    return this
  }

  public setGroups(groups: GroupDescriptor<T>[]) {
    this._config.groups = groups
    return this
  }

  public updateData(data: T[]) {
    this._data = data
    return this
  }

  public getConfig(): ViewConfig<T> {
    return { ...this._config }
  }

  // --- The Core Processing Pipeline ---

  public process(): ViewNode<T>[] {
    // 1. 筛选 (Filtering)
    let processed = this.applyFilters(this._data)

    // 2. 排序 (Sorting) - 对扁平数据先排序
    // 注意：如果有分组，通常先按分组字段排序，再按组内字段排序
    processed = this.applySorting(processed)

    // 3. 分组 (Grouping) - 转换为树形结构
    if (this._config.groups.length > 0) {
      return this.applyGrouping(processed, 0)
    }

    // 如果没有分组，直接包装成 RowNode
    return processed.map((item, index) => ({
      type: ViewNodeType.Row,
      id: `row_${index}`, // 实际项目中建议使用 item 中的唯一 ID
      data: item,
      originalIndex: index
    }))
  }

  // --- Internal Implementation ---

  private applyFilters(data: T[]): T[] {
    const activeFilters = this._config.filters.filter(f => f.enabled)
    if (activeFilters.length === 0) return data

    return data.filter(item => {
      // 默认逻辑：所有 Filter 之间是 AND 关系
      return activeFilters.every(filter => {
        const val = getNestedValue(item, filter.field as string)
        const target = filter.value

        switch (filter.operator) {
          case 'eq':
            return val == target
          case 'neq':
            return val != target
          case 'contains':
            return String(val).toLowerCase().includes(String(target).toLowerCase())
          case 'gt':
            return val > target
          case 'lt':
            return val < target
          case 'gte':
            return val >= target
          case 'lte':
            return val <= target
          case 'in':
            return Array.isArray(target) && target.includes(val)
          default:
            return true
        }
      })
    })
  }

  private applySorting(data: T[]): T[] {
    // 排序优先级：先按分组字段排，再按用户配置的排序排
    const effectiveSorts: { field: string; direction: 'asc' | 'desc' }[] = [
      ...this._config.groups.map(g => ({ field: g.field as string, direction: g.direction })),
      ...this._config.sorts.map(s => ({ field: s.field as string, direction: s.direction }))
    ]

    if (effectiveSorts.length === 0) return [...data]

    return [...data].sort((a, b) => {
      for (const sort of effectiveSorts) {
        const valA = getNestedValue(a, sort.field)
        const valB = getNestedValue(b, sort.field)
        const result = compareValues(valA, valB, sort.direction)
        if (result !== 0) return result
      }
      return 0
    })
  }

  private applyGrouping(data: T[], groupIndex: number): ViewNode<T>[] {
    const groupDesc = this._config.groups[groupIndex]
    if (!groupDesc) {
      // 递归终止：没有更多分组层级，返回数据行
      return data.map((item, idx) => ({
        type: ViewNodeType.Row,
        id: `row_${idx}`, // 实际应取 item.id
        data: item,
        originalIndex: idx
      }))
    }

    // 1. 按当前分组字段将数据分桶
    const groups = new Map<any, T[]>()
    const field = groupDesc.field as string

    for (const item of data) {
      const key = getNestedValue(item, field)
      // 处理 key 为对象的情况 (需要序列化或引用检查，这里简化为 String)
      const mapKey = key === null || key === undefined ? 'null' : key

      if (!groups.has(mapKey)) {
        groups.set(mapKey, [])
      }
      groups.get(mapKey)!.push(item)
    }

    // 2. 将 Map 转换为 GroupNode 列表
    const result: ViewGroupNode<T>[] = []

    groups.forEach((groupData, key) => {
      // 递归处理下一级
      const children = this.applyGrouping(groupData, groupIndex + 1)

      result.push({
        type: ViewNodeType.Group,
        id: `group_${groupDesc.field}_${key}`,
        field: groupDesc.field,
        value: key,
        depth: groupIndex,
        count: groupData.length,
        children,
        isExpanded: true // 默认展开，实际可由外部控制
      })
    })

    // 3. 对分组本身进行排序 (根据 groupDesc.direction)
    // 注意：前面的 applySorting 已经保证了数据内部顺序，这里是保证 GroupNode 的顺序
    result.sort((a, b) => compareValues(a.value, b.value, groupDesc.direction))

    return result
  }
}

{
  interface User {
    id: number
    profile: {
      name: string
      age: number
    }
    role: 'admin' | 'user' | 'guest'
    active: boolean
  }

  const users: User[] = [
    { id: 1, profile: { name: 'Alice', age: 30 }, role: 'admin', active: true },
    { id: 2, profile: { name: 'Bob', age: 25 }, role: 'user', active: true },
    { id: 3, profile: { name: 'Charlie', age: 35 }, role: 'user', active: false },
    { id: 4, profile: { name: 'David', age: 30 }, role: 'admin', active: false }
  ]

  const manager = new ViewConfigManager(users)

  // IDE 会自动补全 'profile.age', 'role' 等
  manager.setFilters([
    {
      id: 'f1',
      field: 'profile.age', // 类型安全：只能输入 User 的合法路径
      operator: 'gte',
      value: 25,
      enabled: true
    }
  ])
  // 先按 Role 分组，组内按 Name 排序
  manager.setGroups([{ id: 'g1', field: 'role', direction: 'asc' }])
  manager.setSorts([{ id: 's1', field: 'profile.name', direction: 'asc' }])

  const viewNodes = manager.process()
  console.log(JSON.stringify(viewNodes, null, 2))

  // 输出结构大致如下:
  // [
  //   {
  //     type: 'GROUP',
  //     value: 'admin',
  //     children: [
  //       { type: 'ROW', data: { name: 'Alice'... } },
  //       { type: 'ROW', data: { name: 'David'... } }
  //     ]
  //   },
  //   {
  //     type: 'GROUP',
  //     value: 'user',
  //     children: [...]
  //   }
  // ]
}
