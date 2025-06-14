package cloudnet

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"bvtc/client"
	"bvtc/log"
	"bvtc/response"

	"github.com/chaunsin/netease-cloud-music/api/types"
	"github.com/chaunsin/netease-cloud-music/api/weapi"
	"github.com/gin-gonic/gin"
)

type ShowPlaylistResp struct {
	PName string `json:"pname"`
	Pid   int64  `json:"pid"`
}

func ShowPlaylist(ctx *gin.Context) {
	api, err := client.GetNetcloudApi()
	if err != nil {
		log.Logger.Error("client fail to init", log.Any("err : ", err))
		ctx.JSON(http.StatusInternalServerError, response.FailMsg(err.Error()))
		return
	}

	//	通过判断用户名找出自己的歌单，存入歌单id和名字
	userinfo, err := api.GetUserInfo(ctx, &weapi.GetUserInfoReq{})
	username := userinfo.Profile.Nickname
	if err != nil {
		log.Logger.Error("userinfo fail to get", log.Any("err : ", err))
		ctx.JSON(http.StatusInternalServerError, response.FailMsg(err.Error()))
		return
	}
	playlist, err := api.Playlist(ctx, &weapi.PlaylistReq{Uid: strconv.FormatInt(userinfo.Account.Id, 10)})
	if err != nil {
		log.Logger.Error("fail to get playlist", log.Any("err : ", err))
		ctx.JSON(http.StatusInternalServerError, response.FailMsg(err.Error()))
		return
	}
	var resp []ShowPlaylistResp
	for i := range playlist.Playlist {
		if playlist.Playlist[i].Creator.Nickname == username {
			resp = append(resp, ShowPlaylistResp{
				PName: playlist.Playlist[i].Name,
				Pid:   playlist.Playlist[i].Id,
			})
		}
	}
	log.Logger.Info("success to get playlist", log.Any("resp", resp))
	ctx.JSON(http.StatusOK, response.SuccessMsg(resp))
}

type UploadToMusicReq struct {
	Pid      int64
	TrackIds int64
}

func UploadToPlaylist(req UploadToMusicReq) error {
	api, _ := client.GetNetcloudApi()
	resp, err := api.PlaylistAddOrDel(context.Background(), &weapi.PlaylistAddOrDelReq{Op: "add", Pid: req.Pid, TrackIds: types.IntsString{req.TrackIds}, Imme: true})
	if err != nil {
		log.Logger.Error("fail", log.Any("err", err))
		return errors.New("fail to upload to playlist")
	}
	if resp.Code != 200 {
		log.Logger.Error("fail to upload to playlist", log.Any("resp", resp))
		return errors.New("fail to upload to playlist")
	}
	log.Logger.Info("success to upload to playlist", log.Any("resp", resp))
	return nil
}
