## 使用 virtualbox 开发 windows 应用

virtualbox 启动虚机失败解决方案[Fix: Virtual machine has terminated unexpectedly during startup with exit code 1 (0x1)](https://appuals.com/fix-virtual-machine-has-terminated-unexpectedly-during-startup-with-exit-code-1-0x1/)

#### windows 开发环境

1，虚机安装增强工具时勾选 3D 加速，如果虚机是 windows7 那么第 2 步开启自动更新后会自动安装 3D 加速补丁[KB2670838](https://support.microsoft.com/en-us/help/2670838/platform-update-for-windows-7-sp1-and-windows-server-2008-r2-sp1)

> KB2670838 summary:  
> The Platform Update for Windows 7 enables improved features and performance on Windows 7 SP1 and Windows Server 2008 R2 SP1. It includes updates to the following components: Direct2D, DirectWrite, Direct3D, Windows Imaging Component (WIC), Windows Advanced Rasterization Platform (WARP), Windows Animation Manager (WAM), XPS Document API , the H.264 Video Decoder and the JPEG XR Codec.

2，自动更新: 配置 windows 虚机时 cable connected 先不勾选，等进了 windows 系统禁用了自动更新再在右下角点击 vd 的网络图标连接上网，开启自动更新直到更新到最新再做别的事，比如安装 vs2015

#### windows 运行环境 for server

1，禁止更新: 配置 windows 虚机时 cable connected 先不勾选，等进了 windows 系统禁用了自动更新再在右下角点击 vd 的网络图标连接上网  
2，虚拟机安装增强工具时不必勾选 3D 加速

#### windows 运行环境 for client

1，禁止更新: 配置 windows 虚机时 cable connected 先不勾选，等进了 windows 系统禁用了自动更新再在右下角点击 vd 的网络图标连接上网  
2，虚机安装增强工具时勾选 3D 加速，如果虚机是 windows7 那么需要手动安装 3D 加速补丁[KB2670838](https://support.microsoft.com/en-us/help/2670838/platform-update-for-windows-7-sp1-and-windows-server-2008-r2-sp1)并且开启 Aero 主题
