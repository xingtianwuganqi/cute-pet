package generate

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func GetPetActionTypeTable(ctx *context.Context) table.Table {

	config := table.DefaultConfigWithDriver("mysql").SetPrimaryKey("id", db.Bigint)
	petActionType := table.NewDefaultTable(ctx, config)

	info := petActionType.GetInfo().HideFilterArea()

	info.AddField("Id", "id", db.Bigint).
		FieldFilterable()
	info.AddField("Created_at", "created_at", db.Datetime)
	info.AddField("Updated_at", "updated_at", db.Datetime)
	info.AddField("Deleted_at", "deleted_at", db.Datetime)
	info.AddField("Type", "type", db.Bigint)
	info.AddField("Action_name", "action_name", db.Longtext)
	info.AddField("Icon", "icon", db.Longtext)

	info.SetTable("pet_action_type").SetTitle("PetActionType").SetDescription("PetActionType")

	formList := petActionType.GetForm()
	formList.AddField("Id", "id", db.Bigint, form.Default)
	formList.AddField("Created_at", "created_at", db.Datetime, form.Datetime)
	formList.AddField("Updated_at", "updated_at", db.Datetime, form.Datetime)
	formList.AddField("Deleted_at", "deleted_at", db.Datetime, form.Datetime)
	formList.AddField("Type", "type", db.Bigint, form.Number)
	formList.AddField("Action_name", "action_name", db.Longtext, form.RichText)
	formList.AddField("Icon", "icon", db.Longtext, form.RichText)

	formList.SetTable("pet_action_type").SetTitle("PetActionType").SetDescription("PetActionType")

	return petActionType
}
