package xdb

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"gorm.io/gorm"
)

const (
	defaultTableName = "casbin_rule"
)

type customTableKey struct{}

type CasbinRule struct {
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	RuleType string `json:"rule_type" gorm:"size:100"`
	V0       string `json:"v0" gorm:"size:100"`
	V1       string `json:"v1" gorm:"size:100"`
	V2       string `json:"v2" gorm:"size:100"`
	V3       string `json:"v3" gorm:"size:100"`
	V4       string `json:"v4" gorm:"size:100"`
	V5       string `json:"v5" gorm:"size:100"`
}

func (CasbinRule) TableName() string {
	return "casbin_rule"
}

type Filter struct {
	RuleType []string
	V0       []string
	V1       []string
	V2       []string
	V3       []string
	V4       []string
	V5       []string
}

// Adapter represents the Gorm adapter for policy storage.
type Adapter struct {
	driverName     string
	dataSourceName string
	databaseName   string
	tablePrefix    string
	tableName      string
	dbSpecified    bool
	db             *gorm.DB
	isFiltered     bool
}

// finalizer is the destructor for Adapter.
func finalizer(a *Adapter) {
	sqlDB, err := a.db.DB()
	if err != nil {
		panic(err)
	}
	err = sqlDB.Close()
	if err != nil {
		panic(err)
	}
}

func NewAdapter(db *gorm.DB, prefix string, tableName string) (*Adapter, error) {
	if len(tableName) == 0 {
		tableName = defaultTableName
	}
	a := &Adapter{
		tablePrefix: prefix,
		tableName:   tableName,
	}
	a.db = db.Scopes(a.casbinRuleTable()).Session(&gorm.Session{Context: db.Statement.Context})
	err := a.createTable()
	if err != nil {
		return nil, err
	}
	return a, nil
}

// NewAdapterByDB creates gorm-adapter by an existing Gorm instance
func NewAdapterByDB(db *gorm.DB) (*Adapter, error) {
	return NewAdapter(db, "", defaultTableName)
}

func NewAdapterByDBWithCustomTable(db *gorm.DB, t interface{}) (*Adapter, error) {
	ctx := db.Statement.Context
	if ctx == nil {
		ctx = context.Background()
	}
	ctx = context.WithValue(ctx, customTableKey{}, t)
	return NewAdapter(db.WithContext(ctx), "", defaultTableName)
}

func (a *Adapter) close() error {
	a.db = nil
	return nil
}

// getTableInstance return the dynamic table name
func (a *Adapter) getTableInstance() *CasbinRule {
	return &CasbinRule{}
}

func (a *Adapter) getFullTableName() string {
	if a.tablePrefix != "" {
		return a.tablePrefix + "_" + a.tableName
	}
	return a.tableName
}

func (a *Adapter) casbinRuleTable() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		tableName := a.getFullTableName()
		return db.Table(tableName)
	}
}

func (a *Adapter) createTable() error {
	t := a.db.Statement.Context.Value(customTableKey{})

	if t == nil {
		t = a.getTableInstance()
	}

	if err := a.db.AutoMigrate(t); err != nil {
		return err
	}
	tableName := a.getFullTableName()
	index := "idx_" + tableName
	hasIndex := a.db.Migrator().HasIndex(t, index)
	if !hasIndex {
		if err := a.db.Exec(fmt.Sprintf("CREATE UNIQUE INDEX %s ON %s (rule_type,v0,v1,v2,v3,v4,v5)", index, tableName)).Error; err != nil {
			return err
		}
	}
	return nil
}

func (a *Adapter) dropTable() error {
	t := a.db.Statement.Context.Value(customTableKey{})
	if t == nil {
		return a.db.Migrator().DropTable(a.getTableInstance())
	}

	return a.db.Migrator().DropTable(t)
}

func loadPolicyLine(line CasbinRule, model model.Model) {
	var p = []string{line.RuleType,
		line.V0, line.V1, line.V2, line.V3, line.V4, line.V5}

	var lineText string
	if line.V5 != "" {
		lineText = strings.Join(p, ", ")
	} else if line.V4 != "" {
		lineText = strings.Join(p[:6], ", ")
	} else if line.V3 != "" {
		lineText = strings.Join(p[:5], ", ")
	} else if line.V2 != "" {
		lineText = strings.Join(p[:4], ", ")
	} else if line.V1 != "" {
		lineText = strings.Join(p[:3], ", ")
	} else if line.V0 != "" {
		lineText = strings.Join(p[:2], ", ")
	}

	persist.LoadPolicyLine(lineText, model)
}

