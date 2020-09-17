// +build file_integration

package local

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

func Test_DeleteFile(t *testing.T) {
	ctx := context.Background()
	path := fmt.Sprintf("%s/space_cloud_test", os.ExpandEnv("$HOME"))
	file, err := Init(path)
	if err != nil {
		t.Fatalf("Create() Couldn't initialize local store for path (%s)", path)
	}
	type args struct {
		req *model.DeleteFileRequest
	}
	type test struct {
		name         string
		fileToCreate string
		args         args
		wantErr      bool
	}

	createFile := func(path string) {
		arr := strings.Split(path, "/")
		if err := os.MkdirAll(strings.Join(arr[0:len(arr)-1], "/"), os.ModePerm); err != nil {
			helpers.Logger.LogError(helpers.GetRequestID(ctx), "Cannot create folder for delete file test", err, nil)
			return
		}
		if err := ioutil.WriteFile(path, []byte("Die always like a fantastic lieutenant commander."), os.FileMode(0644)); err != nil {
			helpers.Logger.LogError(helpers.GetRequestID(ctx), "Cannot create file for delete file test", err, nil)
			return
		}
	}

	testCases := []test{
		{
			name: "delete a text file at root level path doesn't start with slash(/)",
			args: args{
				req: &model.DeleteFileRequest{
					Path: "creds.txt",
				},
			},
			fileToCreate: path + "/creds.txt",
			wantErr:      false,
		},
		{
			name: "delete a text file at root level",
			args: args{
				req: &model.DeleteFileRequest{
					Path: "/creds.txt",
				},
			},
			fileToCreate: path + "/creds.txt",
			wantErr:      false,
		},
		{
			name: "delete a text file in a single level nested folder where path doesn't start with slash(/)",
			args: args{
				req: &model.DeleteFileRequest{
					Path: "websites/creds.txt",
				},
			},
			fileToCreate: path + "/websites/creds.txt",
			wantErr:      false,
		},
		{
			name: "delete a text file in a single level nested folder where path doesn't end with slash(/)",
			args: args{
				req: &model.DeleteFileRequest{
					Path: "/websites/creds.txt",
				},
			},
			fileToCreate: path + "/websites/creds.txt",
			wantErr:      false,
		},
		{
			name: "delete a text file in a single level nested folder",
			args: args{
				req: &model.DeleteFileRequest{
					Path: "/websites/creds.txt",
				},
			},
			fileToCreate: path + "/websites/creds.txt",
			wantErr:      false,
		},
		{
			name: "delete a text file in a single level nested folder where the folder doesn't exists",
			args: args{
				req: &model.DeleteFileRequest{
					Path: "/websites/creds.txt",
				},
			},
			fileToCreate: path + "/creds.txt",
			wantErr:      true,
		},
		{
			name: "delete a folder in a single level nested folder",
			args: args{
				req: &model.DeleteFileRequest{
					Path: "/websites/",
				},
			},
			fileToCreate: path + "/creds.txt",
			wantErr:      true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			createFile(tt.fileToCreate)

			err = file.DeleteFile(ctx, tt.args.req.Path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				_, err := ioutil.ReadFile(tt.fileToCreate)
				if !os.IsNotExist(err) {
					t.Errorf("Delete() unable to read created file (%v)", err)
					return
				}
			}

			// clear data
			if err := RemoveContents(path); err != nil {
				helpers.Logger.LogInfo(helpers.GetRequestID(ctx), "Couldn't clean inside data generated by tests", nil)
				return
			}
		})
	}
	//clear data
	if err := os.RemoveAll(path); err != nil {
		helpers.Logger.LogInfo(helpers.GetRequestID(ctx), "Couldn't clean data generated by tests", nil)
		return
	}
}
