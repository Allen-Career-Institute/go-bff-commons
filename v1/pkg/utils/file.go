package utils

import (
	"fmt"
	frameworkModels "github.com/Allen-Career-Institute/go-bff-commons/v1/framework/models/commons"
	"io"
	"net/http"

	"github.com/Allen-Career-Institute/go-kratos-commons/utils/excel"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/labstack/echo/v4"

	"github.com/Allen-Career-Institute/go-bff-commons/v1/intrnl/datasources"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/logger"
)

func GetDataFromExcel(c echo.Context, l logger.Logger) ([]byte, error) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return nil, err
	}

	fileSrc, err := fileHeader.Open()
	if err != nil {
		l.WithContext(c).Errorf("GetDataFromExcel Error: while opening file :: %v", err)
		return nil, err
	}
	defer fileSrc.Close()

	fileContent, err := io.ReadAll(fileSrc)
	if err != nil {
		l.WithContext(c).Errorf("GetDataFromExcel Error: while reading file :: %v", err)
		return nil, err
	}

	return fileContent, nil
}

func GetParsedData(c echo.Context, fileData []byte, sheetName string, l logger.Logger) ([][]string, error) {
	helper, helperErr := excel.NewExcelHelper(excel.WithBufferOption(fileData))
	if helperErr != nil {
		l.WithContext(c).Errorf("GetParsedData Error: while creating excel helper :: %v", helperErr)
		return nil, helperErr
	}

	sheet, err := helper.ReadSheet(sheetName, 1)
	if err != nil {
		l.WithContext(c).Errorf("GetParsedData Error: while reading %v Sheet:: %v", sheetName, err)
		return nil, errors.New(http.StatusBadRequest, "Sheet Missing", fmt.Sprintf("%v Sheet is missing", sheetName))
	}

	return sheet, nil
}

func PopulateErrorDataRowWise(parsedData [][]string, errs map[int]string) [][]string {
	for rowIdx, err := range errs {
		parsedData[rowIdx] = append(parsedData[rowIdx], err)
	}

	return parsedData
}

func PopulateErrorData(parsedData [][]string, errs []string) [][]string {
	for _, err := range errs {
		errorData := make([][]string, 1)
		errorData[0] = []string{err}
		parsedData = append(parsedData, errorData...)
	}

	return parsedData
}

func AddValidationErrorDataToResponse(c echo.Context, sheetName string, data []byte, parsedDataWithError [][]string, l logger.Logger) ([]byte, error) {
	helper, helperErr := excel.NewExcelHelper(excel.WithBufferOption(data))
	if helperErr != nil {
		return nil, helperErr
	}

	err := helper.AddData(sheetName, 1, parsedDataWithError)
	if err != nil {
		l.WithContext(c).Errorf("AddValidationErrorDataToResponse Error: while adding data to excel :: %v", err)
		return nil, err
	}

	data, err = helper.ToByteArray()
	if err != nil {
		l.WithContext(c).Errorf("AddValidationErrorDataToResponse Error: while converting excel to byte array :: %v", err)
		return nil, err
	}

	return data, nil
}

func AddErrorDataToExcel(c echo.Context, fileName string, data []byte, l logger.Logger, errorResponse error) (frameworkModels.DSResponse, error) {
	l.WithContext(c).Info("addErrorDataToExcel call")
	c.Response().Header().Set(echo.HeaderContentDisposition, "attachment; filename="+fileName)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEOctetStream)

	// Send the byte array as a response
	c.Response().WriteHeader(http.StatusBadRequest)
	_, err := c.Response().Write(data)

	if err != nil {
		return datasources.PopulateResponse(http.StatusInternalServerError, GenericError, err), nil
	}

	if errorResponse != nil {
		return HandleErrorAndConvertResponse(errorResponse)
	}

	return datasources.PopulateResponse(http.StatusInternalServerError, GenericError, nil), nil
}
