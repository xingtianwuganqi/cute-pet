package generate

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func GetPetCustomTypeTable(ctx *context.Context) table.Table {

	petCustomType := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("mysql").SetPrimaryKey("id", db.Bigint))

	info := petCustomType.GetInfo().HideFilterArea()

	info.AddField("Id", "id", db.Bigint).
		FieldFilterable()
	info.AddField("Created_at", "created_at", db.Datetime)
	info.AddField("Updated_at", "updated_at", db.Datetime)
	info.AddField("Deleted_at", "deleted_at", db.Datetime)
	info.AddField("User_id", "user_id", db.Bigint)
	info.AddField("Custom_name", "custom_name", db.Varchar)
	info.AddField("Custom_icon", "custom_icon", db.Varchar)

	info.SetTable("pet_custom_type").SetTitle("PetCustomType").SetDescription("PetCustomType")

	formList := petCustomType.GetForm()
	formList.AddField("Id", "id", db.Bigint, form.Default)
	formList.AddField("Created_at", "created_at", db.Datetime, form.Datetime)
	formList.AddField("Updated_at", "updated_at", db.Datetime, form.Datetime)
	formList.AddField("Deleted_at", "deleted_at", db.Datetime, form.Datetime)
	formList.AddField("User_id", "user_id", db.Bigint, form.Number)
	formList.AddField("Custom_name", "custom_name", db.Varchar, form.Text)
	formList.AddField("Custom_icon", "custom_icon", db.Varchar, form.Text)

	formList.SetTable("pet_custom_type").SetTitle("PetCustomType").SetDescription("PetCustomType")

	return petCustomType
}
