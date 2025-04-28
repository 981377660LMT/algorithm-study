import {
  useAntdTable,
  useCountDown,
  useCounter,
  useDynamicList,
  useHistoryTravel,
  useInfiniteScroll,
  useNetwork,
  usePagination,
  useRequest,
  useSelections,
  useTextSelection,
  useTheme,
  useVirtualList,
  useWebSocket
} from 'ahooks'
import React from 'react'

interface IFooProps {}

const Foo: React.FC<IFooProps> = props => {
  useRequest
  useAntdTable
  useInfiniteScroll
  usePagination
  useDynamicList
  useVirtualList
  useHistoryTravel
  useNetwork
  useSelections
  useCountDown
  useCounter
  useTextSelection
  useWebSocket
  useTheme

  return <></>
}

export default Foo
