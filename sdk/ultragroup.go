package sdk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/httplib"
)

const (
	UGUnPushLevelAllMessage        = -1 // UGUnPushLevelAllMessage 全部消息通知
	UGUnPushLevelNotSet            = 0  // UGUnPushLevelNotSet 未设置
	UGUnPushLevelAtMessage         = 1  // UGUnPushLevelAtMessage 仅@消息通知
	UGUnPushLevelAtUser            = 2  // UGUnPushLevelAtUser @指定用户通知
	UGUnPushLevelAtAllGroupMembers = 4  // UGUnPushLevelAtAllGroupMembers @群全员通知
	UGUnPushLevelNotRecv           = 5  // UGUnPushLevelNotRecv 不接收通知
)

// api 返回结果, data 数组
type RespDataArray struct {
	Code int                                 `json:"code"`
	Data map[string][]map[string]interface{} `json:"data"`
}

// api 返回结果, data
type RespDataKV struct {
	Code int                    `json:"code"`
	Data map[string]interface{} `json:"data"`
}

// 超级群 群组信息
type UGGroupInfo struct {
	GroupId   string `json:"group_id"`
	GroupName string `json:"group_name"`
}

// 超级群 用户信息
type UGUserInfo struct {
	Id        string `json:"id"`
	MutedTime string `json:"time"`
}

// 超级群 频道信息
type UGChannelInfo struct {
	ChannelId  string `json:"channel_id"`
	CreateTime string `json:"create_time"`
}

// 超级群 消息结构
type UGMessage struct {
	FromUserId          string   `json:"from_user_id"`
	ToGroupIds          []string `json:"to_group_ids"`
	ToUserIds           []string `json:"to_user_ids,omitempty"`
	ObjectName          string   `json:"object_name"`
	Content             string   `json:"content"`
	PushContent         string   `json:"push_content,omitempty"`
	PushData            string   `json:"push_data,omitempty"`
	IncludeSenderEnable bool     `json:"include_sender_enable,omitempty"`
	StoreFlag           bool     `json:"store_flag"`
	MentionedFlag       bool     `json:"mentioned_flag,omitempty"`
	SilencePush         bool     `json:"silence_push,omitempty"`
	PushExt             string   `json:"push_ext,omitempty"`
	BusChannel          string   `json:"bus_channel,omitempty"`
}

// 创建群组
func (rc *RongCloud) UGGroupCreate(userId, groupId, groupName string) (err error, requestId string) {
	if userId == "" {
		return RCErrorNewV2(1002, "param 'userId' is required"), ""
	}

	if groupId == "" {
		return RCErrorNewV2(1002, "param 'groupId' is required"), ""
	}

	if groupName == "" {
		return RCErrorNewV2(1002, "param 'groupName' is required"), ""
	}

	url := fmt.Sprintf("%s/v2/ultragroups", rc.rongCloudURI)
	req := httplib.Post(url)
	req.SetTimeout(time.Second*rc.timeout, time.Second*rc.timeout)
	requestId = rc.fillHeaderV2(req)

	// json body
	postBody := map[string]interface{}{"user_id": userId,
		"group_id":   groupId,
		"group_name": groupName,
	}
	_, err = req.JSONBody(postBody)
	if err != nil {
		return err, ""
	}

	// http
	_, err = rc.doV2(req)

	return err, requestId
}

// 解散群组
func (rc *RongCloud) UGGroupDismiss(groupId string) (err error, requestId string) {

	if groupId == "" {
		return RCErrorNewV2(1002, "param 'groupId' is required"), ""
	}

	url := fmt.Sprintf("%s/v2/ultragroups/%s", rc.rongCloudURI, groupId)
	req := httplib.Delete(url)
	req.SetTimeout(time.Second*rc.timeout, time.Second*rc.timeout)
	requestId = rc.fillHeaderV2(req)

	// http
	_, err = rc.doV2(req)

	return err, requestId
}

// 加入群组
func (rc *RongCloud) UGGroupJoin(userId, groupId string) (err error, requestId string) {
	if userId == "" {
		return RCErrorNewV2(1002, "param 'userId' is required"), ""
	}

	if groupId == "" {
		return RCErrorNewV2(1002, "param 'groupId' is required"), ""
	}

	url := fmt.Sprintf("%s/v2/ultragroups/%s/users/%s", rc.rongCloudURI, groupId, userId)
	req := httplib.Post(url)
	req.SetTimeout(time.Second*rc.timeout, time.Second*rc.timeout)
	requestId = rc.fillHeaderV2(req)

	// http
	_, err = rc.doV2(req)

	return err, requestId
}

// 退出群组
func (rc *RongCloud) UGGroupQuit(userId, groupId string) (err error, requestId string) {
	if userId == "" {
		return RCErrorNewV2(1002, "param 'userId' is required"), ""
	}

	if groupId == "" {
		return RCErrorNewV2(1002, "param 'groupId' is required"), ""
	}

	url := fmt.Sprintf("%s/v2/ultragroups/%s/users/%s", rc.rongCloudURI, groupId, userId)
	req := httplib.Delete(url)
	req.SetTimeout(time.Second*rc.timeout, time.Second*rc.timeout)
	requestId = rc.fillHeaderV2(req)

	// http
	_, err = rc.doV2(req)

	return err, requestId
}

// 刷新群组信息
func (rc *RongCloud) UGGroupUpdate(groupId, groupName string) (err error, requestId string) {
	if groupId == "" {
		return RCErrorNewV2(1002, "param 'groupId' is required"), ""
	}
	if groupName == "" {
		return RCErrorNewV2(1002, "param 'groupName' is required"), ""
	}

	url := fmt.Sprintf("%s/v2/ultragroups/%s", rc.rongCloudURI, groupId)
	req := httplib.Put(url)
	req.SetTimeout(time.Second*rc.timeout, time.Second*rc.timeout)
	requestId = rc.fillHeaderV2(req)

	// json body
	postBody := map[string]interface{}{
		"group_name": groupName,
	}
	_, err = req.JSONBody(postBody)
	if err != nil {
		return err, ""
	}

	// http
	_, err = rc.doV2(req)

	return err, requestId
}

