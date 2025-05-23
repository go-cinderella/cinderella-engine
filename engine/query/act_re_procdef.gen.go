// Code generated by github.com/wubin1989/gen. DO NOT EDIT.
// Code generated by github.com/wubin1989/gen. DO NOT EDIT.
// Code generated by github.com/wubin1989/gen. DO NOT EDIT.

package query

import (
	"context"

	gormgen "github.com/wubin1989/gen"
	"github.com/wubin1989/gorm"
	"github.com/wubin1989/gorm/clause"
	"github.com/wubin1989/gorm/schema"

	"github.com/wubin1989/gen/field"

	"github.com/wubin1989/dbresolver"

	"github.com/go-cinderella/cinderella-engine/engine/model"
)

func newActReProcdef(db *gorm.DB, opts ...gormgen.DOOption) actReProcdef {
	_actReProcdef := actReProcdef{}

	_actReProcdef.actReProcdefDo.UseDB(db, opts...)
	_actReProcdef.actReProcdefDo.UseModel(&model.ActReProcdef{})

	tableName := _actReProcdef.actReProcdefDo.TableName()
	_actReProcdef.ALL = field.NewAsterisk(tableName)
	_actReProcdef.ID = field.NewString(tableName, "id_")
	_actReProcdef.Rev = field.NewInt32(tableName, "rev_")
	_actReProcdef.Category = field.NewString(tableName, "category_")
	_actReProcdef.Name = field.NewString(tableName, "name_")
	_actReProcdef.Key = field.NewString(tableName, "key_")
	_actReProcdef.Version = field.NewInt32(tableName, "version_")
	_actReProcdef.DeploymentID = field.NewString(tableName, "deployment_id_")
	_actReProcdef.ResourceName = field.NewString(tableName, "resource_name_")
	_actReProcdef.DgrmResourceName = field.NewString(tableName, "dgrm_resource_name_")
	_actReProcdef.Description = field.NewString(tableName, "description_")
	_actReProcdef.HasStartFormKey = field.NewBool(tableName, "has_start_form_key_")
	_actReProcdef.HasGraphicalNotation = field.NewBool(tableName, "has_graphical_notation_")
	_actReProcdef.SuspensionState = field.NewInt32(tableName, "suspension_state_")
	_actReProcdef.TenantID = field.NewString(tableName, "tenant_id_")
	_actReProcdef.DerivedFrom = field.NewString(tableName, "derived_from_")
	_actReProcdef.DerivedFromRoot = field.NewString(tableName, "derived_from_root_")
	_actReProcdef.DerivedVersion = field.NewInt32(tableName, "derived_version_")
	_actReProcdef.EngineVersion = field.NewString(tableName, "engine_version_")
	_actReProcdef.ProcessID = field.NewString(tableName, "process_id_")
	_actReProcdef.CreatedBy = field.NewString(tableName, "created_by_")
	_actReProcdef.CreatedByName = field.NewString(tableName, "created_by_name_")

	_actReProcdef.fillFieldMap()

	return _actReProcdef
}

type actReProcdef struct {
	actReProcdefDo

	ALL                  field.Asterisk
	ID                   field.String
	Rev                  field.Int32
	Category             field.String
	Name                 field.String
	Key                  field.String
	Version              field.Int32
	DeploymentID         field.String
	ResourceName         field.String
	DgrmResourceName     field.String
	Description          field.String
	HasStartFormKey      field.Bool
	HasGraphicalNotation field.Bool
	SuspensionState      field.Int32
	TenantID             field.String
	DerivedFrom          field.String
	DerivedFromRoot      field.String
	DerivedVersion       field.Int32
	EngineVersion        field.String
	ProcessID            field.String
	CreatedBy            field.String
	CreatedByName        field.String

	fieldMap map[string]field.Expr
}

func (a actReProcdef) Table(newTableName string) *actReProcdef {
	a.actReProcdefDo.UseTable(newTableName)
	return a.updateTableName(newTableName)
}

func (a actReProcdef) As(alias string) *actReProcdef {
	a.actReProcdefDo.DO = *(a.actReProcdefDo.As(alias).(*gormgen.DO))
	return a.updateTableName(alias)
}

