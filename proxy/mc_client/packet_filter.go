package mc_client

import (
	"Eulogist/core/minecraft/protocol/packet"
	RaknetConnection "Eulogist/core/raknet"
	"encoding/json"
	"fmt"
)

/*
数据包过滤器过滤来自 Minecraft 客户端的多个数据包，
然后并将过滤后的多个数据包抄送至网易租赁服。

如果需要，
将根据实际情况由本处的桥接直接发送回应。

writePacketsToServer 指代
用于向客户端抄送数据包的函数。

syncFunc 用于将数据同步到网易租赁服，
它会在每个数据包被过滤处理后执行一次。

返回的 []error 是一个列表，
分别对应 packets 中每一个数据包的处理成功情况
*/
func (m *MinecraftClient) FiltePacketsAndSendCopy(
	packets []RaknetConnection.MinecraftPacket,
	writePacketsToServer func(packets []RaknetConnection.MinecraftPacket),
	syncFunc func() error,
) []error {
	// 初始化
	sendCopy := make([]RaknetConnection.MinecraftPacket, 0)
	doNotSendCopy := make([]bool, len(packets))
	errResults := make([]error, len(packets))
	// 处理每个数据包
	for index, minecraftPacket := range packets {
		// 如果传入的数据包为空
		if minecraftPacket.Packet == nil {
			continue
		}
		// 根据数据包的类型进行不同的处理
		switch pk := minecraftPacket.Packet.(type) {
		case *packet.PyRpc:
			doNotSendCopy[index], errResults[index] = m.OnPyRpc(pk)
			if err := errResults[index]; err != nil {
				errResults[index] = fmt.Errorf("FiltePacketsAndSendCopy: %v", err)
			}
		case *packet.NeteaseJson:
			// 解码 pk.Data 为 JSON 格式
			var jsonMap map[string]any
			err := json.Unmarshal(pk.Data, &jsonMap)
			if err != nil {
				errResults[index] = fmt.Errorf("FiltePacketsAndSendCopy: %v", err)
				break
			}
			// Login UID 已由赞颂者在先前发送，
			// 所以此处不必重复发送
			if eventName, ok := jsonMap["eventName"].(string); ok {
				if eventName == "LOGIN_UID" {
					doNotSendCopy[index] = true
					break
				}
			}
			// 将 NetEase UID 修正为真实值
			if _, ok := jsonMap["uid"]; ok {
				jsonMap["uid"] = m.GetNeteaseUID()
			}
			// 将 JSON 重新编码到 pk.Data
			pk.Data, err = json.Marshal(jsonMap)
			if err != nil {
				errResults[index] = fmt.Errorf("FiltePacketsAndSendCopy: %v", err)
				break
			}
			// 要求该数据包需要经编码后发送
			packets[index].Bytes = nil
		default:
			// 默认情况下，我们需要将
			// 数据包同步到网易租赁服
		}
		// 同步数据到网易租赁服
		if err := syncFunc(); err != nil {
			errResults[index] = fmt.Errorf("FiltePacketsAndSendCopy: %v", err)
		}
	}
	// 抄送数据包
	for index, pk := range packets {
		if doNotSendCopy[index] {
			continue
		}
		sendCopy = append(sendCopy, pk)
	}
	writePacketsToServer(sendCopy)
	// 返回值
	return errResults
}
