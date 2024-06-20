// This file is generated by GoAdmin CLI adm.
package generate

import "github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"

// The key of Generators is the prefix of table info url.
// The corresponding value is the Form and Table data.
//
// http://{{config.Domain}}:{{Port}}/{{config.Prefix}}/info/{{key}}
//
// example:
//
// "pet_action_type" => http://localhost:9033/admin/info/pet_action_type
// "pet_custom_type" => http://localhost:9033/admin/info/pet_custom_type
// "pet_info" => http://localhost:9033/admin/info/pet_info
// "record_lists" => http://localhost:9033/admin/info/record_lists
// "suggestion" => http://localhost:9033/admin/info/suggestion
// "user_info" => http://localhost:9033/admin/info/user_info
//
// example end
var Generators = map[string]table.Generator{

	"pet_action_type": GetPetActionTypeTable,
	"pet_custom_type": GetPetCustomTypeTable,
	"pet_info":        GetPetInfoTable,
	"record_lists":    GetRecordListsTable,
	"suggestion":      GetSuggestionTable,
	"user_info":       GetUserInfoTable,

	// generators end
}