func (a *actReProcdef) updateTableName(table string) *actReProcdef {
	a.ALL = field.NewAsterisk(table)
	a.ID = field.NewString(table, "id_")
	a.Rev = field.NewInt32(table, "rev_")
	a.Category = field.NewString(table, "category_")
	a.Name = field.NewString(table, "name_")
	a.Key = field.NewString(table, "key_")
	a.Version = field.NewInt32(table, "version_")
	a.DeploymentID = field.NewString(table, "deployment_id_")
	a.ResourceName = field.NewString(table, "resource_name_")
	a.DgrmResourceName = field.NewString(table, "dgrm_resource_name_")
	a.Description = field.NewString(table, "description_")
	a.HasStartFormKey = field.NewBool(table, "has_start_form_key_")
	a.HasGraphicalNotation = field.NewBool(table, "has_graphical_notation_")
	a.SuspensionState = field.NewInt32(table, "suspension_state_")
	a.TenantID = field.NewString(table, "tenant_id_")
	a.DerivedFrom = field.NewString(table, "derived_from_")
	a.DerivedFromRoot = field.NewString(table, "derived_from_root_")
	a.DerivedVersion = field.NewInt32(table, "derived_version_")
	a.EngineVersion = field.NewString(table, "engine_version_")
	a.ProcessID = field.NewString(table, "process_id_")
	a.CreatedBy = field.NewString(table, "created_by_")
	a.CreatedByName = field.NewString(table, "created_by_name_")

	a.fillFieldMap()

	return a
}

func (a *actReProcdef) WithContext(ctx context.Context) IActReProcdefDo {
	return a.actReProcdefDo.WithContext(ctx)
}

func (a actReProcdef) TableName() string { return a.actReProcdefDo.TableName() }

func (a actReProcdef) Alias() string { return a.actReProcdefDo.Alias() }

func (a actReProcdef) Columns(cols ...field.Expr) gormgen.Columns {
	return a.actReProcdefDo.Columns(cols...)
}

func (a *actReProcdef) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := a.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (a *actReProcdef) GetFieldExprByName(fieldName string) (field.Expr, bool) {
	_f, ok := a.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	return _f, ok
}

func (a *actReProcdef) fillFieldMap() {
	a.fieldMap = make(map[string]field.Expr, 18)
	a.fieldMap["id_"] = a.ID
	a.fieldMap["rev_"] = a.Rev
	a.fieldMap["category_"] = a.Category
	a.fieldMap["name_"] = a.Name
	a.fieldMap["key_"] = a.Key
	a.fieldMap["version_"] = a.Version
	a.fieldMap["deployment_id_"] = a.DeploymentID
	a.fieldMap["resource_name_"] = a.ResourceName
	a.fieldMap["dgrm_resource_name_"] = a.DgrmResourceName
	a.fieldMap["description_"] = a.Description
	a.fieldMap["has_start_form_key_"] = a.HasStartFormKey
	a.fieldMap["has_graphical_notation_"] = a.HasGraphicalNotation
	a.fieldMap["suspension_state_"] = a.SuspensionState
	a.fieldMap["tenant_id_"] = a.TenantID
	a.fieldMap["derived_from_"] = a.DerivedFrom
	a.fieldMap["derived_from_root_"] = a.DerivedFromRoot
	a.fieldMap["derived_version_"] = a.DerivedVersion
	a.fieldMap["engine_version_"] = a.EngineVersion
	a.fieldMap["process_id_"] = a.ProcessID
	a.fieldMap["created_by_"] = a.CreatedBy
	a.fieldMap["created_by_name_"] = a.CreatedByName
}

func (a actReProcdef) Clone(db *gorm.DB) actReProcdef {
	a.actReProcdefDo.ReplaceConnPool(db.Statement.ConnPool)
	return a
}

func (a actReProcdef) ReplaceDB(db *gorm.DB) actReProcdef {
	a.actReProcdefDo.ReplaceDB(db)
	return a
}

type actReProcdefDo struct{ gormgen.DO }