// 查询用户所在群组(P1)
func (rc *RongCloud) UGQueryUserGroups(userId string, page, size int) (groups []UGGroupInfo, err error, requestId string) {
	if userId == "" {
		return nil, RCErrorNewV2(1002, "param 'userId' is required"), ""
	}

	url := fmt.Sprintf("%s/v2/ultragroups/users/%s/groups", rc.rongCloudURI, userId)
	req := httplib.Get(url)
	req.SetTimeout(time.Second*rc.timeout, time.Second*rc.timeout)
	requestId = rc.fillHeaderV2(req)

	req.Param("page", strconv.Itoa(page))
	req.Param("size", strconv.Itoa(size))

	// http
	respBody, err := rc.doV2(req)
	if err != nil {
		return groups, err, requestId
	}

	// 处理返回结果
	var respJson RespDataArray
	if err := json.Unmarshal(respBody, &respJson); err != nil {
		return groups, err, requestId
	}

	if respJson.Data != nil {
		if gs, ok := respJson.Data["groups"]; ok {
			if gs != nil {
				for _, v := range gs {
					t := UGGroupInfo{GroupId: fmt.Sprint(v["group_id"]),
						GroupName: fmt.Sprint(v["group_name"])}
					groups = append(groups, t)
				}
			}
		}
	}

	return groups, err, requestId
}

// 查询群成员(P1)
func (rc *RongCloud) UGQueryGroupUsers(groupId string, page, size int) (users []UGUserInfo, err error, requestId string) {
	if groupId == "" {
		return nil, RCErrorNewV2(1002, "param 'groupId' is required"), ""
	}

	url := fmt.Sprintf("%s/v2/ultragroups/%s/users", rc.rongCloudURI, groupId)
	req := httplib.Get(url)
	req.SetTimeout(time.Second*rc.timeout, time.Second*rc.timeout)
	requestId = rc.fillHeaderV2(req)

	req.Param("page", strconv.Itoa(page))
	req.Param("size", strconv.Itoa(size))

	// http
	respBody, err := rc.doV2(req)
	if err != nil {
		return users, err, requestId
	}

	// 处理返回结果
	var respJson RespDataArray
	if err := json.Unmarshal(respBody, &respJson); err != nil {
		return users, err, requestId
	}

	if respJson.Data != nil {
		if gs, ok := respJson.Data["users"]; ok {
			if gs != nil {
				for _, v := range gs {
					t := UGUserInfo{Id: fmt.Sprint(v["id"])}
					users = append(users, t)
				}
			}
		}
	}

	return users, err, requestId
}

// 消息发送 普通消息
func (rc *RongCloud) UGGroupSend(msg UGMessage) (err error, requestId string) {
	if msg.FromUserId == "" {
		return RCErrorNewV2(1002, "Paramer 'FromUserId' is required"), ""
	}

	if len(msg.ToGroupIds) == 0 {
		return RCErrorNewV2(1002, "Paramer 'ToGroupIds' is required"), ""
	}

	url := fmt.Sprintf("%s/v2/message/ultragroup/send", rc.rongCloudURI)
	req := httplib.Post(url)
	req.SetTimeout(time.Second*rc.timeout, time.Second*rc.timeout)
	requestId = rc.fillHeaderV2(req)

	// json body
	_, err = req.JSONBody(msg)
	if err != nil {
		return err, ""
	}

	// http
	_, err = rc.doV2(req)

	return err, requestId
}

// 添加禁言成员
func (rc *RongCloud) UGGroupMuteMembersAdd(groupId string, userIds []string) (err error, requestId string) {
	if groupId == "" {
		return RCErrorNewV2(1002, "Paramer 'groupId' is required"), ""
	}

	if len(userIds) == 0 {
		return RCErrorNewV2(1002, "Paramer 'userIds' is required"), ""
	}

	url := fmt.Sprintf("%s/v2/ultragroups/%s/muted-users", rc.rongCloudURI, groupId)
	req := httplib.Post(url)
	req.SetTimeout(time.Second*rc.timeout, time.Second*rc.timeout)
	requestId = rc.fillHeaderV2(req)

	// json body
	postBody := map[string]interface{}{
		"user_ids": userIds,
	}
	_, err = req.JSONBody(postBody)
	if err != nil {
		return err, ""
	}

	// http
	_, err = rc.doV2(req)

	return err, requestId
}

// 移除禁言成员
func (rc *RongCloud) UGGroupMuteMembersRemove(groupId string, userIds []string) (err error, requestId string) {
	if groupId == "" {
		return RCErrorNewV2(1002, "Paramer 'groupId' is required"), ""
	}

	if len(userIds) == 0 {
		return RCErrorNewV2(1002, "Paramer 'userIds' is required"), ""
	}

	url := fmt.Sprintf("%s/v2/ultragroups/%s/muted-users", rc.rongCloudURI, groupId)
	req := httplib.Delete(url)
	req.SetTimeout(time.Second*rc.timeout, time.Second*rc.timeout)
	requestId = rc.fillHeaderV2(req)

	// json body
	postBody := map[string]interface{}{
		"user_ids": userIds,
	}
	_, err = req.JSONBody(postBody)
	if err != nil {
		return err, ""
	}

	// http
	_, err = rc.doV2(req)

	return err, requestId
}

// 获取禁言成员
func (rc *RongCloud) UGGroupMuteMembersGetList(groupId string) (users []UGUserInfo, err error, requestId string) {
	if groupId == "" {
		return nil, RCErrorNewV2(1002, "param 'groupId' is required"), ""
	}

	url := fmt.Sprintf("%s/v2/ultragroups/%s/muted-users", rc.rongCloudURI, groupId)
	req := httplib.Get(url)
	req.SetTimeout(time.Second*rc.timeout, time.Second*rc.timeout)
	requestId = rc.fillHeaderV2(req)

	// http
	respBody, err := rc.doV2(req)
	if err != nil {
		return users, err, requestId
	}

	// 处理返回结果
	var respJson RespDataArray
	if err := json.Unmarshal(respBody, &respJson); err != nil {
		return users, err, requestId
	}

	if respJson.Data != nil {
		if gs, ok := respJson.Data["users"]; ok {
			if gs != nil {
				for _, v := range gs {
					t := UGUserInfo{Id: fmt.Sprint(v["id"]),
						MutedTime: fmt.Sprint(v["time"])}
					users = append(users, t)
				}
			}
		}
	}

	return users, err, requestId
}

