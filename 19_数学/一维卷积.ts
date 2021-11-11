function convolve(volume: number[], kernel: number[]): number[] {
  if (volume.length === 0 || kernel.length === 0) throw new Error('数组不能为空')

  const res = Array(volume.length + kernel.length - 1).fill(0)
  for (let i = 0; i < volume.length; i++) {
    for (let j = 0; j < kernel.length; j++) {
      res[i + j] += volume[i] * kernel[j]
    }
  }

  return res
}

if (require.main === module) {
  console.log(convolve([1, 2, 3], [0, 1, 0.5]))
}
