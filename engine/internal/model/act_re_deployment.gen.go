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

// ActReDeployment mapped from table <act_re_deployment>
type ActReDeployment struct {
	ID_                 string     `gorm:"column:ID_;type:varchar(64);primaryKey" json:"id_"`
	Name_               *string    `gorm:"column:NAME_;type:varchar(255)" json:"name_"`
	Category_           *string    `gorm:"column:CATEGORY_;type:varchar(255)" json:"category_"`
	Key_                *string    `gorm:"column:KEY_;type:varchar(255)" json:"key_"`
	TenantID_           *string    `gorm:"column:TENANT_ID_;type:varchar(255)" json:"tenant_id_"`
	DeployTime_         *time.Time `gorm:"column:DEPLOY_TIME_;type:timestamp without time zone" json:"deploy_time_"`
	DerivedFrom_        *string    `gorm:"column:DERIVED_FROM_;type:varchar(64)" json:"derived_from_"`
	DerivedFromRoot_    *string    `gorm:"column:DERIVED_FROM_ROOT_;type:varchar(64)" json:"derived_from_root_"`
	ParentDeploymentID_ *string    `gorm:"column:PARENT_DEPLOYMENT_ID_;type:varchar(255)" json:"parent_deployment_id_"`
	EngineVersion_      *string    `gorm:"column:ENGINE_VERSION_;type:varchar(255)" json:"engine_version_"`
	ProcessID_          string     `gorm:"column:process_id_;type:varchar(255)" json:"process_id_"`
}

// TableName ActReDeployment's table name
func (*ActReDeployment) TableName() string {
	var TableNameActReDeployment string

	if stringutils.IsNotEmpty(config.G_Config.Db.Name) {
		TableNameActReDeployment = fmt.Sprintf("%s.act_re_deployment", config.G_Config.Db.Name)
	} else {
		TableNameActReDeployment = "act_re_deployment"
	}

	return TableNameActReDeployment
}