import React from 'react'
import useUrlState from '@ahooksjs/use-url-state'
import {
  useAntdTable,
  useAsyncEffect,
  useBoolean,
  useClickAway,
  useControllableValue,
  useCookieState,
  useCountDown,
  useCounter,
  useCreation,
  useDebounce,
  useDebounceEffect,
  useDebounceFn,
  useDeepCompareEffect,
  useDeepCompareLayoutEffect,
  useDocumentVisibility,
  useDrag,
  useDrop,
  useDynamicList,
  useEventEmitter,
  useEventListener,
  useEventTarget,
  useExternal,
  useFavicon,
  useFocusWithin,
  useFullscreen,
  useGetState,
  useHistoryTravel,
  useHover,
  useInfiniteScroll,
  useInterval,
  useInViewport,
  useIsomorphicLayoutEffect,
  useKeyPress,
  useLatest,
  useLocalStorageState,
  useLockFn,
  useLongPress,
  useMap,
  useMemoizedFn,
  useMount,
  useMouse,
  useMutationObserver,
  useNetwork,
  usePagination,
  usePrevious,
  useRafInterval,
  useRafState,
  useRafTimeout,
  useReactive,
  useRequest,
  useResetState,
  useResponsive,
  useSafeState,
  useScroll,
  useSelections,
  useSessionStorageState,
  useSet,
  useSetState,
  useSize,
  useTextSelection,
  useTheme,
  useThrottle,
  useThrottleEffect,
  useThrottleFn,
  useTimeout,
  useTitle,
  useToggle,
  useTrackedEffect,
  useUnmount,
  useUnmountedRef,
  useUpdate,
  useUpdateEffect,
  useUpdateLayoutEffect,
  useVirtualList,
  useWebSocket,
  useWhyDidYouUpdate
} from 'ahooks'

interface IFooProps {}

const Foo: React.FC<IFooProps> = props => {
  // UseRequest
  useRequest

  // Scene
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

  // Lifecycle
  useMount
  useUnmount
  useUnmountedRef

  // State
  useSetState
  useBoolean
  useToggle
  useUrlState // 这个在源码里独立于其他hooks
  useCookieState
  useLocalStorageState
  useSessionStorageState
  useDebounce
  useThrottle
  useMap
  useSet
  usePrevious
  useRafState
  useSafeState
  useGetState
  useResetState

  // Effect
  useUpdateEffect
  useUpdateLayoutEffect
  useAsyncEffect
  useDebounceEffect
  useDebounceFn
  useThrottleFn
  useThrottleEffect
  useDeepCompareEffect
  useDeepCompareLayoutEffect
  useInterval
  useRafInterval
  useTimeout
  useRafTimeout
  useLockFn
  useUpdate

  // Dom
  useEventListener
  useClickAway
  useDocumentVisibility
  useDrop
  useDrag
  useEventTarget
  useExternal
  useTitle
  useFavicon
  useFullscreen
  useHover
  useMutationObserver
  useInViewport
  useKeyPress
  useLongPress
  useMouse
  useResponsive
  useSize
  useScroll
  useFocusWithin

  // Advanced
  useControllableValue
  useCreation
  useEventEmitter
  useIsomorphicLayoutEffect
  useLatest
  useMemoizedFn
  useReactive

  // Dev
  useTrackedEffect
  useWhyDidYouUpdate

  return <></>
}

export default Foo