// 设置全体成员禁言
func (rc *RongCloud) UGGroupMuted(groupId string, status bool) (err error, requestId string) {
	if groupId == "" {
		return RCErrorNewV2(1002, "Paramer 'groupId' is required"), ""
	}

	url := fmt.Sprintf("%s/v2/ultragroups/%s/muted-status", rc.rongCloudURI, groupId)
	req := httplib.Put(url)
	req.SetTimeout(time.Second*rc.timeout, time.Second*rc.timeout)
	requestId = rc.fillHeaderV2(req)

	// json body
	postBody := map[string]interface{}{
		"status": status,
	}
	_, err = req.JSONBody(postBody)
	if err != nil {
		return err, ""
	}

	// http
	_, err = rc.doV2(req)

	return err, requestId
}

// 获取群全体成员禁言状态
func (rc *RongCloud) UGGroupMutedQuery(groupId string) (status bool, err error, requestId string) {
	if groupId == "" {
		return status, RCErrorNewV2(1002, "Paramer 'groupId' is required"), ""
	}

	url := fmt.Sprintf("%s/v2/ultragroups/%s/muted-status", rc.rongCloudURI, groupId)
	req := httplib.Get(url)
	req.SetTimeout(time.Second*rc.timeout, time.Second*rc.timeout)
	requestId = rc.fillHeaderV2(req)

	// http
	respBody, err := rc.doV2(req)
	if err != nil {
		return status, err, requestId
	}

	// 处理返回结果
	var respJson RespDataKV
	if err := json.Unmarshal(respBody, &respJson); err != nil {
		return status, err, requestId
	}

	if respJson.Data != nil {
		if s, ok := respJson.Data["status"]; ok {
			status, err = strconv.ParseBool(fmt.Sprint(s))
		}
	}

	return status, err, requestId
}

// 添加禁言白名单
func (rc *RongCloud) UGGroupMutedWhitelistAdd(groupId string, userIds []string) (err error, requestId string) {
	if groupId == "" {
		return RCErrorNewV2(1002, "Paramer 'groupId' is required"), ""
	}

	if len(userIds) == 0 {
		return RCErrorNewV2(1002, "Paramer 'userIds' is required"), ""
	}

	url := fmt.Sprintf("%s/v2/ultragroups/%s/allowed-users", rc.rongCloudURI, groupId)
	req := httplib.Post(url)
	req.SetTimeout(time.Second*rc.timeout, time.Second*rc.timeout)
	requestId = rc.fillHeaderV2(req)

	// json body
	postBody := map[string]interface{}{
		"user_ids": userIds,
	}
	_, err = req.JSONBody(postBody)
	if err != nil {
		return err, ""
	}

	// http
	_, err = rc.doV2(req)

	return err, requestId
}

// 移除禁言白名单
func (rc *RongCloud) UGGroupMutedWhitelistRemove(groupId string, userIds []string) (err error, requestId string) {
	if groupId == "" {
		return RCErrorNewV2(1002, "Paramer 'groupId' is required"), ""
	}

	if len(userIds) == 0 {
		return RCErrorNewV2(1002, "Paramer 'userIds' is required"), ""
	}

	url := fmt.Sprintf("%s/v2/ultragroups/%s/allowed-users", rc.rongCloudURI, groupId)
	req := httplib.Delete(url)
	req.SetTimeout(time.Second*rc.timeout, time.Second*rc.timeout)
	requestId = rc.fillHeaderV2(req)

	// json body
	postBody := map[string]interface{}{
		"user_ids": userIds,
	}
	_, err = req.JSONBody(postBody)
	if err != nil {
		return err, ""
	}

	// http
	_, err = rc.doV2(req)

	return err, requestId
}

// 获取禁言白名单
func (rc *RongCloud) UGGroupMutedWhitelistQuery(groupId string) (users []UGUserInfo, err error, requestId string) {
	if groupId == "" {
		return nil, RCErrorNewV2(1002, "param 'groupId' is required"), ""
	}

	url := fmt.Sprintf("%s/v2/ultragroups/%s/allowed-users", rc.rongCloudURI, groupId)
	req := httplib.Get(url)
	req.SetTimeout(time.Second*rc.timeout, time.Second*rc.timeout)
	requestId = rc.fillHeaderV2(req)

	// http
	respBody, err := rc.doV2(req)
	if err != nil {
		return users, err, requestId
	}

	// 处理返回结果
	var respJson RespDataArray
	if err := json.Unmarshal(respBody, &respJson); err != nil {
		return users, err, requestId
	}

	if respJson.Data != nil {
		if gs, ok := respJson.Data["users"]; ok {
			if gs != nil {
				for _, v := range gs {
					t := UGUserInfo{Id: fmt.Sprint(v["id"])}
					users = append(users, t)
				}
			}
		}
	}

	return users, err, requestId
}

// 创建群频道
func (rc *RongCloud) UGChannelCreate(groupId, channelId string) (err error, requestId string) {
	if groupId == "" {
		return RCErrorNewV2(1002, "param 'groupId' is required"), ""
	}

	if channelId == "" {
		return RCErrorNewV2(1002, "param 'channelId' is required"), ""
	}

	url := fmt.Sprintf("%s/v2/ultragroups/channels", rc.rongCloudURI)
	req := httplib.Post(url)
	req.SetTimeout(time.Second*rc.timeout, time.Second*rc.timeout)
	requestId = rc.fillHeaderV2(req)

	body := map[string]interface{}{
		"group_id":   groupId,
		"channel_id": channelId,
	}

	if _, err = req.JSONBody(body); err != nil {
		return err, ""
	}

	if _, err = rc.doV2(req); err != nil {
		return err, ""
	}

	return nil, requestId
}

