package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/goburrow/modbus"
	"os"
	"time"
)

/*
!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
!!!!!!!!!!!! VERSION !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
*/
const version = "0.0.2"

/*
!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
!!!!!!!!!!!! VERSION !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
*/

type par1 struct {
	reg     uint16
	regName string
}
type par2 struct {
	reg     uint16
	regtype int
	regName string
}
type param1 []par1

type param2 []par2

var errPar = param1{
	{1, "1_HiMaiVol"},
	{2, "2_LoMaiVol"},
	{3, "3_CharFail"},
	{4, "4_HiBatVol"},
	{5, "5_LoBatVol"},
	{6, "6_HiDCVolt"},
	{7, "7_LoDCVolt"},
	{8, "8_GrFault+"},
	{9, "9_GrFault-"},
	{10, "10_Spa1"},
	{11, "11_Spa2"},
	{12, "12_Spa3"},
	{13, "13_Spa4"},
	{14, "14_Spa5"},
	{15, "15_Spa6"},
	{16, "16_Spa7"},
	{17, "17_Spa8"},
	{18, "18_ChaCurLimInd"},
	{19, "19_BatCurLimInd"},
	{20, "20_HiCharCur"},
	{21, "21_HiBattCur"},
	{22, "22_HiTemp"},
	{23, "23_TempSensErr"},
	{24, "24_IntCommErr"},
	{25, "25_BatTestFail"},
	{26, "26_BatTestAbor"},
	{27, "27_HiBatTemp"},
	{28, "28_HiFloCur"},
	{29, "29_LonCharTime"},
	{30, "30_NoPowSupVol"},
	{31, "31_BatInOper"},
	{32, "32_BatSymFault"},
}

var mesPar = param2{
	{99, 1, "99_MainsVol"},
	{100, 1, "100_BatVol"},
	{101, 1, "101_LoadVol"},
	{102, 1, "102_CharCur"},
	{103, 2, "103_BatCur"},
	{104, 3, "104_AmbTemp"},
	{105, 3, "105_BatTemp"},
	{106, 1, "106_BatSymVol"},
	{107, 1, "107_ComAlmRel"},
	{108, 2, "108_EarFauImp"},
	{109, 4, "109_CharStat"},
	{110, 5, "110_RemCharTime"},
	{111, 5, "111_AhMeter"},
}

type respStruct struct {
	namePar  string
	valuePar string
}

var respons []respStruct

func main() {
	serialPort := flag.String("serial", "/dev/ttyRS485-1", "a string")
	serialSpeed := flag.Int("speed", 9600, "a int")
	slaveID := flag.Int("id", 1, "an int")
	timeout := flag.Int("t", 3000, "an int mSec")
	requestType := flag.Bool("rtype", true, "a bool")
	addressIP := flag.String("ip", "localhost", "a string")
	tcpPort := flag.String("port", "502", "a string")
	flag.Parse()

	tcpServerParam := fmt.Sprint(*addressIP, ":", *tcpPort)

	if *requestType == true {
		resultsErr := readModbusTcp(tcpServerParam, byte(*slaveID), 1, 32, int16(*timeout))
		printErrResult(resultsErr)
		time.Sleep(500 * time.Millisecond)
		resultsMes := readModbusTcp(tcpServerParam, byte(*slaveID), 99, 13, int16(*timeout))
		printMesResult(resultsMes)
		printJson(respons)
	} else {
		resultsErr := readModbusSerial(*serialPort, byte(*slaveID), 1, *serialSpeed, 32, int16(*timeout))
		printErrResult(resultsErr)
		time.Sleep(500 * time.Millisecond)
		resultsMes := readModbusSerial(*serialPort, byte(*slaveID), 99, *serialSpeed, 13, int16(*timeout))
		printMesResult(resultsMes)
		printJson(respons)
	}
}

func readModbusSerial(serialPort string, slaveID byte, startReg uint16, serialSpeed int, regQuan uint16, timeout int16) []byte {
	handler := modbus.NewRTUClientHandler(fmt.Sprint(serialPort))
	handler.BaudRate = serialSpeed
	handler.SlaveId = slaveID
	handler.Timeout = time.Duration(timeout) * time.Millisecond

	defer handler.Close()
	client := modbus.NewClient(handler)

	res, err := client.ReadHoldingRegisters(startReg, regQuan)
	if err != nil {
		printError(err)
	}
	handler.Close()
	return res
}

