package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

type Racetrack struct {
	RaceTrackType       string `json:"raceTrackType"`
	VehicleType         string `json:"vehicleType"`
	NoOfVehiclesAllowed int    `json:"noOfVehiclesAllowed"`
	CostPerHour         int    `json:"costPerHour"`
}

func main() {
	var cost int
	err := isProperTime()
	if err != nil {
		fmt.Errorf(err.Error())
		return
	}

	data, err := ioutil.ReadFile("racetrack-management.json")
	if err != nil {
		fmt.Println("error reading file: ", err.Error())
	}

	rt := []Racetrack{}
	err = json.Unmarshal(data, &rt)
	if err != nil {
		fmt.Println("error parsing racingtrack management json", err.Error())
	}

	for _, r := range rt {
		// race track has to be booked for minimum of 3 hrs per vehicle
		minCostPerVehicle := (r.CostPerHour) / (r.NoOfVehiclesAllowed)
		//fmt.Println("minCostPerVehicle: ", minCostPerVehicle)
		if minCostPerVehicle <= 3 {
			fmt.Println("race track " + r.RaceTrackType + " " + r.VehicleType + " booked for " +
				string(minCostPerVehicle) + " hrs per vehicle")
			continue
		}
		cost = cost + (r.CostPerHour * r.NoOfVehiclesAllowed)

		iEAT, err := isExceedingAllocatedTime()
		if err != nil {
			fmt.Errorf(err.Error())
		}
		if iEAT == true {
			cost = cost + r.CostPerHour
		}
	}

	irta, err := isRegularTrackAvailable()
	if err != nil {
		fmt.Errorf(err.Error())
	}
	if irta != true {
		cost += 50
	}

	fmt.Println("Total cost : ", cost)

}

// type struct to store current data like current time, regular track availability, extra time.
type CurrentData struct {
	CurrentTime              int  `json:"currentTime"`
	IsRegularTrackAvailable  bool `json:"isRegularTrackAvailable"`
	ExtraTimeBeyondAllocated int  `json:"extraTimeBeyondAllocated"`
}

// isProperTime is to check if the current time is between 1 and 5.
func isProperTime() error {
	data, err := ioutil.ReadFile("currentData.json")
	if err != nil {
		return err
	}

	var cd CurrentData
	err = json.Unmarshal(data, &cd)
	if err != nil {
		fmt.Println("error parsing json", err.Error())
	}
	currentTime := cd.CurrentTime
	fmt.Println("currentTime :", currentTime)

	if currentTime <= 1 && currentTime >= 5 {
		fmt.Println("Current Time is ", currentTime, " Cannot proceed now please provide currentTime from currentTime file")
		return errors.New("invalid time set")
	}
	return nil
}

// isRegularTrackAvailable is to check if there is regular track available.
func isRegularTrackAvailable() (bool, error) {
	data, err := ioutil.ReadFile("currentData.json")
	if err != nil {
		return false, err
	}

	cd := CurrentData{}
	err = json.Unmarshal(data, &cd)
	if err != nil {
		return false, err
	}
	if cd.IsRegularTrackAvailable == false {
		fmt.Println("regular track not available")
		return false, nil
	}
	return true, nil
}

// isExceedingAllocatedTime is to check if time taken is exceeding allocated time.
func isExceedingAllocatedTime() (bool, error) {
	data, err := ioutil.ReadFile("currentData.json")
	if err != nil {
		return false, err
	}

	cd := CurrentData{}
	err = json.Unmarshal(data, &cd)
	if err != nil {
		return false, err
	}
	if cd.ExtraTimeBeyondAllocated >= 15 {
		return true, nil
	}
	return false, nil
}