// 删除群频道
func (rc *RongCloud) UGChannelDelete(groupId, channelId string) (err error, requestId string) {
	if groupId == "" {
		return RCErrorNewV2(1002, "param 'groupId' is required"), ""
	}

	if channelId == "" {
		return RCErrorNewV2(1002, "param 'channelId' is required"), ""
	}

	url := fmt.Sprintf("%s/v2/ultragroups/%s/channels/%s", rc.rongCloudURI, groupId, channelId)
	req := httplib.Delete(url)
	req.SetTimeout(time.Second*rc.timeout, time.Second*rc.timeout)
	requestId = rc.fillHeaderV2(req)

	// http
	_, err = rc.doV2(req)

	return err, requestId
}

// 查询群频道列表
func (rc *RongCloud) UGChannelQuery(groupId string, page, size int) (channels []UGChannelInfo, err error, requestId string) {
	if groupId == "" {
		return nil, RCErrorNewV2(1002, "param 'groupId' is required"), ""
	}

	url := fmt.Sprintf("%s/v2/ultragroups/%s/channels", rc.rongCloudURI, groupId)
	req := httplib.Get(url)
	req.SetTimeout(time.Second*rc.timeout, time.Second*rc.timeout)
	requestId = rc.fillHeaderV2(req)

	req.Param("page", strconv.Itoa(page))
	req.Param("limit", strconv.Itoa(size))

	// http
	respBody, err := rc.doV2(req)
	if err != nil {
		return channels, err, requestId
	}

	// 处理返回结果
	var respJson RespDataArray
	if err := json.Unmarshal(respBody, &respJson); err != nil {
		return channels, err, requestId
	}

	if respJson.Data != nil {
		if gs, ok := respJson.Data["channel_list"]; ok {
			if gs != nil {
				for _, v := range gs {
					t := UGChannelInfo{ChannelId: fmt.Sprint(v["channel_id"]),
						CreateTime: fmt.Sprint(v["create_time"])}
					channels = append(channels, t)
				}
			}
		}
	}

	return channels, err, requestId
}

// UGMessageExpansionSet 设置扩展
func (rc *RongCloud) UGMessageExpansionSet(groupId, userId, msgUID, busChannel string, extra map[string]string) error {
	if groupId == "" {
		return RCErrorNewV2(1002, "param 'groupId' is required")
	}

	if userId == "" {
		return RCErrorNewV2(1002, "param 'userId' is required")
	}

	if msgUID == "" {
		return RCErrorNewV2(1002, "param 'msgUID' is required")
	}

	if extra == nil {
		return RCErrorNewV2(1002, "param 'extra' is required")
	}

	if len(extra) > 100 {
		return RCErrorNewV2(1002, "param 'extra' is too long")
	}

	encExtra, err := json.Marshal(extra)
	if err != nil {
		return err
	}

	req := httplib.Post(rc.rongCloudURI + "/ultragroup/message/expansion/set." + ReqType)
	req.SetTimeout(time.Second*rc.timeout, time.Second*rc.timeout)
	rc.fillHeader(req)

	req.Param("msgUID", msgUID)
	req.Param("userId", userId)
	req.Param("groupId", groupId)
	req.Param("extraKeyVal", string(encExtra))

	if busChannel != "" {
		req.Param("busChannel", busChannel)
	}

	if _, err = rc.doV2(req); err != nil {
		return err
	}

	return nil
}

// UGMessageExpansionDelete 删除扩展
func (rc *RongCloud) UGMessageExpansionDelete(groupId, userId, msgUID, busChannel string, keys ...string) error {
	if groupId == "" {
		return RCErrorNewV2(1002, "param 'groupId' is required")
	}

	if userId == "" {
		return RCErrorNewV2(1002, "param 'userId' is required")
	}

	if msgUID == "" {
		return RCErrorNewV2(1002, "param 'msgUID' is required")
	}

	if klens := len(keys); klens <= 0 || klens > 100 {
		return RCErrorNewV2(1002, "invalid param keys")
	}

	encKeys, err := json.Marshal(keys)
	if err != nil {
		return err
	}

	req := httplib.Post(rc.rongCloudURI + "/ultragroup/message/expansion/delete." + ReqType)
	req.SetTimeout(time.Second*rc.timeout, time.Second*rc.timeout)
	rc.fillHeader(req)

	req.Param("msgUID", msgUID)
	req.Param("userId", userId)
	req.Param("groupId", groupId)
	req.Param("extraKey", string(encKeys))

	if busChannel != "" {
		req.Param("busChannel", busChannel)
	}

	if _, err = rc.doV2(req); err != nil {
		return err
	}

	return nil
}

type UGMessageExpansionItem struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	Timestamp int64  `json:"timestamp"`
}

// UGMessageExpansionQuery 获取扩展信息
func (rc *RongCloud) UGMessageExpansionQuery(groupId, msgUID, busChannel string) ([]UGMessageExpansionItem, error) {
	if groupId == "" {
		return nil, RCErrorNewV2(1002, "param 'groupId' is required")
	}

	if msgUID == "" {
		return nil, RCErrorNewV2(1002, "param 'msgUID' is required")
	}

	req := httplib.Post(rc.rongCloudURI + "/ultragroup/message/expansion/query." + ReqType)
	req.SetTimeout(time.Second*rc.timeout, time.Second*rc.timeout)
	rc.fillHeader(req)

	req.Param("msgUID", msgUID)
	req.Param("groupId", groupId)

	if busChannel != "" {
		req.Param("busChannel", busChannel)
	}

	body, err := rc.doV2(req)

	if err != nil {
		return nil, err
	}

	resp := struct {
		Code         int                               `json:"code"`
		ExtraContent map[string]map[string]interface{} `json:"extraContent"`
	}{}
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		return nil, fmt.Errorf("Response error. code: %d", resp.Code)
	}

	var data []UGMessageExpansionItem
	for key, val := range resp.ExtraContent {
		item := UGMessageExpansionItem{
			Key: key,
		}

		if v, ok := val["v"]; ok {
			item.Value = v.(string)
		}

		if ts, ok := val["ts"]; ok {
			item.Timestamp = int64(ts.(float64))
		}

		data = append(data, item)
	}

	return data, nil
}

