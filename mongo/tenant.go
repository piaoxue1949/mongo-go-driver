package mongo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
)

type TenantHookInfo struct {
	// 获取租户字段名的钩子函数
	GetTenantFieldHook func(context.Context) string
	// 获取租户id的钩子函数
	GetTenantIdHook func(context.Context) string
}

var gTenantHookInfo TenantHookInfo

func RegisterTenantHook(hookInfo TenantHookInfo) error {
	if hookInfo.GetTenantFieldHook == nil {
		return errors.New("get tenant field hook is nil")
	}
	if hookInfo.GetTenantIdHook == nil {
		return errors.New("get tenant id hook is nil")
	}
	gTenantHookInfo = hookInfo
	return nil
}

func DeRegisterGetTenantIdHook() {
	gTenantHookInfo.GetTenantFieldHook = nil
	gTenantHookInfo.GetTenantIdHook = nil
}

func GetTenantFilter(ctx context.Context, filter interface{}) interface{} {
	if gTenantHookInfo.GetTenantFieldHook == nil || gTenantHookInfo.GetTenantIdHook == nil {
		return filter
	}
	if filterD, ok := filter.(bson.D); ok {
		tenantField := gTenantHookInfo.GetTenantFieldHook(ctx)
		tenantId := gTenantHookInfo.GetTenantIdHook(ctx)
		if len(filterD) == 0 {
			filter = bson.D{{tenantField, tenantId}}
		} else {
			filter = append(filterD, bson.E{tenantField, tenantId})
		}
	}
	return filter
}