// LoadPolicy loads policy from database.
func (a *Adapter) LoadPolicy(model model.Model) error {
	var lines []CasbinRule
	if err := a.db.Order("ID").Find(&lines).Error; err != nil {
		return err
	}

	for _, line := range lines {
		loadPolicyLine(line, model)
	}

	return nil
}

// LoadFilteredPolicy loads only policy rules that match the filter.
func (a *Adapter) LoadFilteredPolicy(model model.Model, filter interface{}) error {
	var lines []CasbinRule

	filterValue, ok := filter.(Filter)
	if !ok {
		return errors.New("invalid filter type")
	}

	if err := a.db.Scopes(a.filterQuery(a.db, filterValue)).Order("ID").Find(&lines).Error; err != nil {
		return err
	}

	for _, line := range lines {
		loadPolicyLine(line, model)
	}
	a.isFiltered = true

	return nil
}

// IsFiltered returns true if the loaded policy has been filtered.
func (a *Adapter) IsFiltered() bool {
	return a.isFiltered
}

// filterQuery builds the gorm query to match the rule filter to use within a scope.
func (a *Adapter) filterQuery(db *gorm.DB, filter Filter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(filter.RuleType) > 0 {
			db = db.Where("rule_type in (?)", filter.RuleType)
		}
		if len(filter.V0) > 0 {
			db = db.Where("v0 in (?)", filter.V0)
		}
		if len(filter.V1) > 0 {
			db = db.Where("v1 in (?)", filter.V1)
		}
		if len(filter.V2) > 0 {
			db = db.Where("v2 in (?)", filter.V2)
		}
		if len(filter.V3) > 0 {
			db = db.Where("v3 in (?)", filter.V3)
		}
		if len(filter.V4) > 0 {
			db = db.Where("v4 in (?)", filter.V4)
		}
		if len(filter.V5) > 0 {
			db = db.Where("v5 in (?)", filter.V5)
		}
		return db
	}
}

func (a *Adapter) savePolicyLine(ruleType string, rule []string) CasbinRule {
	line := a.getTableInstance()
	line.RuleType = ruleType
	if len(rule) > 0 {
		line.V0 = rule[0]
	}
	if len(rule) > 1 {
		line.V1 = rule[1]
	}
	if len(rule) > 2 {
		line.V2 = rule[2]
	}
	if len(rule) > 3 {
		line.V3 = rule[3]
	}
	if len(rule) > 4 {
		line.V4 = rule[4]
	}
	if len(rule) > 5 {
		line.V5 = rule[5]
	}
	return *line
}

// SavePolicy saves policy to database.
func (a *Adapter) SavePolicy(model model.Model) (err error) {
	if err = a.dropTable(); err != nil {
		return
	}
	if err = a.createTable(); err != nil {
		return
	}
	for ruleType, ast := range model["p"] {
		for _, rule := range ast.Policy {
			line := a.savePolicyLine(ruleType, rule)
			if err = a.db.Create(&line).Error; err != nil {
				return err
			}
		}
	}
	for ruleType, ast := range model["g"] {
		for _, rule := range ast.Policy {
			line := a.savePolicyLine(ruleType, rule)
			if err = a.db.Create(&line).Error; err != nil {
				return err
			}
		}
	}
	return
}

// AddPolicy adds a policy rule to the storage.
func (a *Adapter) AddPolicy(ses string, ruleType string, rule []string) error {
	line := a.savePolicyLine(ruleType, rule)
	err := a.db.Create(&line).Error
	return err
}

// RemovePolicy removes a policy rule from the storage.
func (a *Adapter) RemovePolicy(sec string, ruleType string, rule []string) error {
	line := a.savePolicyLine(ruleType, rule)
	err := a.rawDelete(a.db, line) //can't use db.Delete as we're not using primary key http://jinzhu.me/gorm/crud.html#delete
	return err
}

