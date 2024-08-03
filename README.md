# Demo

## 运行

记得自行加上运行环境，gobot要运行在arm架构下

```shell
GOOS=linux GOARCH=arm GOARM=6 go run xxx.go
```

## 异常

```
fork/exec /var/folders/b1/0fd1b6hs7lz0fm_mh346lybm0000gn/T/go-build844665855/b001/exe/toggle-every-one-second: exec format error
```


### Perhaps you need to flash your Arduino with Firmata?

* 解决方案
```text
Open arduino IDE and go to File > Examples > Firmata > StandardFirmata and open it. Select the appriate port for your arduino and click upload.
```

## 虚拟机
 
* 安装 vm tools
```shell
sudo apt update
sudo apt install open-vm-tools
sudo apt install open-vm-tools-desktop
```

## 共享目录

> [参考](https://askubuntu.com/questions/29284/how-do-i-mount-shared-folders-in-ubuntu-using-vmware-tools)

```shell
mkdir /mnt/hgfs/
sudo vmhgfs-fuse .host:/ /mnt/hgfs/ -o allow_other -o uid=1000
sudo vmhgfs-fuse .host:/ /mnt/ -o allow_other -o uid=1000


# 如果报没权限
sudo vmhgfs-fuse .host:/ /mnt/hgfs -o subtype=vmhgfs-fuse,allow_other
```
