package attentiveness

import (
	"attentiveness/models"
	"github.com/zabawaba99/firego"
	"log"
)

func constructFirebaseRef(firebaseURL Models.FirebaseURL, resourcePath string) (firebaseRef *firego.Firebase) {
	firebaseRef = firego.New(firebaseURL.URL+resourcePath, nil)
	firebaseRef.Auth(firebaseURL.AuthKey)
	return firebaseRef
}

func constructFirebasePollingRef(firebaseURL Models.FirebaseURL, webinarId string) (firebasePollingRef *firego.Firebase) {
	firebasePollingRef = constructFirebaseRef(firebaseURL, webinarId+"/polling")
	return firebasePollingRef
}

func constructFirebasePollingUserRef(firebaseURL Models.FirebaseURL, webinarId string, userId string) (firebasePollingUserRef *firego.Firebase) {
	firebasePollingUserRef = constructFirebaseRef(firebaseURL, webinarId+"/polling/"+userId)
	return firebasePollingUserRef
}

func constructFirebaseAttentivenessRef(firebaseURL Models.FirebaseURL, webinarId string) (firebaseAttentivenessRef *firego.Firebase) {
	firebaseAttentivenessRef = constructFirebaseRef(firebaseURL, webinarId+"/attentiveness/")
	return firebaseAttentivenessRef
}

func setAttentivenessInFirebase(firebaseAttentivenessRef *firego.Firebase, timestamp int, averageAttentiveness Models.Attentiveness) error {
	log.Println("Setting Attentiveness:", firebaseAttentivenessRef, timestamp)
	averageAttentivenessDatapoint := averageAttentiveness.ToMap()
	averageAttentivenessDatapoint["timestamp"] = float64(timestamp) //throw the timestamp into the map too
	if err := firebaseAttentivenessRef.Set(averageAttentivenessDatapoint); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
