package core

import (
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"kructer.com/config"
)

type ModelRegistry struct {
	*gorm.DB
	models map[string]reflect.Value
	isOpen bool
}

func NewRegistry() *ModelRegistry {
	return &ModelRegistry{
		models: make(map[string]reflect.Value),
		isOpen: false,
	}
}

func (m *ModelRegistry) IsOpen() bool {
	return m.isOpen
}

func (m *ModelRegistry) OpenWithConfig(config *config.Configuration) error {
	db, err := gorm.Open(config.Dialect, config.ConnectionString)

	if err != nil {
		return err
	}

	db.DB().SetConnMaxLifetime(time.Minute * 5)
	db.DB().SetMaxIdleConns(0)
	db.DB().SetMaxOpenConns(20)

	m.DB = db
	m.isOpen = true
	return nil
}

func (m *ModelRegistry) Register(values ...interface{}) error {
	models := make(map[string]reflect.Value)
	if len(values) > 0 {
		for _, val := range values {
			rVal := reflect.ValueOf(val)
			if rVal.Kind() == reflect.Ptr {
				rVal = rVal.Elem()
			}
			switch rVal.Kind() {
			case reflect.Struct:
				models[getTypeName(rVal.Type())] = reflect.New(rVal.Type())
			default:
				return errors.New("models must be structs")
			}
		}
	}
	for k, v := range models {
		m.models[k] = v
	}
	return nil
}

func (m *ModelRegistry) AutoMigrateAll() {
	for _, v := range m.models {
		m.AutoMigrate(v.Interface())
	}
}

func (m *ModelRegistry) AutoDropAll() {
	for _, v := range m.models {
		m.DropTableIfExists(v.Interface())
	}
}

func getTypeName(typ reflect.Type) string {
	if typ.Name() != "" {
		return typ.Name()
	}
	split := strings.Split(typ.String(), ".")
	return split[len(split)-1]
}
