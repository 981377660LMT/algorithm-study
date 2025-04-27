## useRequest

### Example

```tsx
import { useRequest } from 'ahooks'
import Mock from 'mockjs'
import React from 'react'

function getUsername(): Promise<string> {
  return new Promise(resolve => {
    setTimeout(() => {
      resolve(Mock.mock('@name'))
    }, 1000)
  })
}

export default () => {
  const { data, error, loading } = useRequest(getUsername)

  if (error) {
    return <div>failed to load</div>
  }
  if (loading) {
    return <div>loading...</div>
  }
  return <div>Username: {data}</div>
}
```

### 解析

```tsx
// TParams：请求参数
// TData：请求返回数据
// service：请求函数
// options：请求配置
const {
  loading: boolean,
  data?: TData,
  error?: Error,
  params: TParams || [],
  run: (...params: TParams) => void,
  runAsync: (...params: TParams) => Promise<TData>,
  refresh: () => void,
  refreshAsync: () => Promise<TData>,
  mutate: (data?: TData | ((oldData?: TData) => (TData | undefined))) => void,
  cancel: () => void,
} = useRequest<TData, TParams>(
  service: (...args: TParams) => Promise<TData>,
  {
    manual?: boolean,
    defaultParams?: TParams,
    onBefore?: (params: TParams) => void,
    onSuccess?: (data: TData, params: TParams) => void,
    onError?: (e: Error, params: TParams) => void,
    onFinally?: (params: TParams, data?: TData, e?: Error) => void,
  }
);
```

- 自动请求/手动请求
  manual + run/runAsync

  - run 是一个普通的同步函数，我们会自动捕获异常，你可以通过 options.onError 来处理异常时的行为。
  - runAsync 是一个返回 Promise 的异步函数，如果使用 runAsync 来调用，则意味着你需要自己捕获异常。

- 生命周期

  - onBefore
  - onSuccess
  - onError
  - onFinally

- 刷新（重复上一次请求）
  refresh/refreshAsync

- 立即变更数据
  mutate 乐观更新视图，失败后自动回滚

- 取消
  忽略当前 promise 返回的数据和错误；
  组件卸载时/竞态取消。

- 参数管理
  useRequest 返回的 params 会记录当次调用 service 的参数数组。比如你触发了 run(1, 2, 3)，则 params 等于 [1, 2, 3] 。

- 轮询
  https://ahooks.js.org/zh-CN/hooks/use-request/polling
- ready
  控制请求是否发出

- 依赖刷新
  通过设置 options.refreshDeps，在依赖变化时， useRequest 会自动调用 refresh 方法，实现刷新（重复上一次请求）的效果
  ```tsx
  const [userId, setUserId] = useState('1')
  const { data, run } = useRequest(() => getUserSchool(userId), {
    refreshDeps: [userId]
  })
  ```
- 防抖
- 节流
- 屏幕聚焦重新请求
  refreshOnWindowFocus
  在浏览器窗口 `refocus 和 revisible` 时，会重新发起请求
  如果和上一次请求间隔大于 focusTimespan ms，则会重新请求一次
- 错误重试
  - retryCount
    错误重试次数。如果设置为 -1，则无限次重试。
  - retryInterval
    如果不设置，默认采用简易的指数退避算法，取 `1000 * 2 ** retryCount`，也就是`第一次重试等待 2s`，第二次重试等待 4s，以此类推，如果大于 30s，`则取 30s`
- loading delay
  延迟 loading 变成 true 的时间，有效防止闪烁
