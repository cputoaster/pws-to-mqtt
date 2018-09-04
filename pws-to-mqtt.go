package main

import (
	"net/http"
	"fmt"
	"log"
	"github.com/eclipse/paho.mqtt.golang"
	"encoding/json"
	"strconv"
)

var client mqtt.Client;

type WeatherSensor struct {
	ID             string
	IndoorTempF    float64
	TempF          float64
	DewpPtF         float64
	WindchillF     float64
	IndoorHumidity int
	Humidity       int
	WindspeedMph   float64
	WindgustMph    float64
	WindDir        int
	AbsBaroMin     float64
	BaroMin        float64
	Rainin         float64
	DailyRainin    float64
	WeeklyRainin   float64
	MonthlyRainin  float64
	YearlyRainin   float64
	SolarRadiation float64
	UV             int
	DateUtc        string
}

func handler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	ws := new(WeatherSensor);
	ws.ID = q["ID"][0]
	ws.IndoorTempF, _ = strconv.ParseFloat(q["indoortempf"][0], 64)
	ws.TempF, _ = strconv.ParseFloat(q["tempf"][0], 64)
	ws.DewpPtF, _ = strconv.ParseFloat(q["dewptf"][0], 64)
	ws.WindchillF, _ = strconv.ParseFloat(q["windchillf"][0], 64)
	ws.IndoorHumidity, _ = strconv.Atoi(q["indoorhumidity"][0])
	ws.Humidity, _ = strconv.Atoi(q["humidity"][0])
	ws.WindspeedMph, _ = strconv.ParseFloat(q["windspeedmph"][0], 64)
	ws.WindgustMph, _ = strconv.ParseFloat(q["windgustmph"][0], 64)
	ws.WindDir, _ = strconv.Atoi(q["winddir"][0])
	ws.AbsBaroMin, _ = strconv.ParseFloat(q["absbaromin"][0], 64)
	ws.BaroMin, _ = strconv.ParseFloat(q["baromin"][0], 64)
	ws.Rainin, _ = strconv.ParseFloat(q["rainin"][0], 64)
	ws.DailyRainin, _ = strconv.ParseFloat(q["dailyrainin"][0], 64)
	ws.WeeklyRainin, _ = strconv.ParseFloat(q["weeklyrainin"][0], 64)
	ws.MonthlyRainin, _ = strconv.ParseFloat(q["monthlyrainin"][0], 64)
	ws.YearlyRainin, _ = strconv.ParseFloat(q["yearlyrainin"][0], 64)
	ws.SolarRadiation, _ = strconv.ParseFloat(q["solarradiation"][0], 64)
	ws.UV, _ = strconv.Atoi(q["UV"][0])
	ws.DateUtc = q["dateutc"][0]


	b, _ := json.Marshal(ws);
	fmt.Fprintf(w, "got: %s", string(b))

	if token := client.Publish("tele/FineOffset00/SENSOR", 0, false,
		string(b)); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
}

func main() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://shibuya:1883")
	opts.SetClientID("pws-to-mqtt")
	client = mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	client.IsConnected();
	log.Print("Starting web server")
	http.HandleFunc("/weatherstation/updateweatherstation.php", handler)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
