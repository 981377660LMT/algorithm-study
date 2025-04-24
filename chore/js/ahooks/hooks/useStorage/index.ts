import { createUseStorageState } from './createUseStorage'

/**
 * readeMe: https://ahooks.js.org/zh-CN/hooks/use-local-storage-state
 */

export const useLocalStorageState = createUseStorageState(() => localStorage)

export const useSessionStorageState = createUseStorageState(() => sessionStorage)
