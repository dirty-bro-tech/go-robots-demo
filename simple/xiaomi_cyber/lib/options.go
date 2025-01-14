package lib

// FeatureCode defines a struct to represent the parameters
type FeatureCode struct {
	Name string
	Code int16
	Typ  string
}

type featureEnum struct {
	MotorOverTemp FeatureCode
	OverTempTime  FeatureCode
	LimitTorque   FeatureCode
	CurKp         FeatureCode
	CurKi         FeatureCode
	SpdKp         FeatureCode
	SpdKi         FeatureCode
	LocKp         FeatureCode
	SpdFiltGain   FeatureCode
	LimitSpd      FeatureCode
	LimitCur      FeatureCode
}

type FeatureParam struct {
	Name   string
	Index  int
	Format string
}

var (
	featureEnums = featureEnum{
		MotorOverTemp: FeatureCode{Name: "motorOverTemp", Code: 0x200D, Typ: "int16"},
		OverTempTime:  FeatureCode{Name: "overTempTime", Code: 0x200E, Typ: "int32"},
		LimitTorque:   FeatureCode{Name: "limit_torque", Code: 0x2007, Typ: "float"},
		CurKp:         FeatureCode{Name: "cur_kp", Code: 0x2012, Typ: "float"},
		CurKi:         FeatureCode{Name: "cur_ki", Code: 0x2013, Typ: "float"},
		SpdKp:         FeatureCode{Name: "spd_kp", Code: 0x2014, Typ: "float"},
		SpdKi:         FeatureCode{Name: "spd_ki", Code: 0x2015, Typ: "float"},
		LocKp:         FeatureCode{Name: "loc_kp", Code: 0x2016, Typ: "float"},
		SpdFiltGain:   FeatureCode{Name: "spd_filt_gain", Code: 0x2017, Typ: "float"},
		LimitSpd:      FeatureCode{Name: "limit_spd", Code: 0x2018, Typ: "float"},
		LimitCur:      FeatureCode{Name: "limit_cur", Code: 0x2019, Typ: "float"},
	}

	featureParams = map[string]FeatureParam{
		"run_mode":      {Name: "run_mode", Index: 0x7005, Format: "u8"},
		"iq_ref":        {Name: "iq_ref", Index: 0x7006, Format: "f"},
		"spd_ref":       {Name: "spd_ref", Index: 0x700A, Format: "f"},
		"limit_torque":  {Name: "limit_torque", Index: 0x700B, Format: "f"},
		"cur_kp":        {Name: "cur_kp", Index: 0x7010, Format: "f"},
		"cur_ki":        {Name: "cur_ki", Index: 0x7011, Format: "f"},
		"cur_filt_gain": {Name: "cur_filt_gain", Index: 0x7014, Format: "f"},
		"loc_ref":       {Name: "loc_ref", Index: 0x7016, Format: "f"},
		"limit_spd":     {Name: "limit_spd", Index: 0x7017, Format: "f"},
		"limit_cur":     {Name: "limit_cur", Index: 0x7018, Format: "f"},
	}

	tpyMap = map[string]int{
		"float": 0x06,
		"int16": 0x03,
		"int32": 0x04,
	}

	twoBytesBits uint = 16
)
