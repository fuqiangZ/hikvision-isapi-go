package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func MockFromFile() []byte {

	bs, err := ioutil.ReadFile("console.log")
	if err != nil {
		fmt.Println("read file error:", err)
		return nil
	}
	fmt.Println(string(bs[0:100]))
	nums := strings.Split(string(bs), ",")
	byteArr := make([]byte, len(nums))

	for i, numStr := range nums {
		var num byte
		fmt.Sscanf(numStr, "%d", &num) // 将数字字符串转换为字节
		byteArr[i] = num
	}
	// fmt.Println(string(b)) // 打印字节数组

	// arr := [4]byte{bs[0], bs[1], bs[2], bs[3]} // 将字节切片转换为字节数组
	// str := string(arr[:])                      // 将字节数组转换为字符串
	// fmt.Println(str)
	return byteArr
}

func Mockdata() string {
	return `---------------------------7e13971310878
	Content-Disposition: form-data; name="anpr.xml";filename="anpr.xml"
	Content-Type: application/xml; charset="UTF-8"
	Content-Length: 3686
	
	<?xml version="1.0" encoding="utf-8"?>
	<EventNotificationAlert version="2.0" xmlns="http://www.isapi.org/ver20/XMLSchema">
	<ipAddress>192.168.100.248</ipAddress>
	<ipv6Address>::</ipv6Address>
	<portNo>20480</portNo>
	<protocol>HTTP</protocol>
	<macAddress>44:a6:42:f7:3a:21</macAddress>
	<channelID>1</channelID>
	<dateTime>2034-01-04T15:57:03+08:00</dateTime>
	<activePostCount>17</activePostCount>
	<eventType>ANPR</eventType>
	<eventState>active</eventState>
	<eventDescription>ANPR</eventDescription>
	<channelName>IP CAPTURE CAMERA</channelName>
	<deviceID>123456789</deviceID>
	<ANPR>
	<country>59</country>
	<licensePlate>蓝苏A8HK52</licensePlate>
	<line>1</line>
	<direction>reverse</direction>
	<confidenceLevel>78</confidenceLevel>
	<plateType>92TypeCivil</plateType>
	<plateColor>blue</plateColor>
	<licenseBright>0</licenseBright>
	<pilotsafebelt>unknown</pilotsafebelt>
	<vicepilotsafebelt>unknown</vicepilotsafebelt>
	<pilotsunvisor>unknown</pilotsunvisor>
	<vicepilotsunvisor>unknown</vicepilotsunvisor>
	<envprosign>unknown</envprosign>
	<dangmark>unknown</dangmark>
	<uphone>unknown</uphone>
	<pendant>unknown</pendant>
	<tissueBox>unknown</tissueBox>
	<frontChild>unknown</frontChild>
	<label>unknown</label>
	<smoking>unknown</smoking>
	<decoration>unknown</decoration>
	<plateCharBelieve>99,99,99,99,99,99,99</plateCharBelieve>
	<speedLimit>0</speedLimit>
	<illegalInfo>
	<illegalCode>0</illegalCode>
	<illegalName>正常</illegalName>
	<illegalDescription>正常</illegalDescription>
	</illegalInfo>
	<vehicleType>buggy</vehicleType>
	<featurePicFileName>1</featurePicFileName>
	<detectDir>2</detectDir>
	<relaLaneDirectionType>0</relaLaneDirectionType>
	<detectType>3</detectType>
	<barrierGateCtrlType>0</barrierGateCtrlType>
	<alarmDataType>0</alarmDataType>
	<dwIllegalTime>0</dwIllegalTime>
	<vehicleInfo>
	<index>563</index>
	<vehicleType>1</vehicleType>
	<colorDepth>0</colorDepth>
	<color>unknown</color>
	<speed>0</speed>
	<length>0</length>
	<vehicleLogoRecog>0</vehicleLogoRecog>
	<vehileSubLogoRecog>0</vehileSubLogoRecog>
	<vehileModel>0</vehileModel>
	<CarWindowFeature>
	<tempPlate>unknown</tempPlate>
	<passCard>unknown</passCard>
	<carCard>unknown</carCard>
	</CarWindowFeature>
	<CarBodyFeature>
	<sparetire>unknown</sparetire>
	<rack>unknown</rack>
	<sunRoof>unknown</sunRoof>
	<words>unknown</words>
	</CarBodyFeature>
	<vehicleUseType>unknown</vehicleUseType>
	</vehicleInfo>
	<pictureInfoList>
	<pictureInfo>
	<fileName>licensePlatePicture.jpg</fileName>
	<type>licensePlatePicture</type>
	<dataType>0</dataType>
	</pictureInfo>
	<pictureInfo>
	<fileName>detectionPicture.jpg</fileName>
	<type>detectionPicture</type>
	<dataType>0</dataType>
	<plateRect>
	<X>375</X>
	<Y>645</Y>
	<width>227</width>
	<height>118</height>
	</plateRect>
	<vehicelRect>
	<X>241</X>
	<Y>86</Y>
	<width>547</width>
	<height>765</height>
	</vehicelRect>
	<PilotRect>
	<x>0</x>
	<y>0</y>
	<width>0</width>
	<height>0</height>
	</PilotRect>
	<VicepilotRect>
	<x>0</x>
	<y>0</y>
	<width>0</width>
	<height>0</height>
	</VicepilotRect>
	<VehicelWindowRect>
	<x>0</x>
	<y>0</y>
	<width>0</width>
	<height>0</height>
	</VehicelWindowRect>
	<capturePicSecurityCode></capturePicSecurityCode>
	</pictureInfo>
	</pictureInfoList>
	<listType>temporary</listType>
	</ANPR>
	<UUID>46dd0716-1dd2-11b2-821f-8ba17e6f56b7</UUID>
	<picNum>2</picNum>
	<pilotStandardSafebelt>unknown</pilotStandardSafebelt>
	<vicepilotStandardSafebelt>unknown</vicepilotStandardSafebelt>
	<VehicleGATInfo>
	<palteTypeByGAT>2</palteTypeByGAT>
	<plateColorByGAT>2</plateColorByGAT>
	<vehicleTypeByGAT>H20</vehicleTypeByGAT>
	<colorByGAT>K</colorByGAT>
	</VehicleGATInfo>
	</EventNotificationAlert>
	
	---------------------------7e13971310878
	Content-Disposition: form-data;name="licensePlatePicture.jpg";filename="licensePlatePicture.jpg"
	Content-Type: image/jpeg
	Content-Length: 78280
	
	efwefwefwef
	`
}
