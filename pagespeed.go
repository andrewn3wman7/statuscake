package statuscake

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"
)


type PageSpeed struct {
	ID              int           `json:"id"              querystring:"id" querystringoptions:"omitempty"`
	Name            string        `json:"name"            querystring:"name" querystringoptions:"omitempty"`
	Website_url     string        `json:"website_url"     querystring:"website_url" querystringoptions:"omitempty"`
	Location_iso    string        `json:"location_iso"    querystring:"location_iso" querystringoptions:"omitempty"`
	Checkrate       int           `json:"checkrate"       querystring:"checkrate" querystringoptions:"omitempty"`
	Location        string        `json:"location"        querystring:"location" querystringoptions:"omitempty"`
	ContactGroup  []string        `json:"contact_groups"  querystring:"contact_groups"`
	AlertSmaller    int           `json:"alert_smaller"   querystring:"alert_smaller"`
	AlertBigger     int           `json:"alert_bigger"    querystring:"alert_bigger"`
	AlertSlower     int           `json:"alert_slower"    querystring:"alert_slower"`
}

	
type Reply struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		ID            int      `json:"id"`
		Name          string   `json:"name"`
		WebsiteURL    string   `json:"website_url"`
		Location      string   `json:"location"`
		LocationIso   string   `json:"location_iso"`
		Checkrate     int      `json:"checkrate"`
		ContactGroups []string `json:"contact_groups"`
		AlertSmaller  int      `json:"alert_smaller"`
		AlertBigger   int      `json:"alert_bigger"`
		AlertSlower   int      `json:"alert_slower"`
		LatestStats   struct {
			LoadtimeMs  int     `json:"Loadtime_ms"`
			FilesizeKb  float64 `json:"Filesize_kb"`
			Requests    int     `json:"Requests"`
			HasIssue    bool    `json:"has_issue"`
			LatestIssue string  `json:"latest_issue"`
		} `json:"latest_stats"`
	} `json:"data"`
}

type ReplyAll struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    []struct {
		ID            int      `json:"ID"`
		Title         string   `json:"Title"`
		URL           string   `json:"URL"`
		Location      string   `json:"Location"`
		LocationISO   string   `json:"Location_ISO"`
		ContactGroups []string `json:"ContactGroups"`
		LatestStats   struct {
			LoadtimeMs int     `json:"Loadtime_ms"`
			FilesizeKb float64 `json:"Filesize_kb"`
			Requests   int     `json:"Requests"`
		} `json:"LatestStats"`
	} `json:"data"`
}

type pagespeeds struct {
	client apiClient
}

//NewPageSpeeds return a new pagespeeds
func NewPageSpeeds(c apiClient) PageSpeeds {
	return &pagespeeds{
		client: c,
	}
}

//PageSpeeds represent the actions done with the API
type PageSpeeds interface {
	All() (*ReplyAll, error)
	Detail(int) (*PageSpeed, error)
	Update(*PageSpeed) (*PageSpeed, error)
	Delete(ID int) error
	Create(*PageSpeed) (*PageSpeed, error)
}

//All return a list of all the pagespeed from the API
func (tt *pagespeeds) All() (*ReplyAll, error) {
	rawResponse, err := tt.client.get("/Pagespeed", nil)
	if err != nil {
		return nil, fmt.Errorf("Error getting StatusCake Pagespeed: %s", err.Error())
	}
	var getResponse *ReplyAll
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
func (tt *pagespeeds) Detail(id int) (*PageSpeed, error) {
	resp, err := tt.client.get("/Pagespeed", url.Values{"id": {fmt.Sprint(id)}})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var dr *Reply

	err = json.NewDecoder(resp.Body).Decode(&dr)
	if err != nil {
		return nil, err
	}

	return dr.pagespeed(), nil
}


//Update the API with params
func (tt *pagespeeds) Update(t *PageSpeed) (*PageSpeed, error) {
	fmt.Println(t.ToURLValuesPg())
	resp, err := tt.client.post("/Pagespeed/Update", t.ToURLValuesPg())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ur updatePageSpeedResponse
	err = json.NewDecoder(resp.Body).Decode(&ur)
	if err != nil {
		return nil, err
	}

	if !ur.Success {
		return nil, &updateError{Message: ur.Message}
	}

	t2 := *t

	return &t2, err
}


// Create pagespeed test with params
func (tt *pagespeeds) Create(t *PageSpeed) (*PageSpeed, error) {
	fmt.Println(t.ToURLValuesPg())
	resp, err := tt.client.post("/Pagespeed/Update", t.ToURLValuesPg())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ur createPageSpeedResponse
	err = json.NewDecoder(resp.Body).Decode(&ur)
	if err != nil {
		return nil, err
	}

	if !ur.Success {
		return nil, &updateError{Message: ur.Message}
	}

	t2 := *t
	t2.ID = ur.Data.NewId
	fmt.Println(ur.Data)

	return &t2, err
}


// ToURLValuesPg returns url.Values of all fields required to create/update a Test.
func (t PageSpeed) ToURLValuesPg() url.Values {
	values := make(url.Values)
	st := reflect.TypeOf(t)
	sv := reflect.ValueOf(t)
	for i := 0; i < st.NumField(); i++ {
		sf := st.Field(i)
		tag := sf.Tag.Get(queryStringTag)
		ft := sf.Type
		if ft.Name() == "" && ft.Kind() == reflect.Ptr {
			// Follow pointer.
			ft = ft.Elem()
		}

		v := sv.Field(i)
		options := sf.Tag.Get("querystringoptions")
		omit := options == "omitempty" && isEmptyValuePg(v)

		if tag != "" && !omit {
			values.Set(tag, valueToQueryStringValuePg(v))
		}
	}

	return values
}

func isEmptyValuePg(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}

	return false
}

func valueToQueryStringValuePg(v reflect.Value) string {
	if v.Type().Name() == "bool" {
		if v.Bool() {
			return "1"
		}

		return "0"
	}

	if v.Type().Kind() == reflect.Slice {
		if ss, ok := v.Interface().([]string); ok {
			return strings.Join(ss, ",")
		}
	}

	return fmt.Sprint(v)
}

func (d *Reply) pagespeed() *PageSpeed {
	c := &PageSpeed{
	}
	c.ID = d.Data.ID
	c.Name = d.Data.Name
	c.Website_url = d.Data.WebsiteURL
	c.Checkrate = d.Data.Checkrate
	c.AlertSmaller = d.Data.AlertSmaller
	c.AlertBigger= d.Data.AlertBigger
	c.AlertSlower= d.Data.AlertSlower
	c.Location_iso = d.Data.LocationIso
	c.ContactGroup = d.Data.ContactGroups
	return c
}