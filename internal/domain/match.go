package domain

// VlrMatch represents a match from vlr.gg
type VlrMatch struct {
	Id            int
	Name          string
	StartDate     string
	StartTime     string
	BestOf        int
	Teams         []Team
	PagePath      string
	EventPagePath string
}

// Match represents a match
type Match struct {
	Id               int    `json:"id" dynamodbav:"id"`               // inherit from vlrMatch
	Name             string `json:"matchName" dynamodbav:"matchName"` // inherit from vlrMatch
	StartDate        string `json:"startDate" dynamodbav:"startDate"` // inherit from vlrMatch
	StartTime        string `json:"startTime" dynamodbav:"startTime"` // inherit from vlrMatch
	BestOf           int    `json:"bestOf" dynamodbav:"bestOf"`       // inherit from vlrMatch
	Teams            []Team `json:"teams" dynamodbav:"teams"`         // inherit from vlrMatch
	PagePath         string `json:"pagePath" dynamodbav:"pagePath"`   // inherit from vlrMatch
	EventName        string `json:"eventName" dynamodbav:"eventName"`
	EventCountryFlag string `json:"eventCountryFlag" dynamodbav:"eventCountryFlag"`
}

// Team represents a team
type Team struct {
	Name string `json:"title" dynamodbav:"title"`
}

// NewMatch creates a new match
func NewMatch(m VlrMatch, e VlrEvent) Match {
	return Match{
		Id:               m.Id,
		Name:             m.Name,
		StartDate:        m.StartDate,
		StartTime:        m.StartTime,
		BestOf:           m.BestOf,
		Teams:            m.Teams,
		PagePath:         m.PagePath,
		EventName:        e.Name,
		EventCountryFlag: e.CountryFlag,
	}
}
