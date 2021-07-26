package ref

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Plan
type SmartContract struct {
	contractapi.Contract
}

// Plan describes basic details of what makes up a simple plan
type Plan struct {
	ID               string   `json:"id"`                 // ID
	UnitOrIndividual string   `json:"unit_or_individual"` // 执行单位或个人
	Contact          int64    `json:"contact"`            // 联系方式
	Aircraft         Aircraft `json:"aircraft"`           // 航空相关信息
	State            string   `json:"state"`              // 状态
	Applicant        string   `json:"applicant"`          // 申请人
	ReadSet          []string `json:"readset"`            // ReadSet
	// Pilot            Pilot    `json:"pilot"`            // 飞行人员相关信息
	// Airport          string   `json:"airport`           // 使用机场
	// TakeOff          string   `json:"takeoff"`          // 起飞地点
	// Landing          string   `json:"landing`           // 降落地点
	// PlanDate         int64    `json:"plandate"`         // 计划执行时间
	// TakeOffTime      int64    `json:"takeofftime"`      // 起飞时间
	// LandingTime      int64    `json:"landingtime"`      // 降落时间
	// FlyTime          int      `json:"flytime"`          // 飞行次数
}

type Aircraft struct {
	Nationality  string `json:"nationality"`  // 航空器国籍
	Type         string `json:"type"`         // 型别
	Number       int    `json:"number"`       // 架数
	CallSign     int    `json:"callsign"`     // 呼号
	Registration int    `json:"registration"` // 注册号
}

// type Pilot struct {
// 	Nationality string `json:"nationality"` // 国籍
// 	Number      int    `json:"number"`      // 机组人数
// }

// InitLedger adds a base set of plans to the ledger
// 初始化账本，主要用于测试
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	plans := []Plan{
		{ID: "plan1", UnitOrIndividual: "医疗", Contact: 18018595217, Aircraft: Aircraft{Nationality: "China", Type: "波音747", Number: 1, CallSign: 111, Registration: 1}, State: "appling", ReadSet: []string{}},
		{ID: "plan2", UnitOrIndividual: "救火", Contact: 18018595217, Aircraft: Aircraft{Nationality: "China", Type: "空客111", Number: 1, CallSign: 112, Registration: 2}, State: "appling", ReadSet: []string{}},
		{ID: "plan3", UnitOrIndividual: "救人", Contact: 18018595217, Aircraft: Aircraft{Nationality: "China", Type: "波音747", Number: 1, CallSign: 113, Registration: 3}, State: "appling", ReadSet: []string{}},
		{ID: "plan4", UnitOrIndividual: "救灾", Contact: 18018595217, Aircraft: Aircraft{Nationality: "China", Type: "波音747", Number: 1, CallSign: 114, Registration: 4}, State: "appling", ReadSet: []string{}},
	}

	for _, plan := range plans {
		planJSON, err := json.Marshal(plan)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(plan.ID, planJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// CreatePlan issues a new plan to the world state with given details.
// 创建飞行计划，只有ApplicantMSP有权限创建飞行计划
func (s *SmartContract) CreatePlan(ctx contractapi.TransactionContextInterface, id string, unitorindividual string, contact int64, aircraft Aircraft) error {
	// 只有ApplicantMSP能创建飞行计划
	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return err
	}
	if mspID != "ApplicantMSP" {
		return fmt.Errorf("you have no right to create a plan")
	}

	exists, err := s.PlanExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the plan %s already exists", id)
	}

	// Get ID of submitting client identity
	clientID, err := s.GetSubmittingClientIdentity(ctx)
	if err != nil {
		return err
	}

	plan := Plan{
		ID:               id,
		UnitOrIndividual: unitorindividual,
		Contact:          contact,
		Aircraft:         aircraft,
		State:            "applying",
		Applicant:        clientID,
		ReadSet:          []string{},
	}
	planJSON, err := json.Marshal(plan)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, planJSON)
}

// ReadPlan returns the plan stored in the world state with given id.
// 读计划，如果是applicant和approver可以直接读取，如果是user，则判断是否在ReadSet中
func (s *SmartContract) ReadPlan(ctx contractapi.TransactionContextInterface, id string) (*Plan, error) {
	// 如果是UserMSP要判断是否在ReadSet中
	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, err
	}
	if mspID == "UserMSP" {
		plan, err := s.ReadPlan(ctx, id)
		if err != nil {
			return nil, err
		}

		clientID, err := s.GetSubmittingClientIdentity(ctx)
		if err != nil {
			return nil, err
		}

		if _, res := InString(plan.ReadSet, clientID); !res {
			return nil, fmt.Errorf("you have no right to read this plan")
		}
	}

	planJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if planJSON == nil {
		return nil, fmt.Errorf("the plan %s does not exist", id)
	}

	var plan Plan
	err = json.Unmarshal(planJSON, &plan)
	if err != nil {
		return nil, err
	}

	return &plan, nil
}

