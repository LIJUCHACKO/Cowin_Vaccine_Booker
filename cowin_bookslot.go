package main

import (
	"bufio"
	//	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	//	"net/url"
	"os"
	"strconv"
	"strings"

	"bytes"
	"crypto/sha256"
	"time"

	"github.com/dbatbold/beep"
)

var Sessiontoken string
var mobileNo string

type District struct {
	Id   int    `json:"district_id"`
	Name string `json:"district_name"`
}

type DistrictList struct {
	Entiities []District `json:"districts"`
	Ttl       int        `json:"ttl"`
}
type State struct {
	Id   int    `json:"state_id"`
	Name string `json:"state_name"`
}

type StateList struct {
	Entiities []State `json:"states"`
	Ttl       int     `json:"ttl"`
}

var Statelist StateList
var Districtlist DistrictList

func PrepareStatelist() {
repeat:
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://cdn-api.co-vin.in/api/v2/admin/location/states", nil)
	req.Header.Add("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8`)
	req.Header.Add("User-Agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11`)
	resp, err := client.Do(req)
	if err != nil {
		goto repeat //connection problem
		log.Fatalln(err)
	}

	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {

		log.Fatalln(err)
	}
	if resp.StatusCode != 200 {
		fmt.Printf("PrepareStatelist -Status code %d", resp.StatusCode)
		os.Exit(1)
	}
	//var data = []byte(`{"states":[{"state_id":1,"state_name":"Andaman and Nicobar Islands"},{"state_id":2,"state_name":"Andhra Pradesh"}]}`)
	//var cont StateList
	if err := json.Unmarshal(body, &Statelist); err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("%+v\n", cont)
	for _, eachstate := range Statelist.Entiities {
		//Statelist = append(Statelist, State{Id: eachstate.Id, Name: eachstate.Name})
		fmt.Printf("State id-%d, Statename-%s\n", eachstate.Id, eachstate.Name)

	}

}
func stateid(statename string) int {
	stateid := -1
	for _, eachstate := range Statelist.Entiities {
		//fmt.Printf("\nState id-%d, Statename-%s", eachstate.Id, eachstate.Name)
		if statename == eachstate.Name {
			stateid = eachstate.Id
		}

	}
	return stateid
	//fmt.Printf("%+v\n", cont)
}
func PrepareDistrictlist(stateid int) {
repeat:
	td := strconv.Itoa(stateid)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://cdn-api.co-vin.in/api/v2/admin/location/districts/"+td, nil)
	req.Header.Add("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8`)
	req.Header.Add("User-Agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11`)
	resp, err := client.Do(req)
	if err != nil {
		goto repeat //connection problem
		log.Fatalln(err)
	}

	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {

		log.Fatalln(err)
	}
	//fmt.Println(string(body))
	//body = []byte(`{"districts":[{"district_id":34,"district_name":"Ahmednagar"}],"ttl":24}`)
	//var cont DistrictList
	if err := json.Unmarshal(body, &Districtlist); err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		fmt.Printf("PrepareDistrictlist -Status code %d", resp.StatusCode)
		os.Exit(1)
	}
	//	fmt.Printf("%+v\n", Districtlist)
	for _, eachdistrict := range Districtlist.Entiities {
		//Districtlist = append(Districtlist, District{Id: eachdistrict.Id, Name: eachdistrict.Name})
		fmt.Printf("District id-%d, District name-%s\n", eachdistrict.Id, eachdistrict.Name)

	}

}
func districtid(districtname string) int {
	districtid := -1
	for _, eachdistrict := range Districtlist.Entiities {
		//fmt.Printf("\nState id-%d, Statename-%s", eachstate.Id, eachstate.Name)
		if districtname == eachdistrict.Name {
			districtid = eachdistrict.Id
		}

	}
	return districtid
	//fmt.Printf("%+v\n", cont)
}

type Session struct {
	Center_Id                int      `json:"center_id"`
	Name                     string   `json:"name"`
	Address                  string   `json:"address"`
	StateName                string   `json:"state_name"`
	DistrictName             string   `json:"district_name"`
	BlockName                string   `json:"block_name"`
	Pincode                  int      `json:"pincode"`
	From                     string   `json:"from"`
	To                       string   `json:"to"`
	Lat                      int      `json:"lat"`
	Long                     int      `json:"long"`
	Fee_type                 string   `json:"fee_type"`
	Session_id               string   `json:"session_id"`
	Date                     string   `json:"date"`
	Available_capacity       int      `json:"available_capacity"`
	Available_capacity_dose1 int      `json:"available_capacity_dose1"`
	Available_capacity_dose2 int      `json:"available_capacity_dose2"`
	Fee                      string   `json:"fee"`
	Min_age_limit            int      `json:"min_age_limit"`
	Vaccine                  string   `json:"vaccine"`
	Slots                    []string `json:"slots"`
}
type SessionList struct {
	Sessions []Session `json:"sessions"`
}

var Sessions SessionList

func PrepareSessionlist(pincodeorDistrictid string, date string) {
repeat:
	//fmt.Printf("\n PrepareSessionlist %s %s", pincode, date)
	pincodeorDistrictid = strings.TrimSpace(pincodeorDistrictid)
	client := &http.Client{}
	url := ""
	if len(pincodeorDistrictid) >= 6 {
		url = "https://cdn-api.co-vin.in/api/v2/appointment/sessions/public/findByPin?pincode=" + pincodeorDistrictid + "&date=" + date
	} else {
		url = "https://cdn-api.co-vin.in/api/v2/appointment/sessions/public/findByDistrict?district_id=" + pincodeorDistrictid + "&date=" + date
	}
	req, err := http.NewRequest("GET", url, nil)

	req.Header.Add("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8`)
	req.Header.Add("User-Agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11`)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		fmt.Println(" connection problem")
		goto repeat //connection problem
		log.Fatalln(err)
	}

	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//	fmt.Println("sessions")

	if err := json.Unmarshal(body, &Sessions); err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		fmt.Printf("PrepareSessionlist -Status code %d", resp.StatusCode)
		os.Exit(1)
	}
	//fmt.Printf("%+v\n", Sessions)
	// for _, eachsession := range Sessions.Sessions {
	// 	fmt.Printf("Center_Id-%d,  name-%s\n", eachsession.Center_Id, eachsession.Name)

	// }

}

var session_ids []string

func indexFunc(slice []string, item string) int {
	for i := range slice {
		if strings.TrimSpace(slice[i]) == item {
			return i
		}
	}
	return -1
}

type Month int

const (
	January Month = 1 + iota
	February
	March
	April
	May
	June
	July
	August
	September
	October
	November
	December
)

func play_music() {
	music := beep.NewMusic("") // output can be a file "music.wav"
	volume := 100

	if err := beep.OpenSoundDevice("default"); err != nil {
		return
		log.Fatal(err)
	}
	if err := beep.InitSoundDevice(); err != nil {
		return
		log.Fatal(err)
	}
	//beep.PrintSheet = true
	musicScore := `
	        VP SA8 SR9
	        A9HRDE cc DScszs|DEc DQzDE[|cc DScszs|DEc DQz DE[|vv DSvcsc|DEvs ]v|cc DScszs|VN
	        A3HLDE [n z,    |cHRq HLz, |[n z,    |cHRq HLz,  |sl z,    |]m   pb|z, ]m    |

	        A9HRDE cz [c|ss DSsz]z|DEs] ps|DSsz][ z][p|DEpDQ[ [|VN
	        A3HLDE [n ov|]m [n    |  pb ic|  n,   lHRq|HLnc DQ[|
	    `

	reader := bufio.NewReader(strings.NewReader(musicScore))
	defer beep.CloseSoundDevice()
	//fmt.Println("playing music")
	go music.Play(reader, volume)
	music.Wait()
	beep.FlushSoundBuffer()
}

type generateOTPtnx struct {
	TnxID string `json:"txnId"`
}

func generateOTP(mobile string) string {
repeat:
	var jsonStr = []byte(`{"mobile":"` + mobile + `","secret":"U2FsdGVkX1+z/4Nr9nta+2DrVJSv7KS6VoQUSQ1ZXYDx/CJUkWxFYG6P3iM/VW+6jLQ9RDQVzp/RcZ8kbT41xw==" }`)

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://cdn-api.co-vin.in/api/v2/auth/generateMobileOTP", bytes.NewBuffer(jsonStr))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "Mozilla/5.0")
	req.Header.Add("origin", "https://selfregistration.cowin.gov.in/")
	req.Header.Add("referer", "https://selfregistration.cowin.gov.in/")
	//req.Header.Add("Authorization", `Bearer `)
	resp, err := client.Do(req)
	//fmt.Println("test")
	if err != nil {
		fmt.Println(err)
		fmt.Println(" connection problem")
		goto repeat //connection problem
		//log.Fatalln(err)
	}
	//fmt.Println("test2")
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var trxidjson generateOTPtnx
	fmt.Println(string(body))

	if resp.StatusCode != 200 {
		fmt.Printf("generateOTP -Status code %d", resp.StatusCode)
		os.Exit(1)
	}
	if err := json.Unmarshal(body, &trxidjson); err != nil {
		log.Fatal(err)
	}
	return trxidjson.TnxID

}

type ConfirmOTPtnx struct {
	Token        string `json:"token"`
	IsNewAccount string `json:"isNewAccount"`
}

func ConfirmOTP(OTP string, txnId string) string {
repeat:
	OTP = strings.TrimSpace(OTP)
	otp_shasum := fmt.Sprintf("%x", sha256.Sum256([]byte(OTP)))
	var jsonStr = []byte(`{"txnId":"` + txnId + `","otp":"` + otp_shasum + `" }`)

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://cdn-api.co-vin.in/api/v2/auth/validateMobileOtp", bytes.NewBuffer(jsonStr))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "Mozilla/5.0")
	req.Header.Add("origin", "https://selfregistration.cowin.gov.in/")
	req.Header.Add("referer", "https://selfregistration.cowin.gov.in/")
	//req.Header.Add("Authorization", `Bearer `)
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		goto repeat //connection problem
		//log.Fatalln(err)
	}

	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var trxidjson ConfirmOTPtnx

	if resp.StatusCode != 200 {
		fmt.Printf("generateOTP -Status code %d", resp.StatusCode)
		return ""
	} else {
		fmt.Println(string(body))
	}
	if err := json.Unmarshal(body, &trxidjson); err != nil {
		log.Fatal(err)
	}
	return trxidjson.Token

}

type Beneficiary struct {
	Reference_Id  string `json:"beneficiary_reference_id"`
	Name          string `json:"name"`
	Birth_Year    string `json:"birth_year"`
	Gender        string `json:"gender"`
	Mobile_number string `json:"mobile_number"`
	Photo_id_type string `json:"photo_id_type"`
}

type BeneficiaryList struct {
	Persons []Beneficiary `json:"beneficiaries"`
}

var BeneficiaryListData BeneficiaryList
var selectedBeneficiaryIds []Beneficiary
var requiredDoses []string

//var slotspriority []string{}

func GetBeneficiariesIds() {
	//data := url.Values{}
repeat:
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://cdn-api.co-vin.in/api/v2/appointment/beneficiaries", nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "Mozilla/5.0")
	req.Header.Add("origin", "https://selfregistration.cowin.gov.in/")
	req.Header.Add("referer", "https://selfregistration.cowin.gov.in/")
	req.Header.Add("Authorization", "Bearer "+Sessiontoken)
	resp, err := client.Do(req)

	if err != nil {
		goto repeat //connection problem
		log.Fatalln(err)
	}

	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(body))

	if resp.StatusCode == 401 {
		fmt.Printf("GetBeneficiariesIds Unauthenticated Access -Status code %d", resp.StatusCode)
		GetSessionToken()
		goto repeat
	}
	if resp.StatusCode != 200 {
		fmt.Printf("GetBeneficiariesIds -Status code %d", resp.StatusCode)
		os.Exit(1)
	}
	if err := json.Unmarshal(body, &BeneficiaryListData); err != nil {
		log.Fatal(err)
	}
	// var list []string
	// for _, eachbeneficary := range BeneficiaryListData.Persons {
	// 	fmt.Printf("\n %s -%s", eachbeneficary.Reference_Id, eachbeneficary.Name)
	// 	list = append(list, eachbeneficary.Reference_Id)
	// }
	// return trxidjson.Token

}
func SelectBeneficiaries() {
	for _, eachbeneficary := range BeneficiaryListData.Persons {
		fmt.Printf("\n %s -%s", eachbeneficary.Reference_Id, eachbeneficary.Name)
		fmt.Printf("\n Do you want to include-(y/n)")
		var selection string
		fmt.Scanln(&selection)
		if strings.TrimSpace(selection) == "y" {
			selectedBeneficiaryIds = append(selectedBeneficiaryIds, eachbeneficary)
			fmt.Printf("\n Dose no ")
			var selection2 string
			fmt.Scanln(&selection2)
			requiredDose := "1"
			if strings.TrimSpace(selection) == "2" {
				requiredDose = "2"
			}
			requiredDoses = append(requiredDoses, requiredDose)
		}
	}

	fmt.Printf("\n Selected beneficiaries \n")
	for index, item := range selectedBeneficiaryIds {
		fmt.Printf(" %s -Dose %s", item.Name, requiredDoses[index])
	}
}
func GetSessionToken() {
	Sessiontoken = ""
	first := true
	for len(Sessiontoken) == 0 {
		trx_id := generateOTP(mobileNo)
		if first {
			go play_music()
			first = false
		}

		fmt.Println("Enter OTP: ")
		var OTP string
		fmt.Scanln(&OTP)
		Sessiontoken = ConfirmOTP(OTP, trx_id)

	}

}

var BookedBeneficiaryIds []string

func bookSlotforperson(personid string, requiredDose string, sessionid string, slot string) int {
repeat:
	//fmt.Printf("\n func bookSlotforperson -%s", sessionid)
	var jsonStr = []byte(`{"dose":` + requiredDose + `,"session_id":"` + sessionid + `","slot":"` + slot + `","beneficiaries":["` + personid + `"]}`)
	fmt.Println(string(jsonStr))
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://cdn-api.co-vin.in/api/v2/appointment/schedule", bytes.NewBuffer(jsonStr))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "Mozilla/5.0")
	req.Header.Add("origin", "https://selfregistration.cowin.gov.in/")
	req.Header.Add("referer", "https://selfregistration.cowin.gov.in/")
	req.Header.Add("Authorization", "Bearer "+Sessiontoken)
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		goto repeat //connection problem
		//log.Fatalln(err)
	}

	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	if resp.StatusCode != 409 && resp.StatusCode != 200 {
		fmt.Println(string(body))
		fmt.Printf("bookSlotforperson -Status code %d", resp.StatusCode)
	}

	if resp.StatusCode == 401 {
		fmt.Printf("bookSlotforperson Unauthenticated Access -Status code %d", resp.StatusCode)
		GetSessionToken()
		goto repeat
	}
	if resp.StatusCode == 200 {
		fmt.Println(string(body))
	}
	return resp.StatusCode
}

