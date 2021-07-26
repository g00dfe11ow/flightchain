package model

// Plan describes basic details of what makes up a simple plan
type Plan struct {
	ID               string   `json:"id"`                 // ID
	UnitOrIndividual string   `json:"unit_or_individual"` // 执行单位或个人
	Contact          string   `json:"contact"`            // 联系方式
	Aircraft         Aircraft `json:"aircraft"`           // 航空相关信息
	State            string   `json:"state"`              // 状态
	Applicant        string   `json:"applicant"`          // 申请人
	ReadSet          []string `json:"readset"`            // ReadSet
}

type Aircraft struct {
	Nationality  string `json:"nationality"`  // 航空器国籍
	Type         string `json:"type"`         // 型别
	Number       int    `json:"number"`       // 架数
	CallSign     int    `json:"callsign"`     // 呼号
	Registration int    `json:"registration"` // 注册号
}

type NewPlan struct {
	ID               string   `json:"id"`                 // ID
	UnitOrIndividual string   `json:"unit_or_individual"` // 执行单位或个人
	Contact          string   `json:"contact"`            // 联系方式
	Aircraft         Aircraft `json:"aircraft"`           // 航空相关信息
}

type QueryID struct {
	ID string `json:"id"`
}

type ApprovalType struct {
	ID    string `json:"id"`
	State string `json:"state"`
}

type Readers struct {
	ID     string `json:"id"`
	Reader string `json:"reader"`
}
