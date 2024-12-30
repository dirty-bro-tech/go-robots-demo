package lib

type CyberGearController struct {
	Bus      interface{} // todo
	MotorId  int
	MainCNId int
	PMin     float64
	PMax     float64
	VMin     float64
	VMax     float64
	TMin     float64
	TMax     float64
	KpMin    float64 // 0.0 ~
	KpMax    float64 // ~ 500.0
	KdMin    float64 // 0.0 ~
	KdMax    float64 //    ~ 5.0
}

func NewController() CyberGearController {

}
