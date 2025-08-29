## 安装 Rust 和 Cargo

您可以通过 `rustup`（Rust 工具链安装程序）来安装 Rust 和 Cargo。`rustup` 会安装 `rustc`（编译器）、`cargo`（包管理器和构建工具）以及其他标准库。

在您的终端中运行以下命令：

```bash
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
```

这个命令会下载并运行官方安装脚本。它会引导您完成安装过程，通常选择默认选项即可。

安装完成后，您需要配置当前 shell 的环境，以便能找到 Cargo 的路径。运行以下命令：

```bash
source "$HOME/.cargo/env"
```

或者，您可以重新启动终端。

最后，通过运行以下命令来验证安装是否成功：

```bash
rustc --version
cargo --version
```

如果安装成功，您将看到 `rustc` 和 `cargo` 的版本信息。

