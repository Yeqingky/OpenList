package fs

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/OpenListTeam/OpenList/v4/internal/driver"
	"github.com/OpenListTeam/OpenList/v4/internal/model"
	"github.com/OpenListTeam/OpenList/v4/internal/op"
)

// the param named path of functions in this package is a mount path
// So, the purpose of this package is to convert mount path to actual path
// then pass the actual path to the op package

type ListArgs struct {
	Refresh            bool
	NoLog              bool
	WithStorageDetails bool
}

func List(ctx context.Context, path string, args *ListArgs) ([]model.Obj, error) {
	res, err := list(ctx, path, args)
	if err != nil {
		if !args.NoLog {
			log.Errorf("failed list %s: %+v", path, err)
		}
		return nil, err
	}
	return res, nil
}

type GetArgs struct {
	NoLog              bool
	WithStorageDetails bool
}

func Get(ctx context.Context, path string, args *GetArgs) (model.Obj, error) {
	res, err := get(ctx, path, args)
	if err != nil {
		if !args.NoLog {
			log.Warnf("failed get %s: %s", path, err)
		}
		return nil, err
	}
	return res, nil
}

func Link(ctx context.Context, path string, args model.LinkArgs) (*model.Link, model.Obj, error) {
	res, file, err := link(ctx, path, args)
	if err != nil {
		log.Errorf("failed link %s: %+v", path, err)
		return nil, nil, err
	}
	return res, file, nil
}

func Remove(ctx context.Context, path string) error {
	err := remove(ctx, path)
	if err != nil {
		log.Errorf("failed remove %s: %+v", path, err)
	}
	return err
}

func PutDirectly(ctx context.Context, dstDirPath string, file model.FileStreamer, skipHook ...bool) error {
	err := putDirectly(ctx, dstDirPath, file, skipHook...)
	if err != nil {
		log.Errorf("failed put %s: %+v", dstDirPath, err)
	}
	return err
}

type GetStoragesArgs struct {
}

func GetStorage(path string, args *GetStoragesArgs) (driver.Driver, error) {
	storageDriver, _, err := op.GetStorageAndActualPath(path)
	if err != nil {
		return nil, err
	}
	return storageDriver, nil
}

func Other(ctx context.Context, args model.FsOtherArgs) (interface{}, error) {
	res, err := other(ctx, args)
	if err != nil {
		log.Errorf("failed get other %s: %+v", args.Path, err)
	}
	return res, err
}

func GetDirectUploadInfo(ctx context.Context, tool, path, dstName string, fileSize int64, overwrite bool) (any, error) {
	info, err := getDirectUploadInfo(ctx, tool, path, dstName, fileSize, overwrite)
	if err != nil {
		log.Errorf("failed get %s direct upload info for %s(%d bytes): %+v", path, dstName, fileSize, err)
	}
	return info, err
}
