package api

import (
	"github.com/elgatito/elementum/cloudhole"
	"github.com/elgatito/elementum/config"
	"github.com/elgatito/elementum/library"
	"github.com/elgatito/elementum/xbmc"
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
)

var cmdLog = logging.MustGetLogger("cmd")

// ClearCache ...
func ClearCache(ctx *gin.Context) {
	key := ctx.Params.ByName("key")
	if ctx != nil {
		ctx.Abort()
	}
	if key != "" {
		library.ClearCacheKey(key)
	} else {
		library.ClearPageCache()
	}
	xbmc.Notify("Elementum", "LOCALIZE[30200]", config.AddonIcon())
}

// ClearPageCache ...
func ClearPageCache(ctx *gin.Context) {
	if ctx != nil {
		ctx.Abort()
	}
	library.ClearPageCache()
}

// ClearTraktCache ...
func ClearTraktCache(ctx *gin.Context) {
	if ctx != nil {
		ctx.Abort()
	}
	library.ClearTraktCache()
}

// ClearTmdbCache ...
func ClearTmdbCache(ctx *gin.Context) {
	if ctx != nil {
		ctx.Abort()
	}
	library.ClearTmdbCache()
}

// ResetClearances ...
func ResetClearances(ctx *gin.Context) {
	cloudhole.ResetClearances()
	xbmc.Notify("Elementum", "LOCALIZE[30264]", config.AddonIcon())
}

// SetViewMode ...
func SetViewMode(ctx *gin.Context) {
	contentType := ctx.Params.ByName("content_type")
	viewName := xbmc.InfoLabel("Container.Viewmode")
	viewMode := xbmc.GetCurrentView()
	cmdLog.Noticef("ViewMode: %s (%s)", viewName, viewMode)
	if viewMode != "0" {
		xbmc.SetSetting("viewmode_"+contentType, viewMode)
	}
	ctx.String(200, "")
}
