package pagehandler

import (
	pbTypes "github.com/Allen-Career-Institute/common-protos/page_service/v1/types"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/framework/datasource"
	commonModels "github.com/Allen-Career-Institute/go-bff-commons/v1/framework/models/commons"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/intrnl/models/page"
	internalUtils "github.com/Allen-Career-Institute/go-bff-commons/v1/intrnl/utils"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/otel"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/utils"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/server/routes"
	"github.com/labstack/echo/v4"
	"runtime/debug"
	"sync"
)

func createDsNameToWidgetData(c echo.Context, widgets []*page.WidgetData, dsToWidgetDataMap map[string][]*page.WidgetData) {
	for _, w := range widgets {
		dsName := w.DataSource
		if _, ok := dsToWidgetDataMap[dsName]; !ok {
			dsToWidgetDataMap[dsName] = make([]*page.WidgetData, 0)
		}
		dsToWidgetDataMap[dsName] = append(dsToWidgetDataMap[dsName], w)
	}
}

// processDataSources processes data sources concurrently(Scatter - Gather)
func (pdh *pageDataHandler) processDataSources(c echo.Context, pageResp *page.CommonPageResponse, dsNamesConfiguredInPage []string, wPosList []string, resolvedWidgetsMap map[string]bool) (*page.CommonPageResponse, error) {
	var err error
	existingDataSources := pdh.GetDSList(c, dsNamesConfiguredInPage)

	//TODO: optimise number of go-routines (check if we can use global pool of go-routines)
	tc := utils.GetMin(len(existingDataSources), pdh.cnf.GoPool.MaxConcurrentRoutines)
	dsResults := make([]*commonModels.DSResponse, tc)

	//pdh.logger.Infof("Echo Addr in processDataSources : %s", reflect.ValueOf(c).Pointer())
	// Mind the "loop variable capture" problem
	// (created local copies of index(id) and datasource(ds) inside loop to make sure expected values are passed to goroutine)
	var wg sync.WaitGroup
	ctxDetailsMap, ok := internalUtils.GetValueFromContext[map[int]*page.WidgetData](c, utils.WidgetIndexToWidgetDataMap)
	if !ok {
		pdh.logger.Errorf("Unable to parse ctxDetailsMap")
	}

	for i, dataSource := range existingDataSources {
		wg.Add(1)
		id := i
		ds := dataSource
		ctxNew := pdh.eutil.CloneContext(c)
		go func(id int, ds *datasource.DataSource) { //using index to map
			defer func() {
				if r := recover(); r != nil {
					// Handle the panic (log, recover, etc.)
					stackTrace := debug.Stack()
					pdh.logger.WithContext(c).Errorf("Panic occurred in goroutine %d, ds: %s, %v\n%s", id, ds.Name(), r, stackTrace)
				}
				wg.Done()
			}()

			if ctxDetailsMap != nil {
				widgetData := ctxDetailsMap[id]
				if widgetData != nil {
					ctxNew.Set(utils.WidgetData, widgetData)
				}
			}

			dsResults[id] = pdh.worker(ctxNew, uint32(id), ds)
		}(id, ds)
	}
	// Wait for all workers to finish
	wg.Wait()
	pdh.logger.WithContext(c).Infof("dsResults : %+v", dsResults)
	// Now you can process the results in orderds
	for i, dsResult := range dsResults {
		var widgetId string
		if ctxDetailsMap != nil {
			if _, ok := ctxDetailsMap[i]; ok {
				widgetId = ctxDetailsMap[i].ConstWidgetID
			}
		}

		//Now map response back into pageResponse
		//pdh.logger.Infof("mapping DS resp into page for ds : %s, pos: %s, resp : %s ", existingDataSources[i].Name(), wPosList[i], dsResult)
		pageResp, err = pdh.prm.MapDataSourceRespToLP(c, existingDataSources[i].Name(), wPosList[i], widgetId, pageResp, dsResult, resolvedWidgetsMap)
		if err != nil {
			pdh.logger.WithContext(c).Errorf("Error while mapping DS resp into page for ds : %s, err: %v", existingDataSources[i].Name(), err)
		}
	}

	var olActions []*pbTypes.Action
	for _, ow := range pageResp.PageContent.OnloadWidgets {
		olAction := &pbTypes.Action{Type: ow.ViewType,
			Data: ow.Data,
		}
		olActions = append(olActions, olAction)
	}
	pageResp.PageInfo.OnloadActions = olActions
	return pageResp, nil
}

func (pdh *pageDataHandler) processPreloadDataSources(c *echo.Context, preloadDS []string) {
	pdh.logger.WithContext(*c).Info("processing preload datasources...")
	ec := *c
	dependentDSList := pdh.GetDSList(*c, preloadDS)
	dsData := make(map[string]*commonModels.DSResponse)

	existingDSMap, ok := ec.Get(utils.SharedDataSource).(map[string]*commonModels.DSResponse)
	if existingDSMap != nil && ok {
		dsData = existingDSMap
	}

	var (
		wg sync.WaitGroup
		mu sync.Mutex
	)

	for i, ds := range dependentDSList {
		if _, ok := dsData[ds.Name()]; ok {
			continue
		}
		ctx := pdh.eutil.CloneContext(*c)
		wg.Add(1)
		go func(id int, ds *datasource.DataSource) {
			defer wg.Done()
			res := pdh.worker(ctx, uint32(id), ds)
			if res != nil {
				mu.Lock()
				pdh.logger.WithContext(*c).Infof("setting data in dsData for ds: %s, res: %v", ds.Name(), res)
				dsData[ds.Name()] = res
				mu.Unlock()
			}
		}(i, ds)
	}

	wg.Wait()
	ec.Set(utils.SharedDataSource, dsData)
	pdh.logger.WithContext(*c).Debugf("CONTEXT %+v", ec)
}

func (pdh *pageDataHandler) worker(c echo.Context, taskID uint32, ds *datasource.DataSource) *commonModels.DSResponse {
	c, span := otel.Trace(c, ds.Name())
	defer span.End()

	resp, err := routes.NewDataSourceExecutor(ds, pdh.dsm, ds.Name(), &pdh.cnf, pdh.logger, pdh.meter, pdh.m).ExecuteDataSourceFromDS(c)
	if err != nil {
		pdh.logger.WithContext(c).Errorf("Data Source: %s failed to fetch the response, with error: %s, for req id: %d", ds.Name(), err, taskID)
		return nil
	}
	if resp.Status == 200 {
		return resp
	}
	return nil
}
