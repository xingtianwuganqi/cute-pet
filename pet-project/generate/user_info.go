package generate

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func GetUserInfoTable(ctx *context.Context) table.Table {

	userInfo := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("mysql").SetPrimaryKey("id", db.Bigint))

	info := userInfo.GetInfo().HideFilterArea()

	info.AddField("Id", "id", db.Bigint).
		FieldFilterable()
	info.AddField("Created_at", "created_at", db.Datetime)
	info.AddField("Updated_at", "updated_at", db.Datetime)
	info.AddField("Deleted_at", "deleted_at", db.Datetime)
	info.AddField("Phone", "phone", db.Varchar)
	info.AddField("Email", "email", db.Varchar)
	info.AddField("Username", "username", db.Varchar)
	info.AddField("Password", "password", db.Varchar)
	info.AddField("Avatar", "avatar", db.Varchar)
	info.AddField("Wx", "wx", db.Varchar)
	info.AddField("Location", "location", db.Bigint)

	info.SetTable("user_info").SetTitle("UserInfo").SetDescription("UserInfo")

	formList := userInfo.GetForm()
	formList.AddField("Id", "id", db.Bigint, form.Default)
	formList.AddField("Created_at", "created_at", db.Datetime, form.Datetime)
	formList.AddField("Updated_at", "updated_at", db.Datetime, form.Datetime)
	formList.AddField("Deleted_at", "deleted_at", db.Datetime, form.Datetime)
	formList.AddField("Phone", "phone", db.Varchar, form.Text)
	formList.AddField("Email", "email", db.Varchar, form.Email)
	formList.AddField("Username", "username", db.Varchar, form.Text)
	formList.AddField("Password", "password", db.Varchar, form.Password)
	formList.AddField("Avatar", "avatar", db.Varchar, form.Text)
	formList.AddField("Wx", "wx", db.Varchar, form.Text)
	formList.AddField("Location", "location", db.Bigint, form.Number)

	formList.SetTable("user_info").SetTitle("UserInfo").SetDescription("UserInfo")

	return userInfo
}
