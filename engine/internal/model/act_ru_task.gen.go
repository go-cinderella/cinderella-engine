// Code generated by github.com/wubin1989/gen. DO NOT EDIT.
// Code generated by github.com/wubin1989/gen. DO NOT EDIT.
// Code generated by github.com/wubin1989/gen. DO NOT EDIT.

package model

import (
	"fmt"

	"github.com/go-cinderella/cinderella-engine/engine/config"
	"github.com/unionj-cloud/toolkit/stringutils"

	"time"
)

//go:generate go-doudou name --file $GOFILE --case snake

// ActRuTask mapped from table <act_ru_task>
type ActRuTask struct {
	ID_                    string     `gorm:"column:ID_;type:varchar(64);primaryKey" json:"id_"`
	Rev_                   *int32     `gorm:"column:REV_;type:integer" json:"rev_"`
	ExecutionID_           *string    `gorm:"column:EXECUTION_ID_;type:varchar(64);index:act_idx_task_exec,priority:1" json:"execution_id_"`
	ProcInstID_            *string    `gorm:"column:PROC_INST_ID_;type:varchar(64);index:act_idx_task_procinst,priority:1" json:"proc_inst_id_"`
	ProcDefID_             *string    `gorm:"column:PROC_DEF_ID_;type:varchar(64);index:act_idx_task_procdef,priority:1" json:"proc_def_id_"`
	TaskDefID_             *string    `gorm:"column:TASK_DEF_ID_;type:varchar(64)" json:"task_def_id_"`
	ScopeID_               *string    `gorm:"column:scope_id_;type:varchar(255);index:act_idx_task_scope,priority:1" json:"scope_id_"`
	SubScopeID_            *string    `gorm:"column:sub_scope_id_;type:varchar(255);index:act_idx_task_sub_scope,priority:1" json:"sub_scope_id_"`
	ScopeType_             *string    `gorm:"column:scope_type_;type:varchar(255);index:act_idx_task_sub_scope,priority:2;index:act_idx_task_scope,priority:2;index:act_idx_task_scope_def,priority:1" json:"scope_type_"`
	ScopeDefinitionID_     *string    `gorm:"column:scope_definition_id_;type:varchar(255);index:act_idx_task_scope_def,priority:2" json:"scope_definition_id_"`
	PropagatedStageInstID_ *string    `gorm:"column:propagated_stage_inst_id_;type:varchar(255)" json:"propagated_stage_inst_id_"`
	Name_                  *string    `gorm:"column:NAME_;type:varchar(255)" json:"name_"`
	ParentTaskID_          *string    `gorm:"column:PARENT_TASK_ID_;type:varchar(64)" json:"parent_task_id_"`
	Description_           *string    `gorm:"column:DESCRIPTION_;type:varchar(4000)" json:"description_"`
	TaskDefKey_            *string    `gorm:"column:TASK_DEF_KEY_;type:varchar(255)" json:"task_def_key_"`
	Owner_                 *string    `gorm:"column:OWNER_;type:varchar(255)" json:"owner_"`
	Assignee_              *string    `gorm:"column:ASSIGNEE_;type:varchar(255)" json:"assignee_"`
	Delegation_            *string    `gorm:"column:DELEGATION_;type:varchar(64)" json:"delegation_"`
	Priority_              *int32     `gorm:"column:PRIORITY_;type:integer" json:"priority_"`
	CreateTime_            *time.Time `gorm:"column:CREATE_TIME_;type:timestamp without time zone;index:act_idx_task_create,priority:1" json:"create_time_"`
	DueDate_               *time.Time `gorm:"column:DUE_DATE_;type:timestamp without time zone" json:"due_date_"`
	Category_              *string    `gorm:"column:CATEGORY_;type:varchar(255)" json:"category_"`
	SuspensionState_       *int32     `gorm:"column:suspension_state_;type:integer" json:"suspension_state_"`
	TenantID_              *string    `gorm:"column:TENANT_ID_;type:varchar(255)" json:"tenant_id_"`
	FormKey_               *string    `gorm:"column:FORM_KEY_;type:varchar(255)" json:"form_key_"`
	ClaimTime_             *time.Time `gorm:"column:CLAIM_TIME_;type:timestamp without time zone" json:"claim_time_"`
	IsCountEnabled_        *bool      `gorm:"column:is_count_enabled_;type:boolean" json:"is_count_enabled_"`
	VarCount_              *int32     `gorm:"column:var_count_;type:integer" json:"var_count_"`
	IDLinkCount_           *int32     `gorm:"column:id_link_count_;type:integer" json:"id_link_count_"`
	SubTaskCount_          *int32     `gorm:"column:sub_task_count_;type:integer" json:"sub_task_count_"`
}

// TableName ActRuTask's table name
func (*ActRuTask) TableName() string {
	var TableNameActRuTask string

	if stringutils.IsNotEmpty(config.G_Config.Db.Name) {
		TableNameActRuTask = fmt.Sprintf("%s.act_ru_task", config.G_Config.Db.Name)
	} else {
		TableNameActRuTask = "act_ru_task"
	}

	return TableNameActRuTask
}