type IActReProcdefDo interface {
	gormgen.SubQuery
	Debug() IActReProcdefDo
	WithContext(ctx context.Context) IActReProcdefDo
	WithResult(fc func(tx gormgen.Dao)) gormgen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IActReProcdefDo
	WriteDB() IActReProcdefDo
	As(alias string) gormgen.Dao
	Session(config *gorm.Session) IActReProcdefDo
	Columns(cols ...field.Expr) gormgen.Columns
	Clauses(conds ...clause.Expression) IActReProcdefDo
	Not(conds ...gormgen.Condition) IActReProcdefDo
	Or(conds ...gormgen.Condition) IActReProcdefDo
	Select(conds ...field.Expr) IActReProcdefDo
	Where(conds ...gormgen.Condition) IActReProcdefDo
	Order(conds ...field.Expr) IActReProcdefDo
	Distinct(cols ...field.Expr) IActReProcdefDo
	Omit(cols ...field.Expr) IActReProcdefDo
	Join(table schema.Tabler, on ...field.Expr) IActReProcdefDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IActReProcdefDo
	RightJoin(table schema.Tabler, on ...field.Expr) IActReProcdefDo
	Group(cols ...field.Expr) IActReProcdefDo
	Having(conds ...gormgen.Condition) IActReProcdefDo
	Limit(limit int) IActReProcdefDo
	Offset(offset int) IActReProcdefDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gormgen.Dao) gormgen.Dao) IActReProcdefDo
	Unscoped() IActReProcdefDo
	Create(values ...*model.ActReProcdef) error
	CreateInBatches(values []*model.ActReProcdef, batchSize int) error
	Save(values ...*model.ActReProcdef) error
	First() (*model.ActReProcdef, error)
	Take() (*model.ActReProcdef, error)
	Last() (*model.ActReProcdef, error)
	Find() ([]*model.ActReProcdef, error)
	FindInBatch(batchSize int, fc func(tx gormgen.Dao, batch int) error) (results []*model.ActReProcdef, err error)
	FindInBatches(result *[]*model.ActReProcdef, batchSize int, fc func(tx gormgen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.ActReProcdef) (info gormgen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gormgen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gormgen.ResultInfo, err error)
	Updates(value interface{}) (info gormgen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gormgen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gormgen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gormgen.ResultInfo, err error)
	UpdateFrom(q gormgen.SubQuery) gormgen.Dao
	Attrs(attrs ...field.AssignExpr) IActReProcdefDo
	Assign(attrs ...field.AssignExpr) IActReProcdefDo
	Joins(fields ...field.RelationField) IActReProcdefDo
	Preload(fields ...field.RelationField) IActReProcdefDo
	FirstOrInit() (*model.ActReProcdef, error)
	FirstOrCreate() (*model.ActReProcdef, error)
	FindByPage(offset int, limit int) (result []*model.ActReProcdef, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Fetch(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IActReProcdefDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (a actReProcdefDo) Debug() IActReProcdefDo {
	return a.withDO(a.DO.Debug())
}

func (a actReProcdefDo) WithContext(ctx context.Context) IActReProcdefDo {
	return a.withDO(a.DO.WithContext(ctx))
}

func (a actReProcdefDo) ReadDB() IActReProcdefDo {
	return a.Clauses(dbresolver.Read)
}

func (a actReProcdefDo) WriteDB() IActReProcdefDo {
	return a.Clauses(dbresolver.Write)
}

func (a actReProcdefDo) Session(config *gorm.Session) IActReProcdefDo {
	return a.withDO(a.DO.Session(config))
}

func (a actReProcdefDo) Clauses(conds ...clause.Expression) IActReProcdefDo {
	return a.withDO(a.DO.Clauses(conds...))
}

func (a actReProcdefDo) Returning(value interface{}, columns ...string) IActReProcdefDo {
	return a.withDO(a.DO.Returning(value, columns...))
}

func (a actReProcdefDo) Not(conds ...gormgen.Condition) IActReProcdefDo {
	return a.withDO(a.DO.Not(conds...))
}

func (a actReProcdefDo) Or(conds ...gormgen.Condition) IActReProcdefDo {
	return a.withDO(a.DO.Or(conds...))
}

func (a actReProcdefDo) Select(conds ...field.Expr) IActReProcdefDo {
	return a.withDO(a.DO.Select(conds...))
}

func (a actReProcdefDo) Where(conds ...gormgen.Condition) IActReProcdefDo {
	return a.withDO(a.DO.Where(conds...))
}

func (a actReProcdefDo) Order(conds ...field.Expr) IActReProcdefDo {
	return a.withDO(a.DO.Order(conds...))
}

func (a actReProcdefDo) Distinct(cols ...field.Expr) IActReProcdefDo {
	return a.withDO(a.DO.Distinct(cols...))
}

func (a actReProcdefDo) Omit(cols ...field.Expr) IActReProcdefDo {
	return a.withDO(a.DO.Omit(cols...))
}

func (a actReProcdefDo) Join(table schema.Tabler, on ...field.Expr) IActReProcdefDo {
	return a.withDO(a.DO.Join(table, on...))
}

func (a actReProcdefDo) LeftJoin(table schema.Tabler, on ...field.Expr) IActReProcdefDo {
	return a.withDO(a.DO.LeftJoin(table, on...))
}

func (a actReProcdefDo) RightJoin(table schema.Tabler, on ...field.Expr) IActReProcdefDo {
	return a.withDO(a.DO.RightJoin(table, on...))
}

func (a actReProcdefDo) Group(cols ...field.Expr) IActReProcdefDo {
	return a.withDO(a.DO.Group(cols...))
}

func (a actReProcdefDo) Having(conds ...gormgen.Condition) IActReProcdefDo {
	return a.withDO(a.DO.Having(conds...))
}

func (a actReProcdefDo) Limit(limit int) IActReProcdefDo {
	return a.withDO(a.DO.Limit(limit))
}

func (a actReProcdefDo) Offset(offset int) IActReProcdefDo {
	return a.withDO(a.DO.Offset(offset))
}

func (a actReProcdefDo) Scopes(funcs ...func(gormgen.Dao) gormgen.Dao) IActReProcdefDo {
	return a.withDO(a.DO.Scopes(funcs...))
}

func (a actReProcdefDo) Unscoped() IActReProcdefDo {
	return a.withDO(a.DO.Unscoped())
}

func (a actReProcdefDo) Create(values ...*model.ActReProcdef) error {
	if len(values) == 0 {
		return nil
	}
	return a.DO.Create(values)
}

func (a actReProcdefDo) CreateInBatches(values []*model.ActReProcdef, batchSize int) error {
	return a.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (a actReProcdefDo) Save(values ...*model.ActReProcdef) error {
	if len(values) == 0 {
		return nil
	}
	return a.DO.Save(values)
}

func (a actReProcdefDo) First() (*model.ActReProcdef, error) {
	if result, err := a.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.ActReProcdef), nil
	}
}

func (a actReProcdefDo) Take() (*model.ActReProcdef, error) {
	if result, err := a.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.ActReProcdef), nil
	}
}

