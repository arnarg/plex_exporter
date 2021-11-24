package plex

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/arnarg/plex_exporter/config"
	"github.com/arnarg/plex_exporter/plex/api"
	"github.com/imdario/mergo"
)

type Server struct {
	ID         string
	Name       string
	Version    string
	Platform   string
	BaseURL    string
	token      string
	httpClient *http.Client
	headers    map[string]string
}

const TestURI = "%s/identity"
const ServerInfoURI = "%s/media/providers"
const StatusURI = "%s/status/sessions"
const LibraryURI = "%s/library/sections"
const SectionURI = "%s/library/sections/%d/all"

func NewServer(c config.PlexServerConfig) (*Server, error) {
	server := &Server{
		BaseURL: c.BaseURL,
		token:   c.Token,
		headers: headers,
		httpClient: &http.Client{
			Timeout: time.Second * 10,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: c.Insecure},
			},
		},
	}
	server.headers["X-Plex-Token"] = c.Token

	serverInfo, err := server.GetServerInfo()
	if err != nil {
		return nil, err
	}

	server.ID = serverInfo.ID
	server.Name = serverInfo.Name
	server.Version = serverInfo.Version
	server.Platform = serverInfo.Platform

	return server, nil
}

func (s *Server) GetServerInfo() (*api.ServerInfoResponse, error) {
	serverInfoResponse := api.ServerInfoResponse{}

	body, err := s.get(fmt.Sprintf(ServerInfoURI, s.BaseURL))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &serverInfoResponse)
	if err != nil {
		return nil, err
	}

	return &serverInfoResponse, nil
}

func (s *Server) GetSessionCount() (int, error) {
	sessionList := api.SessionList{}

	body, err := s.get(fmt.Sprintf(StatusURI, s.BaseURL))
	if err != nil {
		return -1, err
	}

	err = json.Unmarshal(body, &sessionList)
	if err != nil {
		return -1, err
	}

	return sessionList.Size, nil
}

func (s *Server) GetLibrary() (*api.LibraryResponse, error) {
	libraryResponse := api.LibraryResponse{}

	body, err := s.get(fmt.Sprintf(LibraryURI, s.BaseURL))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &libraryResponse)
	if err != nil {
		return nil, err
	}

	return &libraryResponse, nil
}

func (s *Server) GetSectionSize(id int) (int, error) {
	// We don't want to get every item in the library section
	// these headers make sure we only get metadata
	eh := map[string]string{
		"X-Plex-Container-Start": "0",
		"X-Plex-Container-Size":  "0",
	}
	mergo.Merge(&eh, headers)

	sectionResponse := api.SectionResponse{}

	_, body, err := sendRequest("GET", fmt.Sprintf(SectionURI, s.BaseURL, id), eh, s.httpClient)
	if err != nil {
		return -1, err
	}

	err = json.Unmarshal(body, &sectionResponse)
	if err != nil {
		return -1, err
	}

	return sectionResponse.TotalSize, nil
}

func (s *Server) get(url string) ([]byte, error) {
	_, body, err := sendRequest("GET", url, s.headers, s.httpClient)
	return body, err
}

func (s *Server) head(url string) (*http.Response, error) {
	resp, _, err := sendRequest("HEAD", url, s.headers, s.httpClient)
	return resp, err
}
