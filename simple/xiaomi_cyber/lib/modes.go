package lib

type Mode int

func (m Mode) Value() int {
	return int(m)
}

type CMDMode Mode

func (m CMDMode) Mode() Mode {
	return Mode(m)
}

type RunMode Mode

func (r RunMode) Mode() Mode {
	return Mode(r)
}

const (
	CMDModeGetDeviceId       CMDMode = 0
	CMDModeMotorControl      CMDMode = 1
	CMDModeMotorFeedback     CMDMode = 2
	CMDModeMotorEnable       CMDMode = 3
	CMDModeMotorStop         CMDMode = 4
	CMDModeSetMechanicalZero CMDMode = 6
	CMDModeSetMotorCanId     CMDMode = 7
	CMDModeParamTableWrite   CMDMode = 8
	CMDModeSingleParamRead   CMDMode = 17
	CMDModeSingleParamWrite  CMDMode = 18
	CMDModeFaultFeedback     CMDMode = 21
)

const (
	RunModeControlMode  RunMode = 0 // 运控模式
	RunModePositionMode RunMode = 1 // 位置模式
	RunModeSpeedMode    RunMode = 2 // 速度模式
	RunModeCurrentMode  RunMode = 3 // 电流模式
)
