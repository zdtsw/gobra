package logic

// global variable definitions
var (
	version string = "beta"
	author  string = "Wen Zhou"
	project string
	Release bool
	// esIndex string = "bilbo"
	// r       *gin.Engine
)

var render = pageFiller{
	VersionPage:   version,
	ContactAuthor: author,
	EAProject:     project,
}

///////////////////////////////////data strcture /////////////////////////////////////////
type appResponse struct {
	Apps []struct {
		ID        string `json:"id"`
		Container struct {
			Docker struct {
				Image      string `json:"image"`
				Parameters []struct {
					Key   string `json:"key"`
					Value string `json:"value"`
				} `json:"parameters"`
			} `json:"docker"`
		} `json:"container"`
		Env struct {
			JAVAOPTS        string `json:"JAVA_OPTS"`
			PROJECTSEEDROOT string `json:"PROJECT_SEED_ROOT"`
			NAME            string `json:"CLUSTER_NAME"`
		} `json:"env"`
		Labels struct {
			SEEDROOT string `json:"SEED_ROOT"`
			IsTest   string `json:"IS_TEST_INSTANCE"`
			VHOST    string `json:"HAPROXY_0_VHOST"`
		} `json:"labels"`
		TasksRunning int `json:"TasksRunning"`
	} `json:"apps"`
}

type returnAppResp struct {
	Host    string
	URL     string
	Project string
	Live    int
}

type pageFiller struct {
	VersionPage   string
	ContactAuthor string
	EAProject     string
}
