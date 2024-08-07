# Servo Demo

## 舵机原理

* [舵机运动原理](https://fbc-wiki.readthedocs.io/zh/latest/basis_part/steering_gear_control.html)

## [SimpleAngle](./servo/simple_angle.go)

* 说明：
  * 简单For循环转动固定几个角度
* 命令
  * > $ tinygo flash -target=arduino -port=/dev/cu.usbserial-21420 
  * -port 改成自己的设置

## [simple_serial.go](./serial/simple_serial.go)

* 说明：
  * 简单For循环，简单将输入打印

## [serial_servo.go](./serial_servo/serial_servo.go)

* 说明：
  * 简单For循环，简单将输入打印