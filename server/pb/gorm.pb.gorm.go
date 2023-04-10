package pb

import (
	context "context"
	fmt "fmt"
	gorm1 "github.com/infobloxopen/atlas-app-toolkit/gorm"
	errors "github.com/infobloxopen/protoc-gen-gorm/errors"
	gorm "github.com/jinzhu/gorm"
	field_mask "google.golang.org/genproto/protobuf/field_mask"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	strings "strings"
	time "time"
)

type OrderORM struct {
	Cutomer   string `gorm:"not null"`
	Id        uint64 `gorm:"primary_key;not null"`
	Price     float64
	Quantity  uint64
	RequestId uint64 `gorm:"not null"`
	Timestamp *time.Time
}

// TableName overrides the default tablename generated by GORM
func (OrderORM) TableName() string {
	return "order"
}

// ToORM runs the BeforeToORM hook if present, converts the fields of this
// object to ORM format, runs the AfterToORM hook, then returns the ORM object
func (m *Order) ToORM(ctx context.Context) (OrderORM, error) {
	to := OrderORM{}
	var err error
	if prehook, ok := interface{}(m).(OrderWithBeforeToORM); ok {
		if err = prehook.BeforeToORM(ctx, &to); err != nil {
			return to, err
		}
	}
	to.Id = m.Id
	to.Cutomer = m.Cutomer
	to.Quantity = m.Quantity
	to.Price = m.Price
	to.RequestId = m.RequestId
	if m.Timestamp != nil {
		t := m.Timestamp.AsTime()
		to.Timestamp = &t
	}
	if posthook, ok := interface{}(m).(OrderWithAfterToORM); ok {
		err = posthook.AfterToORM(ctx, &to)
	}
	return to, err
}

// ToPB runs the BeforeToPB hook if present, converts the fields of this
// object to PB format, runs the AfterToPB hook, then returns the PB object
func (m *OrderORM) ToPB(ctx context.Context) (Order, error) {
	to := Order{}
	var err error
	if prehook, ok := interface{}(m).(OrderWithBeforeToPB); ok {
		if err = prehook.BeforeToPB(ctx, &to); err != nil {
			return to, err
		}
	}
	to.Id = m.Id
	to.Cutomer = m.Cutomer
	to.Quantity = m.Quantity
	to.Price = m.Price
	to.RequestId = m.RequestId
	if m.Timestamp != nil {
		to.Timestamp = timestamppb.New(*m.Timestamp)
	}
	if posthook, ok := interface{}(m).(OrderWithAfterToPB); ok {
		err = posthook.AfterToPB(ctx, &to)
	}
	return to, err
}

// The following are interfaces you can implement for special behavior during ORM/PB conversions
// of type Order the arg will be the target, the caller the one being converted from

// OrderBeforeToORM called before default ToORM code
type OrderWithBeforeToORM interface {
	BeforeToORM(context.Context, *OrderORM) error
}

// OrderAfterToORM called after default ToORM code
type OrderWithAfterToORM interface {
	AfterToORM(context.Context, *OrderORM) error
}

// OrderBeforeToPB called before default ToPB code
type OrderWithBeforeToPB interface {
	BeforeToPB(context.Context, *Order) error
}

// OrderAfterToPB called after default ToPB code
type OrderWithAfterToPB interface {
	AfterToPB(context.Context, *Order) error
}

// DefaultCreateOrder executes a basic gorm create call
func DefaultCreateOrder(ctx context.Context, in *Order, db *gorm.DB) (*Order, error) {
	if in == nil {
		return nil, errors.NilArgumentError
	}
	ormObj, err := in.ToORM(ctx)
	if err != nil {
		return nil, err
	}
	if hook, ok := interface{}(&ormObj).(OrderORMWithBeforeCreate_); ok {
		if db, err = hook.BeforeCreate_(ctx, db); err != nil {
			return nil, err
		}
	}
	if err = db.Create(&ormObj).Error; err != nil {
		return nil, err
	}
	if hook, ok := interface{}(&ormObj).(OrderORMWithAfterCreate_); ok {
		if err = hook.AfterCreate_(ctx, db); err != nil {
			return nil, err
		}
	}
	pbResponse, err := ormObj.ToPB(ctx)
	return &pbResponse, err
}

type OrderORMWithBeforeCreate_ interface {
	BeforeCreate_(context.Context, *gorm.DB) (*gorm.DB, error)
}
type OrderORMWithAfterCreate_ interface {
	AfterCreate_(context.Context, *gorm.DB) error
}

