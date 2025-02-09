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

// ActRuExecution mapped from table <act_ru_execution>
type ActRuExecution struct {
	ID_                     string     `gorm:"column:id_;type:varchar(64);primaryKey" json:"id_"`
	Rev_                    *int32     `gorm:"column:rev_;type:integer" json:"rev_"`
	ProcInstID_             *string    `gorm:"column:proc_inst_id_;type:varchar(64);index:act_idx_exe_procinst,priority:1" json:"proc_inst_id_"`
	BusinessKey_            *string    `gorm:"column:business_key_;type:varchar(255);index:act_idx_exec_buskey,priority:1" json:"business_key_"`
	ParentID_               *string    `gorm:"column:parent_id_;type:varchar(64);index:act_idx_exe_parent,priority:1" json:"parent_id_"`
	ProcDefID_              *string    `gorm:"column:proc_def_id_;type:varchar(64);index:act_idx_exe_procdef,priority:1" json:"proc_def_id_"`
	SuperExec_              *string    `gorm:"column:super_exec_;type:varchar(64);index:act_idx_exe_super,priority:1" json:"super_exec_"`
	RootProcInstID_         *string    `gorm:"column:root_proc_inst_id_;type:varchar(64);index:act_idx_exe_root,priority:1" json:"root_proc_inst_id_"`
	ActID_                  *string    `gorm:"column:act_id_;type:varchar(255)" json:"act_id_"`
	IsActive_               *bool      `gorm:"column:is_active_;type:boolean" json:"is_active_"`
	IsConcurrent_           *bool      `gorm:"column:is_concurrent_;type:boolean" json:"is_concurrent_"`
	IsScope_                *bool      `gorm:"column:is_scope_;type:boolean" json:"is_scope_"`
	IsEventScope_           *bool      `gorm:"column:is_event_scope_;type:boolean" json:"is_event_scope_"`
	IsMiRoot_               *bool      `gorm:"column:is_mi_root_;type:boolean" json:"is_mi_root_"`
	SuspensionState_        *int32     `gorm:"column:suspension_state_;type:integer" json:"suspension_state_"`
	CachedEntState_         *int32     `gorm:"column:cached_ent_state_;type:integer" json:"cached_ent_state_"`
	TenantID_               *string    `gorm:"column:tenant_id_;type:varchar(255)" json:"tenant_id_"`
	Name_                   *string    `gorm:"column:name_;type:varchar(255)" json:"name_"`
	StartActID_             *string    `gorm:"column:start_act_id_;type:varchar(255)" json:"start_act_id_"`
	StartTime_              *time.Time `gorm:"column:start_time_;type:timestamp without time zone" json:"start_time_"`
	StartUserID_            *string    `gorm:"column:start_user_id_;type:varchar(255)" json:"start_user_id_"`
	LockTime_               *time.Time `gorm:"column:lock_time_;type:timestamp without time zone" json:"lock_time_"`
	LockOwner_              *string    `gorm:"column:lock_owner_;type:varchar(255)" json:"lock_owner_"`
	IsCountEnabled_         *bool      `gorm:"column:is_count_enabled_;type:boolean" json:"is_count_enabled_"`
	EvtSubscrCount_         *int32     `gorm:"column:evt_subscr_count_;type:integer" json:"evt_subscr_count_"`
	TaskCount_              *int32     `gorm:"column:task_count_;type:integer" json:"task_count_"`
	JobCount_               *int32     `gorm:"column:job_count_;type:integer" json:"job_count_"`
	TimerJobCount_          *int32     `gorm:"column:timer_job_count_;type:integer" json:"timer_job_count_"`
	SuspJobCount_           *int32     `gorm:"column:susp_job_count_;type:integer" json:"susp_job_count_"`
	DeadletterJobCount_     *int32     `gorm:"column:deadletter_job_count_;type:integer" json:"deadletter_job_count_"`
	ExternalWorkerJobCount_ *int32     `gorm:"column:external_worker_job_count_;type:integer" json:"external_worker_job_count_"`
	VarCount_               *int32     `gorm:"column:var_count_;type:integer" json:"var_count_"`
	IDLinkCount_            *int32     `gorm:"column:id_link_count_;type:integer" json:"id_link_count_"`
	CallbackID_             *string    `gorm:"column:callback_id_;type:varchar(255)" json:"callback_id_"`
	CallbackType_           *string    `gorm:"column:callback_type_;type:varchar(255)" json:"callback_type_"`
	ReferenceID_            *string    `gorm:"column:reference_id_;type:varchar(255);index:act_idx_exec_ref_id_,priority:1" json:"reference_id_"`
	ReferenceType_          *string    `gorm:"column:reference_type_;type:varchar(255)" json:"reference_type_"`
	PropagatedStageInstID_  *string    `gorm:"column:propagated_stage_inst_id_;type:varchar(255)" json:"propagated_stage_inst_id_"`
	BusinessStatus_         *string    `gorm:"column:business_status_;type:varchar(255)" json:"business_status_"`
}

// TableName ActRuExecution's table name
func (*ActRuExecution) TableName() string {
	var TableNameActRuExecution string

	if stringutils.IsNotEmpty(config.G_Config.Db.Name) {
		TableNameActRuExecution = fmt.Sprintf("%s.act_ru_execution", config.G_Config.Db.Name)
	} else {
		TableNameActRuExecution = "act_ru_execution"
	}

	return TableNameActRuExecution
}