type PushExt struct {
	Title                string                         `json:"title,omitempty"`
	TemplateId           string                         `json:"templateId,omitempty"`
	ForceShowPushContent int                            `json:"forceShowPushContent,omitempty"`
	PushConfigs          []map[string]map[string]string `json:"pushConfigs,omitempty"`
}

// UGMessagePublish 发送消息时设置扩展
func (rc *RongCloud) UGMessagePublish(fromUserId, objectName, content, pushContent, pushData, isPersisted, isMentioned, contentAvailable, busChannel, extraContent string, expansion, unreadCountFlag bool, pushExt *PushExt, toGroupIds ...string) error {
	if fromUserId == "" {
		return RCErrorNewV2(1002, "param 'fromUserId' is required")
	}

	if objectName == "" {
		return RCErrorNewV2(1002, "param 'objectName' is required")
	}

	if content == "" {
		return RCErrorNewV2(1002, "param 'content' is required")
	}

	if groupLen := len(toGroupIds); groupLen <= 0 || groupLen > 3 {
		return RCErrorNewV2(1002, "invalid 'toGroupIds'")
	}

	req := httplib.Post(rc.rongCloudURI + "/message/ultragroup/publish." + ReqType)
	req.SetTimeout(time.Second*rc.timeout, time.Second*rc.timeout)
	rc.fillHeader(req)

	body := map[string]interface{}{
		"fromUserId": fromUserId,
		"toGroupIds": toGroupIds,
		"objectName": objectName,
		"content":    content,
		"expansion":  expansion,
	}

	if pushContent != "" {
		body["pushContent"] = pushContent
	}

	if pushData != "" {
		body["pushData"] = pushData
	}

	if isPersisted != "" {
		body["isPersisted"] = isPersisted
	}

	if isMentioned != "" {
		body["isMentioned"] = isMentioned
	}

	if isMentioned != "1" {
		body["unreadCountFlag"] = unreadCountFlag
	}

	if contentAvailable != "" {
		body["contentAvailable"] = contentAvailable
	}

	if busChannel != "" {
		body["busChannel"] = busChannel
	}

	if extraContent != "" {
		body["extraContent"] = extraContent
	}

	if pushExt != nil {
		encPushExt, e := json.Marshal(pushExt)
		if e != nil {
			return e
		}
		body["pushExt"] = string(encPushExt)
	}

	var err error
	req, err = req.JSONBody(body)
	if err != nil {
		return err
	}

	if _, err = rc.doV2(req); err != nil {
		return err
	}

	return nil
}

// UGMemberExists 查询用户是否在超级群中
func (rc *RongCloud) UGMemberExists(groupId, userId string) (bool, error) {
	if groupId == "" {
		return false, RCErrorNewV2(1002, "param 'groupId' is required")
	}

	if userId == "" {
		return false, RCErrorNewV2(1002, "param 'userId' is required")
	}

	req := httplib.Post(rc.rongCloudURI + "/ultragroup/member/exist." + ReqType)
	req.SetTimeout(time.Second*rc.timeout, time.Second*rc.timeout)
	rc.fillHeader(req)

	req.Param("groupId", groupId)
	req.Param("userId", userId)

	body, err := rc.doV2(req)
	if err != nil {
		return false, err
	}

	resp := struct {
		Code   int  `json:"code"`
		Status bool `json:"status"`
	}{}

	if err = json.Unmarshal(body, &resp); err != nil {
		return false, err
	}

	if resp.Code != http.StatusOK {
		return false, fmt.Errorf("Response error. code: %d", resp.Code)
	}

	return resp.Status, nil
}

// UGNotDisturbSet 设置群/频道默认免打扰接口
func (rc *RongCloud) UGNotDisturbSet(groupId string, unPushLevel int, busChannel string) error {
	if groupId == "" {
		return RCErrorNewV2(1002, "param 'groupId' is required")
	}

	if unPushLevel != UGUnPushLevelAllMessage && unPushLevel != UGUnPushLevelNotSet && unPushLevel != UGUnPushLevelAtMessage &&
		unPushLevel != UGUnPushLevelAtUser && unPushLevel != UGUnPushLevelAtAllGroupMembers && unPushLevel != UGUnPushLevelNotRecv {
		return RCErrorNewV2(1002, "param 'unPushLevel' was wrong")
	}

	var err error

	req := httplib.Post(rc.rongCloudURI + "/ultragroup/notdisturb/set.json")

	req.Param("groupId", groupId)
	req.Param("unpushLevel", strconv.Itoa(unPushLevel))

	if busChannel != "" {
		req.Param("busChannel", busChannel)
	}

	req.SetTimeout(time.Second*rc.timeout, time.Second*rc.timeout)

	rc.fillHeader(req)

	data, err := rc.doV2(req)
	if err != nil {
		return err
	}

	resp := struct {
		Code int `json:"code"`
	}{}
	if err = json.Unmarshal(data, &resp); err != nil {
		return err
	}

	if resp.Code != http.StatusOK {
		return fmt.Errorf("Response error. code: %d", resp.Code)
	}

	return nil
}

type UGNotDisturbGetResponses struct {
	GroupId     string `json:"groupId"`
	BusChannel  string `json:"busChannel"`
	UnPushLevel int    `json:"unpushLevel"`
}

func (rc *RongCloud) UGNotDisturbGet(groupId, busChannel string) (*UGNotDisturbGetResponses, error) {
	if groupId == "" {
		return nil, RCErrorNewV2(1002, "param 'groupId' is required")
	}

	var err error

	req := httplib.Post(rc.rongCloudURI + "/ultragroup/notdisturb/get.json")

	req.Param("groupId", groupId)

	if busChannel != "" {
		req.Param("busChannel", busChannel)
	}

	req.SetTimeout(time.Second*rc.timeout, time.Second*rc.timeout)

	rc.fillHeader(req)

	data, err := rc.doV2(req)
	if err != nil {
		return nil, err
	}

	resp := struct {
		Code        int    `json:"code"`
		GroupId     string `json:"groupId"`
		BusChannel  string `json:"busChannel"`
		UnPushLevel int    `json:"unpushLevel"`
	}{}
	if err = json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	if resp.Code != http.StatusOK {
		return nil, fmt.Errorf("Response error. code: %d", resp.Code)
	}

	return &UGNotDisturbGetResponses{
		GroupId:     resp.GroupId,
		BusChannel:  resp.BusChannel,
		UnPushLevel: resp.UnPushLevel,
	}, nil
}