func readModbusTcp(tcpServerParam string, slaveID byte, startReg uint16, regQuan uint16, timeout int16) []byte {
	handler := modbus.NewTCPClientHandler(tcpServerParam)
	handler.SlaveId = slaveID
	handler.Timeout = time.Duration(timeout) * time.Millisecond

	defer handler.Close()
	client := modbus.NewClient(handler)

	res, err := client.ReadHoldingRegisters(startReg, regQuan)
	if err != nil {
		printError(err)
	}
	handler.Close()
	return res
}

func printError(err error) {
	fmt.Printf("{\"status\":\"error\", \"error\":\"%s\", \"version\":\"%s\"} \n", err, version)
	os.Exit(1)
}
func printErrResult(data []byte) {
	if len(data) != 0 {
		for i := 0; i < (len(data) / 2); i++ {
			temp1 := errPar[i].regName
			data1 := binary.BigEndian.Uint16(data[i*2 : (i*2)+2])
			temp2 := fmt.Sprintf("%d", data1)
			r := respStruct{namePar: temp1, valuePar: temp2}
			respons = append(respons, r)
		}
		//fmt.Println(len(respons))
		//fmt.Println(respons)
	} else {
		fmt.Printf("{\"status\":\"error\", \"error\":\"lengt of data is 0\", \"version\":\"%s\"} \n", version)
		os.Exit(100)
	}

}
func printMesResult(data []byte) {
	if len(data) != 0 {
		for i := 0; i < (len(data) / 2); i++ {
			if mesPar[i].regtype == 4 {
				temp1 := mesPar[i].regName
				var temp2 string
				a := binary.BigEndian.Uint16(data[i*2 : (i*2)+2])
				switch a {
				case 0:
					temp2 = "\"Float\""
				case 1:
					temp2 = "\"Highrate\""
				case 2:
					temp2 = "\"Commissioning\""
				case 3:
					temp2 = "\"Bat_test\""
				case 4:
					temp2 = "\"Char_Off\""
				}
				r := respStruct{namePar: temp1, valuePar: temp2}
				respons = append(respons, r)
			} else if mesPar[i].regtype == 1 {
				temp1 := mesPar[i].regName
				data1 := float64(binary.BigEndian.Uint16(data[i*2:(i*2)+2])) / 10
				temp2 := fmt.Sprintf("%.2f", data1)
				r := respStruct{namePar: temp1, valuePar: temp2}
				respons = append(respons, r)
			} else if mesPar[i].regtype == 2 {
				temp1 := mesPar[i].regName
				i := binary.BigEndian.Uint16(data[i*2 : (i*2)+2])
				b := float64(int16(i)) / 10
				temp2 := fmt.Sprintf("%.2f", b)
				r := respStruct{namePar: temp1, valuePar: temp2}
				respons = append(respons, r)
			} else if mesPar[i].regtype == 3 {
				temp1 := mesPar[i].regName
				i := binary.BigEndian.Uint16(data[i*2 : (i*2)+2])
				b := float64(int16(i))
				temp2 := fmt.Sprintf("%.2f", b)
				r := respStruct{namePar: temp1, valuePar: temp2}
				respons = append(respons, r)
			} else {
				temp1 := mesPar[i].regName
				data1 := binary.BigEndian.Uint16(data[i*2 : (i*2)+2])
				temp2 := fmt.Sprintf("%d", data1)
				r := respStruct{namePar: temp1, valuePar: temp2}
				respons = append(respons, r)
			}
		}
	} else {
		fmt.Printf("{\"status\":\"error\", \"error\":\"lengt of data is 0\", \"version\":\"%s\"} \n", version)
		os.Exit(100)
	}

}

func printJson(data []respStruct) {
	for l := 0; l < len(data); l++ {
		if l == 0 {
			fmt.Printf("{")
		}
		fmt.Printf("\"%s\":", data[l].namePar)
		fmt.Printf("%s,", data[l].valuePar)
		if l == len(data)-1 {
			fmt.Printf("\"version\":\"%s\"}\n", version)
		}
	}
	os.Exit(0)
}

/* build for rapberry
env GOOS=linux GOARCH=arm GOARM=5 go build
*/
