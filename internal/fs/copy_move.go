package fs

import (
	"context"
	stdpath "path"

	"github.com/OpenListTeam/OpenList/v4/internal/conf"
	"github.com/OpenListTeam/OpenList/v4/internal/driver"
	"github.com/OpenListTeam/OpenList/v4/internal/errs"
	"github.com/OpenListTeam/OpenList/v4/internal/model"
	"github.com/OpenListTeam/OpenList/v4/internal/op"
	"github.com/OpenListTeam/OpenList/v4/internal/stream"
	"github.com/OpenListTeam/OpenList/v4/pkg/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// moveObj moves srcObjPath into dstDirPath synchronously.
// The task based transfer was removed together with the task system, so
// integrations that still need a move (e.g. the FTP server) run it inline.
func moveObj(ctx context.Context, srcObjPath, dstDirPath string, skipHook ...bool) error {
	srcStorage, srcActualPath, err := op.GetStorageAndActualPath(srcObjPath)
	if err != nil {
		return errors.WithMessage(err, "failed get src storage")
	}
	dstStorage, dstActualPath, err := op.GetStorageAndActualPath(dstDirPath)
	if err != nil {
		return errors.WithMessage(err, "failed get dst storage")
	}
	if utils.IsBool(skipHook...) {
		ctx = context.WithValue(ctx, conf.SkipHookKey, struct{}{})
	}

	// fast path: the driver is able to move within a single storage
	if srcStorage.GetStorage() == dstStorage.GetStorage() {
		err = op.Move(ctx, srcStorage, srcActualPath, dstActualPath)
		if !errors.Is(err, errs.NotImplement) && !errors.Is(err, errs.NotSupport) {
			return err
		}
	}
	// cross storage: copy first, then drop the source
	if err := copyObj(ctx, srcStorage, srcActualPath, dstStorage, dstActualPath); err != nil {
		return err
	}
	return removeMovedSource(ctx, srcStorage, srcActualPath, dstStorage, dstActualPath)
}

// copyObj copies srcActualPath into dstParentPath, which must be a directory.
func copyObj(ctx context.Context, srcStorage driver.Driver, srcActualPath string, dstStorage driver.Driver, dstParentPath string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	srcObj, err := op.Get(ctx, srcStorage, srcActualPath)
	if err != nil {
		return errors.WithMessagef(err, "failed get src [%s]", srcActualPath)
	}
	if !srcObj.IsDir() {
		return copyFileObj(ctx, srcStorage, srcActualPath, dstStorage, dstParentPath)
	}
	dstDirPath := stdpath.Join(dstParentPath, srcObj.GetName())
	if err := op.MakeDir(ctx, dstStorage, dstDirPath); err != nil {
		return errors.WithMessagef(err, "failed make dir [%s]", dstDirPath)
	}
	objs, err := op.List(ctx, srcStorage, srcActualPath, model.ListArgs{})
	if err != nil {
		return errors.WithMessagef(err, "failed list src [%s] objs", srcActualPath)
	}
	for _, obj := range objs {
		if err := ctx.Err(); err != nil {
			return err
		}
		err = copyObj(ctx, srcStorage, stdpath.Join(srcActualPath, obj.GetName()), dstStorage, dstDirPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func copyFileObj(ctx context.Context, srcStorage driver.Driver, srcActualPath string, dstStorage driver.Driver, dstParentPath string) error {
	link, srcObj, err := op.Link(ctx, srcStorage, srcActualPath, model.LinkArgs{})
	if err != nil {
		return errors.WithMessagef(err, "failed get [%s] link", srcActualPath)
	}
	// any link provided is seekable
	ss, err := stream.NewSeekableStream(&stream.FileStream{
		Obj: srcObj,
		Ctx: ctx,
	}, link)
	if err != nil {
		_ = link.Close()
		return errors.WithMessagef(err, "failed get [%s] stream", srcActualPath)
	}
	return op.Put(context.WithValue(ctx, conf.SkipHookKey, struct{}{}), dstStorage, dstParentPath, ss, nil)
}

// removeMovedSource removes the source of a finished cross storage move. Every
// object is only removed after its destination counterpart has been verified,
// so a partially failed move never loses data.
func removeMovedSource(ctx context.Context, srcStorage driver.Driver, srcActualPath string, dstStorage driver.Driver, dstParentPath string) error {
	srcObj, err := op.GetUnwrap(ctx, srcStorage, srcActualPath)
	if err != nil {
		return errors.WithMessagef(err, "failed get src [%s]", srcActualPath)
	}
	dstObjPath := stdpath.Join(dstParentPath, srcObj.GetName())
	dstObj, err := op.GetUnwrap(ctx, dstStorage, dstObjPath)
	if err != nil {
		return errors.WithMessagef(err, "failed get dst [%s]", dstObjPath)
	}
	if !dstObj.IsDir() {
		return removeSrcObj(ctx, srcStorage, srcActualPath)
	}

	srcObjs, err := op.List(ctx, srcStorage, srcActualPath, model.ListArgs{})
	if err != nil {
		return errors.WithMessagef(err, "failed list src [%s] objs", srcActualPath)
	}
	hasErr := false
	for _, obj := range srcObjs {
		err := removeMovedSource(ctx, srcStorage, stdpath.Join(srcActualPath, obj.GetName()), dstStorage, dstObjPath)
		if err != nil {
			log.Error(err)
			hasErr = true
		}
	}
	if hasErr {
		return errors.Errorf("some subitems of [%s] failed to verify and remove", srcActualPath)
	}
	return removeSrcObj(ctx, srcStorage, srcActualPath)
}

func removeSrcObj(ctx context.Context, srcStorage driver.Driver, srcActualPath string) error {
	if err := op.Remove(ctx, srcStorage, srcActualPath); err != nil {
		return errors.WithMessagef(err, "failed remove src [%s]", srcActualPath)
	}
	return nil
}
