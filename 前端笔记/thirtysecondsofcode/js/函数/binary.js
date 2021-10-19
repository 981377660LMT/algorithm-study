const binary = fn => (a, b) => fn(a, b)
;['2', '1', '0'].map(binary(Math.max)) // [2, 1, 2]
