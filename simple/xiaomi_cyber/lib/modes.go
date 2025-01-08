package lib

type CMDModel int

const (
	CMDModelGetDeviceId       CMDModel = 0
	CMDModelMotorControl      CMDModel = 1
	CMDModelMotorFeedback     CMDModel = 2
	CMDModelMotorEnable       CMDModel = 3
	CMDModelMotorStop         CMDModel = 4
	CMDModelSetMechanicalZero CMDModel = 6
	CMDModelSetMotorCanId     CMDModel = 7
	CMDModelParamTableWrite   CMDModel = 8
	CMDModelSingleParamRead   CMDModel = 17
	CMDModelSingleParamWrite  CMDModel = 18
	CMDModelFaultFeedback     CMDModel = 21
)

const (
	RunModeControlMode  RunMode = 0 // 运控模式
	RunModePositionMode RunMode = 1 // 位置模式
	RunModeSpeedMode    RunMode = 2 // 速度模式
	RunModeCurrentMode  RunMode = 3 // 电流模式
)

type RunMode int
