package api

//func (a *Api) releaseGroups(context *gin.Context) {
//	groups, err := a.c.ReleaseGroups()
//	if err != nil {
//		context.Writer.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//
//	response := &model.ReleaseGroupsResponse{Groups: groups}
//
//	scheme := "http"
//	if context.Request.TLS != nil || context.GetHeader("X-Homessl-Forwarded") == "true" {
//		scheme = "https"
//	}
//
//	base := fmt.Sprintf("%s://%s", scheme, context.Request.Host+context.Request.URL.Path)
//
//	for _, group := range response.Groups {
//		group.CoverArt = fmt.Sprintf("%s/%s/cover", base, group.ID)
//	}
//
//	context.JSON(http.StatusOK, response)
//}
//
//func (a *Api) releaseGroup(context *gin.Context) {
//	uu, err := uuid.Parse(context.Param("uuid"))
//	if err != nil {
//		context.Writer.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	group, err := a.c.ReleaseGroup(uu)
//	if err != nil {
//		context.AbortWithError(http.StatusInternalServerError, err)
//		return
//	}
//
//	if group == nil {
//		context.Status(http.StatusNotFound)
//		return
//	}
//
//	context.JSON(http.StatusOK, group)
//}