func DefaultReadOrder(ctx context.Context, in *Order, db *gorm.DB) (*Order, error) {
	if in == nil {
		return nil, errors.NilArgumentError
	}
	ormObj, err := in.ToORM(ctx)
	if err != nil {
		return nil, err
	}
	if ormObj.Id == 0 {
		return nil, errors.EmptyIdError
	}
	if hook, ok := interface{}(&ormObj).(OrderORMWithBeforeReadApplyQuery); ok {
		if db, err = hook.BeforeReadApplyQuery(ctx, db); err != nil {
			return nil, err
		}
	}
	if db, err = gorm1.ApplyFieldSelection(ctx, db, nil, &OrderORM{}); err != nil {
		return nil, err
	}
	if hook, ok := interface{}(&ormObj).(OrderORMWithBeforeReadFind); ok {
		if db, err = hook.BeforeReadFind(ctx, db); err != nil {
			return nil, err
		}
	}
	ormResponse := OrderORM{}
	if err = db.Where(&ormObj).First(&ormResponse).Error; err != nil {
		return nil, err
	}
	if hook, ok := interface{}(&ormResponse).(OrderORMWithAfterReadFind); ok {
		if err = hook.AfterReadFind(ctx, db); err != nil {
			return nil, err
		}
	}
	pbResponse, err := ormResponse.ToPB(ctx)
	return &pbResponse, err
}

type OrderORMWithBeforeReadApplyQuery interface {
	BeforeReadApplyQuery(context.Context, *gorm.DB) (*gorm.DB, error)
}
type OrderORMWithBeforeReadFind interface {
	BeforeReadFind(context.Context, *gorm.DB) (*gorm.DB, error)
}
type OrderORMWithAfterReadFind interface {
	AfterReadFind(context.Context, *gorm.DB) error
}

func DefaultDeleteOrder(ctx context.Context, in *Order, db *gorm.DB) error {
	if in == nil {
		return errors.NilArgumentError
	}
	ormObj, err := in.ToORM(ctx)
	if err != nil {
		return err
	}
	if ormObj.Id == 0 {
		return errors.EmptyIdError
	}
	if hook, ok := interface{}(&ormObj).(OrderORMWithBeforeDelete_); ok {
		if db, err = hook.BeforeDelete_(ctx, db); err != nil {
			return err
		}
	}
	err = db.Where(&ormObj).Delete(&OrderORM{}).Error
	if err != nil {
		return err
	}
	if hook, ok := interface{}(&ormObj).(OrderORMWithAfterDelete_); ok {
		err = hook.AfterDelete_(ctx, db)
	}
	return err
}

type OrderORMWithBeforeDelete_ interface {
	BeforeDelete_(context.Context, *gorm.DB) (*gorm.DB, error)
}
type OrderORMWithAfterDelete_ interface {
	AfterDelete_(context.Context, *gorm.DB) error
}

func DefaultDeleteOrderSet(ctx context.Context, in []*Order, db *gorm.DB) error {
	if in == nil {
		return errors.NilArgumentError
	}
	var err error
	keys := []uint64{}
	for _, obj := range in {
		ormObj, err := obj.ToORM(ctx)
		if err != nil {
			return err
		}
		if ormObj.Id == 0 {
			return errors.EmptyIdError
		}
		keys = append(keys, ormObj.Id)
	}
	if hook, ok := (interface{}(&OrderORM{})).(OrderORMWithBeforeDeleteSet); ok {
		if db, err = hook.BeforeDeleteSet(ctx, in, db); err != nil {
			return err
		}
	}
	err = db.Where("id in (?)", keys).Delete(&OrderORM{}).Error
	if err != nil {
		return err
	}
	if hook, ok := (interface{}(&OrderORM{})).(OrderORMWithAfterDeleteSet); ok {
		err = hook.AfterDeleteSet(ctx, in, db)
	}
	return err
}

type OrderORMWithBeforeDeleteSet interface {
	BeforeDeleteSet(context.Context, []*Order, *gorm.DB) (*gorm.DB, error)
}
type OrderORMWithAfterDeleteSet interface {
	AfterDeleteSet(context.Context, []*Order, *gorm.DB) error
}

