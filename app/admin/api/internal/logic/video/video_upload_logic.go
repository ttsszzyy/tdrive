package video

import (
	"T-driver/app/admin/api/internal/svc"
	"T-driver/app/admin/api/internal/types"
	"T-driver/app/user/model"
	"T-driver/app/video/rpc/videoPb"
	"T-driver/common/utils"
	"bytes"
	"context"
	"fmt"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	"github.com/zeromicro/go-zero/core/logx"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type VideoUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVideoUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoUploadLogic {
	return &VideoUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VideoUploadLogic) VideoUpload(req *types.VideoUploadReq) (resp *types.VideoUploadResp, err error) {
	id := l.ctx.Value("id").(int64)
	video, err := l.svcCtx.VideoRpc.CreateVideo(l.ctx, &videoPb.CreateVideoReq{
		Uid:      id,
		Title:    req.Title,
		Filename: req.Filename,
		Desc:     req.Desc,
		Status:   2,
	})
	if err != nil {
		return nil, err
	}
	//文件转成m3u8格式保存到服务器
	go func(svcCtx *svc.ServiceContext, id int64, req *types.VideoUploadReq) {
		//保存到Titan
		reader := bytes.NewReader(req.File)
		progress := func(doneSize int64, totalSize int64) {
			// 如果上传完成，从Redis中删除上传进度记录，并打印成功信息。
			if doneSize == totalSize {
				_, err := svcCtx.Redis.Del(fmt.Sprintf(model.UploadVideoId+"%v", id))
				if err != nil {
					logx.Error("redis删除上传进度失败", err)
				}
			}
			// 更新Redis中的上传进度。
			err := svcCtx.Redis.Set(fmt.Sprintf(model.UploadVideoId+"%v", id), strconv.Itoa(int(float64(doneSize)/float64(totalSize)*100)))
			if err != nil {
				logx.Error("redis更新上传进度失败", err)
			}
		}
		url := ""
		uploadStream, err := svcCtx.Storage.UploadStreamV2(context.TODO(), reader, req.Filename, progress)
		if err != nil {
			logx.Error("上传失败", err)
			svcCtx.VideoRpc.UpdateVideo(context.TODO(), &videoPb.UpdateVideoReq{
				Id:     id,
				Status: 4,
			})
			return
		}
		shareAssetResult, err := svcCtx.Storage.GetURL(context.TODO(), uploadStream.String())
		if err != nil {
			logx.Error("GetURL:", err)
			svcCtx.VideoRpc.UpdateVideo(context.TODO(), &videoPb.UpdateVideoReq{
				Id:     id,
				Status: 4,
			})
			return
		}
		if shareAssetResult != nil && len(shareAssetResult.URLs) > 0 {
			url = shareAssetResult.URLs[0]
		}
		//视频转成m3u8
		fileName, err := utils.GenerateRandomFileName()
		if err != nil {
			logx.Error("Error generating random file name:", err)
			svcCtx.VideoRpc.UpdateVideo(context.TODO(), &videoPb.UpdateVideoReq{
				Id:     id,
				Status: 4,
			})
			return
		}
		name := fileName + ".m3u8"
		//文件转换
		outputDir := "./static/path/hls"
		if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
			logx.Error("failed to create output directory: ", err)
			svcCtx.VideoRpc.UpdateVideo(context.TODO(), &videoPb.UpdateVideoReq{
				Id:     id,
				Status: 4,
			})
			return
		}
		// Perform HLS transcoding
		err = ffmpeg_go.Input(url).Output(filepath.Join(outputDir, name), ffmpeg_go.KwArgs{
			"c:v":           "libx264",
			"c:a":           "aac",
			"hls_time":      5,
			"hls_list_size": 0,
			"f":             "hls",
		}).OverWriteOutput().ErrorToStdOut().Run()
		if err != nil {
			logx.Error("failed to transcode video: ", err)
			svcCtx.VideoRpc.UpdateVideo(context.TODO(), &videoPb.UpdateVideoReq{
				Id:     id,
				Status: 4,
			})
			return
		}

		files, err := ioutil.ReadDir(outputDir)
		if err != nil {
			logx.Error("failed to read output directory: ", err)
			svcCtx.VideoRpc.UpdateVideo(context.TODO(), &videoPb.UpdateVideoReq{
				Id:     id,
				Status: 4,
			})
			return
		}
		OutputUrl := ""
		for _, file := range files {
			if strings.HasPrefix(file.Name(), fileName) {
				//上传文件
				location, err := svcCtx.BucketUpload(filepath.Join(outputDir, file.Name()), file.Name())
				if err != nil {
					logx.Error("Error uploading file:", err)
					svcCtx.VideoRpc.UpdateVideo(context.TODO(), &videoPb.UpdateVideoReq{
						Id:     id,
						Status: 4,
					})
					return
				}
				if strings.HasSuffix(file.Name(), ".m3u8") {
					OutputUrl = location
				}
			}
		}

		/*OutputUrl := svcCtx.Config.AwsS3.EndPoint + "/" + svcCtx.Config.AwsS3.Bucket + "/" + name

		err = ffmpeg_go.Input(url).Output(OutputUrl, ffmpeg_go.KwArgs{
			"hls_time":      "5",
			"hls_list_size": "0",
			"c:v":           "libx264",
			"c:a":           "aac",
			"f":             "hls",
			"aws_config": &aws.Config{
				Credentials: credentials.NewStaticCredentials(svcCtx.Config.AwsS3.AccessKey, svcCtx.Config.AwsS3.SecretKey, ""),
				Endpoint:    aws.String(svcCtx.Config.AwsS3.EndPoint),
				Region:      aws.String(svcCtx.Config.AwsS3.Region),
				DisableSSL:  aws.Bool(true),
			},
		}).OverWriteOutput().Run()
		if err != nil {
			logx.Error("Error converting MP4 to M3U8:", err)
			svcCtx.VideoRpc.UpdateVideo(context.TODO(), &videoPb.UpdateVideoReq{
				Id:     id,
				Status: 4,
			})
			return
		}*/

		_, err = svcCtx.VideoRpc.UpdateVideo(context.TODO(), &videoPb.UpdateVideoReq{
			Id:       id,
			Url:      url,
			Cid:      uploadStream.String(),
			FilePath: OutputUrl,
			Status:   3,
		})
		if err != nil {
			logx.Error("Error updating video record:", err)
			return
		}
	}(l.svcCtx, video.Id, req)
	return &types.VideoUploadResp{
		Id: video.Id,
	}, nil
}
