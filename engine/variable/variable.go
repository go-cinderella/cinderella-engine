package variable

import (
	"fmt"
	"github.com/go-cinderella/cinderella-engine/engine/config"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/unionj-cloud/toolkit/stringutils"
)

type Variable struct {
	ID_          string   `gorm:"column:id_;type:varchar(64);primaryKey" json:"id_"`
	Rev_         *int32   `gorm:"column:rev_;type:integer" json:"rev_"`
	Type_        string   `gorm:"column:type_;type:varchar(255);not null" json:"type_"`
	Name_        string   `gorm:"column:name_;type:varchar(255);not null;uniqueIndex:act_ru_variable__uniq,priority:1" json:"name_"`
	ExecutionID_ *string  `gorm:"column:execution_id_;type:varchar(64);index:act_idx_var_exe,priority:1" json:"execution_id_"`
	ProcInstID_  *string  `gorm:"column:proc_inst_id_;type:varchar(64);uniqueIndex:act_ru_variable__uniq,priority:2;index:act_idx_var_procinst,priority:1" json:"proc_inst_id_"`
	TaskID_      *string  `gorm:"column:task_id_;type:varchar(64);index:act_idx_variable_task_id,priority:1" json:"task_id_"`
	ScopeID_     *string  `gorm:"column:scope_id_;type:varchar(255);index:act_idx_ru_var_scope_id_type,priority:1" json:"scope_id_"`
	SubScopeID_  *string  `gorm:"column:sub_scope_id_;type:varchar(255);index:act_idx_ru_var_sub_id_type,priority:1" json:"sub_scope_id_"`
	ScopeType_   *string  `gorm:"column:scope_type_;type:varchar(255);index:act_idx_ru_var_scope_id_type,priority:2;index:act_idx_ru_var_sub_id_type,priority:2" json:"scope_type_"`
	BytearrayID_ *string  `gorm:"column:bytearray_id_;type:varchar(64);index:act_idx_var_bytearray,priority:1" json:"bytearray_id_"`
	Double_      *float64 `gorm:"column:double_;type:double precision" json:"double_"`
	Long_        *int64   `gorm:"column:long_;type:bigint" json:"long_"`
	Text_        *string  `gorm:"column:text_;type:varchar(4000)" json:"text_"`
}

func (variable Variable) GetName() string {
	return variable.Name_
}

func (variable Variable) GetProcessInstanceId() string {
	return cast.ToString(variable.ProcInstID_)
}

func (variable Variable) GetTaskId() string {
	return cast.ToString(variable.TaskID_)
}

func (variable Variable) GetNumberValue() int {
	return cast.ToInt(variable.Long_)
}

func (variable *Variable) SetNumberValue(value int) {
	variable.Long_ = lo.ToPtr(int64(value))
}

func (variable Variable) GetTextValue() string {
	return cast.ToString(variable.Text_)
}

func (variable *Variable) SetTextValue(value string) {
	variable.Text_ = &value
}

func (variable *Variable) SetValue(value interface{}, variableType VariableType) {
	variableType.SetValue(value, variable)
}

// TableName ActRuVariable's table name
func (variable *Variable) TableName() string {
	var TableNameActRuVariable string

	if stringutils.IsNotEmpty(config.G_Config.Db.Name) {
		TableNameActRuVariable = fmt.Sprintf("%s.act_ru_variable", config.G_Config.Db.Name)
	} else {
		TableNameActRuVariable = "act_ru_variable"
	}

	return TableNameActRuVariable
}

func (variable *Variable) GetValue() interface{} {
	variableTypes := variableManager.VariableTypes
	variableType := variableTypes.GetVariableType(variable.Type_)
	return variableType.GetValue(variable)
}