package lib

type Mode int
type CMDMode Mode

func (m CMDMode) Mode() Mode {
	return Mode(m)
}

type RunMode Mode

func (r RunMode) Mode() Mode {
	return Mode(r)
}

const (
	CMDModelGetDeviceId       CMDMode = 0
	CMDModelMotorControl      CMDMode = 1
	CMDModelMotorFeedback     CMDMode = 2
	CMDModelMotorEnable       CMDMode = 3
	CMDModelMotorStop         CMDMode = 4
	CMDModelSetMechanicalZero CMDMode = 6
	CMDModelSetMotorCanId     CMDMode = 7
	CMDModelParamTableWrite   CMDMode = 8
	CMDModelSingleParamRead   CMDMode = 17
	CMDModelSingleParamWrite  CMDMode = 18
	CMDModelFaultFeedback     CMDMode = 21
)

const (
	RunModeControlMode  RunMode = 0 // 运控模式
	RunModePositionMode RunMode = 1 // 位置模式
	RunModeSpeedMode    RunMode = 2 // 速度模式
	RunModeCurrentMode  RunMode = 3 // 电流模式
)