// DefaultStrictUpdateOrder clears / replaces / appends first level 1:many children and then executes a gorm update call
func DefaultStrictUpdateOrder(ctx context.Context, in *Order, db *gorm.DB) (*Order, error) {
	if in == nil {
		return nil, fmt.Errorf("Nil argument to DefaultStrictUpdateOrder")
	}
	ormObj, err := in.ToORM(ctx)
	if err != nil {
		return nil, err
	}
	lockedRow := &OrderORM{}
	db.Model(&ormObj).Set("gorm:query_option", "FOR UPDATE").Where("id=?", ormObj.Id).First(lockedRow)
	if hook, ok := interface{}(&ormObj).(OrderORMWithBeforeStrictUpdateCleanup); ok {
		if db, err = hook.BeforeStrictUpdateCleanup(ctx, db); err != nil {
			return nil, err
		}
	}
	if hook, ok := interface{}(&ormObj).(OrderORMWithBeforeStrictUpdateSave); ok {
		if db, err = hook.BeforeStrictUpdateSave(ctx, db); err != nil {
			return nil, err
		}
	}
	if err = db.Save(&ormObj).Error; err != nil {
		return nil, err
	}
	if hook, ok := interface{}(&ormObj).(OrderORMWithAfterStrictUpdateSave); ok {
		if err = hook.AfterStrictUpdateSave(ctx, db); err != nil {
			return nil, err
		}
	}
	pbResponse, err := ormObj.ToPB(ctx)
	if err != nil {
		return nil, err
	}
	return &pbResponse, err
}

type OrderORMWithBeforeStrictUpdateCleanup interface {
	BeforeStrictUpdateCleanup(context.Context, *gorm.DB) (*gorm.DB, error)
}
type OrderORMWithBeforeStrictUpdateSave interface {
	BeforeStrictUpdateSave(context.Context, *gorm.DB) (*gorm.DB, error)
}
type OrderORMWithAfterStrictUpdateSave interface {
	AfterStrictUpdateSave(context.Context, *gorm.DB) error
}

// DefaultPatchOrder executes a basic gorm update call with patch behavior
func DefaultPatchOrder(ctx context.Context, in *Order, updateMask *field_mask.FieldMask, db *gorm.DB) (*Order, error) {
	if in == nil {
		return nil, errors.NilArgumentError
	}
	var pbObj Order
	var err error
	if hook, ok := interface{}(&pbObj).(OrderWithBeforePatchRead); ok {
		if db, err = hook.BeforePatchRead(ctx, in, updateMask, db); err != nil {
			return nil, err
		}
	}
	pbReadRes, err := DefaultReadOrder(ctx, &Order{Id: in.GetId()}, db)
	if err != nil {
		return nil, err
	}
	pbObj = *pbReadRes
	if hook, ok := interface{}(&pbObj).(OrderWithBeforePatchApplyFieldMask); ok {
		if db, err = hook.BeforePatchApplyFieldMask(ctx, in, updateMask, db); err != nil {
			return nil, err
		}
	}
	if _, err := DefaultApplyFieldMaskOrder(ctx, &pbObj, in, updateMask, "", db); err != nil {
		return nil, err
	}
	if hook, ok := interface{}(&pbObj).(OrderWithBeforePatchSave); ok {
		if db, err = hook.BeforePatchSave(ctx, in, updateMask, db); err != nil {
			return nil, err
		}
	}
	pbResponse, err := DefaultStrictUpdateOrder(ctx, &pbObj, db)
	if err != nil {
		return nil, err
	}
	if hook, ok := interface{}(pbResponse).(OrderWithAfterPatchSave); ok {
		if err = hook.AfterPatchSave(ctx, in, updateMask, db); err != nil {
			return nil, err
		}
	}
	return pbResponse, nil
}

type OrderWithBeforePatchRead interface {
	BeforePatchRead(context.Context, *Order, *field_mask.FieldMask, *gorm.DB) (*gorm.DB, error)
}
type OrderWithBeforePatchApplyFieldMask interface {
	BeforePatchApplyFieldMask(context.Context, *Order, *field_mask.FieldMask, *gorm.DB) (*gorm.DB, error)
}
type OrderWithBeforePatchSave interface {
	BeforePatchSave(context.Context, *Order, *field_mask.FieldMask, *gorm.DB) (*gorm.DB, error)
}
type OrderWithAfterPatchSave interface {
	AfterPatchSave(context.Context, *Order, *field_mask.FieldMask, *gorm.DB) error
}

