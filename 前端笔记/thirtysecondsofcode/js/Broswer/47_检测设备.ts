// 检测是否正在移动设备或桌面上查看页面。
detectDeviceType() // 'Mobile' or 'Desktop'

function detectDeviceType() {
  return /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent)
    ? 'Mobile'
    : 'Desktop'
}
