package task

import "github.com/go-cinderella/cinderella-engine/engine/internal/model"

type TaskDTO struct {
	model.ActRuTask
	GroupID *string `gorm:"column:group_id_;type:varchar(255);index:act_idx_ident_lnk_group,priority:1" json:"group_id_"`
	UserID  *string `gorm:"column:user_id_;type:varchar(255);index:act_idx_ident_lnk_user,priority:1" json:"user_id_"`
}