func (a actReProcdefDo) Last() (*model.ActReProcdef, error) {
	if result, err := a.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.ActReProcdef), nil
	}
}

func (a actReProcdefDo) Find() ([]*model.ActReProcdef, error) {
	result, err := a.DO.Find()
	return result.([]*model.ActReProcdef), err
}

func (a actReProcdefDo) FindInBatch(batchSize int, fc func(tx gormgen.Dao, batch int) error) (results []*model.ActReProcdef, err error) {
	buf := make([]*model.ActReProcdef, 0, batchSize)
	err = a.DO.FindInBatches(&buf, batchSize, func(tx gormgen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (a actReProcdefDo) FindInBatches(result *[]*model.ActReProcdef, batchSize int, fc func(tx gormgen.Dao, batch int) error) error {
	return a.DO.FindInBatches(result, batchSize, fc)
}

func (a actReProcdefDo) Attrs(attrs ...field.AssignExpr) IActReProcdefDo {
	return a.withDO(a.DO.Attrs(attrs...))
}

func (a actReProcdefDo) Assign(attrs ...field.AssignExpr) IActReProcdefDo {
	return a.withDO(a.DO.Assign(attrs...))
}

func (a actReProcdefDo) Joins(fields ...field.RelationField) IActReProcdefDo {
	for _, _f := range fields {
		a = *a.withDO(a.DO.Joins(_f))
	}
	return &a
}

func (a actReProcdefDo) Preload(fields ...field.RelationField) IActReProcdefDo {
	for _, _f := range fields {
		a = *a.withDO(a.DO.Preload(_f))
	}
	return &a
}

func (a actReProcdefDo) FirstOrInit() (*model.ActReProcdef, error) {
	if result, err := a.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.ActReProcdef), nil
	}
}

func (a actReProcdefDo) FirstOrCreate() (*model.ActReProcdef, error) {
	if result, err := a.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.ActReProcdef), nil
	}
}

func (a actReProcdefDo) FindByPage(offset int, limit int) (result []*model.ActReProcdef, count int64, err error) {
	result, err = a.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = a.Offset(-1).Limit(-1).Count()
	return
}

func (a actReProcdefDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = a.Count()
	if err != nil {
		return
	}

	err = a.Offset(offset).Limit(limit).Scan(result)
	return
}

func (a actReProcdefDo) Scan(result interface{}) (err error) {
	return a.DO.Scan(result)
}

func (a actReProcdefDo) Fetch(result interface{}) (err error) {
	return a.DO.Fetch(result)
}

func (a actReProcdefDo) Delete(models ...*model.ActReProcdef) (result gormgen.ResultInfo, err error) {
	return a.DO.Delete(models)
}

func (a *actReProcdefDo) withDO(do gormgen.Dao) *actReProcdefDo {
	a.DO = *do.(*gormgen.DO)
	return a
}