// UltraGroupCreate 创建超级群
func (rc *RongCloud) UltraGroupCreate(userId, groupId, groupName string) error {
	if userId == "" {
		return RCErrorNew(1002, "param 'userId' is empty")
	}

	if groupId == "" {
		return RCErrorNew(1002, "param 'groupId' is empty")
	}

	if groupName == "" {
		return RCErrorNew(1002, "param 'groupName' is empty")
	}

	req := httplib.Post(rc.rongCloudURI + "/ultragroup/create.json")

	req.SetTimeout(rc.timeout*time.Second, rc.timeout*time.Second)
	rc.fillHeader(req)

	req.Param("userId", userId)
	req.Param("groupId", groupId)
	req.Param("groupName", groupName)

	resp, err := rc.do(req)
	if err != nil {
		return err
	}

	data := struct {
		Code int `json:"code"`
	}{}

	if err = json.Unmarshal(resp, &data); err != nil {
		return err
	}

	if data.Code != 200 {
		return fmt.Errorf("response error. code: %d", data.Code)
	}

	return nil
}

// UltraGroupDis 解散超级群
func (rc *RongCloud) UltraGroupDis(groupId string) error {
	if groupId == "" {
		return RCErrorNew(1002, "param 'groupId' is empty")
	}

	req := httplib.Post(rc.rongCloudURI + "/ultragroup/dis.json")

	req.SetTimeout(rc.timeout*time.Second, rc.timeout*time.Second)
	rc.fillHeader(req)

	req.Param("groupId", groupId)

	resp, err := rc.do(req)
	if err != nil {
		return err
	}

	data := struct {
		Code int `json:"code"`
	}{}

	if err = json.Unmarshal(resp, &data); err != nil {
		return err
	}

	if data.Code != 200 {
		return fmt.Errorf("response error. code: %d", data.Code)
	}

	return nil
}

// UltraGroupJoin 加入超级群
func (rc *RongCloud) UltraGroupJoin(userId, groupId string) error {
	if userId == "" {
		return RCErrorNew(1002, "param 'userId' is empty")
	}

	if groupId == "" {
		return RCErrorNew(1002, "param 'groupId' is empty")
	}

	req := httplib.Post(rc.rongCloudURI + "/ultragroup/join.json")

	req.SetTimeout(rc.timeout*time.Second, rc.timeout*time.Second)
	rc.fillHeader(req)

	req.Param("userId", userId)
	req.Param("groupId", groupId)

	resp, err := rc.do(req)
	if err != nil {
		return err
	}

	data := struct {
		Code int `json:"code"`
	}{}

	if err = json.Unmarshal(resp, &data); err != nil {
		return err
	}

	if data.Code != 200 {
		return fmt.Errorf("response error. code: %d", data.Code)
	}

	return nil
}

// UltraGroupQuit 退出超级群
func (rc *RongCloud) UltraGroupQuit(userId, groupId string) error {
	if userId == "" {
		return RCErrorNew(1002, "param 'userId' is empty")
	}

	if groupId == "" {
		return RCErrorNew(1002, "param 'groupId' is empty")
	}

	req := httplib.Post(rc.rongCloudURI + "/ultragroup/quit.json")

	req.SetTimeout(rc.timeout*time.Second, rc.timeout*time.Second)
	rc.fillHeader(req)

	req.Param("userId", userId)
	req.Param("groupId", groupId)

	resp, err := rc.do(req)
	if err != nil {
		return err
	}

	data := struct {
		Code int `json:"code"`
	}{}

	if err = json.Unmarshal(resp, &data); err != nil {
		return err
	}

	if data.Code != 200 {
		return fmt.Errorf("response error. code: %d", data.Code)
	}

	return nil
}

// UltraGroupRefresh 刷新超级群信息
func (rc *RongCloud) UltraGroupRefresh(groupId, groupName string) error {
	if groupId == "" {
		return RCErrorNew(1002, "param 'groupId' is empty")
	}

	if groupName == "" {
		return RCErrorNew(1002, "param 'groupName' is empty")
	}

	req := httplib.Post(rc.rongCloudURI + "/ultragroup/refresh.json")

	req.SetTimeout(rc.timeout*time.Second, rc.timeout*time.Second)
	rc.fillHeader(req)

	req.Param("groupId", groupId)
	req.Param("groupName", groupName)

	resp, err := rc.do(req)
	if err != nil {
		return err
	}

	data := struct {
		Code int `json:"code"`
	}{}

	if err = json.Unmarshal(resp, &data); err != nil {
		return err
	}

	if data.Code != 200 {
		return fmt.Errorf("response error. code: %d", data.Code)
	}

	return nil
}

// UltraGroupUserBannedAdd 添加禁言成员
func (rc *RongCloud) UltraGroupUserBannedAdd(groupId, busChannel string, userIds ...string) error {
	if groupId == "" {
		return RCErrorNew(1002, "param 'groupId' is empty")
	}

	if userIds == nil || len(userIds) <= 0 {
		return RCErrorNew(1002, "param 'userIds' is empty")
	}

	if len(userIds) > 20 {
		return RCErrorNew(1002, "param 'userIds' is too long")
	}

	req := httplib.Post(rc.rongCloudURI + "/ultragroup/userbanned/add.json")

	req.SetTimeout(rc.timeout*time.Second, rc.timeout*time.Second)
	rc.fillHeader(req)

	req.Param("groupId", groupId)
	req.Param("userIds", strings.Join(userIds, ","))

	if busChannel != "" {
		req.Param("busChannel", busChannel)
	}

	resp, err := rc.do(req)
	if err != nil {
		return err
	}

	data := struct {
		Code int `json:"code"`
	}{}

	if err = json.Unmarshal(resp, &data); err != nil {
		return err
	}

	if data.Code != 200 {
		return fmt.Errorf("response error. code: %d", data.Code)
	}

	return nil
}

