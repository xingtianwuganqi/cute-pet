package generate

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func GetRecordListsTable(ctx *context.Context) table.Table {

	config := table.DefaultConfigWithDriver("mysql").SetPrimaryKey("id", db.Bigint)
	recordLists := table.NewDefaultTable(ctx, config)

	info := recordLists.GetInfo().HideFilterArea()

	info.AddField("Id", "id", db.Bigint).
		FieldFilterable()
	info.AddField("Created_at", "created_at", db.Datetime)
	info.AddField("Updated_at", "updated_at", db.Datetime)
	info.AddField("Deleted_at", "deleted_at", db.Datetime)
	info.AddField("User_id", "user_id", db.Bigint)
	info.AddField("Pet_action_type_id", "pet_action_type_id", db.Bigint)
	info.AddField("Pet_custom_type_id", "pet_custom_type_id", db.Bigint)
	info.AddField("Spend", "spend", db.Float)
	info.AddField("Desc", "desc", db.Longtext)

	info.SetTable("record_lists").SetTitle("RecordLists").SetDescription("RecordLists")

	formList := recordLists.GetForm()
	formList.AddField("Id", "id", db.Bigint, form.Default)
	formList.AddField("Created_at", "created_at", db.Datetime, form.Datetime)
	formList.AddField("Updated_at", "updated_at", db.Datetime, form.Datetime)
	formList.AddField("Deleted_at", "deleted_at", db.Datetime, form.Datetime)
	formList.AddField("User_id", "user_id", db.Bigint, form.Number)
	formList.AddField("Pet_action_type_id", "pet_action_type_id", db.Bigint, form.Number)
	formList.AddField("Pet_custom_type_id", "pet_custom_type_id", db.Bigint, form.Number)
	formList.AddField("Spend", "spend", db.Float, form.Text)
	formList.AddField("Desc", "desc", db.Longtext, form.RichText)

	formList.SetTable("record_lists").SetTitle("RecordLists").SetDescription("RecordLists")

	return recordLists
}
