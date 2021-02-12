package config

import "github.com/templexxx/reedsolomon"

const K int = 6
const M int = 3
const W int = 1
const HOST_IP string = "127.0.0.1"
const ChunkSize int = 1024 * 1024 //1MB
const MaxBatchSize int = 100

type OPType int
type CMDType int

const (
	ListenPort           int = 8977
	NodeListenClientPort int = 8988
	NodeListenMSCMDPort  int = 8999
	CMDPort              int = 9000
	DataPort             int = 9090
	ACKPort              int = 10010
)

//DataNode操作
const (
	UPDT_REQ         OPType = 0
	MoveDataToRoot   OPType = 1 //内部发送数据
	SendDataToParity OPType = 2 //data发送给parity
	DDU              OPType = 3
)

//MS命令
const (
	DataDeltaUpdate CMDType = 0 // DDU(i < j)，发送命令给DataNode，使其转发更新数据给rootParity

)

type Strategy int

const (
	CAU Strategy = 0
)
const IP_BASE string = "192.168.1."
const MSIP string = IP_BASE + "99"

var DataNodeIPs = [K]string{IP_BASE + "100", IP_BASE + "101", IP_BASE + "102", IP_BASE + "103", IP_BASE + "104", IP_BASE + "105"}
var ParityNodeIPs = [M]string{IP_BASE + "106", IP_BASE + "107", IP_BASE + "108"}
var Rack0 = Rack{
	Nodes:        map[string]string{"0": IP_BASE + "100", "1": IP_BASE + "101", "2": IP_BASE + "102"},
	NodeNum:      3,
	CurUpdateNum: 0,
	Stripes:      map[int][]int{},
}
var Rack1 = Rack{
	Nodes:        map[string]string{"3": IP_BASE + "103", "4": IP_BASE + "104", "5": IP_BASE + "105"},
	NodeNum:      3,
	CurUpdateNum: 0,
	Stripes:      map[int][]int{},
}
var Rack2 = Rack{
	Nodes:        map[string]string{"0": IP_BASE + "106", "1": IP_BASE + "107", "2": IP_BASE + "108"},
	NodeNum:      3,
	CurUpdateNum: 0,
	Stripes:      map[int][]int{},
}

//传输数据格式
type TD struct {
	SendSize           int
	OPType             OPType
	StripeID           int
	DataChunkID        int
	UpdateParityID     int
	NumRecvChunkItem   int
	NumRecvChunkParity int
	PortNum            int
	NextIP             string
	SenderIP           string
	FromIP             string
	Buff               []byte
}

type ACKData struct {
	ChunkID int
}

//传输命令格式
type CMD struct {
	SendSize           int
	Type               CMDType
	StripeID           int
	DataChunkID        int
	UpdateParityID     int
	NumRecvChunkItem   int
	NumRecvChunkParity int
	PortNum            int
	NextIP             string
	FromIP             string
	ToIP               string
}
type UpdateReqData struct {
	OPType       OPType
	LocalChunkID int
}

type MetaInfo struct {
	StripeID        int
	DataChunkID     int
	ChunkStoreIndex int //chunkID
	RelatedParities []string
	ChunkIP         string
	DataNodeID      int
	RackID          int
}

type UpdateStripe struct {
	StripeID  int
	DataIDs   []int
	ParityIDs []int
}
type Rack struct {
	Nodes        map[string]string
	NodeNum      int
	CurUpdateNum int
	Stripes      map[int][]int
}

var RS *reedsolomon.RS

//获取数据块（chunkID）对应的IP
func GetRelatedParities(chunkID int) []string {
	var relatedParities []string = make([]string, RS.ParityNum)
	col := chunkID % RS.DataNum

	for i := 0; i < RS.ParityNum; i++ {
		if RS.GenMatrix[i*RS.DataNum+col] > 0 {
			relatedParities = append(relatedParities, ParityNodeIPs[i])
		}
	}
	return relatedParities
}
func getRackID(dataNodeID int) int {
	if dataNodeID < 3 {
		return 0
	} else {
		return 1
	}
}