// AddPolicies adds multiple policy rules to the storage.
func (a *Adapter) AddPolicies(ruleType string, rules [][]string) error {
	return a.db.Transaction(func(tx *gorm.DB) error {
		for _, rule := range rules {
			line := a.savePolicyLine(ruleType, rule)
			if err := tx.Create(&line).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// RemovePolicies removes multiple policy rules from the storage.
func (a *Adapter) RemovePolicies(ruleType string, rules [][]string) error {
	return a.db.Transaction(func(tx *gorm.DB) error {
		for _, rule := range rules {
			line := a.savePolicyLine(ruleType, rule)
			if err := a.rawDelete(tx, line); err != nil { //can't use db.Delete as we're not using primary key http://jinzhu.me/gorm/crud.html#delete
				return err
			}
		}
		return nil
	})
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
func (a *Adapter) RemoveFilteredPolicy(sec string, ruleType string, fieldIndex int, fieldValues ...string) error {
	line := a.getTableInstance()
	line.RuleType = ruleType
	if fieldIndex <= 0 && 0 < fieldIndex+len(fieldValues) {
		line.V0 = fieldValues[0-fieldIndex]
	}
	if fieldIndex <= 1 && 1 < fieldIndex+len(fieldValues) {
		line.V1 = fieldValues[1-fieldIndex]
	}
	if fieldIndex <= 2 && 2 < fieldIndex+len(fieldValues) {
		line.V2 = fieldValues[2-fieldIndex]
	}
	if fieldIndex <= 3 && 3 < fieldIndex+len(fieldValues) {
		line.V3 = fieldValues[3-fieldIndex]
	}
	if fieldIndex <= 4 && 4 < fieldIndex+len(fieldValues) {
		line.V4 = fieldValues[4-fieldIndex]
	}
	if fieldIndex <= 5 && 5 < fieldIndex+len(fieldValues) {
		line.V5 = fieldValues[5-fieldIndex]
	}
	err := a.rawDelete(a.db, *line)
	return err
}

func (a *Adapter) rawDelete(db *gorm.DB, line CasbinRule) error {
	queryArgs := []interface{}{line.RuleType}
	queryStr := "rule_type = ?"
	if line.V0 != "" {
		queryStr += " and v0 = ?"
		queryArgs = append(queryArgs, line.V0)
	}
	if line.V1 != "" {
		queryStr += " and v1 = ?"
		queryArgs = append(queryArgs, line.V1)
	}
	if line.V2 != "" {
		queryStr += " and v2 = ?"
		queryArgs = append(queryArgs, line.V2)
	}
	if line.V3 != "" {
		queryStr += " and v3 = ?"
		queryArgs = append(queryArgs, line.V3)
	}
	if line.V4 != "" {
		queryStr += " and v4 = ?"
		queryArgs = append(queryArgs, line.V4)
	}
	if line.V5 != "" {
		queryStr += " and v5 = ?"
		queryArgs = append(queryArgs, line.V5)
	}
	args := append([]interface{}{queryStr}, queryArgs...)
	err := db.Delete(a.getTableInstance(), args...).Error
	return err
}

func appendWhere(line CasbinRule) (string, []interface{}) {
	queryArgs := []interface{}{line.RuleType}
	queryStr := "rule_type = ?"
	if line.V0 != "" {
		queryStr += " and v0 = ?"
		queryArgs = append(queryArgs, line.V0)
	}
	if line.V1 != "" {
		queryStr += " and v1 = ?"
		queryArgs = append(queryArgs, line.V1)
	}
	if line.V2 != "" {
		queryStr += " and v2 = ?"
		queryArgs = append(queryArgs, line.V2)
	}
	if line.V3 != "" {
		queryStr += " and v3 = ?"
		queryArgs = append(queryArgs, line.V3)
	}
	if line.V4 != "" {
		queryStr += " and v4 = ?"
		queryArgs = append(queryArgs, line.V4)
	}
	if line.V5 != "" {
		queryStr += " and v5 = ?"
		queryArgs = append(queryArgs, line.V5)
	}
	return queryStr, queryArgs
}

// UpdatePolicy updates a new policy rule to DB.
func (a *Adapter) UpdatePolicy(ruleType string, oldRule, newPolicy []string) error {
	oldLine := a.savePolicyLine(ruleType, oldRule)
	queryStr, queryArgs := appendWhere(oldLine)
	newLine := a.savePolicyLine(ruleType, newPolicy)
	err := a.db.Where(queryStr, queryArgs...).Updates(newLine).Error
	return err
}