// ModifyPlan updates an existing plan in the world state with provided parameters.
// 用参数来修改已有的飞行计划
func (s *SmartContract) ModifyPlan(ctx contractapi.TransactionContextInterface, id string, unitorindividual string, contact int64, aircraft Aircraft) error {
	// 只有ApplicantMSP能修改飞行计划
	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return err
	}
	if mspID != "ApplicantMSP" {
		return fmt.Errorf("you have no right to modify this plan")
	}

	plan, err := s.ReadPlan(ctx, id)
	if err != nil {
		return err
	}

	clientID, err := s.GetSubmittingClientIdentity(ctx)
	if err != nil {
		return err
	}

	if clientID != plan.Applicant {
		return fmt.Errorf("submitting client not authorized to modify plan, this is not your plan")
	}

	if plan.State == "applying" {
		return fmt.Errorf("submitting client not authorized to modify plan, the plan state is already changed")
	}

	// overwriting original plan with new plan
	*plan = Plan{
		ID:               id,
		UnitOrIndividual: unitorindividual,
		Contact:          contact,
		Aircraft:         aircraft,
		State:            plan.State,
		Applicant:        clientID,
	}
	planJSON, err := json.Marshal(plan)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, planJSON)
}

// DeletePlan deletes an given plan from the world state.
// 删除已有的飞行计划
func (s *SmartContract) DeletePlan(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.PlanExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the plan %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// PlanExists returns true when plan with given ID exists in world state
// 查询对应ID的飞行计划是否存在
func (s *SmartContract) PlanExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	planJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return planJSON != nil, nil
}

// ApprovalPlan updates the owner field of plan with given id in world state.
// approver审批提交上来的飞行计划
func (s *SmartContract) ApprovalPlan(ctx contractapi.TransactionContextInterface, id string, newState string) error {
	// 只有ApproverMSP能审批飞行计划
	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return err
	}
	if mspID != "ApproverMSP" {
		return fmt.Errorf("you have no right to approval a plan")
	}

	plan, err := s.ReadPlan(ctx, id)
	if err != nil {
		return err
	}

	plan.State = newState
	planJSON, err := json.Marshal(plan)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, planJSON)
}

// AddReader updates the owner field of plan with given id in world state.
// 添加新的用户到ReadSet中
func (s *SmartContract) AddReader(ctx contractapi.TransactionContextInterface, id string, newReader string) error {
	// 只有ApplicantMSP能修改飞行计划的ReadSet
	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return err
	}
	if mspID != "ApplicantMSP" {
		return fmt.Errorf("you have no right to add new reader to this plan")
	}

	plan, err := s.ReadPlan(ctx, id)
	if err != nil {
		return err
	}

	if _, res := InString(plan.ReadSet, newReader); res {
		return fmt.Errorf("the new reader is already in the readset")
	} else {
		plan.ReadSet = append(plan.ReadSet, newReader)
	}

	planJSON, err := json.Marshal(plan)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, planJSON)
}

// RemoveReader updates the owner field of plan with given id in world state.
// 将某用户从ReadSet中删除
func (s *SmartContract) RemoveReader(ctx contractapi.TransactionContextInterface, id string, reader string) error {
	// 只有ApplicantMSP能修改飞行计划的ReadSet
	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return err
	}
	if mspID != "ApplicantMSP" {
		return fmt.Errorf("you have no right to remove reader from this plan")
	}

	plan, err := s.ReadPlan(ctx, id)
	if err != nil {
		return err
	}

	if _, res := InString(plan.ReadSet, reader); !res {
		return fmt.Errorf("the reader is not in the readset")
	} else {
		newset, err := Remove(plan.ReadSet, reader)
		if err != nil {
			return err
		}
		plan.ReadSet = newset
	}

	planJSON, err := json.Marshal(plan)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, planJSON)
}

// GetAllPlans returns all plans found in world state
// 获取所有飞行计划
func (s *SmartContract) GetAllPlans(ctx contractapi.TransactionContextInterface) ([]*Plan, error) {
	// 如果是UserMSP不能使用该方法
	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, err
	}
	if mspID == "UserMSP" {
		return nil, fmt.Errorf("you have no right to read this plan")
	}

	// range query with empty string for startKey and endKey does an
	// open-ended query of all plans in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var plans []*Plan
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var plan Plan
		err = json.Unmarshal(queryResponse.Value, &plan)
		if err != nil {
			return nil, err
		}
		plans = append(plans, &plan)
	}

	return plans, nil
}

// GetSubmittingClientIdentity returns the name and issuer of the identity that
// invokes the smart contract. This function base64 decodes the identity string
// before returning the value to the client or smart contract.
// 获取用户身份
func (s *SmartContract) GetSubmittingClientIdentity(ctx contractapi.TransactionContextInterface) (string, error) {

	b64ID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("failed to read clientID: %v", err)
	}
	decodeID, err := base64.StdEncoding.DecodeString(b64ID)
	if err != nil {
		return "", fmt.Errorf("failed to base64 decode clientID: %v", err)
	}
	return string(decodeID), nil
}

// 判断字符串是否在字符数组中
func InString(slice []string, val string) (int, bool) {
	for index, item := range slice {
		if item == val {
			return index, true
		}
	}
	return -1, false
}

// 删除index位置的字符串
func Remove(slice []string, val string) ([]string, error) {
	if index, res := InString(slice, val); !res {
		return nil, fmt.Errorf("there is no string %v in string slice", val)
	} else {
		slice = append(slice[:index], slice[index+1:]...)
		return slice, nil
	}
}