- 缓存 & SWR(stale-while-revalidate)
  如果设置了 options.cacheKey，useRequest 会将当前请求成功的数据缓存起来。下次组件初始化时，如果有缓存数据，我们会优先返回缓存数据，然后在背后发送新请求，也就是 SWR 的能力。
  你可以通过 options.staleTime 设置数据保持新鲜时间，在该时间内，我们认为数据是新鲜的，不会重新发起请求。
  你也可以通过 options.cacheTime 设置数据缓存时间，超过该时间，我们会清空该条缓存数据。

  - cacheKey
    请求 Promise 共享：相同的 cacheKey 同时只会有一个在发起请求，后发起的会共用同一个请求 Promise(`singleflight`)
    数据同步：当某个 cacheKey 发起请求时，其它相同 cacheKey 的内容均会随之同步
  - cacheTime
    缓存数据回收时间，默认 5 分钟，-1 表示不过期
  - staleTime
    缓存数据新鲜时间，默认 0，-1 表示永远新鲜
  - setCache、getCache、clearCache 自定义缓存
    ```tsx
    type CachedKey = string | number
    export interface CachedData<TData = any, TParams = any> {
      data: TData
      params: TParams
      time: number
    }
    interface RecordData extends CachedData {
      timer: Timer | undefined
    }
    declare const setCache: (key: CachedKey, cacheTime: number, cachedData: CachedData) => void
    declare const getCache: (key: CachedKey) => RecordData | undefined
    declare const clearCache: (key?: string | string[]) => void // 单个、多个、所有
    ```

## useAntdTable

集成了 antd 的 表格和表单

### Example

```tsx
import { Table } from 'antd'
import React from 'react'
import { useAntdTable } from 'ahooks'

interface Item {
  name: {
    last: string
  }
  email: string
  phone: string
  gender: 'male' | 'female'
}

interface Result {
  total: number
  list: Item[]
}

const getTableData = ({ current, pageSize }): Promise<Result> => {
  const query = `page=${current}&size=${pageSize}`

  return fetch(`https://randomuser.me/api?results=55&${query}`)
    .then(res => res.json())
    .then(res => ({
      total: res.info.results,
      list: res.results
    }))
}

export default () => {
  const { tableProps } = useAntdTable(getTableData)

  const columns = [
    {
      title: 'name',
      dataIndex: ['name', 'last']
    },
    {
      title: 'email',
      dataIndex: 'email'
    },
    {
      title: 'phone',
      dataIndex: 'phone'
    },
    {
      title: 'gender',
      dataIndex: 'gender'
    }
  ]

  return <Table columns={columns} rowKey="email" style={{ overflow: 'auto' }} {...tableProps} />
}
```

## useInfiniteScroll

无限滚动、下拉场景使用

### Example

```tsx
import React from 'react'
import { useInfiniteScroll } from 'ahooks'

interface Result {
  list: string[]
  nextId: string | undefined
}

