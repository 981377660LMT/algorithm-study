/**
 * 将 num 拆分成 cands[i] 的线性组合，使得拆分的个数最(多/少).
 * @param num 非负数.
 * @param cands 非负数数组.
 * @param minimize 是否使得拆分的个数最少. 默认为true.
 * @returns [counts, ok] counts是拆分成cands[i]的个数，ok表示是否可以拆分.
 */
function splitTo(num: number, cands: number[], minimize = true): [counts: number, ok: boolean] {}

/**
 * 完全背包.
 */
function _solve1(num: number, cands: number[], minimize = true): [cands: number, ok: boolean] {}

/**
 * 同余最短路.
 */
function _solve2(num: number, cands: number[], minimize = true): [cands: number, ok: boolean] {}

export { splitTo }