var newSession []Session

func BookSlot() {
	booked := false
	for _, eachsession := range newSession {
		slotno := 0

		fmt.Printf("\nsession id -   %s(%s)", eachsession.Session_id, eachsession.Name)
		fmt.Println(eachsession.Slots)
		if len(eachsession.Slots) == 0 {
			fmt.Printf("\n no slots")
			continue
		}
		for indexno, eachbeneficary := range selectedBeneficiaryIds {
			now := time.Now()
			year, _, _ := now.Date()
			birthyear, err := strconv.Atoi(eachbeneficary.Birth_Year)
			if err != nil {
				log.Fatalln(err)
			}
			Beneficiaryage := year - birthyear

			if indexFunc(BookedBeneficiaryIds, eachbeneficary.Reference_Id) >= 0 {
				continue
			}
			if eachsession.Min_age_limit > Beneficiaryage {
				fmt.Printf("\n%s age not covered -%d Min_age_limit-%d", eachbeneficary.Name, Beneficiaryage, eachsession.Min_age_limit)
				continue
			}
			fmt.Printf("\nBooking for -%s age %d", eachbeneficary.Name, Beneficiaryage)
			status := 409
			for slotno < len(eachsession.Slots) && status == 409 {
				eachslot := eachsession.Slots[slotno]
				fmt.Printf("\nTrying slot -%s", eachslot)
				status = bookSlotforperson(eachbeneficary.Reference_Id, requiredDoses[indexno], eachsession.Session_id, eachslot)
				if status == 200 {
					BookedBeneficiaryIds = append(BookedBeneficiaryIds, eachbeneficary.Reference_Id)
					fmt.Printf("\nBooked for -%s  on %s slot%s", eachbeneficary.Name, eachsession.Date, eachslot)
					booked = true
					break
				}
				if status == 409 {
					fmt.Printf("\n %s(%s) slot %s fully booked ", eachsession.Name, eachsession.Session_id, eachslot)
					//break
					slotno = slotno + 1
				}
			}

		}

	}
	if booked {
		go play_music()
	}

}

