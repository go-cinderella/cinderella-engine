package variable

import (
	"fmt"
	"github.com/go-cinderella/cinderella-engine/engine/config"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/unionj-cloud/toolkit/stringutils"
	"time"
)

type HistoricVariable struct {
	ID_              string     `gorm:"column:ID_;type:varchar(64);primaryKey" json:"id_"`
	Rev_             *int32     `gorm:"column:REV_;type:integer;default:1" json:"rev_"`
	ProcInstID_      *string    `gorm:"column:PROC_INST_ID_;type:varchar(64);index:act_idx_hi_procvar_proc_inst,priority:1" json:"proc_inst_id_"`
	ExecutionID_     *string    `gorm:"column:EXECUTION_ID_;type:varchar(64);index:act_idx_hi_procvar_exe,priority:1" json:"execution_id_"`
	TaskID_          *string    `gorm:"column:TASK_ID_;type:varchar(64);index:act_idx_hi_procvar_task_id,priority:1" json:"task_id_"`
	Name_            string     `gorm:"column:NAME_;type:varchar(255);not null;index:act_idx_hi_procvar_name_type,priority:1" json:"name_"`
	VarType_         *string    `gorm:"column:VAR_TYPE_;type:varchar(100);index:act_idx_hi_procvar_name_type,priority:2" json:"var_type_"`
	ScopeID_         *string    `gorm:"column:SCOPE_ID_;type:varchar(255);index:act_idx_hi_var_scope_id_type,priority:1" json:"scope_id_"`
	SubScopeID_      *string    `gorm:"column:SUB_SCOPE_ID_;type:varchar(255);index:act_idx_hi_var_sub_id_type,priority:1" json:"sub_scope_id_"`
	ScopeType_       *string    `gorm:"column:SCOPE_TYPE_;type:varchar(255);index:act_idx_hi_var_scope_id_type,priority:2;index:act_idx_hi_var_sub_id_type,priority:2" json:"scope_type_"`
	BytearrayID_     *string    `gorm:"column:BYTEARRAY_ID_;type:varchar(64)" json:"bytearray_id_"`
	Double_          *float64   `gorm:"column:DOUBLE_;type:double precision" json:"double_"`
	Long_            *int64     `gorm:"column:LONG_;type:bigint" json:"long_"`
	Text_            *string    `gorm:"column:TEXT_;type:varchar(4000)" json:"text_"`
	CreateTime_      *time.Time `gorm:"column:CREATE_TIME_;type:timestamp without time zone" json:"create_time_"`
	LastUpdatedTime_ *time.Time `gorm:"column:LAST_UPDATED_TIME_;type:timestamp without time zone" json:"last_updated_time_"`
}

func (variable HistoricVariable) GetName() string {
	return variable.Name_
}

func (variable HistoricVariable) GetProcessInstanceId() string {
	return cast.ToString(variable.ProcInstID_)
}

func (variable HistoricVariable) GetTaskId() string {
	return cast.ToString(variable.TaskID_)
}

func (variable HistoricVariable) GetNumberValue() int {
	return cast.ToInt(variable.Long_)
}

func (variable *HistoricVariable) SetNumberValue(value int) {
	variable.Long_ = lo.ToPtr(int64(value))
}

func (variable HistoricVariable) GetTextValue() string {
	return cast.ToString(variable.Text_)
}

func (variable *HistoricVariable) SetTextValue(value string) {
	variable.Text_ = &value
}

func (variable *HistoricVariable) SetValue(value interface{}, variableType VariableType) {
	variableType.SetValue(value, variable)
}

func (variable *HistoricVariable) TableName() string {
	var TableNameActHiVarinst string

	if stringutils.IsNotEmpty(config.G_Config.Db.Name) {
		TableNameActHiVarinst = fmt.Sprintf("%s.act_hi_varinst", config.G_Config.Db.Name)
	} else {
		TableNameActHiVarinst = "act_hi_varinst"
	}

	return TableNameActHiVarinst
}
