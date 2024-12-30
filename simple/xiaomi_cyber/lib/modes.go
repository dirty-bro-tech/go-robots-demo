package lib

const (
	CMDModelGetDeviceId       = 0
	CMDModelMotorControl      = 1
	CMDModelMotorFeedback     = 2
	CMDModelMotorEnable       = 3
	CMDModelMotorStop         = 4
	CMDModelSetMechanicalZero = 6
	CMDModelSetMotorCanId     = 7
	CMDModelParamTableWrite   = 8
	CMDModelSingleParamRead   = 17
	CMDModelSingleParamWrite  = 18
	CMDModelFaultFeedback     = 21
)

const (
	RunModeControlMode  = 0 // 运控模式
	RunModePositionMode = 1 // 位置模式
	RunModeSpeedMode    = 2 // 速度模式
	RunModeCurrentMode  = 3 // 电流模式
)
