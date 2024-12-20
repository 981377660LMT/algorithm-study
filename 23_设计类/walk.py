import os

# 遍历指定目录及其所有子目录
for root, dirs, files in os.walk("."):
    print(f"当前目录: {root}")
    print(f"子目录: {dirs}")
    print(f"文件: {files}")
    print("---")