const resultData = ['0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '10', '11', '12', '13']

function getLoadMoreList(nextId: string | undefined, limit: number): Promise<Result> {
  let start = 0
  if (nextId) {
    start = resultData.findIndex(i => i === nextId)
  }
  const end = start + limit
  const list = resultData.slice(start, end)
  const nId = resultData.length >= end ? resultData[end] : undefined
  return new Promise(resolve => {
    setTimeout(() => {
      resolve({
        list,
        nextId: nId
      })
    }, 1000)
  })
}

export default () => {
  const { data, loading, loadMore, loadingMore } = useInfiniteScroll(d => getLoadMoreList(d?.nextId, 4))

  return (
    <div>
      {loading ? (
        <p>loading</p>
      ) : (
        <div>
          {data?.list?.map(item => (
            <div key={item} style={{ padding: 12, border: '1px solid #f5f5f5' }}>
              item-{item}
            </div>
          ))}
        </div>
      )}

      <div style={{ marginTop: 8 }}>
        {data?.nextId && (
          <button type="button" onClick={loadMore} disabled={loadingMore}>
            {loadingMore ? 'Loading more...' : 'Click to load more'}
          </button>
        )}

        {!data?.nextId && <span>No more data</span>}
      </div>
    </div>
  )
}
```

### 解析

Result 与 Options

```tsx
declare const useInfiniteScroll: <TData extends Data>(service: Service<TData>, options?: InfiniteScrollOptions<TData>) => InfiniteScrollResult<TData>

export type Service<TData extends Data> = (currentData?: TData) => Promise<TData>
export interface InfiniteScrollResult<TData extends Data> {
  data: TData
  loading: boolean // 是否正在进行首次请求
  loadingMore: boolean // 是否正在进行更多数据请求
  error?: Error
  noMore: boolean
  loadMore: () => void
  loadMoreAsync: () => Promise<TData>
  reload: () => void
  reloadAsync: () => Promise<TData>
  cancel: () => void
  mutate: (data?: TData) => void
}

export interface InfiniteScrollOptions<TData extends Data> {
  target?: BasicTarget<Element | Document> // 滚动的
  isNoMore?: (data?: TData) => boolean
  direction?: 'bottom' | 'top'

  threshold?: number // 滚动到距离底部多少像素时触发 loadMore

  manual?: boolean
  reloadDeps?: DependencyList

  onBefore?: () => void
  onSuccess?: (data: TData) => void
  onError?: (e: Error) => void
  onFinally?: (data?: TData, e?: Error) => void
}
```

- 滚动自动加载

  - options.target 指定父级元素（父级元素需设置固定高度，且支持内部滚动）
    父级容器，如果存在，则在滚动到底部时，自动触发 loadMore。
    需要配合 isNoMore 使用，以便知道什么时候到最后一页了。 当 target 为 document 时，定义为整个视口
  - options.isNoMore 判断是不是没有更多数据了
  - options.direction 滚动的方向，默认为向下滚动

  ```tsx
  const ref = useRef<HTMLDivElement>(null)
  const { data, loading, loadMore, loadingMore, noMore } = useInfiniteScroll(d => getLoadMoreList(d?.nextId, 4), {
    target: ref,
    isNoMore: d => d?.nextId === undefined
  })
  ```

- 数据重置
  当 reloadDeps 变化时，会自动触发 reload
- 数据突变

## usePagination

分页场景

### Example

```tsx
import { usePagination } from 'ahooks'
import { Pagination } from 'antd'
import Mock from 'mockjs'
import React from 'react'

interface UserListItem {
  id: string
  name: string
  gender: 'male' | 'female'
  email: string
  disabled: boolean
}

const userList = (current: number, pageSize: number) =>
  Mock.mock({
    total: 55,
    [`list|${pageSize}`]: [
      {
        id: '@guid',
        name: '@name',
        'gender|1': ['male', 'female'],
        email: '@email',
        disabled: false
      }
    ]
  })

async function getUserList(params: { current: number; pageSize: number }): Promise<{ total: number; list: UserListItem[] }> {
  return new Promise(resolve => {
    setTimeout(() => {
      resolve(userList(params.current, params.pageSize))
    }, 1000)
  })
}

export default () => {
  const { data, loading, pagination } = usePagination(getUserList)
  return (
    <div>
      {loading ? (
        <p>loading</p>
      ) : (
        <ul>
          {data?.list?.map(item => (
            <li key={item.email}>
              {item.name} - {item.email}
            </li>
          ))}
        </ul>
      )}
      <Pagination
        current={pagination.current}
        pageSize={pagination.pageSize}
        total={data?.total}
        onChange={pagination.onChange}
        onShowSizeChange={pagination.onChange}
        showQuickJumper
        showSizeChanger
        style={{ marginTop: 16, textAlign: 'right' }}
      />
    </div>
  )
}
```

### 解析

- Options 与 Result 继承了 useRequest 的 Options 和 Result
- service 的第一个参数为 { current: number, pageSize: number }
- service 返回的数据结构为 { total: number, list: Item[] }
- 会额外返回 pagination 字段，包含所有分页信息，及操作分页的函数。
- refreshDeps 变化，会重置 current 到第一页，并重新发起请求，一般你可以把 pagination 依赖的条件放这里

```tsx
declare const usePagination: <TData extends Data, TParams extends Params>(
  service: Service<TData, TParams>,
  options?: PaginationOptions<TData, TParams>
) => PaginationResult<TData, TParams>

export type Data = {
  total: number
  list: any[]
}
export type Params = [
  {
    current: number
    pageSize: number
    [key: string]: any
  },
  ...any[]
]
export type Service<TData extends Data, TParams extends Params> = (...args: TParams) => Promise<TData>
export interface PaginationResult<TData extends Data, TParams extends Params> extends Result<TData, TParams> {
  pagination: {
    current: number
    pageSize: number
    total: number
    totalPage: number

    onChange: (current: number, pageSize: number) => void
    changeCurrent: (current: number) => void
    changePageSize: (pageSize: number) => void
  }
}
export interface PaginationOptions<TData extends Data, TParams extends Params> extends Options<TData, TParams> {
  defaultCurrent?: number
  defaultPageSize?: number
}
```

## useDynamicList

一个帮助你管理动态列表状态，并能`生成唯一 key` 的 Hook。

### Example

### 解析

```ts
declare const useDynamicList: <T>(initialList?: T[]) => {
  list: T[]
  insert: (index: number, item: T) => void
  merge: (index: number, items: T[]) => void
  replace: (index: number, item: T) => void
  remove: (index: number) => void
  batchRemove: (indexes: number[]) => void
  getKey: (index: number) => number
  getIndex: (key: number) => number
  move: (oldIndex: number, newIndex: number) => void
  push: (item: T) => void
  pop: () => void
  unshift: (item: T) => void
  shift: () => void
  sortList: (result: T[]) => T[] // 校准排序?
  resetList: (newList: T[]) => void
}
```

## useVirtualList

虚拟渲染。
提供虚拟化列表能力的 Hook，用于解决展示海量数据渲染时首屏渲染缓慢和滚动卡顿问题。
支持动态元素高度。

### Example

```tsx
import React, { useMemo, useRef } from 'react'
import { useVirtualList } from 'ahooks'

export default () => {
  const containerRef = useRef(null)
  const wrapperRef = useRef(null)

  const originalList = useMemo(() => Array.from(Array(99999).keys()), [])

  const [value, onChange] = React.useState<number>(0)

  const [list, scrollTo] = useVirtualList(originalList, {
    containerTarget: containerRef,
    wrapperTarget: wrapperRef,
    itemHeight: i => (i % 2 === 0 ? 42 + 8 : 84 + 8),
    overscan: 10
  })

  return (
    <div>
      <div style={{ textAlign: 'right', marginBottom: 16 }}>
        <input style={{ width: 120 }} placeholder="line number" type="number" value={value} onChange={e => onChange(Number(e.target.value))} />
        <button
          style={{ marginLeft: 8 }}
          type="button"
          onClick={() => {
            scrollTo(Number(value))
          }}
        >
          scroll to
        </button>
      </div>
      <div ref={containerRef} style={{ height: '300px', overflow: 'auto' }}>
        <div ref={wrapperRef}>
          {list.map(ele => (
            <div
              style={{
                height: ele.index % 2 === 0 ? 42 : 84,
                display: 'flex',
                justifyContent: 'center',
                alignItems: 'center',
                border: '1px solid #e8e8e8',
                marginBottom: 8
              }}
              key={ele.index}
            >
              Row: {ele.data} size: {ele.index % 2 === 0 ? 'small' : 'large'}
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
```

### 解析

```tsx
type ItemHeight<T> = (index: number, data: T) => number
export interface Options<T> {
  containerTarget: BasicTarget // 外部容器
  wrapperTarget: BasicTarget // 内部容器
  itemHeight: number | ItemHeight<T> // 行高度，静态高度可以直接写入像素值，动态高度可传入函数
  overscan?: number // 视区上、下额外展示的 DOM 节点数量
}
declare const useVirtualList: <T = any>(
  list: T[], // 包含大量数据的列表。 注意：必须经过 useMemo 处理或者永不变化，否则会死循环
  options: Options<T>
) => readonly [
  {
    index: number
    data: T
  }[],
  (index: number) => void
]
```

## useHistoryTravel

管理状态历史变化记录，方便在历史记录中前进与后退。

### Example

```tsx
import { useHistoryTravel } from 'ahooks'
import React from 'react'

export default () => {
  const maxLength = 3
  const { value, setValue, backLength, forwardLength, back, forward } = useHistoryTravel<string>('', maxLength)

  return (
    <div>
      <div>maxLength: {maxLength}</div>
      <div>backLength: {backLength}</div>
      <div>forwardLength: {forwardLength}</div>
      <input value={value || ''} onChange={e => setValue(e.target.value)} />
      <button disabled={backLength <= 0} onClick={back} style={{ margin: '0 8px' }}>
        back
      </button>
      <button disabled={forwardLength <= 0} onClick={forward}>
        forward
      </button>
    </div>
  )
}
```

### 解析

```ts
export default function useHistoryTravel<T>(
  initialValue?: T,
  maxLength?: number
): {
  value: T | undefined
  backLength: number
  forwardLength: number
  setValue: (val: T) => void
  go: (step: number) => void
  back: () => void
  forward: () => void
  reset: (...params: any[]) => void
}
```
