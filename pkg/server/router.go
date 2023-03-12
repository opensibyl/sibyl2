package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensibyl/sibyl2/pkg/server/service"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func injectV1Group(v1group *gin.RouterGroup) {
	// scope
	scopeGroup := v1group.Group("/")
	scopeGroup.Handle(http.MethodGet, "/repo", service.HandleRepoQuery)
	scopeGroup.Handle(http.MethodDelete, "/repo", service.HandleRepoDelete)
	scopeGroup.Handle(http.MethodGet, "/rev", service.HandleRevQuery)
	scopeGroup.Handle(http.MethodDelete, "/rev", service.HandleRevDelete)
	scopeGroup.Handle(http.MethodGet, "/file", service.HandleFileQuery)
	// upload
	uploadGroup := v1group.Group("/")
	uploadGroup.Handle(http.MethodPost, "/func", service.HandleFunctionUpload)
	uploadGroup.Handle(http.MethodPost, "/funcctx", service.HandleFunctionContextUpload)
	uploadGroup.Handle(http.MethodPost, "/clazz", service.HandleClazzUpload)
	// basic
	basicGroup := v1group.Group("/")
	basicGroup.Handle(http.MethodGet, "/func", service.HandleFunctionsQuery)
	basicGroup.Handle(http.MethodGet, "/funcctx", service.HandleFunctionContextsQuery)
	basicGroup.Handle(http.MethodGet, "/clazz", service.HandleClazzesQuery)

	// query by signature
	signatureGroup := v1group.Group("signature")
	signatureGroup.Handle(http.MethodGet, "/regex/func", service.HandleSignatureRegexFunc)
	signatureGroup.Handle(http.MethodGet, "/func", service.HandleSignatureFunc)
	signatureGroup.Handle(http.MethodGet, "/funcctx", service.HandleSignatureFuncctx)
	signatureGroup.Handle(http.MethodGet, "/funcctx/chain", service.HandleSignatureFuncctxChain)
	signatureGroup.Handle(http.MethodGet, "/funcctx/rchain", service.HandleSignatureFuncctxReverseChain)
	// query by regex
	regexGroup := v1group.Group("regex")
	regexGroup.Handle(http.MethodGet, "/func", service.HandleRegexFunc)
	regexGroup.Handle(http.MethodGet, "/clazz", service.HandleRegexClazz)
	regexGroup.Handle(http.MethodGet, "/funcctx", service.HandleRegexFuncctx)
	// query by reference
	referenceGroup := v1group.Group("reference")
	countGroup := referenceGroup.Group("count")
	countGroup.Handle(http.MethodGet, "/funcctx", service.HandleReferenceCountFuncctx)
	countGroup.Handle(http.MethodGet, "/funcctx/reverse", service.HandleReferenceCountFuncctxReverse)
	// tag
	tagGroup := v1group.Group("tag")
	tagGroup.Handle(http.MethodGet, "/func", service.HandleFuncTagQuery)
	tagGroup.Handle(http.MethodPost, "/func", service.HandleFuncTagCreate)
	// query by stat
	v1group.Handle(http.MethodGet, "/rev/stat", service.HandleRevStatQuery)
}

func injectOpsGroup(opsGroup *gin.RouterGroup) {
	opsGroup.Handle(http.MethodGet, "/ping", service.HandlePing)
	opsGroup.Handle(http.MethodGet, "/monitor/upload", service.HandleStatusUpload)
	opsGroup.Handle(http.MethodGet, "/version", service.HandleVersion)
	opsGroup.Handle(http.MethodGet, "/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	opsGroup.Handle(http.MethodGet, "/metrics", gin.WrapH(promhttp.Handler()))
}