// UltraGroupUserBannedDel 移除禁言成员
func (rc *RongCloud) UltraGroupUserBannedDel(groupId, busChannel string, userIds ...string) error {
	if groupId == "" {
		return RCErrorNew(1002, "param 'groupId' is empty")
	}

	if userIds == nil || len(userIds) <= 0 {
		return RCErrorNew(1002, "param 'userIds' is empty")
	}

	if len(userIds) > 20 {
		return RCErrorNew(1002, "param 'userIds' is too long")
	}

	req := httplib.Post(rc.rongCloudURI + "/ultragroup/userbanned/del.json")

	req.SetTimeout(rc.timeout*time.Second, rc.timeout*time.Second)
	rc.fillHeader(req)

	req.Param("groupId", groupId)
	req.Param("userIds", strings.Join(userIds, ","))

	if busChannel != "" {
		req.Param("busChannel", busChannel)
	}

	resp, err := rc.do(req)
	if err != nil {
		return err
	}

	data := struct {
		Code int `json:"code"`
	}{}

	if err = json.Unmarshal(resp, &data); err != nil {
		return err
	}

	if data.Code != 200 {
		return fmt.Errorf("response error. code: %d", data.Code)
	}

	return nil
}

type UltraGroupUserBannedResponseItem struct {
	Id string `json:"id"`
}

// UltraGroupUserBannedGet 获取禁言成员
func (rc *RongCloud) UltraGroupUserBannedGet(groupId, busChannel string, page, pageSize int) ([]UltraGroupUserBannedResponseItem, error) {
	if groupId == "" {
		return nil, RCErrorNew(1002, "param 'groupId' is empty")
	}

	req := httplib.Post(rc.rongCloudURI + "/ultragroup/userbanned/get.json")

	req.SetTimeout(rc.timeout*time.Second, rc.timeout*time.Second)
	rc.fillHeader(req)

	req.Param("groupId", groupId)

	if busChannel != "" {
		req.Param("busChannel", busChannel)
	}

	if page != 0 {
		req.Param("page", strconv.Itoa(page))
	}

	if pageSize != 0 {
		req.Param("pageSize", strconv.Itoa(pageSize))
	}

	resp, err := rc.do(req)
	if err != nil {
		return nil, err
	}

	data := struct {
		Code  int                                `json:"code"`
		Users []UltraGroupUserBannedResponseItem `json:"users"`
	}{}

	if err = json.Unmarshal(resp, &data); err != nil {
		return nil, err
	}

	if data.Code != 200 {
		return nil, fmt.Errorf("response error. code: %d", data.Code)
	}

	return data.Users, nil
}

// UltraGroupGlobalBannedSet 设置超级群禁言状态
func (rc *RongCloud) UltraGroupGlobalBannedSet(groupId, busChannel string, status bool) error {
	if groupId == "" {
		return RCErrorNew(1002, "param 'groupId' is empty")
	}

	req := httplib.Post(rc.rongCloudURI + "/ultragroup/globalbanned/set.json")

	req.SetTimeout(rc.timeout*time.Second, rc.timeout*time.Second)
	rc.fillHeader(req)

	req.Param("groupId", groupId)
	req.Param("status", strconv.FormatBool(status))

	if busChannel != "" {
		req.Param("busChannel", busChannel)
	}

	resp, err := rc.do(req)
	if err != nil {
		return err
	}

	data := struct {
		Code int `json:"code"`
	}{}

	if err = json.Unmarshal(resp, &data); err != nil {
		return err
	}

	if data.Code != 200 {
		return fmt.Errorf("response error. code: %d", data.Code)
	}

	return nil
}

// UltraGroupGlobalBannedGet 查询超级群禁言状态
func (rc *RongCloud) UltraGroupGlobalBannedGet(groupId, busChannel string) (bool, error) {
	if groupId == "" {
		return false, RCErrorNew(1002, "param 'groupId' is empty")
	}

	req := httplib.Post(rc.rongCloudURI + "/ultragroup/globalbanned/get.json")

	req.SetTimeout(rc.timeout*time.Second, rc.timeout*time.Second)
	rc.fillHeader(req)

	req.Param("groupId", groupId)

	if busChannel != "" {
		req.Param("busChannel", busChannel)
	}

	resp, err := rc.do(req)
	if err != nil {
		return false, err
	}

	data := struct {
		Code   int  `json:"code"`
		Status bool `json:"status"`
	}{}

	if err = json.Unmarshal(resp, &data); err != nil {
		return false, err
	}

	if data.Code != 200 {
		return false, fmt.Errorf("response error. code: %d", data.Code)
	}

	return data.Status, nil
}

// UltraGroupBannedWhiteListAdd 添加禁言白名单
func (rc *RongCloud) UltraGroupBannedWhiteListAdd(groupId, busChannel string, userIds ...string) error {
	if groupId == "" {
		return RCErrorNew(1002, "param 'groupId' is empty")
	}

	if userIds == nil || len(userIds) <= 0 {
		return RCErrorNew(1002, "param 'userIds' is empty")
	}

	if len(userIds) > 20 {
		return RCErrorNew(1002, "param 'userIds' is too long")
	}

	req := httplib.Post(rc.rongCloudURI + "/ultragroup/banned/whitelist/add.json")

	req.SetTimeout(rc.timeout*time.Second, rc.timeout*time.Second)
	rc.fillHeader(req)

	req.Param("groupId", groupId)
	req.Param("userIds", strings.Join(userIds, ","))

	if busChannel != "" {
		req.Param("busChannel", busChannel)
	}

	resp, err := rc.do(req)
	if err != nil {
		return err
	}

	data := struct {
		Code int `json:"code"`
	}{}

	if err = json.Unmarshal(resp, &data); err != nil {
		return err
	}

	if data.Code != 200 {
		return fmt.Errorf("response error. code: %d", data.Code)
	}

	return nil
}

