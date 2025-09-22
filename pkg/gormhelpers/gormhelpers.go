package gormhelpers

import (
	"context"
	"fmt"

	"github.com/thekrauss/beto-shared/pkg/errors"
	"gorm.io/gorm"
)

// insère un nouvel enregistrement
func Create[T any](ctx context.Context, db *gorm.DB, entity *T) error {
	if err := db.WithContext(ctx).Create(entity).Error; err != nil {
		return errors.Wrap(err, errors.CodeDBError, "failed to create entity")
	}
	return nil
}

// récupère un enregistrement par ID
func FindByID[T any](ctx context.Context, db *gorm.DB, id any) (*T, error) {
	var entity T
	err := db.WithContext(ctx).First(&entity, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.CodeDBNotFound, "record not found")
		}
		return nil, errors.Wrap(err, errors.CodeDBError, "failed to find entity by ID")
	}
	return &entity, nil
}

// met à jour un enregistrement complet
func Update[T any](ctx context.Context, db *gorm.DB, entity *T) error {
	if err := db.WithContext(ctx).Save(entity).Error; err != nil {
		return errors.Wrap(err, errors.CodeDBError, "failed to update entity")
	}
	return nil
}

// met à jour certains champs
func UpdateFields[T any](ctx context.Context, db *gorm.DB, id any, fields map[string]interface{}) error {
	if err := db.WithContext(ctx).Model(new(T)).Where("id = ?", id).Updates(fields).Error; err != nil {
		return errors.Wrap(err, errors.CodeDBError, "failed to update fields")
	}
	return nil
}

// supprime un enregistrement
func Delete[T any](ctx context.Context, db *gorm.DB, id any) error {
	if err := db.WithContext(ctx).Delete(new(T), "id = ?", id).Error; err != nil {
		return errors.Wrap(err, errors.CodeDBError, "failed to delete entity")
	}
	return nil
}

// vérifie si un enregistrement existe
func Exists[T any](ctx context.Context, db *gorm.DB, conditions map[string]interface{}) (bool, error) {
	var count int64
	if err := db.WithContext(ctx).Model(new(T)).Where(conditions).Count(&count).Error; err != nil {
		return false, errors.Wrap(err, errors.CodeDBError, "failed to check existence")
	}
	return count > 0, nil
}

// retourne le nombre d’enregistrements
func Count[T any](ctx context.Context, db *gorm.DB, conditions map[string]interface{}) (int64, error) {
	var count int64
	if err := db.WithContext(ctx).Model(new(T)).Where(conditions).Count(&count).Error; err != nil {
		return 0, errors.Wrap(err, errors.CodeDBError, "failed to count records")
	}
	return count, nil
}

// retourne une liste paginée
func FindAllPaginated[T any](ctx context.Context, db *gorm.DB, page, pageSize int, conditions map[string]interface{}) ([]T, int64, error) {
	var results []T
	var total int64

	query := db.WithContext(ctx).Model(new(T)).Where(conditions)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.Wrap(err, errors.CodeDBError, "failed to count results")
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&results).Error; err != nil {
		return nil, 0, errors.Wrap(err, errors.CodeDBError, "failed to fetch paginated results")
	}

	return results, total, nil
}

// recherche avec conditions dynamiques
func FindByConditions[T any](ctx context.Context, db *gorm.DB, conditions map[string]interface{}) ([]T, error) {
	var results []T
	if err := db.WithContext(ctx).Model(new(T)).Where(conditions).Find(&results).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBError, "failed to find by conditions")
	}
	return results, nil
}

// récupère ou crée si inexistant
func FirstOrCreate[T any](ctx context.Context, db *gorm.DB, conds map[string]interface{}, defaults *T) (*T, error) {
	var entity T
	if err := db.WithContext(ctx).Where(conds).FirstOrCreate(&entity, defaults).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBError, "failed to first-or-create entity")
	}
	return &entity, nil
}

// exécute une fonction dans une transaction
func Transaction(ctx context.Context, db *gorm.DB, fn func(tx *gorm.DB) error) error {
	if err := db.WithContext(ctx).Transaction(fn); err != nil {
		return errors.Wrap(err, errors.CodeDBError, "transaction failed")
	}
	return nil
}

// affiche la requête SQL générée (for dev)
func DebugSQL(db *gorm.DB) {
	sql := db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)
	fmt.Println("SQL:", sql)
}
