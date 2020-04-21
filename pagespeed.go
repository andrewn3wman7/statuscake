package statuscake

import (
	"encoding/json"
	"fmt"
	"net/url"
	// "log"
	"strings"

	"github.com/google/go-querystring/query"
)


type PageSpeed struct {
	ID              int           `json:"ID"              url:"id,omitempty"`
	Name            string        `json:"Title"           url:"Title,omitempty"`
	Website_url     string        `json:"URL"             url:"URL,omitempty"`
	Location_iso    string        `json:"location_iso"    url:"location_iso,omitempty"`
	Checkrate       string        `json:"checkrate"       url:"checkrate"`
	Location        string        `json:"location"        url:"location"`
	ContactGroupsC  string         `                      url:"contact_groups,omitempty"`
	ContactGroups []string        `json:"contact_groups"`
	// AlertBigger     string        `json:"alert_bigger"    url:"alert_bigger"`
	// AlertSlower     string        `json:"alert_slower"    url:"alert_slower"`
	// AlertSmaller	string        `json:"alert_smaller"   url:"alert_smaller"`

}

type PageSpeedResponse struct {
    Success bool `json:"success"`
    Message string `json:"message"`
    PageSpeedList []*PageSpeed `json:"data"`
}

//PartialPageSpeed represent a pagespeed test creation or modification
type PartialPageSpeed struct {
    ID              int
    Name            string
	Website_url     string
    Location_iso    string
	Checkrate       string
    ContactGroupsC  string
	// AlertBigger     string
	// AlertSlower	    string
	// AlertSmaller	string

}

type createPageSpeed struct {
	ID              int                 `url:"id,omitempty"`
	Name            string              `url:"name"           json:"name"`
	Website_url     string              `url:"website_url"    json:"website_url"`
	Location_iso    string              `url:"location_iso"   json:"location_iso"`
	Checkrate       jsonNumberString    `url:"checkrate"      json:"checkrate"`
	ContactGroupsC  string              `url:"contact_groups" json:"contact_groups"`
}

type updatePagespeed struct {
	ID             int                 `url:"id"`
	Name           string              `url:"name"           json:"name"`
	Website_url	   string              `url:"website_url"    json:"website_url"`
	Location_iso   string              `url:"location_iso"   json:"location_iso"`
	Checkrate      jsonNumberString    `url:"checkrate"      json:"checkrate"`
	ContactGroupsC string              `url:"contact_groups" json:"contact_groups"`
}

type pagespeedCreateResponse struct {
	Success bool                          `json:"success"`
	Message interface{}                   `json:"message"`
	Data    *pagespeedCreateResponseData  `json:"data"`
}

type pagespeedCreateResponseData struct {
	NewId int                 `json:"new_id"`
}

type pagespeedUpdateRespoonse struct {
	Success bool              `json:"success"`
	Message interface{}       `json:"message"`
}

type pagespeeds struct {
	client apiClient
}

func (cps *createPageSpeed) fromPartial(p *PartialPageSpeed) {
	cps.ID = p.ID
	cps.Name = p.Name
	cps.Checkrate = jsonNumberString(p.Checkrate)
	cps.Website_url = p.Website_url
	cps.Location_iso = p.Location_iso
	cps.ContactGroupsC = p.ContactGroupsC
}

func (cps *createPageSpeed) toPartial(p *PartialPageSpeed) {
	p.ID = cps.ID
	p.Name = cps.Name
	p.Checkrate = string(cps.Checkrate)
	p.Location_iso = cps.Location_iso
	p.ContactGroupsC = cps.ContactGroupsC
}

func (ups *updatePagespeed) fromPartial(p *PartialPageSpeed) {
	ups.ID = p.ID
	ups.Name = p.Name
	ups.Checkrate = jsonNumberString(p.Checkrate)
	ups.Website_url = p.Website_url
	ups.Location_iso = p.Location_iso
	ups.ContactGroupsC = p.ContactGroupsC
}

func (ups *updatePagespeed) toPartial(p *PartialPageSpeed) {
	p.ID = ups.ID
	p.Name = ups.Name
	p.Checkrate = string(ups.Checkrate)
	p.Website_url = ups.Website_url
	p.Location_iso = ups.Location_iso
	p.ContactGroupsC = ups.ContactGroupsC
}

//NewPageSpeeds return a new pagespeeds
func NewPageSpeeds(c apiClient) PageSpeeds {
	return &pagespeeds{
		client: c,
	}
}