// UltraGroupBannedWhiteListDel 移除禁言白名单
func (rc *RongCloud) UltraGroupBannedWhiteListDel(groupId, busChannel string, userIds ...string) error {
	if groupId == "" {
		return RCErrorNew(1002, "param 'groupId' is empty")
	}

	if userIds == nil || len(userIds) <= 0 {
		return RCErrorNew(1002, "param 'userIds' is empty")
	}

	if len(userIds) > 20 {
		return RCErrorNew(1002, "param 'userIds' is too long")
	}

	req := httplib.Post(rc.rongCloudURI + "/ultragroup/banned/whitelist/del.json")

	req.SetTimeout(rc.timeout*time.Second, rc.timeout*time.Second)
	rc.fillHeader(req)

	req.Param("groupId", groupId)
	req.Param("userIds", strings.Join(userIds, ","))

	if busChannel != "" {
		req.Param("busChannel", busChannel)
	}

	resp, err := rc.do(req)
	if err != nil {
		return err
	}

	data := struct {
		Code int `json:"code"`
	}{}

	if err = json.Unmarshal(resp, &data); err != nil {
		return err
	}

	if data.Code != 200 {
		return fmt.Errorf("response error. code: %d", data.Code)
	}

	return nil
}

type UltraGroupBannedWhiteListGetResponseItem struct {
	Id string `json:"id"`
}

// UltraGroupBannedWhiteListGet 获取禁言白名单
func (rc *RongCloud) UltraGroupBannedWhiteListGet(groupId, busChannel string, page, pageSize int) ([]UltraGroupBannedWhiteListGetResponseItem, error) {
	if groupId == "" {
		return nil, RCErrorNew(1002, "param 'groupId' is empty")
	}

	req := httplib.Post(rc.rongCloudURI + "/ultragroup/banned/whitelist/get.json")

	req.SetTimeout(rc.timeout*time.Second, rc.timeout*time.Second)
	rc.fillHeader(req)

	req.Param("groupId", groupId)

	if busChannel != "" {
		req.Param("busChannel", busChannel)
	}

	if page != 0 {
		req.Param("page", strconv.Itoa(page))
	}

	if pageSize != 0 {
		req.Param("pageSize", strconv.Itoa(pageSize))
	}

	resp, err := rc.do(req)
	if err != nil {
		return nil, err
	}

	data := struct {
		Code  int                                        `json:"code"`
		Users []UltraGroupBannedWhiteListGetResponseItem `json:"users"`
	}{}

	if err = json.Unmarshal(resp, &data); err != nil {
		return nil, err
	}

	if data.Code != 200 {
		return nil, fmt.Errorf("response error. code: %d", data.Code)
	}

	return data.Users, nil
}

// UltraGroupChannelCreate 创建频道
func (rc *RongCloud) UltraGroupChannelCreate(groupId, busChannel string) error {
	if groupId == "" {
		return RCErrorNew(1002, "param 'groupId' is empty")
	}

	if busChannel == "" {
		return RCErrorNew(1002, "param 'busChannel' is empty")
	}

	req := httplib.Post(rc.rongCloudURI + "/ultragroup/channel/create.json")

	req.SetTimeout(rc.timeout*time.Second, rc.timeout*time.Second)
	rc.fillHeader(req)

	req.Param("groupId", groupId)
	req.Param("busChannel", busChannel)

	resp, err := rc.do(req)
	if err != nil {
		return err
	}

	data := struct {
		Code int `json:"code"`
	}{}

	if err = json.Unmarshal(resp, &data); err != nil {
		return err
	}

	if data.Code != 200 {
		return fmt.Errorf("response error. code: %d", data.Code)
	}

	return nil
}

// UltraGroupChannelDel 删除频道
func (rc *RongCloud) UltraGroupChannelDel(groupId, busChannel string) error {
	if groupId == "" {
		return RCErrorNew(1002, "param 'groupId' is empty")
	}

	if busChannel == "" {
		return RCErrorNew(1002, "param 'busChannel' is empty")
	}

	req := httplib.Post(rc.rongCloudURI + "/ultragroup/channel/del.json")

	req.SetTimeout(rc.timeout*time.Second, rc.timeout*time.Second)
	rc.fillHeader(req)

	req.Param("groupId", groupId)
	req.Param("busChannel", busChannel)

	resp, err := rc.do(req)
	if err != nil {
		return err
	}

	data := struct {
		Code int `json:"code"`
	}{}

	if err = json.Unmarshal(resp, &data); err != nil {
		return err
	}

	if data.Code != 200 {
		return fmt.Errorf("response error. code: %d", data.Code)
	}

	return nil
}

type UltraGroupChannelGetResponseItem struct {
	ChannelId  string `json:"channelId"`
	CreateTime string `json:"createTime"`
}

// UltraGroupChannelGet 查询频道列表
func (rc *RongCloud) UltraGroupChannelGet(groupId string, page, limit int) ([]UltraGroupChannelGetResponseItem, error) {
	if groupId == "" {
		return nil, RCErrorNew(1002, "param 'groupId' is empty")
	}

	req := httplib.Post(rc.rongCloudURI + "/ultragroup/channel/get.json")

	req.SetTimeout(rc.timeout*time.Second, rc.timeout*time.Second)
	rc.fillHeader(req)

	req.Param("groupId", groupId)

	if page != 0 {
		req.Param("page", strconv.Itoa(page))
	}

	if limit != 0 {
		req.Param("limit", strconv.Itoa(limit))
	}

	resp, err := rc.do(req)
	if err != nil {
		return nil, err
	}

	data := struct {
		Code     int                                `json:"code"`
		Channels []UltraGroupChannelGetResponseItem `json:"channelList"`
	}{}

	if err = json.Unmarshal(resp, &data); err != nil {
		return nil, err
	}

	if data.Code != 200 {
		return nil, fmt.Errorf("response error. code: %d", data.Code)
	}

	return data.Channels, nil
}