// DefaultPatchSetOrder executes a bulk gorm update call with patch behavior
func DefaultPatchSetOrder(ctx context.Context, objects []*Order, updateMasks []*field_mask.FieldMask, db *gorm.DB) ([]*Order, error) {
	if len(objects) != len(updateMasks) {
		return nil, fmt.Errorf(errors.BadRepeatedFieldMaskTpl, len(updateMasks), len(objects))
	}

	results := make([]*Order, 0, len(objects))
	for i, patcher := range objects {
		pbResponse, err := DefaultPatchOrder(ctx, patcher, updateMasks[i], db)
		if err != nil {
			return nil, err
		}

		results = append(results, pbResponse)
	}

	return results, nil
}

// DefaultApplyFieldMaskOrder patches an pbObject with patcher according to a field mask.
func DefaultApplyFieldMaskOrder(ctx context.Context, patchee *Order, patcher *Order, updateMask *field_mask.FieldMask, prefix string, db *gorm.DB) (*Order, error) {
	if patcher == nil {
		return nil, nil
	} else if patchee == nil {
		return nil, errors.NilArgumentError
	}
	var err error
	var updatedTimestamp bool
	for i, f := range updateMask.Paths {
		if f == prefix+"Id" {
			patchee.Id = patcher.Id
			continue
		}
		if f == prefix+"Cutomer" {
			patchee.Cutomer = patcher.Cutomer
			continue
		}
		if f == prefix+"Quantity" {
			patchee.Quantity = patcher.Quantity
			continue
		}
		if f == prefix+"Price" {
			patchee.Price = patcher.Price
			continue
		}
		if f == prefix+"RequestId" {
			patchee.RequestId = patcher.RequestId
			continue
		}
		if !updatedTimestamp && strings.HasPrefix(f, prefix+"Timestamp.") {
			if patcher.Timestamp == nil {
				patchee.Timestamp = nil
				continue
			}
			if patchee.Timestamp == nil {
				patchee.Timestamp = &timestamppb.Timestamp{}
			}
			childMask := &field_mask.FieldMask{}
			for j := i; j < len(updateMask.Paths); j++ {
				if trimPath := strings.TrimPrefix(updateMask.Paths[j], prefix+"Timestamp."); trimPath != updateMask.Paths[j] {
					childMask.Paths = append(childMask.Paths, trimPath)
				}
			}
			if err := gorm1.MergeWithMask(patcher.Timestamp, patchee.Timestamp, childMask); err != nil {
				return nil, nil
			}
		}
		if f == prefix+"Timestamp" {
			updatedTimestamp = true
			patchee.Timestamp = patcher.Timestamp
			continue
		}
	}
	if err != nil {
		return nil, err
	}
	return patchee, nil
}

// DefaultListOrder executes a gorm list call
func DefaultListOrder(ctx context.Context, db *gorm.DB) ([]*Order, error) {
	in := Order{}
	ormObj, err := in.ToORM(ctx)
	if err != nil {
		return nil, err
	}
	if hook, ok := interface{}(&ormObj).(OrderORMWithBeforeListApplyQuery); ok {
		if db, err = hook.BeforeListApplyQuery(ctx, db); err != nil {
			return nil, err
		}
	}
	db, err = gorm1.ApplyCollectionOperators(ctx, db, &OrderORM{}, &Order{}, nil, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	if hook, ok := interface{}(&ormObj).(OrderORMWithBeforeListFind); ok {
		if db, err = hook.BeforeListFind(ctx, db); err != nil {
			return nil, err
		}
	}
	db = db.Where(&ormObj)
	db = db.Order("id")
	ormResponse := []OrderORM{}
	if err := db.Find(&ormResponse).Error; err != nil {
		return nil, err
	}
	if hook, ok := interface{}(&ormObj).(OrderORMWithAfterListFind); ok {
		if err = hook.AfterListFind(ctx, db, &ormResponse); err != nil {
			return nil, err
		}
	}
	pbResponse := []*Order{}
	for _, responseEntry := range ormResponse {
		temp, err := responseEntry.ToPB(ctx)
		if err != nil {
			return nil, err
		}
		pbResponse = append(pbResponse, &temp)
	}
	return pbResponse, nil
}

type OrderORMWithBeforeListApplyQuery interface {
	BeforeListApplyQuery(context.Context, *gorm.DB) (*gorm.DB, error)
}
type OrderORMWithBeforeListFind interface {
	BeforeListFind(context.Context, *gorm.DB) (*gorm.DB, error)
}
type OrderORMWithAfterListFind interface {
	AfterListFind(context.Context, *gorm.DB, *[]OrderORM) error
}
