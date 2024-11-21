package utils

import (
	"fmt"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/logger"
	"io"

	"github.com/Allen-Career-Institute/go-kratos-commons/utils/excel"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/labstack/echo/v4"
)

func GetFileFromRequest(c echo.Context, l logger.Logger) ([]byte, error) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return nil, err
	}

	fileSrc, err := fileHeader.Open()
	if err != nil {
		l.WithContext(c).Errorf("GetFileFromRequest Error: while opening file :: %v", err)
		return nil, err
	}
	defer fileSrc.Close()

	fileContent, err := io.ReadAll(fileSrc)
	if err != nil {
		l.WithContext(c).Errorf("GetFileFromRequest Error: while reading file :: %v", err)
		return nil, err
	}

	return fileContent, nil
}

func GetParsedExcelSheet(c echo.Context, fileData []byte, l logger.Logger, sheetName string) ([][]string, error) {
	helper, err := excel.NewExcelHelper(excel.WithBufferOption(fileData))
	if err != nil {
		l.WithContext(c).Errorf("GetParsedData Error: while creating excel helper: %v", err)
		return nil, err
	}

	sheet, err := helper.ReadSheet(sheetName, 1)
	if err != nil {
		l.WithContext(c).Errorf("GetParsedData Error: while reading User Privileges Sheet: %v", err)
		return nil, errors.New(400, "Sheet Missing", fmt.Sprintf(" %v is missing", sheetName))
	}

	return sheet, nil
}

func ConvertExcelSheetToBytes(c echo.Context, l logger.Logger, sheet [][]string, sheetName string) ([]byte, error) {
	helper, err := excel.NewExcelHelper()
	if err != nil {
		l.WithContext(c).Errorf("ConvertExcelSheetToBytes Error: while creating excel helper :: %v", err)
		return nil, err
	}

	err = helper.AddSheet(excel.SheetData{
		SheetName: sheetName,
		Data:      sheet,
	})
	if err != nil {
		l.WithContext(c).Errorf("ConvertExcelSheetToBytes Error: while adding sheet :: %v", err)
		return nil, err
	}

	return helper.ToByteArray()
}
