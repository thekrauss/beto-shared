package gormhelpers

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// insère un nouvel enregistrement
func Create[T any](ctx context.Context, db *gorm.DB, entity *T) error {
	return db.WithContext(ctx).Create(entity).Error
}

// récupère un enregistrement par ID
func FindByID[T any](ctx context.Context, db *gorm.DB, id any) (*T, error) {
	var entity T
	err := db.WithContext(ctx).First(&entity, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// met à jour un enregistrement complet
func Update[T any](ctx context.Context, db *gorm.DB, entity *T) error {
	return db.WithContext(ctx).Save(entity).Error
}

// met à jour seulement certains champs
func UpdateFields[T any](ctx context.Context, db *gorm.DB, id any, fields map[string]interface{}) error {
	return db.WithContext(ctx).Model(new(T)).Where("id = ?", id).Updates(fields).Error
}

// supprime un enregistrement
func Delete[T any](ctx context.Context, db *gorm.DB, id any) error {
	return db.WithContext(ctx).Delete(new(T), "id = ?", id).Error
}

// vérifie si un enregistrement existe
func Exists[T any](ctx context.Context, db *gorm.DB, conditions map[string]interface{}) (bool, error) {
	var count int64
	err := db.WithContext(ctx).Model(new(T)).Where(conditions).Count(&count).Error
	return count > 0, err
}

// retourne le nombre d’enregistrements
func Count[T any](ctx context.Context, db *gorm.DB, conditions map[string]interface{}) (int64, error) {
	var count int64
	err := db.WithContext(ctx).Model(new(T)).Where(conditions).Count(&count).Error
	return count, err
}

// retourne une liste paginée
func FindAllPaginated[T any](ctx context.Context, db *gorm.DB, page, pageSize int, conditions map[string]interface{}) ([]T, int64, error) {
	var results []T
	var total int64

	query := db.WithContext(ctx).Model(new(T)).Where(conditions)

	// compte total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	//  offset & limit
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&results).Error; err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

// recherche avec conditions dynamiques
func FindByConditions[T any](ctx context.Context, db *gorm.DB, conditions map[string]interface{}) ([]T, error) {
	var results []T
	err := db.WithContext(ctx).Model(new(T)).Where(conditions).Find(&results).Error
	return results, err
}

// récupère ou crée si inexistant
func FirstOrCreate[T any](ctx context.Context, db *gorm.DB, conds map[string]interface{}, defaults *T) (*T, error) {
	var entity T
	err := db.WithContext(ctx).Where(conds).FirstOrCreate(&entity, defaults).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// exécute une fonction dans une transaction
func Transaction(ctx context.Context, db *gorm.DB, fn func(tx *gorm.DB) error) error {
	return db.WithContext(ctx).Transaction(fn)
}

// affiche la requête SQL générée (for dev)
func DebugSQL(db *gorm.DB) {
	sql := db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)
	fmt.Println("SQL:", sql)
}
