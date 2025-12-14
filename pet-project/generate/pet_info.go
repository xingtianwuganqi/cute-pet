package generate

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func GetPetInfoTable(ctx *context.Context) table.Table {

	petInfo := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("mysql").SetPrimaryKey("id", db.Bigint))

	info := petInfo.GetInfo().HideFilterArea()

	info.AddField("Id", "id", db.Bigint).
		FieldFilterable()
	info.AddField("Created_at", "created_at", db.Datetime)
	info.AddField("Updated_at", "updated_at", db.Datetime)
	info.AddField("Deleted_at", "deleted_at", db.Datetime)
	info.AddField("User_id", "user_id", db.Bigint)
	info.AddField("Pet_type", "pet_type", db.Bigint)
	info.AddField("Age", "age", db.Bigint)
	info.AddField("Name", "name", db.Varchar)
	info.AddField("Avatar", "avatar", db.Varchar)
	info.AddField("Birth_day", "birth_day", db.Varchar)

	info.SetTable("pet_info").SetTitle("PetInfo").SetDescription("PetInfo")

	formList := petInfo.GetForm()
	formList.AddField("Id", "id", db.Bigint, form.Default)
	formList.AddField("Created_at", "created_at", db.Datetime, form.Datetime)
	formList.AddField("Updated_at", "updated_at", db.Datetime, form.Datetime)
	formList.AddField("Deleted_at", "deleted_at", db.Datetime, form.Datetime)
	formList.AddField("User_id", "user_id", db.Bigint, form.Number)
	formList.AddField("Pet_type", "pet_type", db.Bigint, form.Number)
	formList.AddField("Age", "age", db.Bigint, form.Number)
	formList.AddField("Name", "name", db.Varchar, form.Text)
	formList.AddField("Avatar", "avatar", db.Varchar, form.Text)
	formList.AddField("Birth_day", "birth_day", db.Varchar, form.Text)

	formList.SetTable("pet_info").SetTitle("PetInfo").SetDescription("PetInfo")

	return petInfo
}
