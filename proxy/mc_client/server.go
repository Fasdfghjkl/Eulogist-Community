package mc_client

import (
	raknet_connection "Eulogist/core/raknet"
	"Eulogist/core/tools/skin_process"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"net"

	"Eulogist/core/standard/protocol"

	"github.com/sandertv/go-raknet"
)

// CreateListener 在 127.0.0.1:19132 上以 Raknet
// 协议侦听 Minecraft 客户端的连接，
// 这意味着您成功创建了一个 Minecraft 数据包代理服务器
func (m *MinecraftClient) CreateListener() error {
	// 创建一个 Raknet 监听器
	listener, err := raknet.Listen("127.0.0.1:19132")
	if err != nil {
		return fmt.Errorf("CreateListener: %v", err)
	}
	// 获取监听器的地址
	address, ok := listener.Addr().(*net.UDPAddr)
	if !ok {
		return fmt.Errorf("CreateListener: Failed to get address for listener")
	}
	// 设置 pong data
	listener.PongData([]byte(
		fmt.Sprintf(
			"MCPE;%v;%v;%v;%v;%v;%v;Gophertunnel;%v;%v;%v;%v;",
			"Eulogist", protocol.CurrentProtocol, protocol.CurrentVersion, "0", "1",
			listener.ID(), "Creative", 1, address.Port, address.Port,
		),
	))
	// 初始化变量
	m.listener = listener
	m.address = address
	m.connected = make(chan struct{}, 1)
	m.Raknet = raknet_connection.NewStandardRaknetWrapper()
	// 返回成功
	return nil
}

// WaitConnect 等待 Minecraft 客户端连接到服务器
func (m *MinecraftClient) WaitConnect() error {
	// 接受客户端连接
	conn, err := m.listener.Accept()
	if err != nil {
		return fmt.Errorf("WaitConnect: %v", err)
	}
	// 初始化变量
	serverKey, _ := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	m.SetConnection(conn, serverKey)
	m.connected <- struct{}{}
	// 返回成功
	return nil
}

// GetServerIP 获取服务器的 IP 地址
func (m *MinecraftClient) GetServerIP() string {
	return m.address.IP.String()
}

// GetServerPort 获取服务器的端口号
func (m *MinecraftClient) GetServerPort() int {
	return m.address.Port
}

// ...
func (m *MinecraftClient) GetNeteaseUID() string {
	return m.neteaseUID
}

// ...
func (m *MinecraftClient) SetNeteaseUID(neteaseUID string) {
	m.neteaseUID = neteaseUID
}

// ...
func (m *MinecraftClient) InitPlayerSkin() {
	m.playerSkin = &skin_process.Skin{}
}

// ...
func (m *MinecraftClient) GetPlayerSkin() *skin_process.Skin {
	return m.playerSkin
}

// ...
func (m *MinecraftClient) SetPlayerSkin(skin *skin_process.Skin) {
	m.playerSkin = skin
}

// ...
func (m *MinecraftClient) GetOutfitInfo() map[string]*int {
	return m.outfitInfo
}

// ...
func (m *MinecraftClient) SetOutfitInfo(outfitInfo map[string]*int) {
	m.outfitInfo = outfitInfo
}

// ...
func (m *MinecraftClient) GetEntityUniqueID() int64 {
	return m.entityUniqueID
}

// ...
func (m *MinecraftClient) SetEntityUniqueID(entityUniqueID int64) {
	m.entityUniqueID = entityUniqueID
}

// ...
func (m *MinecraftClient) GetEntityRuntimeID() uint64 {
	return m.entityRuntimeID
}

// ...
func (m *MinecraftClient) SetEntityRuntimeID(entityRuntimeID uint64) {
	m.entityRuntimeID = entityRuntimeID
}
