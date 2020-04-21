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
const version = "0.0.1b"

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
	{2, "2_HiMaiVol"},
	{3, "3_LoMaiVol"},
	{4, "4_CharFail"},
	{5, "5_HiBatVol"},
	{6, "6_LoBatVol"},
	{7, "7_HiDCVolt"},
	{8, "8_LoDCVolt"},
	{9, "9_GrFault+"},
	{10, "10_GrFault-"},
	{11, "11_Spa1"},
	{12, "12_Spa2"},
	{13, "13_Spa3"},
	{14, "14_Spa4"},
	{15, "15_Spa5"},
	{16, "16_Spa6"},
	{17, "17_Spa7"},
	{18, "18_Spa8"},
	{19, "19_ChaCurLimInd"},
	{20, "20_BatCurLimInd"},
	{21, "21_HiCharCur"},
	{22, "22_HiBattCur"},
	{23, "23_HiTemp"},
	{24, "24_TempSensErr"},
	{25, "25_IntCommErr"},
	{26, "26_BatTestFail"},
	{27, "27_BatTestAbor"},
	{28, "28_HiBatTemp"},
	{29, "29_HiFloCur"},
	{30, "30_LonCharTime"},
	{31, "31_NoPowSupVol"},
	{32, "32_BatInOper"},
}

var mesPar = param2{
	{100, 1, "100_MainsVol"},
	{101, 1, "101_BatVol"},
	{102, 1, "102_LoadVol"},
	{103, 1, "103_CharCur"},
	{104, 2, "104_BatCur"},
	{105, 3, "105_AmbTemp"},
	{106, 3, "106_BatTemp"},
	{107, 1, "107_AnalSpa"},
	{108, 1, "108_Freq"},
	{109, 2, "109_EarFauImp"},
	{110, 4, "110_CharStat"},
	{111, 5, "111_RemCharTime"},
	{112, 5, "112_AhMeter"},
}

type respStruct struct {
	namePar  string
	valuePar string
}

var respons []respStruct

func main() {
	serialPort := flag.String("serial", "/dev/ttyRS485-1", "a string")
	serialSpeed := flag.Int("speed", 9600, "a int")
	slaveID := flag.Int("id", 4, "an int")
	//timeout := flag.Int("t", 3000, "an int")
	flag.Parse()

	//		printResult(resultsErr)
	resultsErr := readModbus(*serialPort, byte(*slaveID),2, *serialSpeed,30)
//	resultsErr := []byte{0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	//fmt.Println(len(resultsErr))
//	fmt.Println(resultsErr)
	printErrResult(resultsErr)

		resultsMes := readModbus(*serialPort, byte(*slaveID),100, *serialSpeed,13)
//	resultsMes := []byte{9, 70, 8, 224, 0, 0, 0, 0, 0, 25, 3, 231, 0, 1, 0, 0, 176, 3, 0, 0, 0, 0, 0, 100, 0, 0}
	//		printResult(resultsMes)
	//fmt.Println(len(resultsMes))
//	fmt.Println(resultsMes)
	printMesResult(resultsMes)
	//		os.Exit(0)
	printJson(respons)
	/*	60
		[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
		26
		[9 70 8 224 0 0 0 0 0 25 3 231 0 1 0 0 176 3 0 0 0 0 0 100 0 0]
	*/
	/*var ttt = []byte{0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 1, 244, 0, 226, 0, 225, 0, 225, 1, 244, 0, 230, 0, 229, 0, 229, 0, 243, 0, 0, 5, 220, 0, 100, 4, 86, 1, 244, 0, 219, 0, 6, 0, 1, 0, 1, 0, 219, 0, 19, 0, 2, 0, 1, 0, 219, 0, 43, 0, 4, 0, 4, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	var ttt = []byte{}
	printResult(ttt)*/
}
func readModbus(serialPort string, slaveID byte, startReg uint16, serialSpeed int, regQuan uint16) []byte {
	handler := modbus.NewRTUClientHandler(fmt.Sprint(serialPort))
	handler.BaudRate = serialSpeed
	handler.SlaveId = slaveID
	handler.Timeout = time.Duration(3000) * time.Millisecond

	defer handler.Close()
	client := modbus.NewClient(handler)

	res, err := client.ReadHoldingRegisters(startReg, regQuan)
	if err != nil {
		printError(err)
	}
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
