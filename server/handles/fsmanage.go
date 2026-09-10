package handles

import (
	"fmt"
	stdpath "path"
	"strings"

	"github.com/OpenListTeam/OpenList/v4/internal/conf"
	"github.com/OpenListTeam/OpenList/v4/internal/errs"
	"github.com/OpenListTeam/OpenList/v4/internal/fs"
	"github.com/OpenListTeam/OpenList/v4/internal/model"
	"github.com/OpenListTeam/OpenList/v4/internal/op"
	"github.com/OpenListTeam/OpenList/v4/internal/sign"
	"github.com/OpenListTeam/OpenList/v4/pkg/utils"
	"github.com/OpenListTeam/OpenList/v4/server/common"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type MkdirOrLinkReq struct {
	Path string `json:"path" form:"path"`
}

func checkRelativePath(path string) error {
	if strings.ContainsAny(path, "/\\") || path == "" || path == "." || path == ".." {
		return errs.RelativePath
	}
	return nil
}

type RemoveReq struct {
	Dir   string   `json:"dir"`
	Names []string `json:"names"`
}

// FsRemove performs batch remove (individual item permission checks skipped for performance).
func FsRemove(c *gin.Context) {
	var req RemoveReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	if len(req.Names) == 0 {
		common.ErrorStrResp(c, "Empty file names", 400)
		return
	}
	user := c.Request.Context().Value(conf.UserKey).(*model.User)
	if !user.CanRemove() {
		common.ErrorResp(c, errs.PermissionDenied, 403)
		return
	}
	reqPath, err := user.JoinPath(req.Dir)
	if err != nil {
		common.ErrorResp(c, err, 403)
		return
	}
	meta, err := op.GetNearestMeta(reqPath)
	if err != nil && !errors.Is(errors.Cause(err), errs.MetaNotFound) {
		common.ErrorResp(c, err, 500, true)
		return
	}
	if !common.CanWrite(user, meta, reqPath) {
		common.ErrorResp(c, errs.PermissionDenied, 403)
		return
	}
	if !strings.HasSuffix(reqPath, "/") {
		reqPath += "/"
	}
	for i, name := range req.Names {
		fullPath := stdpath.Join(reqPath, name)
		if !strings.HasPrefix(fullPath+"/", reqPath) {
			req.Names[i] = ""
			continue
		}
		req.Names[i] = fullPath
	}
	for _, path := range req.Names {
		if path == "" {
			continue
		}
		err := fs.Remove(c.Request.Context(), path)
		if err != nil {
			common.ErrorResp(c, err, 500)
			return
		}
	}
	//fs.ClearCache(req.Dir)
	common.SuccessResp(c)
}

// Link return real link, just for proxy program, it may contain cookie, so just allowed for admin
func Link(c *gin.Context) {
	var req MkdirOrLinkReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	//user := c.Request.Context().Value(conf.UserKey).(*model.User)
	//rawPath := stdpath.Join(user.BasePath, req.Path)
	// why need not join base_path? because it's always the full path
	rawPath := req.Path
	storage, err := fs.GetStorage(rawPath, &fs.GetStoragesArgs{})
	if err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	if storage.Config().NoLinkURL {
		common.SuccessResp(c, model.Link{
			URL: fmt.Sprintf("%s/p%s?d&sign=%s",
				common.GetApiUrl(c),
				utils.EncodePath(rawPath, true),
				sign.Sign(rawPath)),
		})
		return
	}
	link, _, err := fs.Link(c.Request.Context(), rawPath, model.LinkArgs{IP: c.ClientIP(), Header: c.Request.Header, Redirect: true})
	if err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	defer link.Close()
	common.SuccessResp(c, link)
}
