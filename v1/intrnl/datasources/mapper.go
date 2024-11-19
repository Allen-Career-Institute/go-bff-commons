package datasources

import models "github.com/Allen-Career-Institute/go-bff-commons/v1/framework/models/commons"

// PopulateResponse , this method should be used to wrap errors or success responses to
// a generic response which would be propagated to the UI. Status, reason and the actual
// payload/error would all be sent to the UI.
func PopulateResponse(status int, reason string, data interface{}) models.DSResponse {
	return models.DSResponse{Status: status, Reason: reason, Data: data}
}

func PopulateResponseWithWidgetType(status int, reason string, viewType string, data interface{}) models.DSResponse {
	return models.DSResponse{Status: status, Reason: reason, ViewType: viewType, Data: data}
}
