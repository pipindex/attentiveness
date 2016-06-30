package attentiveness

import (
	"attentiveness/models"

	// "github.com/apex/go-apex"
	"github.com/zabawaba99/firego"
	"log"
	"strconv"
	"sync"
	"time"
)

var (
	sampleRatePolesPerSecond = 1
	windowLengthSeconds      = 10
	polesPerWindow           = windowLengthSeconds / sampleRatePolesPerSecond
	wg                       sync.WaitGroup
)

func CalculateAverageAttentivenessForWebinar(firebaseURL string, firebaseToken string, timeNowMilliseconds int, webinarId string) (bool) {
    timeWindow := Models.TimeWindow{
        Start: timeNowMilliseconds - (sampleRatePolesPerSecond*1000)*windowLengthSeconds,
        End:   timeNowMilliseconds,
    }
    wg.Add(1)
    calculateAverageAttentivenessForWebinar(&wg, Models.FirebaseURL{firebaseURL, firebaseToken}, timeWindow, webinarId)

    wg.Wait()
    return true
}

func CalculateAverageAttentivenessForActiveWebinars(firebaseURL Models.FirebaseURL, activeWebinars []string, timeNowMilliseconds int) Models.Response {

	timeWindow := Models.TimeWindow{
		Start: timeNowMilliseconds - (sampleRatePolesPerSecond*1000)*windowLengthSeconds,
		End:   timeNowMilliseconds,
	}

	for _, webinarId := range activeWebinars {
		wg.Add(1)
		go calculateAverageAttentivenessForWebinar(&wg, firebaseURL, timeWindow, webinarId)
	}

	wg.Wait()

	return Models.Response{true, activeWebinars}
}

func calculateAverageAttentivenessForWebinar(wg *sync.WaitGroup, firebaseURL Models.FirebaseURL, timeWindow Models.TimeWindow, webinarId string) {

	defer wg.Done()

	firebasePollingRef := constructFirebasePollingRef(firebaseURL, webinarId)
	pollingUsers := getPollingUsers(firebasePollingRef)

	polesChan := make(chan map[string]map[string]interface{})

	for userId := range pollingUsers {
		firebasePollingUserRef := constructFirebasePollingUserRef(firebaseURL, webinarId, userId)
		go getPollingDataForUser(polesChan, firebasePollingUserRef, userId, timeWindow)
	}

	// For each user, as the data becomes ready, keep adding to total
	allAttentiveness := make([]Models.Attentiveness, 0, len(pollingUsers))

	for range pollingUsers {
		select {
		case poles := <-polesChan:
			userAttentiveness := calculateAttentivenessForPoles(polesPerWindow, poles)
			allAttentiveness = append(allAttentiveness, userAttentiveness)
		case <-time.After(time.Millisecond * 100000):
		}
	}

	averageAttentiveness := averageAttentiveness(allAttentiveness)

    log.Println("WebinarId:", webinarId, "TimeWindow:", timeWindow)
	log.Println("WebinarId:", webinarId, "Users:", len(pollingUsers))
	log.Println("WebinarId:", webinarId, "Calculated:", len(allAttentiveness))
	log.Println("WebinarId:", webinarId, "Average Attentiveness:", averageAttentiveness)

	firebaseAttentivenessRef := constructFirebaseAttentivenessRef(firebaseURL, webinarId)
	if err := setAttentivenessInFirebase(firebaseAttentivenessRef, timeWindow.End, averageAttentiveness); err != nil {
		panic(err)
	}
}

func GetActiveWebinars(firebaseURL Models.FirebaseURL) (activeWebinars []string) {
	var activeWebinarsMap map[string]string
	firebaseRunningWebinarsRef := constructFirebaseRef(firebaseURL, "/active_webinars")
	firebaseRunningWebinarsRef.Value(&activeWebinarsMap)

	activeWebinars = make([]string, len(activeWebinarsMap))

	//convert to array of webinar Ids
	i := 0
	for k := range activeWebinarsMap {
		activeWebinars[i] = k
		i++
	}

	log.Println("Active Webinars:", activeWebinars)
	return activeWebinars
}

func averageAttentiveness(attentivenessList []Models.Attentiveness) (averageAttentiveness Models.Attentiveness) {

	if attentivenessLen := float64(len(attentivenessList)); attentivenessLen > 0 {

		var activeRatioTotal, idleRatioTotal, hiddenRatioTotal, notPresentRatioTotal float64

		for _, attentiveness := range attentivenessList {
			activeRatioTotal += attentiveness.ActiveRatio
			idleRatioTotal += attentiveness.IdleRatio
			hiddenRatioTotal += attentiveness.HiddenRatio
			notPresentRatioTotal += attentiveness.NotPresentRatio
		}

		averageAttentiveness = Models.Attentiveness{
			ActiveRatio:     activeRatioTotal / attentivenessLen,
			IdleRatio:       idleRatioTotal / attentivenessLen,
			HiddenRatio:     hiddenRatioTotal / attentivenessLen,
			NotPresentRatio: notPresentRatioTotal / attentivenessLen,
		}
	} else {
		averageAttentiveness = Models.Attentiveness{0, 0, 0, 1}
	}

	return averageAttentiveness
}

func getPollingUsers(firebasePollingRef *firego.Firebase) (pollingUsers map[string]interface{}) {
	firebasePollingRef.Shallow(true)
	if err := firebasePollingRef.Value(&pollingUsers); err != nil {
		log.Println("error!")
		log.Fatal(err)
	}
	return pollingUsers
}

func getPollingDataForUser(polesChan chan map[string]map[string]interface{}, firebaseRef *firego.Firebase, userId string, timeWindow Models.TimeWindow) {
	log.Println("Getting Polling Data for User", userId)
	var poles map[string]map[string]interface{}
	if err := firebaseRef.OrderBy("timestamp").
		StartAt(strconv.Itoa(timeWindow.Start)).
		EndAt(strconv.Itoa(timeWindow.End)).
		Value(&poles); err != nil {
		log.Println("error!")
		log.Fatal(err)
	}
	polesChan <- poles
}

func calculateAttentivenessForPoles(totalPoles int, poles map[string]map[string]interface{}) Models.Attentiveness {
	poleCount := map[string]int{
		"active": 0,
		"idle":   0,
		"hidden": 0,
	}
	for _, pole := range poles {
		switch pole["status"] {
		case "active":
			poleCount["active"]++
		case "idle":
			poleCount["idle"]++
		case "hidden":
			poleCount["hidden"]++
		}
	}
	return Models.NewAttentiveness(totalPoles, poleCount["active"], poleCount["idle"], poleCount["hidden"])
}
