import type { Mutate, StateCreator, StoreApi, StoreMutatorIdentifier } from './vanilla.mjs';
type ExtractState<S> = S extends {
    getState: () => infer T;
} ? T : never;
type ReadonlyStoreApi<T> = Pick<StoreApi<T>, 'getState' | 'subscribe'>;
type WithReact<S extends ReadonlyStoreApi<unknown>> = S & {
    /** @deprecated please use api.getState() */
    getServerState?: () => ExtractState<S>;
};
export declare function useStore<S extends WithReact<StoreApi<unknown>>>(api: S): ExtractState<S>;
export declare function useStore<S extends WithReact<StoreApi<unknown>>, U>(api: S, selector: (state: ExtractState<S>) => U): U;
/**
 * @deprecated Use `useStoreWithEqualityFn` from 'zustand/traditional'
 * https://github.com/pmndrs/zustand/discussions/1937
 */
export declare function useStore<S extends WithReact<StoreApi<unknown>>, U>(api: S, selector: (state: ExtractState<S>) => U, equalityFn: ((a: U, b: U) => boolean) | undefined): U;
export type UseBoundStore<S extends WithReact<ReadonlyStoreApi<unknown>>> = {
    (): ExtractState<S>;
    <U>(selector: (state: ExtractState<S>) => U): U;
    /**
     * @deprecated Use `createWithEqualityFn` from 'zustand/traditional'
     */
    <U>(selector: (state: ExtractState<S>) => U, equalityFn: (a: U, b: U) => boolean): U;
} & S;
type Create = {
    <T, Mos extends [StoreMutatorIdentifier, unknown][] = []>(initializer: StateCreator<T, [], Mos>): UseBoundStore<Mutate<StoreApi<T>, Mos>>;
    <T>(): <Mos extends [StoreMutatorIdentifier, unknown][] = []>(initializer: StateCreator<T, [], Mos>) => UseBoundStore<Mutate<StoreApi<T>, Mos>>;
    /**
     * @deprecated Use `useStore` hook to bind store
     */
    <S extends StoreApi<unknown>>(store: S): UseBoundStore<S>;
};
export declare const create: Create;
/**
 * @deprecated Use `import { create } from 'zustand'`
 */
declare const _default: Create;
export default _default;
