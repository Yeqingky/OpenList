package op_test

import (
	"context"
	"testing"

	"github.com/OpenListTeam/OpenList/v4/internal/conf"
	"github.com/OpenListTeam/OpenList/v4/internal/db"
	"github.com/OpenListTeam/OpenList/v4/internal/model"
	"github.com/OpenListTeam/OpenList/v4/internal/op"
	"github.com/OpenListTeam/OpenList/v4/pkg/utils"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func init() {
	dB, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	conf.Conf = conf.DefaultConfig("data")
	db.Init(dB)
}

func TestCreateStorage(t *testing.T) {
	var storages = []struct {
		storage model.Storage
		isErr   bool
	}{
		{storage: model.Storage{Driver: "Template", MountPath: "/local", Addition: `{"field":"a"}`}, isErr: false},
		{storage: model.Storage{Driver: "Template", MountPath: "/local", Addition: `{"field":"a"}`}, isErr: true},
		{storage: model.Storage{Driver: "None", MountPath: "/none", Addition: `{"field":"a"}`}, isErr: true},
	}
	for _, storage := range storages {
		_, err := op.CreateStorage(context.Background(), storage.storage)
		if err != nil {
			if !storage.isErr {
				t.Errorf("failed to create storage: %+v", err)
			} else {
				t.Logf("expect failed to create storage: %+v", err)
			}
		}
	}
}

func TestGetStorageVirtualFilesByPath(t *testing.T) {
	setupStorages(t)
	virtualFiles := op.GetStorageVirtualFilesByPath("/a")
	var names []string
	for _, virtualFile := range virtualFiles {
		names = append(names, virtualFile.GetName())
	}
	var expectedNames = []string{"b", "c", "d"}
	if utils.SliceEqual(names, expectedNames) {
		t.Logf("passed")
	} else {
		t.Errorf("expected: %+v, got: %+v", expectedNames, names)
	}
}

func TestGetBalancedStorage(t *testing.T) {
	set := mapset.NewSet[string]()
	for i := 0; i < 5; i++ {
		storage := op.GetBalancedStorage("/a/d/e1")
		set.Add(storage.GetStorage().MountPath)
	}
	expected := mapset.NewSet([]string{"/a/d/e1", "/a/d/e1.balance"}...)
	if !expected.Equal(set) {
		t.Errorf("expected: %+v, got: %+v", expected, set)
	}
}

func setupStorages(t *testing.T) {
	var storages = []model.Storage{
		{Driver: "Template", MountPath: "/a/b", Order: 0, Addition: `{"field":"a"}`},
		{Driver: "Template", MountPath: "/adc", Order: 0, Addition: `{"field":"a"}`},
		{Driver: "Template", MountPath: "/a/c", Order: 1, Addition: `{"field":"a"}`},
		{Driver: "Template", MountPath: "/a/d", Order: 2, Addition: `{"field":"a"}`},
		{Driver: "Template", MountPath: "/a/d/e1", Order: 3, Addition: `{"field":"a"}`},
		{Driver: "Template", MountPath: "/a/d/e", Order: 4, Addition: `{"field":"a"}`},
		{Driver: "Template", MountPath: "/a/d/e1.balance", Order: 4, Addition: `{"field":"a"}`},
	}
	for _, storage := range storages {
		_, err := op.CreateStorage(context.Background(), storage)
		if err != nil {
			t.Fatalf("failed to create storage: %+v", err)
		}
	}
}