var selectFee_type string
var vaccine_type string
var Centerids []string

func CheckSessionAvailabilty() bool {
	result := false
	newSession = newSession[:0]
	for _, eachsession := range Sessions.Sessions {
		proceed := true
		if len(selectFee_type) > 0 {
			if strings.ToLower(eachsession.Fee_type) != strings.TrimSpace(strings.ToLower(selectFee_type)) {
				proceed = false
			}

		}
		if len(Centerids) > 0 {
			s := strconv.Itoa(eachsession.Center_Id)
			if indexFunc(Centerids, s) == -1 {
				proceed = false
			}
		}
		if len(vaccine_type) > 0 {
			if strings.ToLower(vaccine_type) != strings.ToLower(eachsession.Vaccine) {
				proceed = false
			}
		}
		//if eachsession.Fee_type == "Free" {
		if proceed {
			if eachsession.Available_capacity > 0 {
				if indexFunc(session_ids, eachsession.Session_id) == -1 {
					session_ids = append(session_ids, eachsession.Session_id)
					newSession = append(newSession, eachsession)
					fmt.Printf("pincode -%d Center_Id-%d,  name-%s  date -%s\n", eachsession.Pincode, eachsession.Center_Id, eachsession.Name, eachsession.Date)
					result = true
				}
			}

		}

	}
	return result
}
func main() {
	fmt.Println("Enter mobile no: ")
	fmt.Scanln(&mobileNo)
	mobileNo = strings.TrimSpace(mobileNo)

	GetSessionToken()
	GetBeneficiariesIds()

	SelectBeneficiaries()
	now := time.Now()
	fmt.Println("\nCurrent Time in String: ", now.String())
	var after string
reenterfrom:
	fmt.Printf("\nSearch slots from today+(no of days) ")
	fmt.Scanln(&after)
	from, errf := strconv.Atoi(after)
	if errf != nil {
		fmt.Println(errf)
		goto reenterfrom
	}
reentertill:
	fmt.Printf("\nSearch slots till today+(no of days)  ")
	fmt.Scanln(&after)
	NoofDays, errt := strconv.Atoi(after)
	if errt != nil {
		fmt.Println(errt)
		goto reentertill
	}
	//NoofDays := 8

	//for next 7days
	var pinnos string
repeatselection:
	fmt.Println("Do you want to enter pincode(y/n): ")

	var selection string
	fmt.Scanln(&selection)
	if strings.TrimSpace(selection) == "y" {
		fmt.Println("Enter pincodes separated by comma: ")
		fmt.Scanln(&pinnos)
	} else {
		PrepareStatelist()
	reenter:
		fmt.Println("Enter State code: ")
		var selection2 string
		fmt.Scanln(&selection2)
		i, err := strconv.Atoi(selection2)
		if err != nil {
			fmt.Println(err)
			goto reenter
		}
		PrepareDistrictlist(i)
		fmt.Println("Enter District codes separated by comma: ")
		fmt.Scanln(&pinnos)

	}
	var Centerid_str string
	fmt.Println("Enter Center ids(if any particularly) separated by comma: ")
	fmt.Scanln(&Centerid_str)
	if len(strings.TrimSpace(Centerid_str)) > 0 {
		Centerids = strings.Split(Centerid_str, ",")
	}
	fmt.Println("Enter vaccine_type(COVAXIN/COVISHIELD)if any: ")
	fmt.Scanln(&vaccine_type)
	vaccine_type = strings.TrimSpace(vaccine_type)
	fmt.Println("Enter fee_type(Paid/Free/Any): ")
	fmt.Scanln(&selectFee_type)
	if selectFee_type == "Paid" || selectFee_type == "Free" {

	} else {
		selectFee_type = ""
	}

	pincodes := strings.Split(pinnos, ",")
	//pincodes := []string{"400094", "400088"}
	fmt.Println("Checking for pincodes/district_ids")
	fmt.Println(pincodes)
	fmt.Println("Checking for dates")
	for i := from; i < NoofDays; i++ {
		after := now.AddDate(0, 0, i)
		year, month, day := after.Date()
		Day := strconv.Itoa(day)
		Month := strconv.Itoa(int(month))
		Year := strconv.Itoa(year)
		fmt.Println(Day + "-" + Month + "-" + Year)
	}
	fmt.Println("----------------------------------------")
	if len(Centerids) > 0 {
		fmt.Println("Checking only for centerids")
		fmt.Println(Centerids)
	}
	if len(vaccine_type) > 0 {
		fmt.Printf("vaccine type %s\n", vaccine_type)
	} else {
		fmt.Printf("vaccine_type Any\n")
	}
	if len(selectFee_type) > 0 {
		fmt.Printf("Fee type %s\n", selectFee_type)
	} else {
		fmt.Printf("Fee type Free/Paid\n")
	}
	fmt.Println("Do you want to continue(y/n): ")
	var selection4 string
	fmt.Scanln(&selection4)
	if strings.TrimSpace(selection4) != "y" {
		goto repeatselection
	}
	for len(pincodes) > 0 { //while(1)
		now := time.Now()
		for _, pin := range pincodes {
			for i := from; i < NoofDays; i++ {
				//time.Sleep(5 * time.Second)
				after := now.AddDate(0, 0, i)
				year, month, day := after.Date()
				Day := strconv.Itoa(day)
				Month := strconv.Itoa(int(month))
				Year := strconv.Itoa(year)
				//fmt.Println(Day + "-" + Month + "-" + Year)
				PrepareSessionlist(pin, Day+"-"+Month+"-"+Year)
				if CheckSessionAvailabilty() {
					fmt.Println("Current Time in String: ", now.String())
					fmt.Println("----------------------------------------\n")
					BookSlot()
				}
			}
		}
		fmt.Printf("Last update: %s\r", now.String())
		time.Sleep(20 * time.Second)
	}

}
