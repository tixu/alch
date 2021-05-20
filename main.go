package alch

import (
	"fmt"
	"time"

	"github.com/andygrunwald/cachet"
	"github.com/goombaio/namegenerator"
	"github.com/tixu/alch/cachethq"
)

func main() {
	client := createClient("http://cachet:8000/", "7ABqTvQVOwOySelPvY91")
	instance := cachethq.Instance{
		client: client,
	}
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)
	//createComponent(client, nameGenerator.Generate())
	createGroup(client, nameGenerator.Generate())
	//createIncident(client, "201", "201201", 5, cachet.ComponentStatusMajorOutage, cachet.IncidentStatusIdentified)
	// Output: Beer Fridge
	// ID > 0!
	// Status: 200 OK
}

func createClient(url string, token string) *cachet.Client {
	client, _ := cachet.NewClient(url, nil)
	client.Authentication.SetTokenAuth(token)
	return client
}
func createIncident(client *cachet.Client, IncidentName string, IncidentMessage string, ComponentId int, ComponentStatus int, IncidentStatus int) {
	i := &cachet.Incident{
		Name:            IncidentName,
		Message:         IncidentMessage,
		ComponentID:     ComponentId,
		ComponentStatus: ComponentStatus,
		Status:          IncidentStatus,
		Visible:         cachet.IncidentVisibilityPublic,
	}
	fmt.Printf("incident %v", i)
	i, resp, err := client.Incidents.Create(i)
	if err != nil {
		fmt.Printf("Status: %s\n", resp.Status)
		fmt.Printf("incidents %v\n", i)
	}

}

func createComponent(client *cachet.Client, name string) {

	component := &cachet.Component{
		Name:        name,
		Description: "Description",
		Status:      cachet.ComponentStatusOperational,
	}
	component, resp, _ := client.Components.Create(component)
	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Component %v\n", component)
}

func createGroup(client *cachet.Client, name string) {

	group := &cachet.ComponentGroup{
		Name: name,
	}
	component, resp, _ := client.ComponentGroups.Create(group)
	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Component %v\n", component)
}
