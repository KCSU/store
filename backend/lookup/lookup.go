package lookup

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/kcsu/store/config"
	"github.com/kcsu/store/model"
	"gorm.io/gorm"
)

type LookupUser struct {
	Cancelled  string `xml:"cancelled,attr"`
	Identifier struct {
		Text   string `xml:",chardata"`
		Scheme string `xml:"scheme,attr"`
	} `xml:"identifier"`
	DisplayName    string `xml:"displayName"`
	RegisteredName string `xml:"registeredName"`
	Surname        string `xml:"surname"`
	VisibleName    string `xml:"visibleName"`
	MisAffiliation string `xml:"misAffiliation"`
}

type Result struct {
	XMLName xml.Name     `xml:"result"`
	People  []LookupUser `xml:"people>person"`
}

// Fetch XML from an API endpoint
func GetXML(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/xml")
	res, err := http.DefaultClient.Do(req)
	return res, err
}

// Fetch users from the Lookup API
func GetPeople(baseURL *url.URL, group model.Group) ([]LookupUser, error) {
	var path string
	switch group.Type {
	case "inst":
		path = fmt.Sprintf("inst/%s/members", group.Lookup)
	case "group":
		path = fmt.Sprintf("group/%s/members", group.Lookup)
	default:
		return nil, errors.New("invalid group type")
	}
	requestUrl, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	u := baseURL.ResolveReference(requestUrl).String()
	response, err := GetXML(u)
	if err != nil {
		return nil, err
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var result Result
	if err := xml.Unmarshal(responseData, &result); err != nil {
		return nil, err
	}
	return result.People, nil
}

// Convert Lookup API Users to Group users
func GetGroupUsers(people []LookupUser) []model.GroupUser {
	users := make([]model.GroupUser, 0, len(people))
	for _, person := range people {
		if person.Identifier.Scheme != "crsid" {
			continue
		}
		email := person.Identifier.Text + "@cam.ac.uk"
		groupUser := model.GroupUser{UserEmail: email, IsManual: false}
		users = append(users, groupUser)
	}
	return users
}

func Run(c *config.Config, d *gorm.DB) error {
	// Setup API url
	lookupUrl, err := url.Parse(c.LookupApiUrl)
	if err != nil {
		return err
	}
	// Get a list of groups
	var groups []model.Group
	if err := d.Find(&groups).Error; err != nil {
		return err
	}
	for _, group := range groups {
		// Ignore manually assigned groups
		if group.Type == "manual" {
			continue
		}
		// Fetch users from lookup API
		people, err := GetPeople(lookupUrl, group)
		if err != nil {
			return err
		}
		// Convert lookup users -> group users
		users := GetGroupUsers(people)
		// Replace group users in database
		err = d.Model(&group).Association("GroupUsers").Replace(&users)
		if err != nil {
			return err
		}
	}
	return nil
}