//PageSpeeds represent the actions done wit the API
type PageSpeeds interface {
	All() (*PageSpeedResponse, error)
	Detail(int) (*PageSpeed, error)
	Update(*PartialPageSpeed) (*PageSpeed, error)
	UpdatePartial(*PartialPageSpeed) (*PartialPageSpeed, error)
	Delete(ID int) error
	completePageSpeed(*PartialPageSpeed) (*PageSpeed, error)
	CreatePartial(*PartialPageSpeed) (*PartialPageSpeed, error)
	Create(*PartialPageSpeed) (*PageSpeed, error)
}

//Create the pagespeed with the data in s and return the PageSpeed created
func (tt *pagespeeds) Create(s *PartialPageSpeed) (*PageSpeed, error) {
	var err error
	s, err = tt.CreatePartial(s)
	if err != nil {
		return nil, err
	}
	return tt.completePageSpeed(s)
}

func (tt *pagespeeds) completePageSpeed(s *PartialPageSpeed) (*PageSpeed, error) {
	full, err := tt.Detail(((*s).ID))
	if err != nil {
		return nil, err
	}
	(*full).ContactGroups = strings.Split((*s).ContactGroupsC, ",")
	return full, nil
}

//CreatePartial create the pagespeed whith the data in s and return the PartialSsl created
func (tt *pagespeeds) CreatePartial(s *PartialPageSpeed) (*PartialPageSpeed, error) {
	(*s).ID = 0
	var v url.Values
	{
		cps := createPageSpeed{}
		cps.fromPartial(s)
		v, _ = query.Values(cps)
	}

	rawResponse, err := tt.client.post("/Pagespeed/Update", v)
	if err != nil {
		return nil, fmt.Errorf("Error creating StatusCake Pagespeed: %s", err.Error())
	}

	var createResponse pagespeedCreateResponse
	err = json.NewDecoder(rawResponse.Body).Decode(&createResponse)
	if err != nil {
		return nil, err
	}

	if !createResponse.Success {
		return nil, fmt.Errorf("%s", createResponse.Message.(string))
	}

	(*s).ID = int(createResponse.Data.NewId)
	return s, nil
}

//All return a list of all the pagespeed from the API
func (tt *pagespeeds) All() (*PageSpeedResponse, error) {
	rawResponse, err := tt.client.get("/Pagespeed", nil)
	if err != nil {
		return nil, fmt.Errorf("Error getting StatusCake Pagespeed: %s", err.Error())
	}
	var getResponse *PageSpeedResponse
	err = json.NewDecoder(rawResponse.Body).Decode(&getResponse)
	if err != nil {
		return nil, err
	}
	return getResponse, err
}

//Delete delete the pagespeed which ID is id
func (tt *pagespeeds) Delete(id int) error {
	_, err := tt.client.delete("/Pagespeed/Update", url.Values{"id": {fmt.Sprint(id)}})
	if err != nil {
		return err
	}
	return nil
}


//Detail return the pagespeed corresponding to the ID
func (tt *pagespeeds) Detail(ID int) (*PageSpeed, error) {
	responses, err := tt.All()
	if err != nil {
		return nil, err
	}
	myPageSpeed, errF := findPageSpeed(responses, ID)
	if errF != nil {
		return nil, errF
	}
	return myPageSpeed, nil
}

//Update update the API with s and create one if s.ID=0 then return the corresponding Ssl
func (tt *pagespeeds) Update(s *PartialPageSpeed) (*PageSpeed, error) {
	var err error
	s, err = tt.UpdatePartial(s)
	if err != nil {
		return nil, err
	}
	return tt.completePageSpeed(s)
}

//UpdatePartial update the API with s and create one if s.ID=0 then return the corresponding PartialSsl
func (tt *pagespeeds) UpdatePartial(s *PartialPageSpeed) (*PartialPageSpeed, error) {
	if (*s).ID == 0 {
		return tt.CreatePartial(s)
	}
	
	var v url.Values
	{
		us := updatePagespeed{}
		us.fromPartial(s)
		v, _ = query.Values(us)
	}

	rawResponse, err := tt.client.post("/Pagespeed/Update", v)
	if err != nil {
		return nil, fmt.Errorf("Error creating StatusCake PageSpeed: %s", err.Error())
	}

	var updateResponse pagespeedUpdateRespoonse
	err = json.NewDecoder(rawResponse.Body).Decode(&updateResponse)
	if err != nil {
		return nil, err
	}

	if !updateResponse.Success {
		return nil, fmt.Errorf("%s", updateResponse.Message.(string))
	}
	return s, nil
}

func findPageSpeed(responses *PageSpeedResponse, ID int) (*PageSpeed, error) {
	var response *PageSpeed
	for _, elem := range responses.PageSpeedList {
		if (*elem).ID == ID {
			return elem, nil
		}
	}
	return response, fmt.Errorf("%s Not found", ID)
}
