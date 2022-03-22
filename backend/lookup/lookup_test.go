package lookup_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kcsu/store/config"
	"github.com/kcsu/store/lookup"
	mocks "github.com/kcsu/store/mocks/db"
	"github.com/kcsu/store/model"
	"github.com/stretchr/testify/suite"
)

type LookupSuite struct {
	suite.Suite
	response  string
	store     *mocks.GroupStore
	server    *httptest.Server
	lookup    *lookup.ApiLookup
	expectURL string
	numCalls  int
}

const twoUserResponse = `<result xmlns="http://www.lookup.cam.ac.uk"
		xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
		version="1.2"
		xsi:schemaLocation="http://www.lookup.cam.ac.uk https://www.lookup.cam.ac.uk/xsd/ibis-api-1.2.xsd"
	>
		<people>
			<person cancelled="false">
				<identifier scheme="crsid">kt123</identifier>
				<displayName>K. Thrace</displayName>
				<registeredName>K. Thrace</registeredName>
				<surname>Thrace</surname>
				<visibleName>K. Thrace</visibleName>
				<misAffiliation>student</misAffiliation>
			</person>
			<person cancelled="false">
				<identifier scheme="crsid">lr456</identifier>
				<displayName>L. Roslin</displayName>
				<registeredName>L. Roslin</registeredName>
				<surname>Roslin</surname>
				<visibleName>L. Roslin</visibleName>
				<misAffiliation>student</misAffiliation>
			</person>
		</people>
	</result>`

func (s *LookupSuite) SetupTest() {
	s.numCalls = 0
	s.store = new(mocks.GroupStore)
	s.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.numCalls += 1
		s.Equal(s.expectURL, r.URL.Path)
		s.Equal("application/xml", r.Header.Get("Accept"))
		s.Equal(http.MethodGet, r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(s.response))
	}))
	s.lookup = &lookup.ApiLookup{
		Url:   s.server.URL,
		Store: s.store,
	}
}

func (s *LookupSuite) TearDownTest() {
	s.server.Close()
}

func (s *LookupSuite) TestGetPeople() {
	s.response = `<result xmlns="http://www.lookup.cam.ac.uk" 
			xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" 
			version="1.2"
			xsi:schemaLocation="http://www.lookup.cam.ac.uk https://www.lookup.cam.ac.uk/xsd/ibis-api-1.2.xsd"
		>
			<people>
				<person cancelled="false">
					<identifier scheme="crsid">kt123</identifier>
					<displayName>K. Thrace</displayName>
					<registeredName>K. Thrace</registeredName>
					<surname>Thrace</surname>
					<visibleName>K. Thrace</visibleName>
					<misAffiliation>student</misAffiliation>
				</person>
			</people>
		</result>`
	group := model.Group{
		Name:   "My Group",
		Type:   "inst",
		Lookup: "MYGRP",
	}
	s.expectURL = "/inst/MYGRP/members"
	expectUsers := []lookup.LookupUser{
		{
			Cancelled: "false",
			Identifier: lookup.Identifier{
				Text:   "kt123",
				Scheme: "crsid",
			},
			DisplayName:    "K. Thrace",
			RegisteredName: "K. Thrace",
			Surname:        "Thrace",
			VisibleName:    "K. Thrace",
			MisAffiliation: "student",
		},
	}
	users, err := s.lookup.GetPeople(group)
	s.NoError(err)
	s.Equal(expectUsers, users)
	s.Equal(1, s.numCalls)
}

func (s *LookupSuite) TestGetGroupUsers() {
	people := []lookup.LookupUser{
		{
			Cancelled: "false",
			Identifier: lookup.Identifier{
				Text:   "kt123",
				Scheme: "crsid",
			},
			DisplayName:    "K. Thrace",
			RegisteredName: "K. Thrace",
			Surname:        "Thrace",
			VisibleName:    "K. Thrace",
			MisAffiliation: "student",
		},
		{
			Cancelled: "false",
			Identifier: lookup.Identifier{
				Text:   "lr456",
				Scheme: "crsid",
			},
			DisplayName:    "L. Roslin",
			RegisteredName: "L. Roslin",
			Surname:        "Roslin",
			VisibleName:    "L. Roslin",
			MisAffiliation: "student",
		},
	}
	users := s.lookup.GetGroupUsers(people)
	expectUsers := []model.GroupUser{
		{UserEmail: "kt123@cam.ac.uk"},
		{UserEmail: "lr456@cam.ac.uk"},
	}
	s.Equal(expectUsers, users)
}

func (s *LookupSuite) TestProcessGroup() {
	s.response = twoUserResponse
	group := model.Group{
		Name:   "My Group",
		Type:   "group",
		Lookup: "MYGRP2",
	}
	s.expectURL = "/group/MYGRP2/members"
	s.store.On("ReplaceLookupUsers", &group, []model.GroupUser{
		{UserEmail: "kt123@cam.ac.uk"},
		{UserEmail: "lr456@cam.ac.uk"},
	}).Return(nil).Once()
	err := s.lookup.ProcessGroup(group)
	s.NoError(err)
	s.store.AssertExpectations(s.T())
	s.Equal(1, s.numCalls)
}

func (s *LookupSuite) TestRun() {
	s.response = twoUserResponse
	group := model.Group{
		Name:   "My Group",
		Type:   "group",
		Lookup: "MYGRP",
	}
	s.expectURL = "/group/MYGRP/members"
	s.store.On("Get").Return([]model.Group{group}, nil).Once()
	s.store.On("ReplaceLookupUsers", &group, []model.GroupUser{
		{UserEmail: "kt123@cam.ac.uk"},
		{UserEmail: "lr456@cam.ac.uk"},
	}).Return(nil).Once()
	c := config.Config{
		LookupApiUrl: s.server.URL,
	}
	err := lookup.Run(&c, s.store)
	s.NoError(err)
	s.Equal(1, s.numCalls)
	s.store.AssertExpectations(s.T())
}

func TestLookupSuite(t *testing.T) {
	suite.Run(t, new(LookupSuite))
}
