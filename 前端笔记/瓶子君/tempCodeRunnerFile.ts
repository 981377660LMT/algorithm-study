nst flatten1 = (nums: NestedArray): number[] =>
//   nums.reduce<number[]>(
//     (pre, cur) => (Array.isArray(cur) ? [...pre, ...flatten(cur)] : [...pre, cur]),
//     []
//   )