# godep
维护golang项目的依赖包

## 环境

- 本项目编译之前，需用户自行安装`Golang`；
- 编译时间视具体网络环境而定；
- 要求`linux`或`macos`系统即可。

## 说明

该工具依赖glide.yaml文件依次拉去依赖包，不支持间接依赖包的自动获取(需用户把间接依赖的包手动维护到glide.yaml).

## 获取

git clone https://github.com/dongjialong2006/godep.git

## 编译

- make

## 安装

- 将godep放到PATH所引到的bin目录下即可.
- go get -u github.com/dongjialong2006/godep

## 命令

- 在y有glide.yaml文件的目录下执行`godep`即可.

