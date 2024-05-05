package dtos

import (
	"encoding/json"
	"github.com/MKwann7/zgEXCELL-3-Media/app/media/src/code/libraries/builder"
	"github.com/MKwann7/zgEXCELL-3-Media/app/media/src/code/libraries/db"
	"github.com/MKwann7/zgEXCELL-3-Media/app/media/src/code/libraries/helper"
	"github.com/google/uuid"
	"reflect"
)

type Images struct {
	builder builder.Builder
}

func (images *Images) GetById(userId int) (*Image, error) {
	connection := images.getConnection()
	model := Image{}
	interfaceModel, error := images.builder.GetById(userId, connection, reflect.TypeOf(model))

	if error != nil {
		return nil, error
	}

	returnModel := images.assignInterfaceModel(interfaceModel)

	return returnModel, nil
}

func (images *Images) GetByUuid(userUuid uuid.UUID) (*Image, error) {
	connection := images.getConnection()
	model := Image{}
	interfaceModel, error := images.builder.GetByUuid(userUuid, connection, reflect.TypeOf(model))

	if error != nil {
		return nil, error
	}

	returnModel := images.assignInterfaceModel(interfaceModel)

	return returnModel, nil
}

func (images *Images) CreateNew(model Image) (*Image, error) {
	connection := images.getConnection()

	var myModelMap map[string]interface{}
	data, _ := json.Marshal(model)
	json.Unmarshal(data, &myModelMap)

	fields := reflect.TypeOf(model)
	values := reflect.ValueOf(model)
	interfaceModel, error := images.builder.CreateNew(db.MakeSqlExecModel(fields, values, connection.PrimaryKey), connection)

	if error != nil {
		return nil, error
	}

	returnModel := images.assignInterfaceModel(interfaceModel)

	return returnModel, nil
}

// LocalAddr returns the local network address.
func (images *Images) getConnection() db.Connection {
	connection := db.Connection{}
	return connection.GetMedia("image", "image_id", "sys_row_id")
}

func (images *Images) assignInterfaceModel(model map[string]interface{}) *Image {
	returnModel := &Image{}
	returnModel.ImageId = helper.CastAsNullableInt(model["image_id"])
	returnModel.ParentId = helper.CastAsIntWithNull(model["parent_id"])
	returnModel.CompanyId = helper.CastAsNullableString(model["company_id"])
	returnModel.UserId = helper.CastAsNullableInt(model["user_id"])
	returnModel.EntityId = helper.CastAsNullableInt(model["entity_id"])
	returnModel.EntityName = helper.CastAsNullableString(model["entity_name"])
	returnModel.ImageClass = helper.CastAsNullableString(model["image_class"])
	returnModel.Title = helper.CastAsNullableString(model["title"])
	returnModel.Url = helper.CastAsNullableString(model["url"])
	returnModel.Thumb = helper.CastAsNullableString(model["thumb"])
	returnModel.Width = helper.CastAsNullableInt(model["width"])
	returnModel.Height = helper.CastAsNullableInt(model["height"])
	returnModel.Type = helper.CastAsNullableString(model["type"])

	if helper.CastAsNullableUuid(model["sys_row_id"]).Valid {
		returnModel.SysRowId = helper.CastAsNullableUuid(model["sys_row_id"]).Value
	}

	return returnModel
}

type Image struct {
	ImageId     int             `field:"image_id"`
	ParentId    helper.NullInt  `field:"parent_id"`
	CompanyId   string          `field:"company_id"`
	UserId      int             `field:"user_id"`
	EntityId    int             `field:"entity_id"`
	EntityName  string          `field:"entity_name"`
	ImageClass  string          `field:"image_class"`
	Title       string          `field:"title"`
	Url         string          `field:"url"`
	Thumb       string          `field:"thumb"`
	Width       int             `field:"width"`
	Height      int             `field:"height"`
	Type        string          `field:"type"`
	CreatedOn   helper.NullTime `field:"created_on"`
	CreatedBy   int             `field:"created_by"`
	LastUpdated helper.NullTime `field:"last_updated"`
	UpdatedBy   int             `field:"updated_by"`
	SysRowId    uuid.UUID       `field:"sys_row_id"`
}
