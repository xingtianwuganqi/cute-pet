package generate

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func GetSuggestionTable(ctx *context.Context) table.Table {

	suggestion := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("mysql").SetPrimaryKey("id", db.Bigint))

	info := suggestion.GetInfo().HideFilterArea()

	info.AddField("Id", "id", db.Bigint).
		FieldFilterable()
	info.AddField("Created_at", "created_at", db.Datetime)
	info.AddField("Updated_at", "updated_at", db.Datetime)
	info.AddField("Deleted_at", "deleted_at", db.Datetime)
	info.AddField("User_id", "user_id", db.Bigint)
	info.AddField("Contact", "contact", db.Varchar)
	info.AddField("Content", "content", db.Varchar)

	info.SetTable("suggestion").SetTitle("Suggestion").SetDescription("Suggestion")

	formList := suggestion.GetForm()
	formList.AddField("Id", "id", db.Bigint, form.Default)
	formList.AddField("Created_at", "created_at", db.Datetime, form.Datetime)
	formList.AddField("Updated_at", "updated_at", db.Datetime, form.Datetime)
	formList.AddField("Deleted_at", "deleted_at", db.Datetime, form.Datetime)
	formList.AddField("User_id", "user_id", db.Bigint, form.Number)
	formList.AddField("Contact", "contact", db.Varchar, form.Text)
	formList.AddField("Content", "content", db.Varchar, form.Text)

	formList.SetTable("suggestion").SetTitle("Suggestion").SetDescription("Suggestion")

	return suggestion
}
